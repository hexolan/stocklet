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

	"github.com/doug-martin/goqu/v9"
	_ "github.com/doug-martin/goqu/v9/dialect/postgres"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/hexolan/stocklet/internal/pkg/errors"
	pb "github.com/hexolan/stocklet/internal/pkg/protogen/shipping/v1"
	"github.com/hexolan/stocklet/internal/svc/shipping"
)

const (
	pgShipmentBaseQuery      string = "SELECT id, order_id, dispatched, created_at FROM shipments"
	pgShipmentItemsBaseQuery string = "SELECT shipment_id, product_id, quantity FROM shipment_items"
)

type postgresController struct {
	cl *pgxpool.Pool
}

func NewPostgresController(cl *pgxpool.Pool) shipping.StorageController {
	return postgresController{cl: cl}
}

func (c postgresController) GetShipment(ctx context.Context, shipmentId string) (*pb.Shipment, error) {
	return c.getShipment(ctx, nil, shipmentId)
}

func (c postgresController) getShipment(ctx context.Context, tx *pgx.Tx, shipmentId string) (*pb.Shipment, error) {
	// Determine if a db transaction is being used
	var row pgx.Row
	const query = pgShipmentBaseQuery + " WHERE id=$1"
	if tx == nil {
		row = c.cl.QueryRow(ctx, query, shipmentId)
	} else {
		row = (*tx).QueryRow(ctx, query, shipmentId)
	}

	// Scan row to protobuf obj
	shipment, err := scanRowToShipment(row)
	if err != nil {
		return nil, err
	}

	return shipment, nil
}

func (c postgresController) getShipmentByOrderId(ctx context.Context, tx *pgx.Tx, orderId string) (*pb.Shipment, error) {
	// Determine if a db transaction is being used
	var row pgx.Row
	const query = pgShipmentBaseQuery + " WHERE order_id=$1"
	if tx == nil {
		row = c.cl.QueryRow(ctx, query, orderId)
	} else {
		row = (*tx).QueryRow(ctx, query, orderId)
	}

	// Scan row to protobuf obj
	shipment, err := scanRowToShipment(row)
	if err != nil {
		return nil, err
	}

	return shipment, nil
}

func (c postgresController) GetShipmentItems(ctx context.Context, shipmentId string) ([]*pb.ShipmentItem, error) {
	return c.getShipmentItems(ctx, nil, shipmentId)
}

func (c postgresController) getShipmentItems(ctx context.Context, tx *pgx.Tx, shipmentId string) ([]*pb.ShipmentItem, error) {
	// Determine if transaction is being used
	var rows pgx.Rows
	var err error
	if tx == nil {
		rows, err = c.cl.Query(ctx, pgShipmentItemsBaseQuery+" WHERE shipment_id=$1", shipmentId)
	} else {
		rows, err = (*tx).Query(ctx, pgShipmentItemsBaseQuery+" WHERE shipment_id=$1", shipmentId)
	}

	if err != nil {
		return nil, errors.WrapServiceError(errors.ErrCodeService, "query error whilst fetching items", err)
	}

	shipmentItems := []*pb.ShipmentItem{}
	for rows.Next() {
		var shipmentItem pb.ShipmentItem
		shipmentItem.ShipmentId = shipmentId
		err := rows.Scan(
			&shipmentItem.ProductId,
			&shipmentItem.Quantity,
		)
		if err != nil {
			return nil, errors.WrapServiceError(errors.ErrCodeService, "failed to scan an order item", err)
		}

		shipmentItems = append(shipmentItems, &shipmentItem)
	}

	if rows.Err() != nil {
		return nil, errors.WrapServiceError(errors.ErrCodeService, "error whilst scanning order item rows", rows.Err())
	}

	return shipmentItems, nil
}

func (c postgresController) AllocateOrderShipment(ctx context.Context, orderId string, orderMetadata shipping.EventOrderMetadata, productQuantities map[string]int32) error {
	// Begin a DB transaction
	tx, err := c.cl.Begin(ctx)
	if err != nil {
		return errors.WrapServiceError(errors.ErrCodeExtService, "failed to begin transaction", err)
	}
	defer tx.Rollback(ctx)

	// Create shipment
	var shipmentId string
	err = tx.QueryRow(ctx, "INSERT INTO shipments (order_id) VALUES ($1) RETURNING id", orderId).Scan(&shipmentId)
	if err != nil {
		return errors.WrapServiceError(errors.ErrCodeExtService, "failed to create shipment", err)
	}

	// Add shipment items
	vals := [][]interface{}{}
	for productId, quantity := range productQuantities {
		vals = append(
			vals,
			goqu.Vals{shipmentId, productId, quantity},
		)
	}

	statement, args, err := goqu.Dialect("postgres").From("shipment_items").Insert().Cols("shipment_id", "product_id", "quantity").Vals(vals...).Prepared(true).ToSQL()
	if err != nil {
		return errors.WrapServiceError(errors.ErrCodeService, "failed to build statement", err)
	}

	_, err = tx.Exec(ctx, statement, args...)
	if err != nil {
		return errors.WrapServiceError(errors.ErrCodeExtService, "failed to add shipment items", err)
	}

	// Prepare and append shipment allocated event to transaction
	evt, evtTopic, err := shipping.PrepareShipmentAllocationEvent_Allocated(orderId, orderMetadata, shipmentId, productQuantities)
	if err != nil {
		return errors.WrapServiceError(errors.ErrCodeService, "failed to create event", err)
	}

	_, err = tx.Exec(ctx, "INSERT INTO event_outbox (aggregateid, aggregatetype, payload) VALUES ($1, $2, $3)", shipmentId, evtTopic, evt)
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

func (c postgresController) CancelOrderShipment(ctx context.Context, orderId string) error {
	// Begin a DB transaction
	tx, err := c.cl.Begin(ctx)
	if err != nil {
		return errors.WrapServiceError(errors.ErrCodeExtService, "failed to begin transaction", err)
	}
	defer tx.Rollback(ctx)

	// Get the shipment
	shipment, err := c.getShipmentByOrderId(ctx, &tx, orderId)
	if err != nil {
		return errors.WrapServiceError(errors.ErrCodeExtService, "failed to fetch shipment info", err)
	}

	// Get the shipment items
	shipmentItems, err := c.getShipmentItems(ctx, &tx, shipment.Id)
	if err != nil {
		return errors.WrapServiceError(errors.ErrCodeExtService, "failed to fetch shipment manifest", err)
	}

	// Delete shipment
	_, err = tx.Exec(ctx, "DELETE FROM shipments WHERE id=$1", shipment.Id)
	if err != nil {
		return errors.WrapServiceError(errors.ErrCodeExtService, "failed to delete shipment", err)
	}

	// Add the event to the outbox table with the transaction
	evt, evtTopic, err := shipping.PrepareShipmentAllocationEvent_AllocationReleased(orderId, shipment.Id, shipmentItems)
	if err != nil {
		return errors.WrapServiceError(errors.ErrCodeService, "failed to create event", err)
	}

	_, err = tx.Exec(ctx, "INSERT INTO event_outbox (aggregateid, aggregatetype, payload) VALUES ($1, $2, $3)", shipment.Id, evtTopic, evt)
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

// Scan a postgres row to a protobuf object
func scanRowToShipment(row pgx.Row) (*pb.Shipment, error) {
	var shipment pb.Shipment

	// Temporary variables that require conversion
	var tmpCreatedAt pgtype.Timestamp

	err := row.Scan(
		&shipment.Id,
		&shipment.OrderId,
		&shipment.Dispatched,
		&tmpCreatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, errors.WrapServiceError(errors.ErrCodeNotFound, "shipment not found", err)
		} else {
			return nil, errors.WrapServiceError(errors.ErrCodeExtService, "something went wrong scanning object", err)
		}
	}

	// convert postgres timestamps to unix format
	if tmpCreatedAt.Valid {
		shipment.CreatedAt = tmpCreatedAt.Time.Unix()
	} else {
		return nil, errors.NewServiceError(errors.ErrCodeUnknown, "something went wrong scanning object (timestamp conversion error)")
	}

	return &shipment, nil
}
