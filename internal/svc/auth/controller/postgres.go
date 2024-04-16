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

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/hexolan/stocklet/internal/pkg/errors"
	"github.com/hexolan/stocklet/internal/svc/auth"
)

type postgresController struct {
	cl *pgxpool.Pool
}

func NewPostgresController(cl *pgxpool.Pool) auth.StorageController {
	return postgresController{cl: cl}
}

func (c postgresController) getPasswordAuthMethod(ctx context.Context, userId string) (string, error) {
	var hashedPassword string
	err := c.cl.QueryRow(ctx, "SELECT hashed_password FROM auth_methods WHERE user_id=$1", userId).Scan(&hashedPassword)
	if err != nil {
		return "", errors.WrapServiceError(errors.ErrCodeNotFound, "unknown user id", err)
	}

	return hashedPassword, nil
}

func (c postgresController) SetPassword(ctx context.Context, userId string, password string) error {
	hashedPassword, err := auth.HashPassword(password)
	if err != nil {
		return errors.WrapServiceError(errors.ErrCodeInvalidArgument, "unable to hash password", err)
	}

	// Check if auth method already exists
	var statement string
	_, err = c.getPasswordAuthMethod(ctx, userId)
	if err != nil {
		// Auth method does not exist
		statement = "INSERT INTO auth_methods (user_id, hashed_password) VALUES ($1, $2)"
	} else {
		// Auth method already exists
		statement = "UPDATE auth_methods SET hashed_password=$2 WHERE user_id=$1"
	}

	// Execute statement
	result, err := c.cl.Exec(ctx, statement, userId, hashedPassword)
	if err != nil {
		return errors.WrapServiceError(errors.ErrCodeExtService, "failed to set password", err)
	}

	// Ensure a row was affected
	if result.RowsAffected() != 1 {
		return errors.WrapServiceError(errors.ErrCodeExtService, "auth methods unaffected", err)
	}

	return nil
}

func (c postgresController) VerifyPassword(ctx context.Context, userId string, password string) (bool, error) {
	hashedPassword, err := c.getPasswordAuthMethod(ctx, userId)
	if err != nil {
		return false, err
	}

	match := auth.CompareHashAndPassword(password, hashedPassword)
	if !match {
		return false, errors.WrapServiceError(errors.ErrCodeForbidden, "invalid user id or password", err)
	}

	return true, nil
}

func (c postgresController) DeleteAuthMethods(ctx context.Context, userId string) error {
	_, err := c.cl.Exec(ctx, "DELETE FROM auth_methods WHERE user_id=$1", userId)
	if err != nil {
		return errors.WrapServiceError(errors.ErrCodeExtService, "failed to delete", err)
	}

	return nil
}
