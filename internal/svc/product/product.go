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

package product

import (
	"context"

	"github.com/bufbuild/protovalidate-go"
	"github.com/rs/zerolog/log"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/hexolan/stocklet/internal/pkg/errors"
	"github.com/hexolan/stocklet/internal/pkg/messaging"
	commonpb "github.com/hexolan/stocklet/internal/pkg/protogen/common/v1"
	eventpb "github.com/hexolan/stocklet/internal/pkg/protogen/events/v1"
	pb "github.com/hexolan/stocklet/internal/pkg/protogen/product/v1"
)

// Interface for the service
type ProductService struct {
	pb.UnimplementedProductServiceServer

	store StorageController
	pbVal *protovalidate.Validator
}

// Interface for database methods
// Flexibility for implementing seperate controllers for different databases (e.g. Postgres, MongoDB, etc)
type StorageController interface {
	GetProduct(ctx context.Context, productId string) (*pb.Product, error)
	GetProducts(ctx context.Context) ([]*pb.Product, error)

	UpdateProductPrice(ctx context.Context, productId string, price float32) (*pb.Product, error)
	DeleteProduct(ctx context.Context, productId string) error

	PriceOrderProducts(ctx context.Context, orderId string, customerId string, productQuantities map[string]int32) error
}

// Interface for event consumption
// Flexibility for seperate controllers for different messaging systems (e.g. Kafka, NATS, etc)
type ConsumerController interface {
	messaging.ConsumerController

	Attach(svc pb.ProductServiceServer)
}

// Create the product service
func NewProductService(cfg *ServiceConfig, store StorageController) *ProductService {
	// Initialise the protobuf validator
	pbVal, err := protovalidate.New()
	if err != nil {
		log.Panic().Err(err).Msg("failed to initialise protobuf validator")
	}

	// Initialise the service
	return &ProductService{
		store: store,
		pbVal: pbVal,
	}
}

func (svc ProductService) ServiceInfo(ctx context.Context, req *commonpb.ServiceInfoRequest) (*commonpb.ServiceInfoResponse, error) {
	return &commonpb.ServiceInfoResponse{
		Name:          "product",
		Source:        "https://github.com/hexolan/stocklet",
		SourceLicense: "AGPL-3.0",
	}, nil
}

func (svc ProductService) ViewProduct(ctx context.Context, req *pb.ViewProductRequest) (*pb.ViewProductResponse, error) {
	// Validate the request args
	if err := svc.pbVal.Validate(req); err != nil {
		// Provide the validation error to the user.
		return nil, errors.NewServiceError(errors.ErrCodeInvalidArgument, "invalid request: "+err.Error())
	}

	// Get product from DB
	product, err := svc.store.GetProduct(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return &pb.ViewProductResponse{Product: product}, nil
}

func (svc ProductService) ViewProducts(ctx context.Context, req *pb.ViewProductsRequest) (*pb.ViewProductsResponse, error) {
	// todo: pagination mechanism
	products, err := svc.store.GetProducts(ctx)
	if err != nil {
		return nil, err
	}

	return &pb.ViewProductsResponse{Products: products}, nil
}

func (svc ProductService) ProcessOrderCreatedEvent(ctx context.Context, req *eventpb.OrderCreatedEvent) (*emptypb.Empty, error) {
	err := svc.store.PriceOrderProducts(ctx, req.OrderId, req.CustomerId, req.ItemQuantities)
	if err != nil {
		return nil, errors.WrapServiceError(errors.ErrCodeExtService, "error processing event", err)
	}

	return &emptypb.Empty{}, nil
}
