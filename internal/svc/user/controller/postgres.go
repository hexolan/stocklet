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

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/hexolan/stocklet/internal/pkg/errors"
	authpb "github.com/hexolan/stocklet/internal/pkg/protogen/auth/v1"
	pb "github.com/hexolan/stocklet/internal/pkg/protogen/user/v1"
	"github.com/hexolan/stocklet/internal/svc/user"
)

const pgUserBaseQuery string = "SELECT id, first_name, last_name, email, created_at, updated_at FROM users"

type postgresController struct {
	cl          *pgxpool.Pool
	serviceOpts *user.ServiceConfigOpts
}

func NewPostgresController(cl *pgxpool.Pool, serviceOpts *user.ServiceConfigOpts) user.StorageController {
	return postgresController{cl: cl, serviceOpts: serviceOpts}
}

func (c postgresController) GetUser(ctx context.Context, userId string) (*pb.User, error) {
	return c.getUser(ctx, nil, userId)
}

func (c postgresController) getUser(ctx context.Context, tx *pgx.Tx, userId string) (*pb.User, error) {
	// Determine if a db transaction is being used
	var row pgx.Row
	const query = pgUserBaseQuery + " WHERE id=$1"
	if tx == nil {
		row = c.cl.QueryRow(ctx, query, userId)
	} else {
		row = (*tx).QueryRow(ctx, query, userId)
	}

	// Scan row to protobuf obj
	userObj, err := scanRowToUser(row)
	if err != nil {
		return nil, err
	}

	return userObj, nil
}

func (c postgresController) RegisterUser(ctx context.Context, email string, password string, firstName string, lastName string) (*pb.User, error) {
	// Establish connection with auth service
	authConn, err := grpc.Dial(
		c.serviceOpts.AuthServiceGrpc,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithStatsHandler(otelgrpc.NewClientHandler()),
	)
	if err != nil {
		return nil, errors.WrapServiceError(errors.ErrCodeExtService, "failed to establish connection to auth service", err)
	}
	defer authConn.Close()
	authCl := authpb.NewAuthServiceClient(authConn)

	// Begin a DB transaction
	tx, err := c.cl.Begin(ctx)
	if err != nil {
		return nil, errors.WrapServiceError(errors.ErrCodeExtService, "failed to begin db transaction", err)
	}
	defer tx.Rollback(ctx)

	// Create user in database
	var userId string
	err = tx.QueryRow(
		ctx,
		"INSERT INTO users (first_name, last_name, email) VALUES ($1, $2, $3) RETURNING id",
		firstName,
		lastName,
		email,
	).Scan(&userId)
	if err != nil {
		return nil, errors.WrapServiceError(errors.ErrCodeExtService, "failed to create user", err)
	}

	// Get the newly created user obj
	userObj, err := c.getUser(ctx, &tx, userId)
	if err != nil {
		return nil, errors.WrapServiceError(errors.ErrCodeExtService, "failed to register user", err)
	}

	// Prepare user created event and append to transaction
	evt, evtTopic, err := user.PrepareUserCreatedEvent(userObj)
	if err != nil {
		return nil, errors.WrapServiceError(errors.ErrCodeService, "failed to create event", err)
	}

	_, err = tx.Exec(ctx, "INSERT INTO event_outbox (aggregateid, aggregatetype, payload) VALUES ($1, $2, $3)", userObj.Id, evtTopic, evt)
	if err != nil {
		return nil, errors.WrapServiceError(errors.ErrCodeExtService, "failed to insert event", err)
	}

	// Attempt to add auth method for user
	authCtx, authCtxCancel := context.WithTimeout(context.Background(), time.Second*10)
	defer authCtxCancel()
	_, err = authCl.SetPassword(authCtx, &authpb.SetPasswordRequest{UserId: userObj.Id, Password: password})
	if err != nil {
		return nil, errors.WrapServiceError(errors.ErrCodeExtService, "error registering user: failed to set auth method", err)
	}

	// Commit the transaction
	err = tx.Commit(ctx)
	if err != nil {
		return nil, errors.WrapServiceError(errors.ErrCodeExtService, "failed to commit transaction", err)
	}

	return userObj, nil
}

func (c postgresController) UpdateUserEmail(ctx context.Context, userId string, email string) error {
	// Begin a DB transaction
	tx, err := c.cl.Begin(ctx)
	if err != nil {
		return errors.WrapServiceError(errors.ErrCodeExtService, "failed to begin transaction", err)
	}
	defer tx.Rollback(ctx)

	// Update email
	_, err = tx.Exec(ctx, "UPDATE users SET email=$1 WHERE id=$2", email, userId)
	if err != nil {
		return errors.WrapServiceError(errors.ErrCodeExtService, "failed to update email", err)
	}

	// Add the event to the outbox table with the transaction
	evt, evtTopic, err := user.PrepareUserEmailUpdatedEvent(userId, email)
	if err != nil {
		return errors.WrapServiceError(errors.ErrCodeService, "failed to create event", err)
	}

	_, err = tx.Exec(ctx, "INSERT INTO event_outbox (aggregateid, aggregatetype, payload) VALUES ($1, $2, $3)", userId, evtTopic, evt)
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

func (c postgresController) DeleteUser(ctx context.Context, userId string) (*pb.User, error) {
	// Begin a DB transaction
	tx, err := c.cl.Begin(ctx)
	if err != nil {
		return nil, errors.WrapServiceError(errors.ErrCodeExtService, "failed to begin transaction", err)
	}
	defer tx.Rollback(ctx)

	// Get the user
	userObj, err := c.getUser(ctx, &tx, userId)
	if err != nil {
		return nil, errors.WrapServiceError(errors.ErrCodeExtService, "failed to fetch user", err)
	}

	// Delete user
	_, err = tx.Exec(ctx, "DELETE FROM users WHERE id=$1", userId)
	if err != nil {
		return nil, errors.WrapServiceError(errors.ErrCodeExtService, "failed to delete user", err)
	}

	// Add the event to the outbox table with the transaction
	evt, evtTopic, err := user.PrepareUserDeletedEvent(userObj)
	if err != nil {
		return nil, errors.WrapServiceError(errors.ErrCodeService, "failed to create event", err)
	}

	_, err = tx.Exec(ctx, "INSERT INTO event_outbox (aggregateid, aggregatetype, payload) VALUES ($1, $2, $3)", userObj.Id, evtTopic, evt)
	if err != nil {
		return nil, errors.WrapServiceError(errors.ErrCodeExtService, "failed to insert event", err)
	}

	// Commit the transaction
	err = tx.Commit(ctx)
	if err != nil {
		return nil, errors.WrapServiceError(errors.ErrCodeExtService, "failed to commit transaction", err)
	}

	return userObj, nil
}

// Scan a postgres row to a protobuf object
func scanRowToUser(row pgx.Row) (*pb.User, error) {
	var user pb.User

	// Temporary variables that require conversion
	var tmpCreatedAt pgtype.Timestamp
	var tmpUpdatedAt pgtype.Timestamp

	err := row.Scan(
		&user.Id,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&tmpCreatedAt,
		&tmpUpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, errors.WrapServiceError(errors.ErrCodeNotFound, "user not found", err)
		} else {
			return nil, errors.WrapServiceError(errors.ErrCodeExtService, "failed to scan object from database", err)
		}
	}

	// convert postgres timestamps to unix format
	if tmpCreatedAt.Valid {
		user.CreatedAt = tmpCreatedAt.Time.Unix()
	} else {
		return nil, errors.NewServiceError(errors.ErrCodeUnknown, "failed to scan object from database (timestamp conversion)")
	}

	if tmpUpdatedAt.Valid {
		unixUpdated := tmpUpdatedAt.Time.Unix()
		user.UpdatedAt = &unixUpdated
	}

	return &user, nil
}
