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

package payment

import (
	"context"

	"github.com/bufbuild/protovalidate-go"
	"github.com/rs/zerolog/log"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/hexolan/stocklet/internal/pkg/errors"
	"github.com/hexolan/stocklet/internal/pkg/messaging"
	commonpb "github.com/hexolan/stocklet/internal/pkg/protogen/common/v1"
	eventpb "github.com/hexolan/stocklet/internal/pkg/protogen/events/v1"
	pb "github.com/hexolan/stocklet/internal/pkg/protogen/payment/v1"
)

// Interface for the service
type PaymentService struct {
	pb.UnimplementedPaymentServiceServer

	store StorageController
	pbVal *protovalidate.Validator
}

// Interface for database methods
// Flexibility for implementing seperate controllers for different databases (e.g. Postgres, MongoDB, etc)
type StorageController interface {
	GetBalance(ctx context.Context, customerId string) (*pb.CustomerBalance, error)
	GetTransaction(ctx context.Context, transactionId string) (*pb.Transaction, error)

	CreateBalance(ctx context.Context, customerId string) error
	CreditBalance(ctx context.Context, customerId string, amount float32) error
	DebitBalance(ctx context.Context, customerId string, amount float32, orderId *string) (*pb.Transaction, error)
	CloseBalance(ctx context.Context, customerId string) error

	PaymentForOrder(ctx context.Context, orderId string, customerId string, amount float32) error
}

// Interface for event consumption
// Flexibility for seperate controllers for different messaging systems (e.g. Kafka, NATS, etc)
type ConsumerController interface {
	messaging.ConsumerController

	Attach(svc pb.PaymentServiceServer)
}

// Create the payment service
func NewPaymentService(cfg *ServiceConfig, store StorageController) *PaymentService {
	// Initialise the protobuf validator
	pbVal, err := protovalidate.New()
	if err != nil {
		log.Panic().Err(err).Msg("failed to initialise protobuf validator")
	}

	// Initialise the service
	return &PaymentService{
		store: store,
		pbVal: pbVal,
	}
}

func (svc PaymentService) ServiceInfo(ctx context.Context, req *commonpb.ServiceInfoRequest) (*commonpb.ServiceInfoResponse, error) {
	return &commonpb.ServiceInfoResponse{
		Name:          "payment",
		Source:        "https://github.com/hexolan/stocklet",
		SourceLicense: "AGPL-3.0",
	}, nil
}

func (svc PaymentService) ViewTransaction(ctx context.Context, req *pb.ViewTransactionRequest) (*pb.ViewTransactionResponse, error) {
	// Attempt to get the transaction from the db
	transaction, err := svc.store.GetTransaction(ctx, req.TransactionId)
	if err != nil {
		return nil, err
	}

	return &pb.ViewTransactionResponse{Transaction: transaction}, nil
}

func (svc PaymentService) ViewBalance(ctx context.Context, req *pb.ViewBalanceRequest) (*pb.ViewBalanceResponse, error) {
	// todo: permission checking

	// Attempt to get the balance from the db
	balance, err := svc.store.GetBalance(ctx, req.CustomerId)
	if err != nil {
		return nil, err
	}

	return &pb.ViewBalanceResponse{Balance: balance}, nil
}

func (svc PaymentService) ProcessUserCreatedEvent(ctx context.Context, req *eventpb.UserCreatedEvent) (*emptypb.Empty, error) {
	err := svc.store.CreateBalance(ctx, req.UserId)
	if err != nil {
		return nil, errors.WrapServiceError(errors.ErrCodeExtService, "error processing event", err)
	}

	return &emptypb.Empty{}, nil
}

func (svc PaymentService) ProcessUserDeletedEvent(ctx context.Context, req *eventpb.UserDeletedEvent) (*emptypb.Empty, error) {
	err := svc.store.CloseBalance(ctx, req.UserId)
	if err != nil {
		return nil, errors.WrapServiceError(errors.ErrCodeExtService, "error processing event", err)
	}

	return &emptypb.Empty{}, nil
}

func (svc PaymentService) ProcessShipmentAllocationEvent(ctx context.Context, req *eventpb.ShipmentAllocationEvent) (*emptypb.Empty, error) {
	err := svc.store.PaymentForOrder(ctx, req.OrderId, req.OrderMetadata.CustomerId, req.OrderMetadata.TotalPrice)
	if err != nil {
		return nil, errors.WrapServiceError(errors.ErrCodeExtService, "error processing event", err)
	}

	return &emptypb.Empty{}, nil
}
