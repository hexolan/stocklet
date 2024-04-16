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

package user

import (
	"context"

	"github.com/bufbuild/protovalidate-go"
	"github.com/rs/zerolog/log"

	"github.com/hexolan/stocklet/internal/pkg/errors"
	"github.com/hexolan/stocklet/internal/pkg/messaging"
	commonpb "github.com/hexolan/stocklet/internal/pkg/protogen/common/v1"
	pb "github.com/hexolan/stocklet/internal/pkg/protogen/user/v1"
)

// Interface for the service
type UserService struct {
	pb.UnimplementedUserServiceServer

	store StorageController
	pbVal *protovalidate.Validator
}

// Interface for database methods
// Flexibility for implementing seperate controllers for different databases (e.g. Postgres, MongoDB, etc)
type StorageController interface {
	GetUser(ctx context.Context, userId string) (*pb.User, error)

	RegisterUser(ctx context.Context, email string, password string, firstName string, lastName string) (*pb.User, error)
	UpdateUserEmail(ctx context.Context, userId string, email string) error

	DeleteUser(ctx context.Context, userId string) (*pb.User, error)
}

// Interface for event consumption
// Flexibility for seperate controllers for different messaging systems (e.g. Kafka, NATS, etc)
type ConsumerController interface {
	messaging.ConsumerController

	Attach(svc pb.UserServiceServer)
}

// Create the shipping service
func NewUserService(cfg *ServiceConfig, store StorageController) *UserService {
	// Initialise the protobuf validator
	pbVal, err := protovalidate.New()
	if err != nil {
		log.Panic().Err(err).Msg("failed to initialise protobuf validator")
	}

	// Initialise the service
	return &UserService{
		store: store,
		pbVal: pbVal,
	}
}

func (svc UserService) ServiceInfo(ctx context.Context, req *commonpb.ServiceInfoRequest) (*commonpb.ServiceInfoResponse, error) {
	return &commonpb.ServiceInfoResponse{
		Name:          "user",
		Source:        "https://github.com/hexolan/stocklet",
		SourceLicense: "AGPL-3.0",
	}, nil
}

func (svc UserService) ViewUser(ctx context.Context, req *pb.ViewUserRequest) (*pb.ViewUserResponse, error) {
	// Validate the request args
	if err := svc.pbVal.Validate(req); err != nil {
		// Provide the validation error to the user.
		return nil, errors.NewServiceError(errors.ErrCodeInvalidArgument, "invalid request: "+err.Error())
	}

	// Get user from DB
	user, err := svc.store.GetUser(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return &pb.ViewUserResponse{User: user}, nil
}

func (svc UserService) RegisterUser(ctx context.Context, req *pb.RegisterUserRequest) (*pb.RegisterUserResponse, error) {
	// Validate the request args
	if err := svc.pbVal.Validate(req); err != nil {
		// Provide the validation error to the user.
		return nil, errors.NewServiceError(errors.ErrCodeInvalidArgument, "invalid request: "+err.Error())
	}

	// Attempt to register the user
	// This process involves calling the auth service to add an auth method for the user
	user, err := svc.store.RegisterUser(ctx, req.Email, req.Password, req.FirstName, req.LastName)
	if err != nil {
		return nil, err
	}

	return &pb.RegisterUserResponse{User: user}, nil
}
