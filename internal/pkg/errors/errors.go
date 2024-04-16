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

package errors

import (
	"fmt"

	"google.golang.org/grpc/status"
)

type ServiceError struct {
	code    ErrorCode
	msg     string
	wrapped error
}

func NewServiceError(code ErrorCode, msg string) error {
	return &ServiceError{
		code: code,
		msg:  msg,
	}
}

func NewServiceErrorf(code ErrorCode, msg string, args ...interface{}) error {
	return NewServiceError(code, fmt.Sprintf(msg, args...))
}

func WrapServiceError(code ErrorCode, msg string, wrapped error) error {
	return &ServiceError{
		code:    code,
		msg:     msg,
		wrapped: wrapped,
	}
}

func (e ServiceError) Error() string {
	if e.wrapped != nil {
		return fmt.Sprintf("%s: %s", e.msg, e.wrapped.Error())
	}

	return e.msg
}

func (e ServiceError) Code() ErrorCode {
	return e.code
}

// Set the gRPC status to only expose the top error message.
//
// This is to prevent any full error contexts (from wrapped errors) being exposed to users by the gateway.
// e.g. "{"code":2,"message":"something went wrong scanning order: failed to connect to `host=postgres user=postgres database=postgres`: hostname resolving error (lookup postgres on 127.0.0.11:53: server misbehaving)","details":[]}"
func (e ServiceError) GRPCStatus() *status.Status {
	return status.New(e.Code().GRPCCode(), e.msg)
}

func (e ServiceError) Unwrap() error {
	return e.wrapped
}
