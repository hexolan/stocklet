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

package order

import (
	"context"

	"github.com/bufbuild/protovalidate-go"
	"github.com/rs/zerolog/log"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/hexolan/stocklet/internal/pkg/errors"
	"github.com/hexolan/stocklet/internal/pkg/gwauth"
	"github.com/hexolan/stocklet/internal/pkg/messaging"
	commonpb "github.com/hexolan/stocklet/internal/pkg/protogen/common/v1"
	eventpb "github.com/hexolan/stocklet/internal/pkg/protogen/events/v1"
	pb "github.com/hexolan/stocklet/internal/pkg/protogen/order/v1"
)

// Interface for the service
type OrderService struct {
	pb.UnimplementedOrderServiceServer

	store StorageController
	pbVal *protovalidate.Validator
}

// Interface for database methods
// Flexibility for implementing seperate controllers for different databases (e.g. Postgres, MongoDB, etc)
type StorageController interface {
	GetOrder(ctx context.Context, orderId string) (*pb.Order, error)
	GetCustomerOrders(ctx context.Context, customerId string) ([]*pb.Order, error)

	CreateOrder(ctx context.Context, order *pb.Order) (*pb.Order, error)
	ApproveOrder(ctx context.Context, orderId string, transactionId string) (*pb.Order, error)
	ProcessOrder(ctx context.Context, orderId string, itemsPrice float32) (*pb.Order, error)
	RejectOrder(ctx context.Context, orderId string) (*pb.Order, error)
	SetOrderShipmentId(ctx context.Context, orderId string, shippingId string) error
}

// Interface for event consumption
// Flexibility for seperate controllers for different messaging systems (e.g. Kafka, NATS, etc)
type ConsumerController interface {
	messaging.ConsumerController

	Attach(svc pb.OrderServiceServer)
}

// Create the order service
func NewOrderService(cfg *ServiceConfig, store StorageController) *OrderService {
	// Initialise the protobuf validator
	pbVal, err := protovalidate.New()
	if err != nil {
		log.Panic().Err(err).Msg("failed to initialise protobuf validator")
	}

	// Initialise the service
	return &OrderService{
		store: store,
		pbVal: pbVal,
	}
}

func (svc OrderService) ServiceInfo(ctx context.Context, req *commonpb.ServiceInfoRequest) (*commonpb.ServiceInfoResponse, error) {
	return &commonpb.ServiceInfoResponse{
		Name:          "order",
		Source:        "https://github.com/hexolan/stocklet",
		SourceLicense: "AGPL-3.0",
	}, nil
}

func (svc OrderService) ViewOrder(ctx context.Context, req *pb.ViewOrderRequest) (*pb.ViewOrderResponse, error) {
	// Validate the request args
	if err := svc.pbVal.Validate(req); err != nil {
		// Provide the validation error to the user.
		return nil, errors.NewServiceError(errors.ErrCodeInvalidArgument, "invalid request: "+err.Error())
	}

	// Get the order from the DB
	order, err := svc.store.GetOrder(ctx, req.OrderId)
	if err != nil {
		return nil, err
	}

	return &pb.ViewOrderResponse{Order: order}, nil
}

func (svc OrderService) ViewOrders(ctx context.Context, req *pb.ViewOrdersRequest) (*pb.ViewOrdersResponse, error) {
	// Validate the request args
	if err := svc.pbVal.Validate(req); err != nil {
		// provide validation err context to user
		return nil, errors.NewServiceError(errors.ErrCodeInvalidArgument, "invalid request: "+err.Error())
	}

	// Get the orders from the storage controller
	orders, err := svc.store.GetCustomerOrders(ctx, req.CustomerId)
	if err != nil {
		return nil, err
	}

	return &pb.ViewOrdersResponse{Orders: orders}, nil
}

func (svc OrderService) PlaceOrder(ctx context.Context, req *pb.PlaceOrderRequest) (*pb.PlaceOrderResponse, error) {
	// If the request is through the gateway, then substitute req.CustomerId for current user
	gatewayRequest, gwMd := gwauth.IsGatewayRequest(ctx)
	if gatewayRequest {
		// ensure user is authenticated
		claims, err := gwauth.GetGatewayUser(gwMd)
		if err != nil {
			return nil, err
		}

		req.CustomerId = claims.Subject
	}

	// Validate the request args
	if err := svc.pbVal.Validate(req); err != nil {
		// provide validation err context to user
		return nil, errors.NewServiceError(errors.ErrCodeInvalidArgument, "invalid request: "+err.Error())
	}

	// Create the order.
	//
	// This will initiate a SAGA process involving
	// all the services required to create the order
	order, err := svc.store.CreateOrder(
		ctx,
		&pb.Order{
			Status:     pb.OrderStatus_ORDER_STATUS_PROCESSING,
			Items:      req.Cart,
			CustomerId: req.CustomerId,
		},
	)
	if err != nil {
		return nil, errors.WrapServiceError(errors.ErrCodeUnknown, "failed to create order", err)
	}

	// Return the pending order
	return &pb.PlaceOrderResponse{Order: order}, nil
}

func (svc OrderService) ProcessProductPriceQuoteEvent(ctx context.Context, req *eventpb.ProductPriceQuoteEvent) (*emptypb.Empty, error) {
	if req.Type == eventpb.ProductPriceQuoteEvent_TYPE_AVALIABLE {
		// Set order status to processing (from pending)
		// Dispatch OrderProcessingEvent
		_, err := svc.store.ProcessOrder(ctx, req.OrderId, req.TotalPrice)
		if err != nil {
			return nil, errors.WrapServiceError(errors.ErrCodeExtService, "failed to update in response to event", err)
		}

	} else if req.Type == eventpb.ProductPriceQuoteEvent_TYPE_UNAVALIABLE {
		// Set order status to rejected (from pending)
		// Dispatch OrderRejectedEvent
		_, err := svc.store.RejectOrder(ctx, req.OrderId)
		if err != nil {
			return nil, errors.WrapServiceError(errors.ErrCodeExtService, "failed to update in response to event", err)
		}
	} else {
		return nil, errors.NewServiceError(errors.ErrCodeInvalidArgument, "invalid event type")
	}

	return &emptypb.Empty{}, nil
}

func (svc OrderService) ProcessStockReservationEvent(ctx context.Context, req *eventpb.StockReservationEvent) (*emptypb.Empty, error) {
	if req.Type == eventpb.StockReservationEvent_TYPE_INSUFFICIENT_STOCK {
		// Set order status to rejected (from processing)
		// Dispatch OrderRejectedEvent
		_, err := svc.store.RejectOrder(ctx, req.OrderId)
		if err != nil {
			return nil, errors.WrapServiceError(errors.ErrCodeExtService, "failed to update in response to event", err)
		}
	}

	return &emptypb.Empty{}, nil
}

func (svc OrderService) ProcessShipmentAllocationEvent(ctx context.Context, req *eventpb.ShipmentAllocationEvent) (*emptypb.Empty, error) {
	if req.Type == eventpb.ShipmentAllocationEvent_TYPE_FAILED {
		// Set order status to rejected (from processing)
		// Dispatch OrderRejectedEvent
		_, err := svc.store.RejectOrder(ctx, req.OrderId)
		if err != nil {
			return nil, errors.WrapServiceError(errors.ErrCodeExtService, "failed to update in response to event", err)
		}
	} else if req.Type == eventpb.ShipmentAllocationEvent_TYPE_ALLOCATED {
		// Append shipment id to order
		err := svc.store.SetOrderShipmentId(ctx, req.OrderId, req.ShipmentId)
		if err != nil {
			return nil, errors.WrapServiceError(errors.ErrCodeExtService, "failed to update in response to event", err)
		}
	}

	return &emptypb.Empty{}, nil
}

func (svc OrderService) ProcessPaymentProcessedEvent(ctx context.Context, req *eventpb.PaymentProcessedEvent) (*emptypb.Empty, error) {
	if req.Type == eventpb.PaymentProcessedEvent_TYPE_SUCCESS {
		// Set order status to approved (from processing)
		// Dispatch OrderApprovedEvent
		_, err := svc.store.ApproveOrder(ctx, req.OrderId, *req.TransactionId)
		if err != nil {
			return nil, errors.WrapServiceError(errors.ErrCodeExtService, "failed to update in response to event", err)
		}
	} else if req.Type == eventpb.PaymentProcessedEvent_TYPE_FAILED {
		// Set order status to rejected (from processing)
		// Dispatch OrderRejectedEvent
		_, err := svc.store.RejectOrder(ctx, req.OrderId)
		if err != nil {
			return nil, errors.WrapServiceError(errors.ErrCodeExtService, "failed to update in response to event", err)
		}
	}

	return &emptypb.Empty{}, nil
}
