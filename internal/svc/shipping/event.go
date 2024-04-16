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
	"github.com/hexolan/stocklet/internal/pkg/messaging"
	eventspb "github.com/hexolan/stocklet/internal/pkg/protogen/events/v1"
	pb "github.com/hexolan/stocklet/internal/pkg/protogen/shipping/v1"
)

type EventOrderMetadata struct {
	CustomerId string
	ItemsPrice float32
	TotalPrice float32
}

func PrepareShipmentAllocationEvent_Failed(orderId string, orderMetadata EventOrderMetadata, productQuantities map[string]int32) ([]byte, string, error) {
	topic := messaging.Shipping_Shipment_Allocation_Topic
	event := &eventspb.ShipmentAllocationEvent{
		Revision: 1,

		Type:    eventspb.ShipmentAllocationEvent_TYPE_FAILED,
		OrderId: orderId,
		OrderMetadata: &eventspb.ShipmentAllocationEvent_OrderMetadata{
			CustomerId: orderMetadata.CustomerId,
			ItemsPrice: orderMetadata.ItemsPrice,
			TotalPrice: orderMetadata.TotalPrice,
		},
		ProductQuantities: productQuantities,
	}

	return messaging.MarshalEvent(event, topic)
}

func PrepareShipmentAllocationEvent_Allocated(orderId string, orderMetadata EventOrderMetadata, shipmentId string, productQuantities map[string]int32) ([]byte, string, error) {
	topic := messaging.Shipping_Shipment_Allocation_Topic
	event := &eventspb.ShipmentAllocationEvent{
		Revision: 1,

		Type:    eventspb.ShipmentAllocationEvent_TYPE_ALLOCATED,
		OrderId: orderId,
		OrderMetadata: &eventspb.ShipmentAllocationEvent_OrderMetadata{
			CustomerId: orderMetadata.CustomerId,
			ItemsPrice: orderMetadata.ItemsPrice,
			TotalPrice: orderMetadata.TotalPrice,
		},
		ShipmentId:        shipmentId,
		ProductQuantities: productQuantities,
	}

	return messaging.MarshalEvent(event, topic)
}

func PrepareShipmentAllocationEvent_AllocationReleased(orderId string, shipmentId string, shipmentItems []*pb.ShipmentItem) ([]byte, string, error) {
	productQuantities := make(map[string]int32)
	for _, item := range shipmentItems {
		productQuantities[item.ProductId] = item.Quantity
	}

	topic := messaging.Shipping_Shipment_Allocation_Topic
	event := &eventspb.ShipmentAllocationEvent{
		Revision: 1,

		Type:              eventspb.ShipmentAllocationEvent_TYPE_ALLOCATION_RELEASED,
		OrderId:           orderId,
		ShipmentId:        shipmentId,
		ProductQuantities: productQuantities,
	}

	return messaging.MarshalEvent(event, topic)
}

func PrepareShipmentDispatchedEvent(orderId string, shipmentId string, productQuantities map[string]int32) ([]byte, string, error) {
	topic := messaging.Shipping_Shipment_Dispatched_Topic
	event := &eventspb.ShipmentDispatchedEvent{
		Revision: 1,

		OrderId:           orderId,
		ShipmentId:        shipmentId,
		ProductQuantities: productQuantities,
	}

	return messaging.MarshalEvent(event, topic)
}
