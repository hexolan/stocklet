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
	"github.com/hexolan/stocklet/internal/pkg/messaging"
	eventspb "github.com/hexolan/stocklet/internal/pkg/protogen/events/v1"
	pb "github.com/hexolan/stocklet/internal/pkg/protogen/product/v1"
)

func PrepareProductCreatedEvent(product *pb.Product) ([]byte, string, error) {
	topic := messaging.Product_State_Created_Topic
	event := &eventspb.ProductCreatedEvent{
		Revision: 1,

		ProductId:   product.Id,
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
	}

	return messaging.MarshalEvent(event, topic)
}

func PrepareProductPriceUpdatedEvent(product *pb.Product) ([]byte, string, error) {
	topic := messaging.Product_Attribute_Price_Topic
	event := &eventspb.ProductPriceUpdatedEvent{
		Revision: 1,

		ProductId: product.Id,
		Price:     product.Price,
	}

	return messaging.MarshalEvent(event, topic)
}

func PrepareProductDeletedEvent(product *pb.Product) ([]byte, string, error) {
	topic := messaging.Product_State_Deleted_Topic
	event := &eventspb.ProductDeletedEvent{
		Revision: 1,

		ProductId: product.Id,
	}

	return messaging.MarshalEvent(event, topic)
}

func PrepareProductPriceQuoteEvent_Avaliable(orderId string, productQuantities map[string]int32, productPrices map[string]float32, totalPrice float32) ([]byte, string, error) {
	topic := messaging.Product_PriceQuotation_Topic
	event := &eventspb.ProductPriceQuoteEvent{
		Revision: 1,

		Type:              eventspb.ProductPriceQuoteEvent_TYPE_AVALIABLE,
		OrderId:           orderId,
		ProductQuantities: productQuantities,
		ProductPrices:     productPrices,
		TotalPrice:        totalPrice,
	}

	return messaging.MarshalEvent(event, topic)
}

func PrepareProductPriceQuoteEvent_Unavaliable(orderId string) ([]byte, string, error) {
	topic := messaging.Product_PriceQuotation_Topic
	event := &eventspb.ProductPriceQuoteEvent{
		Revision: 1,

		Type:    eventspb.ProductPriceQuoteEvent_TYPE_UNAVALIABLE,
		OrderId: orderId,
	}

	return messaging.MarshalEvent(event, topic)
}
