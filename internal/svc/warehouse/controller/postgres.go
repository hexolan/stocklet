// Copyright (C) 2024 Declan Teevan
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package controller

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/hexolan/stocklet/internal/pkg/errors"
	pb "github.com/hexolan/stocklet/internal/pkg/protogen/warehouse/v1"
	"github.com/hexolan/stocklet/internal/svc/warehouse"
)

const (
	pgProductStockBaseQuery     string = "SELECT product_id, quantity FROM product_stock"
	pgReservationBaseQuery      string = "SELECT id, order_id, created_at FROM reservations"
	pgReservationItemsBaseQuery string = "SELECT product_id, quantity FROM reservation_items"
)

type postgresController struct {
	cl *pgxpool.Pool
}

func NewPostgresController(cl *pgxpool.Pool) warehouse.StorageController {
	return postgresController{cl: cl}
}

func (c postgresController) GetProductStock(ctx context.Context, productId string) (*pb.ProductStock, error) {
	return c.getProductStock(ctx, nil, productId)
}

func (c postgresController) getProductStock(ctx context.Context, tx *pgx.Tx, productId string) (*pb.ProductStock, error) {
	// Determine if a db transaction is being used
	var row pgx.Row
	const query = pgProductStockBaseQuery + " WHERE product_id=$1"
	if tx == nil {
		row = c.cl.QueryRow(ctx, query, productId)
	} else {
		row = (*tx).QueryRow(ctx, query, productId)
	}

	// Scan row to protobuf obj
	stock, err := scanRowToProductStock(row)
	if err != nil {
		return nil, err
	}

	return stock, nil
}

func (c postgresController) GetReservation(ctx context.Context, reservationId string) (*pb.Reservation, error) {
	return c.getReservation(ctx, nil, reservationId)
}

func (c postgresController) getReservation(ctx context.Context, tx *pgx.Tx, reservationId string) (*pb.Reservation, error) {
	// Determine if a db transaction is being used
	var row pgx.Row
	const query = pgReservationBaseQuery + " WHERE id=$1"
	if tx == nil {
		row = c.cl.QueryRow(ctx, query, reservationId)
	} else {
		row = (*tx).QueryRow(ctx, query, reservationId)
	}

	// Scan row to protobuf obj
	reservation, err := scanRowToReservation(row)
	if err != nil {
		return nil, err
	}

	// Append items to reservation
	reservation, err = c.appendItemsToReservation(ctx, tx, reservation)
	if err != nil {
		return nil, err
	}

	return reservation, nil
}

func (c postgresController) getReservationByOrderId(ctx context.Context, tx *pgx.Tx, orderId string) (*pb.Reservation, error) {
	// Determine if a db transaction is being used
	var row pgx.Row
	const query = pgReservationBaseQuery + " WHERE order_id=$1"
	if tx == nil {
		row = c.cl.QueryRow(ctx, query, orderId)
	} else {
		row = (*tx).QueryRow(ctx, query, orderId)
	}

	// Scan row to protobuf obj
	reservation, err := scanRowToReservation(row)
	if err != nil {
		return nil, err
	}

	// Append items to reservation
	reservation, err = c.appendItemsToReservation(ctx, tx, reservation)
	if err != nil {
		return nil, err
	}

	return reservation, nil
}

func (c postgresController) CreateProductStock(ctx context.Context, productId string, startingQuantity int32) error {
	// Begin a DB transaction
	tx, err := c.cl.Begin(ctx)
	if err != nil {
		return errors.WrapServiceError(errors.ErrCodeExtService, "failed to begin transaction", err)
	}
	defer tx.Rollback(ctx)

	// Create stock
	_, err = tx.Exec(ctx, "INSERT INTO product_stock (product_id, quantity) VALUES ($1, $2)", productId, startingQuantity)
	if err != nil {
		return errors.WrapServiceError(errors.ErrCodeExtService, "failed to create product stock", err)
	}

	// Add the event to the outbox table with the transaction
	evt, evtTopic, err := warehouse.PrepareStockCreatedEvent(&pb.ProductStock{ProductId: productId, Quantity: startingQuantity})
	if err != nil {
		return errors.WrapServiceError(errors.ErrCodeService, "failed to create event", err)
	}

	_, err = tx.Exec(ctx, "INSERT INTO event_outbox (aggregateid, aggregatetype, payload) VALUES ($1, $2, $3)", productId, evtTopic, evt)
	if err != nil {
		return errors.WrapServiceError(errors.ErrCodeExtService, "failed to insert event", err)
	}

	// Commit the transaction
	err = tx.Commit(ctx)
	if err != nil {
		return errors.WrapServiceError(errors.ErrCodeExtService, "failed to commit transaction", err)
	}

	return nil
}

func (c postgresController) ReserveOrderStock(ctx context.Context, orderId string, orderMetadata warehouse.EventOrderMetadata, productQuantities map[string]int32) error {
	// Begin a DB transaction
	tx, err := c.cl.Begin(ctx)
	if err != nil {
		return errors.WrapServiceError(errors.ErrCodeExtService, "failed to begin transaction", err)
	}
	defer tx.Rollback(ctx)

	// Create reservation
	var reservationId string
	err = tx.QueryRow(ctx, "INSERT INTO reservations (order_id) VALUES ($1) RETURNING id", orderId).Scan(&reservationId)
	if err != nil {
		return errors.WrapServiceError(errors.ErrCodeExtService, "failed to create reservation", err)
	}

	// Reserve the items
	insufficientStockProductIds := []string{}
	for productId, quantity := range productQuantities {
		err = c.reserveStock(ctx, &tx, reservationId, productId, quantity)
		if err != nil {
			insufficientStockProductIds = append(insufficientStockProductIds, productId)
		}
	}

	// Ensure that all of the stock was reserved
	if len(insufficientStockProductIds) > 0 {
		// Add the event to the outbox table with the transaction
		evt, evtTopic, err := warehouse.PrepareStockReservationEvent_Failed(orderId, orderMetadata, insufficientStockProductIds)
		if err != nil {
			return errors.WrapServiceError(errors.ErrCodeService, "failed to create event", err)
		}

		_, err = c.cl.Exec(ctx, "INSERT INTO event_outbox (aggregateid, aggregatetype, payload) VALUES ($1, $2, $3)", orderId, evtTopic, evt)
		if err != nil {
			return errors.WrapServiceError(errors.ErrCodeExtService, "failed to insert event", err)
		}

		return nil
	}

	// Add the event to the outbox table with the transaction
	evt, evtTopic, err := warehouse.PrepareStockReservationEvent_Reserved(orderId, orderMetadata, reservationId, productQuantities)
	if err != nil {
		return errors.WrapServiceError(errors.ErrCodeService, "failed to create event", err)
	}

	_, err = tx.Exec(ctx, "INSERT INTO event_outbox (aggregateid, aggregatetype, payload) VALUES ($1, $2, $3)", reservationId, evtTopic, evt)
	if err != nil {
		return errors.WrapServiceError(errors.ErrCodeExtService, "failed to insert event", err)
	}

	// Commit the transaction
	err = tx.Commit(ctx)
	if err != nil {
		return errors.WrapServiceError(errors.ErrCodeExtService, "failed to commit transaction", err)
	}

	return nil
}

func (c postgresController) reserveStock(ctx context.Context, tx *pgx.Tx, reservationId string, productId string, quantity int32) error {
	// Determine if a transaction has already been provided
	var (
		funcTx pgx.Tx
		err    error
	)
	if tx != nil {
		funcTx = *tx
	} else {
		funcTx, err = c.cl.Begin(ctx)
		if err != nil {
			return errors.WrapServiceError(errors.ErrCodeExtService, "failed to begin transaction", err)
		}
		defer funcTx.Rollback(ctx)
	}

	// Subtract from quantity
	_, err = funcTx.Exec(ctx, "UPDATE product_stock SET quantity = quantity - $1 WHERE product_id=$2", quantity, productId)
	if err != nil {
		return errors.WrapServiceError(errors.ErrCodeExtService, "failed to update product stock", err)
	}

	// Get updated stock
	stock, err := c.getProductStock(ctx, tx, productId)
	if err != nil {
		return err
	}

	// Ensure that the stock is not negative
	if stock.Quantity < 0 {
		return errors.NewServiceError(errors.ErrCodeInvalidArgument, "insufficient stock")
	}

	// Add as reservation item
	_, err = funcTx.Exec(ctx, "INSERT INTO reservation_items (reservation_id, product_id, quantity) VALUES ($1, $2, $3)", reservationId, productId, quantity)
	if err != nil {
		return errors.WrapServiceError(errors.ErrCodeExtService, "failed to add as reservation item", err)
	}

	// Commit the transaction (if created in this func)
	if tx == nil {
		err = funcTx.Commit(ctx)
		if err != nil {
			return errors.WrapServiceError(errors.ErrCodeExtService, "failed to commit transaction", err)
		}
	}

	return nil
}

func (c postgresController) ReturnReservedOrderStock(ctx context.Context, orderId string) error {
	// Begin a DB transaction
	tx, err := c.cl.Begin(ctx)
	if err != nil {
		return errors.WrapServiceError(errors.ErrCodeExtService, "failed to begin transaction", err)
	}
	defer tx.Rollback(ctx)

	// Get the reservation
	reservation, err := c.getReservationByOrderId(ctx, &tx, orderId)
	if err != nil {
		return errors.WrapServiceError(errors.ErrCodeExtService, "failed to locate order reservation", err)
	}

	// Return all of the reserved stock
	for _, reservedStock := range reservation.ReservedStock {
		err = c.returnReservedStock(ctx, &tx, reservation.Id, reservedStock.ProductId, reservedStock.Quantity)
		if err != nil {
			return errors.WrapServiceError(errors.ErrCodeExtService, "failed to return reserved stock", err)
		}
	}

	// Delete the reservation
	_, err = tx.Exec(ctx, "DELETE FROM reservations WHERE id=$1", reservation.Id)
	if err != nil {
		return errors.WrapServiceError(errors.ErrCodeExtService, "failed to delete reservation", err)
	}

	// Prepare and add reservation consumed event to outbox
	evt, evtTopic, err := warehouse.PrepareStockReservationEvent_Returned(reservation.OrderId, reservation.Id, reservation.ReservedStock)
	if err != nil {
		return errors.WrapServiceError(errors.ErrCodeService, "failed to create event", err)
	}

	_, err = tx.Exec(ctx, "INSERT INTO event_outbox (aggregateid, aggregatetype, payload) VALUES ($1, $2, $3)", reservation.Id, evtTopic, evt)
	if err != nil {
		return errors.WrapServiceError(errors.ErrCodeExtService, "failed to insert event", err)
	}

	// Commit the transaction
	err = tx.Commit(ctx)
	if err != nil {
		return errors.WrapServiceError(errors.ErrCodeExtService, "failed to commit transaction", err)
	}

	return nil
}

func (c postgresController) returnReservedStock(ctx context.Context, tx *pgx.Tx, reservationId string, productId string, quantity int32) error {
	// Determine if a transaction has already been provided
	var (
		funcTx pgx.Tx
		err    error
	)
	if tx != nil {
		funcTx = *tx
	} else {
		funcTx, err = c.cl.Begin(ctx)
		if err != nil {
			return errors.WrapServiceError(errors.ErrCodeExtService, "failed to begin transaction", err)
		}
		defer funcTx.Rollback(ctx)
	}

	// Add back to stock quantity
	_, err = funcTx.Exec(ctx, "UPDATE product_stock SET quantity = quantity + $1 WHERE product_id=$2", quantity, productId)
	if err != nil {
		return errors.WrapServiceError(errors.ErrCodeExtService, "failed to update product stock", err)
	}

	// Delete reservation item
	_, err = funcTx.Exec(ctx, "DELETE FROM reservation_items WHERE reservation_id=$1 AND product_id=$2", reservationId, productId)
	if err != nil {
		return errors.WrapServiceError(errors.ErrCodeExtService, "failed to add as reservation item", err)
	}

	// Commit the transaction (if created in this func)
	if tx == nil {
		err = funcTx.Commit(ctx)
		if err != nil {
			return errors.WrapServiceError(errors.ErrCodeExtService, "failed to commit transaction", err)
		}
	}

	return nil
}

func (c postgresController) ConsumeReservedOrderStock(ctx context.Context, orderId string) error {
	// Begin a DB transaction
	tx, err := c.cl.Begin(ctx)
	if err != nil {
		return errors.WrapServiceError(errors.ErrCodeExtService, "failed to begin transaction", err)
	}
	defer tx.Rollback(ctx)

	// Get the reservation
	reservation, err := c.getReservationByOrderId(ctx, &tx, orderId)
	if err != nil {
		return errors.WrapServiceError(errors.ErrCodeExtService, "failed to locate order reservation", err)
	}

	// Delete the reservation
	_, err = tx.Exec(ctx, "DELETE FROM reservations WHERE id=$1", reservation.Id)
	if err != nil {
		return errors.WrapServiceError(errors.ErrCodeExtService, "failed to delete reservation", err)
	}

	// Dispatch stock removed events
	for _, reservedStock := range reservation.ReservedStock {
		evt, evtTopic, err := warehouse.PrepareStockRemovedEvent(reservedStock.ProductId, reservedStock.Quantity, &reservation.Id)
		if err != nil {
			return errors.WrapServiceError(errors.ErrCodeService, "failed to create event", err)
		}

		_, err = tx.Exec(ctx, "INSERT INTO event_outbox (aggregateid, aggregatetype, payload) VALUES ($1, $2, $3)", reservedStock.ProductId, evtTopic, evt)
		if err != nil {
			return errors.WrapServiceError(errors.ErrCodeExtService, "failed to insert event", err)
		}
	}

	// Prepare and add reservation consumed event to outbox
	evt, evtTopic, err := warehouse.PrepareStockReservationEvent_Consumed(reservation.OrderId, reservation.Id, reservation.ReservedStock)
	if err != nil {
		return errors.WrapServiceError(errors.ErrCodeService, "failed to create event", err)
	}

	_, err = tx.Exec(ctx, "INSERT INTO event_outbox (aggregateid, aggregatetype, payload) VALUES ($1, $2, $3)", reservation.Id, evtTopic, evt)
	if err != nil {
		return errors.WrapServiceError(errors.ErrCodeExtService, "failed to insert event", err)
	}

	// Commit the transaction
	err = tx.Commit(ctx)
	if err != nil {
		return errors.WrapServiceError(errors.ErrCodeExtService, "failed to commit transaction", err)
	}

	return nil
}

func (c postgresController) getReservationItems(ctx context.Context, tx *pgx.Tx, reservationId string) ([]*pb.ReservationStock, error) {
	// Determine if transaction is being used.
	var rows pgx.Rows
	var err error
	const query = pgReservationItemsBaseQuery + " WHERE reservation_id=$1"
	if tx == nil {
		rows, err = c.cl.Query(ctx, query, reservationId)
	} else {
		rows, err = (*tx).Query(ctx, query, reservationId)
	}

	if err != nil {
		return nil, errors.WrapServiceError(errors.ErrCodeService, "query error whilst fetching reserved items", err)
	}

	items := []*pb.ReservationStock{}
	for rows.Next() {
		var reservStock pb.ReservationStock

		err := rows.Scan(
			&reservStock.ProductId,
			&reservStock.Quantity,
		)
		if err != nil {
			return nil, errors.WrapServiceError(errors.ErrCodeService, "failed to scan a reservation item", err)
		}

		items = append(items, &reservStock)
	}

	if rows.Err() != nil {
		return nil, errors.WrapServiceError(errors.ErrCodeService, "error scanning item rows", rows.Err())
	}

	return items, nil
}

// Append items to the reservation
func (c postgresController) appendItemsToReservation(ctx context.Context, tx *pgx.Tx, reservation *pb.Reservation) (*pb.Reservation, error) {
	reservedItems, err := c.getReservationItems(ctx, tx, reservation.Id)
	if err != nil {
		return nil, err
	}

	// Add the items to the reservation protobuf
	reservation.ReservedStock = reservedItems

	return reservation, nil
}

// Scan a postgres row to a protobuf object
func scanRowToProductStock(row pgx.Row) (*pb.ProductStock, error) {
	var stock pb.ProductStock

	err := row.Scan(
		&stock.ProductId,
		&stock.Quantity,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, errors.WrapServiceError(errors.ErrCodeNotFound, "stock not found", err)
		} else {
			return nil, errors.WrapServiceError(errors.ErrCodeExtService, "something went wrong scanning object", err)
		}
	}

	return &stock, nil
}

// Scan a postgres row to a protobuf object
func scanRowToReservation(row pgx.Row) (*pb.Reservation, error) {
	var reservation pb.Reservation

	// Temporary variables that require conversion
	var tmpCreatedAt pgtype.Timestamp

	err := row.Scan(
		&reservation.Id,
		&reservation.OrderId,
		&tmpCreatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, errors.WrapServiceError(errors.ErrCodeNotFound, "reservation not found", err)
		} else {
			return nil, errors.WrapServiceError(errors.ErrCodeExtService, "something went wrong scanning object", err)
		}
	}

	// convert postgres timestamps to unix format
	if tmpCreatedAt.Valid {
		reservation.CreatedAt = tmpCreatedAt.Time.Unix()
	}

	return &reservation, nil
}
