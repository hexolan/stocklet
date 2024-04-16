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
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/hexolan/stocklet/internal/pkg/errors"
	pb "github.com/hexolan/stocklet/internal/pkg/protogen/payment/v1"
	"github.com/hexolan/stocklet/internal/svc/payment"
)

const (
	pgTransactionBaseQuery     string = "SELECT id, order_id, customer_id, amount, reversed_at, processed_at FROM transactions"
	pgCustomerBalanceBaseQuery string = "SELECT customer_id, balance FROM customer_balances"
)

type postgresController struct {
	cl *pgxpool.Pool
}

func NewPostgresController(cl *pgxpool.Pool) payment.StorageController {
	return postgresController{cl: cl}
}

func (c postgresController) GetBalance(ctx context.Context, customerId string) (*pb.CustomerBalance, error) {
	return c.getBalance(ctx, nil, customerId)
}

func (c postgresController) getBalance(ctx context.Context, tx *pgx.Tx, customerId string) (*pb.CustomerBalance, error) {
	// Determine if a db transaction is being used
	var row pgx.Row
	const query = pgCustomerBalanceBaseQuery + " WHERE customer_id=$1"
	if tx == nil {
		row = c.cl.QueryRow(ctx, query, customerId)
	} else {
		row = (*tx).QueryRow(ctx, query, customerId)
	}

	// Scan row to protobuf obj
	balance, err := scanRowToCustomerBalance(row)
	if err != nil {
		return nil, err
	}

	return balance, nil
}

func (c postgresController) GetTransaction(ctx context.Context, transactionId string) (*pb.Transaction, error) {
	return c.getTransaction(ctx, nil, transactionId)
}

func (c postgresController) getTransaction(ctx context.Context, tx *pgx.Tx, transactionId string) (*pb.Transaction, error) {
	// Determine if a db transaction is being used
	var row pgx.Row
	const query = pgTransactionBaseQuery + " WHERE id=$1"
	if tx == nil {
		row = c.cl.QueryRow(ctx, query, transactionId)
	} else {
		row = (*tx).QueryRow(ctx, query, transactionId)
	}

	// Scan row to protobuf obj
	transaction, err := scanRowToTransaction(row)
	if err != nil {
		return nil, err
	}

	return transaction, nil
}

func (c postgresController) CreateBalance(ctx context.Context, customerId string) error {
	// Begin a DB transaction
	tx, err := c.cl.Begin(ctx)
	if err != nil {
		return errors.WrapServiceError(errors.ErrCodeExtService, "failed to begin transaction", err)
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(
		ctx,
		"INSERT INTO customer_balances (customer_id, balance) VALUES ($1, $2)",
		customerId,
		0.00,
	)
	if err != nil {
		return errors.WrapServiceError(errors.ErrCodeExtService, "failed to create balance", err)
	}

	// Add the event to the outbox table with the transaction
	evt, evtTopic, err := payment.PrepareBalanceCreatedEvent(&pb.CustomerBalance{CustomerId: customerId, Balance: 0.00})
	if err != nil {
		return errors.WrapServiceError(errors.ErrCodeService, "failed to create event", err)
	}

	_, err = tx.Exec(ctx, "INSERT INTO event_outbox (aggregateid, aggregatetype, payload) VALUES ($1, $2, $3)", customerId, evtTopic, evt)
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

func (c postgresController) CreditBalance(ctx context.Context, customerId string, amount float32) error {
	// Begin a DB transaction
	tx, err := c.cl.Begin(ctx)
	if err != nil {
		return errors.WrapServiceError(errors.ErrCodeExtService, "failed to begin transaction", err)
	}
	defer tx.Rollback(ctx)

	// Add to balance
	_, err = tx.Exec(
		ctx,
		"UPDATE customer_balances SET balance = balance + $1 WHERE customer_id=$2",
		amount,
		customerId,
	)
	if err != nil {
		return errors.WrapServiceError(errors.ErrCodeExtService, "failed to update balance", err)
	}

	// Get updated balance
	balance, err := c.getBalance(ctx, &tx, customerId)
	if err != nil {
		return err
	}

	// Add the event to the outbox table with the transaction
	evt, evtTopic, err := payment.PrepareBalanceCreditedEvent(
		balance.CustomerId,
		amount,
		balance.Balance,
	)
	if err != nil {
		return errors.WrapServiceError(errors.ErrCodeService, "failed to create event", err)
	}

	_, err = tx.Exec(ctx, "INSERT INTO event_outbox (aggregateid, aggregatetype, payload) VALUES ($1, $2, $3)", balance.CustomerId, evtTopic, evt)
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

func (c postgresController) DebitBalance(ctx context.Context, customerId string, amount float32, orderId *string) (*pb.Transaction, error) {
	return c.debitBalance(ctx, nil, customerId, amount, orderId)
}

func (c postgresController) debitBalance(ctx context.Context, tx *pgx.Tx, customerId string, amount float32, orderId *string) (*pb.Transaction, error) {
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
			return nil, errors.WrapServiceError(errors.ErrCodeExtService, "failed to begin transaction", err)
		}
		defer funcTx.Rollback(ctx)
	}

	// Subtract from balance
	_, err = funcTx.Exec(
		ctx,
		"UPDATE customer_balances SET balance = balance - $1 WHERE customer_id=$2",
		amount,
		customerId,
	)
	if err != nil {
		return nil, errors.WrapServiceError(errors.ErrCodeExtService, "failed to update balance", err)
	}

	// Get updated balance
	balance, err := c.getBalance(ctx, &funcTx, customerId)
	if err != nil {
		return nil, err
	}

	// Ensure that the balance is not negative
	if balance.Balance < 0.00 {
		return nil, errors.NewServiceError(errors.ErrCodeInvalidArgument, "insufficient balance")
	}

	// Add the balance event to the outbox table with the transaction
	evt, evtTopic, err := payment.PrepareBalanceDebitedEvent(
		balance.CustomerId,
		amount,
		balance.Balance,
	)
	if err != nil {
		return nil, errors.WrapServiceError(errors.ErrCodeService, "failed to create event", err)
	}

	_, err = funcTx.Exec(ctx, "INSERT INTO event_outbox (aggregateid, aggregatetype, payload) VALUES ($1, $2, $3)", balance.CustomerId, evtTopic, evt)
	if err != nil {
		return nil, errors.WrapServiceError(errors.ErrCodeExtService, "failed to insert event", err)
	}

	// Create a payment transaction record
	transaction, err := c.createTransaction(ctx, &funcTx, orderId, customerId, amount)
	if err != nil {
		return nil, err
	}

	// Commit the transaction (if created in this func)
	if tx == nil {
		err = funcTx.Commit(ctx)
		if err != nil {
			return nil, errors.WrapServiceError(errors.ErrCodeExtService, "failed to commit transaction", err)
		}
	}

	return transaction, nil
}

func (c postgresController) CloseBalance(ctx context.Context, customerId string) error {
	// Begin a DB transaction
	tx, err := c.cl.Begin(ctx)
	if err != nil {
		return errors.WrapServiceError(errors.ErrCodeExtService, "failed to begin transaction", err)
	}
	defer tx.Rollback(ctx)

	// Get current balance
	balance, err := c.getBalance(ctx, &tx, customerId)
	if err != nil {
		return err
	}

	// Delete balance
	_, err = tx.Exec(ctx, "DELETE FROM customer_balances WHERE customer_id=$1", customerId)
	if err != nil {
		return errors.WrapServiceError(errors.ErrCodeExtService, "failed to delete balance", err)
	}

	// Add the event to the outbox table with the transaction
	evt, evtTopic, err := payment.PrepareBalanceClosedEvent(balance)
	if err != nil {
		return errors.WrapServiceError(errors.ErrCodeService, "failed to create event", err)
	}

	_, err = tx.Exec(ctx, "INSERT INTO event_outbox (aggregateid, aggregatetype, payload) VALUES ($1, $2, $3)", balance.CustomerId, evtTopic, evt)
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

func (c postgresController) PaymentForOrder(ctx context.Context, orderId string, customerId string, amount float32) error {
	// Begin a DB transaction
	tx, err := c.cl.Begin(ctx)
	if err != nil {
		return errors.WrapServiceError(errors.ErrCodeExtService, "failed to begin transaction", err)
	}
	defer tx.Rollback(ctx)

	// Attempt to debit balance for the order
	transaction, err := c.debitBalance(ctx, &tx, customerId, amount, &orderId)
	if err != nil {
		// check that error is not a result of insufficient balance
		// - or the customer not having a balance
		errText := err.Error()
		if errText != "insufficient balance" && !strings.HasPrefix(errText, "failed to update balance") {
			return err
		}
	}

	// Prepare response event
	var (
		evt      []byte
		evtTopic string
	)
	if transaction != nil {
		// Succesful
		evt, evtTopic, err = payment.PreparePaymentProcessedEvent_Success(transaction)
	} else {
		// Failure
		// - result of insufficient/non-existent balance
		evt, evtTopic, err = payment.PreparePaymentProcessedEvent_Failure(orderId, customerId, amount)
	}

	// Ensure the event was prepared succesfully
	if err != nil {
		return errors.WrapServiceError(errors.ErrCodeService, "failed to create event", err)
	}

	// Add the event to the outbox table
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

func (c postgresController) createTransaction(ctx context.Context, tx *pgx.Tx, orderId *string, customerId string, amount float32) (*pb.Transaction, error) {
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
			return nil, errors.WrapServiceError(errors.ErrCodeExtService, "failed to begin transaction", err)
		}
		defer funcTx.Rollback(ctx)
	}

	// Insert the transaction
	var transactionId string
	err = funcTx.QueryRow(
		ctx,
		"INSERT INTO transactions (order_id, customer_id, amount) VALUES ($1, $2, $3) RETURNING id",
		orderId,
		customerId,
		amount,
	).Scan(&transactionId)
	if err != nil {
		return nil, errors.WrapServiceError(errors.ErrCodeExtService, "failed to insert transaction", err)
	}

	// Get the transaction obj
	transaction, err := c.getTransaction(ctx, &funcTx, transactionId)
	if err != nil {
		return nil, err
	}

	// Add the event to the outbox table with the transaction
	evt, evtTopic, err := payment.PrepareTransactionLoggedEvent(transaction)
	if err != nil {
		return nil, errors.WrapServiceError(errors.ErrCodeService, "failed to create event", err)
	}

	_, err = funcTx.Exec(ctx, "INSERT INTO event_outbox (aggregateid, aggregatetype, payload) VALUES ($1, $2, $3)", transaction.Id, evtTopic, evt)
	if err != nil {
		return nil, errors.WrapServiceError(errors.ErrCodeExtService, "failed to insert event", err)
	}

	// Commit the transaction (if created in this func)
	if tx == nil {
		err = funcTx.Commit(ctx)
		if err != nil {
			return nil, errors.WrapServiceError(errors.ErrCodeExtService, "failed to commit transaction", err)
		}
	}

	return transaction, nil
}

// Scan a postgres row to a protobuf object
func scanRowToTransaction(row pgx.Row) (*pb.Transaction, error) {
	var transaction pb.Transaction

	// Temporary variables that require conversion
	var tmpProcessedAt pgtype.Timestamp
	var tmpReversedAt pgtype.Timestamp

	err := row.Scan(
		&transaction.Id,
		&transaction.OrderId,
		&transaction.CustomerId,
		&transaction.Amount,
		&tmpReversedAt,
		&tmpProcessedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, errors.WrapServiceError(errors.ErrCodeNotFound, "transaction not found", err)
		} else {
			return nil, errors.WrapServiceError(errors.ErrCodeExtService, "failed to scan object from database", err)
		}
	}

	// Convert the temporary variables
	// - converting postgres timestamps to unix format
	if tmpProcessedAt.Valid {
		transaction.ProcessedAt = tmpProcessedAt.Time.Unix()
	} else {
		return nil, errors.NewServiceError(errors.ErrCodeUnknown, "failed to scan object from database (timestamp conversion)")
	}

	if tmpReversedAt.Valid {
		unixReversed := tmpReversedAt.Time.Unix()
		transaction.ReversedAt = &unixReversed
	}

	return &transaction, nil
}

// Scan a postgres row to a protobuf object
func scanRowToCustomerBalance(row pgx.Row) (*pb.CustomerBalance, error) {
	var balance pb.CustomerBalance

	err := row.Scan(
		&balance.CustomerId,
		&balance.Balance,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, errors.WrapServiceError(errors.ErrCodeNotFound, "balance not found", err)
		} else {
			return nil, errors.WrapServiceError(errors.ErrCodeExtService, "failed to scan object from database", err)
		}
	}

	return &balance, nil
}
