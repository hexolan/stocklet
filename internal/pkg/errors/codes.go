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
	"google.golang.org/grpc/codes"
)

type ErrorCode int32

const (
	ErrCodeUnknown ErrorCode = iota

	ErrCodeService
	ErrCodeExtService

	ErrCodeNotFound
	ErrCodeAlreadyExists

	ErrCodeForbidden
	ErrCodeUnauthorised

	ErrCodeInvalidArgument
)

// Maps the custom service error codes
// to their gRPC status code equivalents.
func (c ErrorCode) GRPCCode() codes.Code {
	codeMap := map[ErrorCode]codes.Code{
		ErrCodeUnknown: codes.Unknown,

		ErrCodeService:    codes.Internal,
		ErrCodeExtService: codes.Unavailable,

		ErrCodeNotFound:      codes.NotFound,
		ErrCodeAlreadyExists: codes.AlreadyExists,

		ErrCodeForbidden:    codes.PermissionDenied,
		ErrCodeUnauthorised: codes.PermissionDenied,

		ErrCodeInvalidArgument: codes.InvalidArgument,
	}

	grpcCode, mapped := codeMap[c]
	if mapped {
		return grpcCode
	}
	return codes.Unknown
}
