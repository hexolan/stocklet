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
	"github.com/hexolan/stocklet/internal/pkg/messaging"
	eventspb "github.com/hexolan/stocklet/internal/pkg/protogen/events/v1"
	pb "github.com/hexolan/stocklet/internal/pkg/protogen/order/v1"
)

func PrepareOrderCreatedEvent(order *pb.Order) ([]byte, string, error) {
	topic := messaging.Order_State_Created_Topic
	event := &eventspb.OrderCreatedEvent{
		Revision: 1,

		OrderId:        order.Id,
		CustomerId:     order.CustomerId,
		ItemQuantities: order.Items,
	}

	return messaging.MarshalEvent(event, topic)
}

func PrepareOrderPendingEvent(order *pb.Order) ([]byte, string, error) {
	topic := messaging.Order_State_Pending_Topic
	event := &eventspb.OrderPendingEvent{
		Revision: 1,

		OrderId:        order.Id,
		CustomerId:     order.CustomerId,
		ItemQuantities: order.Items,
	}

	return messaging.MarshalEvent(event, topic)
}

func PrepareOrderRejectedEvent(order *pb.Order) ([]byte, string, error) {
	topic := messaging.Order_State_Rejected_Topic
	event := &eventspb.OrderRejectedEvent{
		Revision: 1,

		OrderId:       order.Id,
		TransactionId: order.TransactionId,
		ShippingId:    order.ShippingId,
	}

	return messaging.MarshalEvent(event, topic)
}

func PrepareOrderApprovedEvent(order *pb.Order) ([]byte, string, error) {
	topic := messaging.Order_State_Approved_Topic
	event := &eventspb.OrderApprovedEvent{
		Revision: 1,

		OrderId:       order.Id,
		TransactionId: order.GetTransactionId(),
		ShippingId:    order.GetShippingId(),
	}

	return messaging.MarshalEvent(event, topic)
}
