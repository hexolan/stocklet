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

package warehouse

import (
	"context"

	"github.com/bufbuild/protovalidate-go"
	"github.com/rs/zerolog/log"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/hexolan/stocklet/internal/pkg/errors"
	"github.com/hexolan/stocklet/internal/pkg/messaging"
	commonpb "github.com/hexolan/stocklet/internal/pkg/protogen/common/v1"
	eventpb "github.com/hexolan/stocklet/internal/pkg/protogen/events/v1"
	pb "github.com/hexolan/stocklet/internal/pkg/protogen/warehouse/v1"
)

// Interface for the service
type WarehouseService struct {
	pb.UnimplementedWarehouseServiceServer

	store StorageController
	pbVal *protovalidate.Validator
}

// Interface for database methods
// Flexibility for implementing seperate controllers for different databases (e.g. Postgres, MongoDB, etc)
type StorageController interface {
	GetProductStock(ctx context.Context, productId string) (*pb.ProductStock, error)
	GetReservation(ctx context.Context, reservationId string) (*pb.Reservation, error)

	CreateProductStock(ctx context.Context, productId string, startingQuantity int32) error

	ReserveOrderStock(ctx context.Context, orderId string, orderMetadata EventOrderMetadata, productQuantities map[string]int32) error
	ReturnReservedOrderStock(ctx context.Context, orderId string) error
	ConsumeReservedOrderStock(ctx context.Context, orderId string) error
}

// Interface for event consumption
// Flexibility for seperate controllers for different messaging systems (e.g. Kafka, NATS, etc)
type ConsumerController interface {
	messaging.ConsumerController

	Attach(svc pb.WarehouseServiceServer)
}

// Create the shipping service
func NewWarehouseService(cfg *ServiceConfig, store StorageController) *WarehouseService {
	// Initialise the protobuf validator
	pbVal, err := protovalidate.New()
	if err != nil {
		log.Panic().Err(err).Msg("failed to initialise protobuf validator")
	}

	// Initialise the service
	return &WarehouseService{
		store: store,
		pbVal: pbVal,
	}
}

func (svc WarehouseService) ServiceInfo(ctx context.Context, req *commonpb.ServiceInfoRequest) (*commonpb.ServiceInfoResponse, error) {
	return &commonpb.ServiceInfoResponse{
		Name:          "warehouse",
		Source:        "https://github.com/hexolan/stocklet",
		SourceLicense: "AGPL-3.0",
	}, nil
}

func (svc WarehouseService) ViewProductStock(ctx context.Context, req *pb.ViewProductStockRequest) (*pb.ViewProductStockResponse, error) {
	// Validate the request args
	if err := svc.pbVal.Validate(req); err != nil {
		// Provide the validation error to the user.
		return nil, errors.NewServiceError(errors.ErrCodeInvalidArgument, "invalid request: "+err.Error())
	}

	// Get stock from db
	stock, err := svc.store.GetProductStock(ctx, req.ProductId)
	if err != nil {
		return nil, err
	}

	return &pb.ViewProductStockResponse{Stock: stock}, nil
}

func (svc WarehouseService) ViewReservation(ctx context.Context, req *pb.ViewReservationRequest) (*pb.ViewReservationResponse, error) {
	// Validate the request args
	if err := svc.pbVal.Validate(req); err != nil {
		// Provide the validation error to the user.
		return nil, errors.NewServiceError(errors.ErrCodeInvalidArgument, "invalid request: "+err.Error())
	}

	// Get reservation from db
	reservation, err := svc.store.GetReservation(ctx, req.ReservationId)
	if err != nil {
		return nil, err
	}

	return &pb.ViewReservationResponse{Reservation: reservation}, nil
}

func (svc WarehouseService) ProcessProductCreatedEvent(ctx context.Context, req *eventpb.ProductCreatedEvent) (*emptypb.Empty, error) {
	err := svc.store.CreateProductStock(ctx, req.ProductId, 0)
	if err != nil {
		return nil, errors.WrapServiceError(errors.ErrCodeExtService, "error processing event", err)
	}

	return &emptypb.Empty{}, nil
}

func (svc WarehouseService) ProcessOrderPendingEvent(ctx context.Context, req *eventpb.OrderPendingEvent) (*emptypb.Empty, error) {
	err := svc.store.ReserveOrderStock(
		ctx,
		req.OrderId,
		EventOrderMetadata{
			CustomerId: req.CustomerId,
			ItemsPrice: req.ItemsPrice,
			TotalPrice: req.TotalPrice,
		},
		req.ItemQuantities,
	)
	if err != nil {
		return nil, errors.WrapServiceError(errors.ErrCodeExtService, "error processing event", err)
	}

	return &emptypb.Empty{}, nil
}

func (svc WarehouseService) ProcessShipmentAllocationEvent(ctx context.Context, req *eventpb.ShipmentAllocationEvent) (*emptypb.Empty, error) {
	if req.Type == eventpb.ShipmentAllocationEvent_TYPE_FAILED {
		err := svc.store.ReturnReservedOrderStock(ctx, req.OrderId)
		if err != nil {
			return nil, errors.WrapServiceError(errors.ErrCodeExtService, "error processing event", err)
		}
	}

	return &emptypb.Empty{}, nil
}

func (svc WarehouseService) ProcessPaymentProcessedEvent(ctx context.Context, req *eventpb.PaymentProcessedEvent) (*emptypb.Empty, error) {
	if req.Type == eventpb.PaymentProcessedEvent_TYPE_FAILED {
		err := svc.store.ReturnReservedOrderStock(ctx, req.OrderId)
		if err != nil {
			return nil, errors.WrapServiceError(errors.ErrCodeExtService, "error processing event", err)
		}
	} else if req.Type == eventpb.PaymentProcessedEvent_TYPE_SUCCESS {
		err := svc.store.ConsumeReservedOrderStock(ctx, req.OrderId)
		if err != nil {
			return nil, errors.WrapServiceError(errors.ErrCodeExtService, "error processing event", err)
		}
	}

	return &emptypb.Empty{}, nil
}
