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
	"golang.org/x/exp/maps"

	"github.com/doug-martin/goqu/v9"
	_ "github.com/doug-martin/goqu/v9/dialect/postgres"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/hexolan/stocklet/internal/pkg/errors"
	pb "github.com/hexolan/stocklet/internal/pkg/protogen/product/v1"
	"github.com/hexolan/stocklet/internal/svc/product"
)

const pgProductBaseQuery string = "SELECT id, name, description, price, created_at, updated_at FROM products"

type postgresController struct {
	cl *pgxpool.Pool
}

func NewPostgresController(cl *pgxpool.Pool) product.StorageController {
	return postgresController{cl: cl}
}

func (c postgresController) GetProduct(ctx context.Context, productId string) (*pb.Product, error) {
	return c.getProduct(ctx, nil, productId)
}

func (c postgresController) getProduct(ctx context.Context, tx *pgx.Tx, productId string) (*pb.Product, error) {
	// Determine if a db transaction is being used
	var row pgx.Row
	const query = pgProductBaseQuery + " WHERE id=$1"
	if tx == nil {
		row = c.cl.QueryRow(ctx, query, productId)
	} else {
		row = (*tx).QueryRow(ctx, query, productId)
	}

	// Scan row to protobuf obj
	product, err := scanRowToProduct(row)
	if err != nil {
		return nil, err
	}

	return product, nil
}

// todo: implementing pagination mechanism
func (c postgresController) GetProducts(ctx context.Context) ([]*pb.Product, error) {
	rows, err := c.cl.Query(ctx, pgProductBaseQuery+" LIMIT 10")
	if err != nil {
		return nil, errors.WrapServiceError(errors.ErrCodeService, "query error", err)
	}

	products := []*pb.Product{}
	for rows.Next() {
		productObj, err := scanRowToProduct(rows)
		if err != nil {
			return nil, err
		}

		products = append(products, productObj)
	}

	if rows.Err() != nil {
		return nil, errors.WrapServiceError(errors.ErrCodeService, "error whilst scanning rows", rows.Err())
	}

	return products, nil
}

// Update a product price.
func (c postgresController) UpdateProductPrice(ctx context.Context, productId string, price float32) (*pb.Product, error) {
	// Begin a DB transaction
	tx, err := c.cl.Begin(ctx)
	if err != nil {
		return nil, errors.WrapServiceError(errors.ErrCodeExtService, "failed to begin transaction", err)
	}
	defer tx.Rollback(ctx)

	// Update product price
	_, err = tx.Exec(ctx, "UPDATE products SET price=$1 WHERE id=$2", price, productId)
	if err != nil {
		return nil, errors.WrapServiceError(errors.ErrCodeExtService, "failed to update product price", err)
	}

	// Get updated product
	productObj, err := c.getProduct(ctx, &tx, productId)
	if err != nil {
		return nil, err
	}

	// Add the event to the outbox table with the transaction
	evt, evtTopic, err := product.PrepareProductPriceUpdatedEvent(productObj)
	if err != nil {
		return nil, errors.WrapServiceError(errors.ErrCodeService, "failed to create event", err)
	}

	_, err = tx.Exec(ctx, "INSERT INTO event_outbox (aggregateid, aggregatetype, payload) VALUES ($1, $2, $3)", productObj.Id, evtTopic, evt)
	if err != nil {
		return nil, errors.WrapServiceError(errors.ErrCodeExtService, "failed to insert event", err)
	}

	// Commit the transaction
	err = tx.Commit(ctx)
	if err != nil {
		return nil, errors.WrapServiceError(errors.ErrCodeExtService, "failed to commit transaction", err)
	}

	return productObj, nil
}

// Delete a product by its specified id.
func (c postgresController) DeleteProduct(ctx context.Context, productId string) error {
	// Begin a DB transaction
	tx, err := c.cl.Begin(ctx)
	if err != nil {
		return errors.WrapServiceError(errors.ErrCodeExtService, "failed to begin transaction", err)
	}
	defer tx.Rollback(ctx)

	// Get product
	productObj, err := c.getProduct(ctx, &tx, productId)
	if err != nil {
		return err
	}

	// Delete product
	_, err = tx.Exec(ctx, "DELETE FROM products WHERE id=$1", productId)
	if err != nil {
		return errors.WrapServiceError(errors.ErrCodeExtService, "failed to delete product", err)
	}

	// Add the event to the outbox table with the transaction
	evt, evtTopic, err := product.PrepareProductDeletedEvent(productObj)
	if err != nil {
		return errors.WrapServiceError(errors.ErrCodeService, "failed to create event", err)
	}

	_, err = tx.Exec(ctx, "INSERT INTO event_outbox (aggregateid, aggregatetype, payload) VALUES ($1, $2, $3)", productObj.Id, evtTopic, evt)
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

func (c postgresController) PriceOrderProducts(ctx context.Context, orderId string, customerId string, productQuantities map[string]int32) error {
	// Begin a DB transaction
	tx, err := c.cl.Begin(ctx)
	if err != nil {
		return errors.WrapServiceError(errors.ErrCodeExtService, "failed to begin transaction", err)
	}
	defer tx.Rollback(ctx)

	// Get prices of all specified products (in productQuantities)
	productIds := maps.Keys(productQuantities)
	statement, args, err := goqu.Dialect("postgres").From("products").Select("id", "price").Where(goqu.C("id").In(productIds)).Prepared(true).ToSQL()
	if err != nil {
		return errors.WrapServiceError(errors.ErrCodeService, "failed to build statement", err)
	}

	rows, err := tx.Query(ctx, statement, args...)
	if err != nil {
		return errors.WrapServiceError(errors.ErrCodeExtService, "failed to fetch price quotes", err)
	}

	var productPrices map[string]float32
	for rows.Next() {
		var productId string
		var productPrice float32
		err := rows.Scan(&productId, &productPrice)
		if err != nil {
			return errors.WrapServiceError(errors.ErrCodeNotFound, "failed to fetch price quotes: error whilst scanning row", err)
		}

		productPrices[productId] = productPrice
	}

	if rows.Err() != nil {
		return errors.WrapServiceError(errors.ErrCodeService, "failed to fetch price quotes: error whilst scanning rows", rows.Err())
	}

	// Calculate total price
	// Also ensuring that all items in the itemQuantities have a fetched price
	var totalPrice float32
	for productId, quantity := range productQuantities {
		productPrice, ok := productPrices[productId]
		if !ok {
			// Prepare and dispatch failure product pricing event
			evt, evtTopic, err := product.PrepareProductPriceQuoteEvent_Unavaliable(orderId)
			if err != nil {
				return errors.WrapServiceError(errors.ErrCodeService, "failed to create event", err)
			}

			_, err = c.cl.Exec(ctx, "INSERT INTO event_outbox (aggregateid, aggregatetype, payload) VALUES ($1, $2, $3)", orderId, evtTopic, evt)
			if err != nil {
				return errors.WrapServiceError(errors.ErrCodeExtService, "failed to insert event", err)
			}

			return nil
		}

		// Add to total price
		totalPrice += productPrice * float32(quantity)
	}

	// Prepare and dispatch succesful product pricing event
	evt, evtTopic, err := product.PrepareProductPriceQuoteEvent_Avaliable(
		orderId,
		productQuantities,
		productPrices,
		totalPrice,
	)
	if err != nil {
		return errors.WrapServiceError(errors.ErrCodeService, "failed to create event", err)
	}

	_, err = tx.Exec(ctx, "INSERT INTO event_outbox (aggregateid, aggregatetype, payload) VALUES ($1, $2, $3)", orderId, evtTopic, evt)
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
func scanRowToProduct(row pgx.Row) (*pb.Product, error) {
	var productObj pb.Product

	// Temporary variables that require conversion
	var tmpCreatedAt pgtype.Timestamp
	var tmpUpdatedAt pgtype.Timestamp

	err := row.Scan(
		&productObj.Id,
		&productObj.Name,
		&productObj.Description,
		&productObj.Price,
		&tmpCreatedAt,
		&tmpUpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, errors.WrapServiceError(errors.ErrCodeNotFound, "product not found", err)
		} else {
			return nil, errors.WrapServiceError(errors.ErrCodeExtService, "failed to scan object from database", err)
		}
	}

	// convert postgres timestamps to unix format
	if tmpCreatedAt.Valid {
		productObj.CreatedAt = tmpCreatedAt.Time.Unix()
	} else {
		return nil, errors.NewServiceError(errors.ErrCodeUnknown, "failed to scan object from database (timestamp conversion)")
	}

	if tmpUpdatedAt.Valid {
		unixUpdated := tmpUpdatedAt.Time.Unix()
		productObj.UpdatedAt = &unixUpdated
	}

	return &productObj, nil
}
