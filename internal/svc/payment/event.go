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
	"github.com/hexolan/stocklet/internal/pkg/messaging"
	eventspb "github.com/hexolan/stocklet/internal/pkg/protogen/events/v1"
	pb "github.com/hexolan/stocklet/internal/pkg/protogen/payment/v1"
)

func PrepareBalanceCreatedEvent(bal *pb.CustomerBalance) ([]byte, string, error) {
	topic := messaging.Payment_Balance_Created_Topic
	event := &eventspb.BalanceCreatedEvent{
		Revision: 1,

		CustomerId: bal.CustomerId,
		Balance:    bal.Balance,
	}

	return messaging.MarshalEvent(event, topic)
}

func PrepareBalanceCreditedEvent(customerId string, amount float32, newBalance float32) ([]byte, string, error) {
	topic := messaging.Payment_Balance_Credited_Topic
	event := &eventspb.BalanceCreditedEvent{
		Revision: 1,

		CustomerId: customerId,
		Amount:     amount,
		NewBalance: newBalance,
	}

	return messaging.MarshalEvent(event, topic)
}

func PrepareBalanceDebitedEvent(customerId string, amount float32, newBalance float32) ([]byte, string, error) {
	topic := messaging.Payment_Balance_Debited_Topic
	event := &eventspb.BalanceDebitedEvent{
		Revision: 1,

		CustomerId: customerId,
		Amount:     amount,
		NewBalance: newBalance,
	}

	return messaging.MarshalEvent(event, topic)
}

func PrepareBalanceClosedEvent(bal *pb.CustomerBalance) ([]byte, string, error) {
	topic := messaging.Payment_Balance_Closed_Topic
	event := &eventspb.BalanceClosedEvent{
		Revision: 1,

		CustomerId: bal.CustomerId,
		Balance:    bal.Balance,
	}

	return messaging.MarshalEvent(event, topic)
}

func PrepareTransactionLoggedEvent(transaction *pb.Transaction) ([]byte, string, error) {
	topic := messaging.Payment_Transaction_Created_Topic
	event := &eventspb.TransactionLoggedEvent{
		Revision: 1,

		TransactionId: transaction.Id,
		Amount:        transaction.Amount,
		OrderId:       transaction.OrderId,
		CustomerId:    transaction.CustomerId,
	}

	return messaging.MarshalEvent(event, topic)
}

func PrepareTransactionReversedEvent(transaction *pb.Transaction) ([]byte, string, error) {
	topic := messaging.Payment_Transaction_Reversed_Topic
	event := &eventspb.TransactionReversedEvent{
		Revision: 1,

		TransactionId: transaction.Id,
		Amount:        transaction.Amount,
		OrderId:       transaction.OrderId,
		CustomerId:    transaction.CustomerId,
	}

	return messaging.MarshalEvent(event, topic)
}

func PreparePaymentProcessedEvent_Success(transaction *pb.Transaction) ([]byte, string, error) {
	topic := messaging.Payment_Processing_Topic
	event := &eventspb.PaymentProcessedEvent{
		Revision: 1,

		Type:          eventspb.PaymentProcessedEvent_TYPE_SUCCESS,
		OrderId:       transaction.OrderId,
		CustomerId:    transaction.CustomerId,
		Amount:        transaction.Amount,
		TransactionId: &transaction.Id,
	}

	return messaging.MarshalEvent(event, topic)
}

func PreparePaymentProcessedEvent_Failure(orderId string, customerId string, amount float32) ([]byte, string, error) {
	topic := messaging.Payment_Processing_Topic
	event := &eventspb.PaymentProcessedEvent{
		Revision: 1,

		Type:       eventspb.PaymentProcessedEvent_TYPE_FAILED,
		OrderId:    orderId,
		CustomerId: customerId,
		Amount:     amount,
	}

	return messaging.MarshalEvent(event, topic)
}
