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

package shipping

import (
	"context"

	"github.com/bufbuild/protovalidate-go"
	"github.com/rs/zerolog/log"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/hexolan/stocklet/internal/pkg/errors"
	"github.com/hexolan/stocklet/internal/pkg/messaging"
	commonpb "github.com/hexolan/stocklet/internal/pkg/protogen/common/v1"
	eventpb "github.com/hexolan/stocklet/internal/pkg/protogen/events/v1"
	pb "github.com/hexolan/stocklet/internal/pkg/protogen/shipping/v1"
)

// Interface for the service
type ShippingService struct {
	pb.UnimplementedShippingServiceServer

	store StorageController
	pbVal *protovalidate.Validator
}

// Interface for database methods
// Flexibility for implementing seperate controllers for different databases (e.g. Postgres, MongoDB, etc)
type StorageController interface {
	GetShipment(ctx context.Context, shipmentId string) (*pb.Shipment, error)
	GetShipmentItems(ctx context.Context, shipmentId string) ([]*pb.ShipmentItem, error)

	AllocateOrderShipment(ctx context.Context, orderId string, orderMetadata EventOrderMetadata, productQuantities map[string]int32) error
	CancelOrderShipment(ctx context.Context, orderId string) error
}

// Interface for event consumption
// Flexibility for seperate controllers for different messaging systems (e.g. Kafka, NATS, etc)
type ConsumerController interface {
	messaging.ConsumerController

	Attach(svc pb.ShippingServiceServer)
}

// Create the shipping service
func NewShippingService(cfg *ServiceConfig, store StorageController) *ShippingService {
	// Initialise the protobuf validator
	pbVal, err := protovalidate.New()
	if err != nil {
		log.Panic().Err(err).Msg("failed to initialise protobuf validator")
	}

	// Initialise the service
	return &ShippingService{
		store: store,
		pbVal: pbVal,
	}
}

func (svc ShippingService) ServiceInfo(ctx context.Context, req *commonpb.ServiceInfoRequest) (*commonpb.ServiceInfoResponse, error) {
	return &commonpb.ServiceInfoResponse{
		Name:          "shipping",
		Source:        "https://github.com/hexolan/stocklet",
		SourceLicense: "AGPL-3.0",
	}, nil
}

func (svc ShippingService) ViewShipment(ctx context.Context, req *pb.ViewShipmentRequest) (*pb.ViewShipmentResponse, error) {
	// Validate the request args
	if err := svc.pbVal.Validate(req); err != nil {
		// Provide the validation error to the user.
		return nil, errors.NewServiceError(errors.ErrCodeInvalidArgument, "invalid request: "+err.Error())
	}

	// todo: permission checking?

	// Get shipment from DB
	shipment, err := svc.store.GetShipment(ctx, req.ShipmentId)
	if err != nil {
		return nil, err
	}

	return &pb.ViewShipmentResponse{Shipment: shipment}, nil
}

func (svc ShippingService) ViewShipmentManifest(ctx context.Context, req *pb.ViewShipmentManifestRequest) (*pb.ViewShipmentManifestResponse, error) {
	// Validate the request args
	if err := svc.pbVal.Validate(req); err != nil {
		// Provide the validation error to the user.
		return nil, errors.NewServiceError(errors.ErrCodeInvalidArgument, "invalid request: "+err.Error())
	}

	// todo: permission checking?

	shipmentItems, err := svc.store.GetShipmentItems(ctx, req.ShipmentId)
	if err != nil {
		return nil, err
	}

	return &pb.ViewShipmentManifestResponse{Manifest: shipmentItems}, nil
}

func (svc ShippingService) ProcessStockReservationEvent(ctx context.Context, req *eventpb.StockReservationEvent) (*emptypb.Empty, error) {
	if req.Type == eventpb.StockReservationEvent_TYPE_STOCK_RESERVED {
		err := svc.store.AllocateOrderShipment(
			ctx,
			req.OrderId,
			EventOrderMetadata{
				CustomerId: req.OrderMetadata.CustomerId,
				ItemsPrice: req.OrderMetadata.ItemsPrice,
				TotalPrice: req.OrderMetadata.TotalPrice,
			},
			req.ReservationStock,
		)
		if err != nil {
			return nil, errors.WrapServiceError(errors.ErrCodeExtService, "failed to update in response to event", err)
		}

	}

	return &emptypb.Empty{}, nil
}

func (svc ShippingService) ProcessPaymentProcessedEvent(ctx context.Context, req *eventpb.PaymentProcessedEvent) (*emptypb.Empty, error) {
	if req.Type == eventpb.PaymentProcessedEvent_TYPE_FAILED {
		err := svc.store.CancelOrderShipment(ctx, req.OrderId)
		if err != nil {
			return nil, errors.WrapServiceError(errors.ErrCodeExtService, "failed to update in response to event", err)
		}

	}

	return &emptypb.Empty{}, nil
}
