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
	"github.com/hexolan/stocklet/internal/pkg/messaging"
	eventspb "github.com/hexolan/stocklet/internal/pkg/protogen/events/v1"
	pb "github.com/hexolan/stocklet/internal/pkg/protogen/warehouse/v1"
)

type EventOrderMetadata struct {
	CustomerId string
	ItemsPrice float32
	TotalPrice float32
}

func PrepareStockCreatedEvent(productStock *pb.ProductStock) ([]byte, string, error) {
	topic := messaging.Warehouse_Stock_Created_Topic
	event := &eventspb.StockCreatedEvent{
		Revision: 1,

		ProductId: productStock.ProductId,
		Quantity:  productStock.Quantity,
	}

	return messaging.MarshalEvent(event, topic)
}

func PrepareStockAddedEvent(productId string, amount int32, reservationId *string) ([]byte, string, error) {
	topic := messaging.Warehouse_Stock_Added_Topic
	event := &eventspb.StockAddedEvent{
		Revision: 1,

		ProductId:     productId,
		Amount:        amount,
		ReservationId: reservationId,
	}

	return messaging.MarshalEvent(event, topic)
}

func PrepareStockRemovedEvent(productId string, amount int32, reservationId *string) ([]byte, string, error) {
	topic := messaging.Warehouse_Stock_Removed_Topic
	event := &eventspb.StockRemovedEvent{
		Revision: 1,

		ProductId:     productId,
		Amount:        amount,
		ReservationId: reservationId,
	}

	return messaging.MarshalEvent(event, topic)
}

func PrepareStockReservationEvent_Failed(orderId string, orderMetadata EventOrderMetadata, insufficientStockProductIds []string) ([]byte, string, error) {
	topic := messaging.Warehouse_Reservation_Failed_Topic
	event := &eventspb.StockReservationEvent{
		Revision: 1,

		Type:    eventspb.StockReservationEvent_TYPE_INSUFFICIENT_STOCK,
		OrderId: orderId,
		OrderMetadata: &eventspb.StockReservationEvent_OrderMetadata{
			CustomerId: orderMetadata.CustomerId,
			ItemsPrice: orderMetadata.ItemsPrice,
			TotalPrice: orderMetadata.TotalPrice,
		},
		InsufficientStock: insufficientStockProductIds,
	}

	return messaging.MarshalEvent(event, topic)
}

func PrepareStockReservationEvent_Reserved(orderId string, orderMetadata EventOrderMetadata, reservationId string, reservationStock map[string]int32) ([]byte, string, error) {
	topic := messaging.Warehouse_Reservation_Reserved_Topic
	event := &eventspb.StockReservationEvent{
		Revision: 1,

		Type:    eventspb.StockReservationEvent_TYPE_STOCK_RESERVED,
		OrderId: orderId,
		OrderMetadata: &eventspb.StockReservationEvent_OrderMetadata{
			CustomerId: orderMetadata.CustomerId,
			ItemsPrice: orderMetadata.ItemsPrice,
			TotalPrice: orderMetadata.TotalPrice,
		},
		ReservationId:    reservationId,
		ReservationStock: reservationStock,
	}

	return messaging.MarshalEvent(event, topic)
}

func PrepareStockReservationEvent_Returned(orderId string, reservationId string, reservedStock []*pb.ReservationStock) ([]byte, string, error) {
	reservationStock := make(map[string]int32)
	for _, item := range reservedStock {
		reservationStock[item.ProductId] = item.Quantity
	}

	topic := messaging.Warehouse_Reservation_Returned_Topic
	event := &eventspb.StockReservationEvent{
		Revision: 1,

		Type:             eventspb.StockReservationEvent_TYPE_STOCK_RESERVED,
		OrderId:          orderId,
		ReservationId:    reservationId,
		ReservationStock: reservationStock,
	}

	return messaging.MarshalEvent(event, topic)
}

func PrepareStockReservationEvent_Consumed(orderId string, reservationId string, reservedStock []*pb.ReservationStock) ([]byte, string, error) {
	reservationStock := make(map[string]int32)
	for _, item := range reservedStock {
		reservationStock[item.ProductId] = item.Quantity
	}

	topic := messaging.Warehouse_Reservation_Consumed_Topic
	event := &eventspb.StockReservationEvent{
		Revision: 1,

		Type:             eventspb.StockReservationEvent_TYPE_STOCK_CONSUMED,
		OrderId:          orderId,
		ReservationId:    reservationId,
		ReservationStock: reservationStock,
	}

	return messaging.MarshalEvent(event, topic)
}
