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

package messaging

import (
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

type ConsumerController interface {
	Start()
	Stop()
}

// Topic Definitions
const (
	// Order Topics
	Order_State_Topic          = "order.state"
	Order_State_Created_Topic  = Order_State_Topic + ".created"
	Order_State_Pending_Topic  = Order_State_Topic + ".pending"
	Order_State_Rejected_Topic = Order_State_Topic + ".rejected"
	Order_State_Approved_Topic = Order_State_Topic + ".approved"

	// Payment Topics
	Payment_Balance_Topic          = "payment.balance"
	Payment_Balance_Created_Topic  = Payment_Balance_Topic + ".created"
	Payment_Balance_Credited_Topic = Payment_Balance_Topic + ".credited"
	Payment_Balance_Debited_Topic  = Payment_Balance_Topic + ".debited"
	Payment_Balance_Closed_Topic   = Payment_Balance_Topic + ".closed"

	Payment_Transaction_Topic          = "payment.transaction"
	Payment_Transaction_Created_Topic  = Payment_Transaction_Topic + ".created"
	Payment_Transaction_Reversed_Topic = Payment_Transaction_Topic + ".reversed"

	Payment_Processing_Topic = "payment.processing"

	// Product Topics
	Product_State_Topic         = "product.state"
	Product_State_Created_Topic = Product_State_Topic + ".created"
	Product_State_Deleted_Topic = Product_State_Topic + ".deleted"

	Product_Attribute_Topic       = "product.attr"
	Product_Attribute_Price_Topic = Product_Attribute_Topic + ".price"

	Product_PriceQuotation_Topic = "product.pricequotation"

	// Shipping Topics
	Shipping_Shipment_Topic            = "shipping.shipment"
	Shipping_Shipment_Allocation_Topic = Shipping_Shipment_Topic + ".allocation"
	Shipping_Shipment_Dispatched_Topic = Shipping_Shipment_Topic + ".dispatched"

	// User Topics
	User_State_Topic         = "user.state"
	User_State_Created_Topic = User_State_Topic + ".created"
	User_State_Deleted_Topic = User_State_Topic + ".deleted"

	User_Attribute_Topic       = "user.attr"
	User_Attribute_Email_Topic = User_Attribute_Topic + ".email"

	// Warehouse Topics
	Warehouse_Stock_Topic         = "warehouse.stock"
	Warehouse_Stock_Created_Topic = Warehouse_Stock_Topic + ".created"
	Warehouse_Stock_Added_Topic   = Warehouse_Stock_Topic + ".added"
	Warehouse_Stock_Removed_Topic = Warehouse_Stock_Topic + ".removed"

	Warehouse_Reservation_Topic          = "warehouse.reservation"
	Warehouse_Reservation_Failed_Topic   = Warehouse_Reservation_Topic + ".failed"
	Warehouse_Reservation_Reserved_Topic = Warehouse_Reservation_Topic + ".reserved"
	Warehouse_Reservation_Returned_Topic = Warehouse_Reservation_Topic + ".returned"
	Warehouse_Reservation_Consumed_Topic = Warehouse_Reservation_Topic + ".consumed"
)

func MarshalEvent(event protoreflect.ProtoMessage, topic string) ([]byte, string, error) {
	wireEvent, err := proto.Marshal(event)
	if err != nil {
		return []byte{}, "", err
	}

	return wireEvent, topic, nil
}
