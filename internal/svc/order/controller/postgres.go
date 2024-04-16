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
	"time"

	"github.com/doug-martin/goqu/v9"
	_ "github.com/doug-martin/goqu/v9/dialect/postgres"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/hexolan/stocklet/internal/pkg/errors"
	pb "github.com/hexolan/stocklet/internal/pkg/protogen/order/v1"
	"github.com/hexolan/stocklet/internal/svc/order"
)

const (
	pgOrderBaseQuery      string = "SELECT id, status, customer_id, shipping_id, transaction_id, created_at, updated_at FROM orders"
	pgOrderItemsBaseQuery string = "SELECT product_id, quantity FROM order_items"
)

// The postgres controller is responsible for implementing the StorageController interface
// to store and retrieve the requested items from the Postgres database.
//
// Other controllers can be implemented to interface with different database systems.
type postgresController struct {
	cl *pgxpool.Pool
}

// Creates a new postgresController that implements the StorageController interface.
func NewPostgresController(cl *pgxpool.Pool) order.StorageController {
	return postgresController{cl: cl}
}

// Internal method - Validation is assumed to have taken place already
// Gets an order by its specified id from the database
func (c postgresController) GetOrder(ctx context.Context, orderId string) (*pb.Order, error) {
	return c.getOrder(ctx, nil, orderId)
}

func (c postgresController) getOrder(ctx context.Context, tx *pgx.Tx, orderId string) (*pb.Order, error) {
	// Query order
	var row pgx.Row
	if tx == nil {
		row = c.cl.QueryRow(ctx, pgOrderBaseQuery+" WHERE id=$1", orderId)
	} else {
		row = (*tx).QueryRow(ctx, pgOrderBaseQuery+" WHERE id=$1", orderId)
	}

	// Scan row to protobuf order
	order, err := scanRowToOrder(row)
	if err != nil {
		return nil, err
	}

	// Append the order items
	order, err = c.appendItemsToOrderObj(ctx, tx, order)
	if err != nil {
		return nil, err
	}

	return order, nil
}

// Internal method - Inputs are assumed valid
// Create a new order in the database
func (c postgresController) CreateOrder(ctx context.Context, orderObj *pb.Order) (*pb.Order, error) {
	// Begin a DB transaction
	tx, err := c.cl.Begin(ctx)
	if err != nil {
		return nil, errors.WrapServiceError(errors.ErrCodeExtService, "failed to begin transaction", err)
	}
	defer tx.Rollback(ctx)

	// Prepare and perform insert query
	newOrder := pb.Order{
		Items:      orderObj.Items,
		Status:     orderObj.Status,
		CustomerId: orderObj.CustomerId,
		CreatedAt:  time.Now().Unix(),
	}
	err = tx.QueryRow(
		ctx,
		"INSERT INTO orders (status, customer_id) VALUES ($1, $2) RETURNING id",
		newOrder.Status,
		newOrder.CustomerId,
	).Scan(&newOrder.Id)
	if err != nil {
		return nil, errors.WrapServiceError(errors.ErrCodeExtService, "failed to create order", err)
	}

	// Create records for any order items
	err = c.createOrderItems(ctx, tx, newOrder.Id, newOrder.Items)
	if err != nil {
		// The deffered rollback will be called (so the transaction will not be commited)
		return nil, err
	}

	// Prepare a created event.
	//
	// Then add the event to the outbox table with the transaction
	// to ensure that the event will be dispatched if
	// the transaction succeeds.
	evt, evtTopic, err := order.PrepareOrderCreatedEvent(&newOrder)
	if err != nil {
		return nil, errors.WrapServiceError(errors.ErrCodeService, "failed to create order event", err)
	}

	_, err = tx.Exec(ctx, "INSERT INTO event_outbox (aggregateid, aggregatetype, payload) VALUES ($1, $2, $3)", newOrder.Id, evtTopic, evt)
	if err != nil {
		return nil, errors.WrapServiceError(errors.ErrCodeExtService, "failed to insert order event", err)
	}

	// Commit the transaction
	err = tx.Commit(ctx)
	if err != nil {
		return nil, errors.WrapServiceError(errors.ErrCodeExtService, "failed to commit transaction", err)
	}

	return &newOrder, nil
}

// Get all orders related to a specified customer.
// TODO: implement pagination
func (c postgresController) GetCustomerOrders(ctx context.Context, customerId string) ([]*pb.Order, error) {
	rows, err := c.cl.Query(ctx, pgOrderBaseQuery+" WHERE customer_id=$1 LIMIT 10", customerId)
	if err != nil {
		return nil, errors.WrapServiceError(errors.ErrCodeService, "query error whilst fetching customer orders", err)
	}

	orders := []*pb.Order{}
	for rows.Next() {
		orderObj, err := scanRowToOrder(rows)
		if err != nil {
			return nil, err
		}

		// Append the order's items
		orderObj, err = c.appendItemsToOrderObj(ctx, nil, orderObj)
		if err != nil {
			return nil, err
		}

		orders = append(orders, orderObj)
	}

	if rows.Err() != nil {
		return nil, errors.WrapServiceError(errors.ErrCodeService, "error whilst scanning order rows", rows.Err())
	}

	return orders, nil
}

// Set order status to approved
// Dispatch OrderApprovedEvent
func (c postgresController) ApproveOrder(ctx context.Context, orderId string, transactionId string) (*pb.Order, error) {
	// Begin a DB transaction
	tx, err := c.cl.Begin(ctx)
	if err != nil {
		return nil, errors.WrapServiceError(errors.ErrCodeExtService, "failed to begin transaction", err)
	}
	defer tx.Rollback(ctx)

	// Execute update query
	_, err = tx.Exec(
		ctx,
		"UPDATE orders SET status = $1, transaction_id = $2 WHERE id = $3",
		pb.OrderStatus_ORDER_STATUS_APPROVED,
		transactionId,
		orderId,
	)
	if err != nil {
		return nil, errors.WrapServiceError(errors.ErrCodeExtService, "failed to approve order", err)
	}

	orderObj, err := c.getOrder(ctx, &tx, orderId)
	if err != nil {
		return nil, errors.WrapServiceError(errors.ErrCodeExtService, "failed to approve order", err)
	}

	// Then add the event to the outbox table with the transaction.
	evt, evtTopic, err := order.PrepareOrderApprovedEvent(orderObj)
	if err != nil {
		return nil, errors.WrapServiceError(errors.ErrCodeService, "failed to create event", err)
	}

	_, err = tx.Exec(ctx, "INSERT INTO event_outbox (aggregateid, aggregatetype, payload) VALUES ($1, $2, $3)", orderObj.Id, evtTopic, evt)
	if err != nil {
		return nil, errors.WrapServiceError(errors.ErrCodeExtService, "failed to insert event", err)
	}

	// Commit the transaction
	err = tx.Commit(ctx)
	if err != nil {
		return nil, errors.WrapServiceError(errors.ErrCodeExtService, "failed to commit transaction", err)
	}

	return orderObj, nil
}

// Set order status to processing
// Dispatch OrderProcessingEvent
func (c postgresController) ProcessOrder(ctx context.Context, orderId string, itemsPrice float32) (*pb.Order, error) {
	// Begin a DB transaction
	tx, err := c.cl.Begin(ctx)
	if err != nil {
		return nil, errors.WrapServiceError(errors.ErrCodeExtService, "failed to begin transaction", err)
	}
	defer tx.Rollback(ctx)

	// Execute update query
	_, err = tx.Exec(
		ctx,
		"UPDATE orders SET status = $1, items_price = $2, total_price = $3 WHERE id = $4",
		pb.OrderStatus_ORDER_STATUS_PROCESSING,
		itemsPrice,
		itemsPrice,
		orderId,
	)
	if err != nil {
		return nil, errors.WrapServiceError(errors.ErrCodeExtService, "failed to update order", err)
	}

	orderObj, err := c.getOrder(ctx, &tx, orderId)
	if err != nil {
		return nil, errors.WrapServiceError(errors.ErrCodeExtService, "failed to update order", err)
	}

	// Then add the event to the outbox table with the transaction.
	// todo: fix name discrepency (mixed up processing and pending in my wording)
	evt, evtTopic, err := order.PrepareOrderPendingEvent(orderObj)
	if err != nil {
		return nil, errors.WrapServiceError(errors.ErrCodeService, "failed to create event", err)
	}

	_, err = tx.Exec(ctx, "INSERT INTO event_outbox (aggregateid, aggregatetype, payload) VALUES ($1, $2, $3)", orderObj.Id, evtTopic, evt)
	if err != nil {
		return nil, errors.WrapServiceError(errors.ErrCodeExtService, "failed to insert event", err)
	}

	// Commit the transaction
	err = tx.Commit(ctx)
	if err != nil {
		return nil, errors.WrapServiceError(errors.ErrCodeExtService, "failed to commit transaction", err)
	}

	return orderObj, nil
}

// Set order status to rejected (from processing)
// Dispatch OrderRejectedEvent
func (c postgresController) RejectOrder(ctx context.Context, orderId string) (*pb.Order, error) {
	// Begin a DB transaction
	tx, err := c.cl.Begin(ctx)
	if err != nil {
		return nil, errors.WrapServiceError(errors.ErrCodeExtService, "failed to begin transaction", err)
	}
	defer tx.Rollback(ctx)

	// Execute update query
	_, err = tx.Exec(
		ctx,
		"UPDATE orders SET status = $1 WHERE id = $2",
		pb.OrderStatus_ORDER_STATUS_REJECTED,
		orderId,
	)
	if err != nil {
		return nil, errors.WrapServiceError(errors.ErrCodeExtService, "failed to approve order", err)
	}

	orderObj, err := c.getOrder(ctx, &tx, orderId)
	if err != nil {
		return nil, errors.WrapServiceError(errors.ErrCodeExtService, "failed to approve order", err)
	}

	// Then add the event to the outbox table with the transaction.
	evt, evtTopic, err := order.PrepareOrderRejectedEvent(orderObj)
	if err != nil {
		return nil, errors.WrapServiceError(errors.ErrCodeService, "failed to create event", err)
	}

	_, err = tx.Exec(ctx, "INSERT INTO event_outbox (aggregateid, aggregatetype, payload) VALUES ($1, $2, $3)", orderObj.Id, evtTopic, evt)
	if err != nil {
		return nil, errors.WrapServiceError(errors.ErrCodeExtService, "failed to insert event", err)
	}

	// Commit the transaction
	err = tx.Commit(ctx)
	if err != nil {
		return nil, errors.WrapServiceError(errors.ErrCodeExtService, "failed to commit transaction", err)
	}

	return orderObj, nil
}

// Append shipment id to order
func (c postgresController) SetOrderShipmentId(ctx context.Context, orderId string, shippingId string) error {
	// Execute update query
	_, err := c.cl.Exec(
		ctx,
		"UPDATE orders SET shipment_id = $1 WHERE id = $2",
		shippingId,
		orderId,
	)
	if err != nil {
		return errors.WrapServiceError(errors.ErrCodeExtService, "failed to approve order", err)
	}

	return nil
}

// Build and exec an insert statement for a map of order items
func (c postgresController) createOrderItems(ctx context.Context, tx pgx.Tx, orderId string, items map[string]int32) error {
	// check there are items to add
	if len(items) > 1 {
		vals := [][]interface{}{}
		for productId, quantity := range items {
			vals = append(
				vals,
				goqu.Vals{orderId, productId, quantity},
			)
		}

		statement, args, err := goqu.Dialect("postgres").From(
			"order_items",
		).Insert().Cols(
			"order_id",
			"product_id",
			"quantity",
		).Vals(
			vals...,
		).Prepared(
			true,
		).ToSQL()
		if err != nil {
			return errors.WrapServiceError(errors.ErrCodeService, "failed to build SQL statement", err)
		}

		// Execute the statement on the transaction
		_, err = tx.Exec(ctx, statement, args...)
		if err != nil {
			return errors.WrapServiceError(errors.ErrCodeExtService, "failed to add items to order", err)
		}
	}

	return nil
}

func (c postgresController) getOrderItems(ctx context.Context, tx *pgx.Tx, orderId string) (*map[string]int32, error) {
	// Determine if transaction is being used.
	var rows pgx.Rows
	var err error
	if tx == nil {
		rows, err = c.cl.Query(ctx, pgOrderItemsBaseQuery+" WHERE order_id=$1", orderId)
	} else {
		rows, err = (*tx).Query(ctx, pgOrderItemsBaseQuery+" WHERE order_id=$1", orderId)
	}

	if err != nil {
		return nil, errors.WrapServiceError(errors.ErrCodeService, "query error whilst fetching order items", err)
	}

	items := make(map[string]int32)
	for rows.Next() {
		var (
			itemId   string
			quantity int32
		)
		err := rows.Scan(
			&itemId,
			&quantity,
		)
		if err != nil {
			return nil, errors.WrapServiceError(errors.ErrCodeService, "failed to scan an order item", err)
		}

		items[itemId] = quantity
	}

	if rows.Err() != nil {
		return nil, errors.WrapServiceError(errors.ErrCodeService, "error whilst scanning order item rows", rows.Err())
	}

	return &items, nil
}

// Appends order items to an order object.
func (c postgresController) appendItemsToOrderObj(ctx context.Context, tx *pgx.Tx, orderObj *pb.Order) (*pb.Order, error) {
	// Load the order items
	orderItems, err := c.getOrderItems(ctx, tx, orderObj.Id)
	if err != nil {
		return nil, err
	}

	// Add the order items to the order protobuf
	orderObj.Items = *orderItems

	// Return the order
	return orderObj, nil
}

// Scan a postgres row to a protobuf order object.
func scanRowToOrder(row pgx.Row) (*pb.Order, error) {
	var order pb.Order

	// Temporary variables that require conversion
	var tmpCreatedAt pgtype.Timestamp
	var tmpUpdatedAt pgtype.Timestamp

	err := row.Scan(
		&order.Id,
		&order.Status,
		&order.CustomerId,
		&order.ShippingId,
		&order.TransactionId,
		&tmpCreatedAt,
		&tmpUpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, errors.WrapServiceError(errors.ErrCodeNotFound, "order not found", err)
		} else {
			return nil, errors.WrapServiceError(errors.ErrCodeExtService, "something went wrong scanning order", err)
		}
	}

	// Convert the temporary variables
	//
	// This includes converting postgres timestamps to unix format
	if tmpCreatedAt.Valid {
		order.CreatedAt = tmpCreatedAt.Time.Unix()
	} else {
		return nil, errors.NewServiceError(errors.ErrCodeUnknown, "failed to convert order (created_at) timestamp")
	}

	if tmpUpdatedAt.Valid {
		unixUpdated := tmpUpdatedAt.Time.Unix()
		order.UpdatedAt = &unixUpdated
	}

	return &order, nil
}
