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

package controller

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/twmb/franz-go/pkg/kgo"
	"google.golang.org/protobuf/proto"

	"github.com/hexolan/stocklet/internal/pkg/messaging"
	eventpb "github.com/hexolan/stocklet/internal/pkg/protogen/events/v1"
	pb "github.com/hexolan/stocklet/internal/pkg/protogen/order/v1"
	"github.com/hexolan/stocklet/internal/svc/order"
)

type kafkaController struct {
	cl *kgo.Client

	svc pb.OrderServiceServer

	ctx       context.Context
	ctxCancel context.CancelFunc
}

func NewKafkaController(cl *kgo.Client) order.ConsumerController {
	// Create a cancellable context for the consumer
	ctx, ctxCancel := context.WithCancel(context.Background())

	// Ensure the required Kafka topics exist
	err := messaging.EnsureKafkaTopics(
		cl,

		messaging.Order_State_Created_Topic,
		messaging.Order_State_Pending_Topic,
		messaging.Order_State_Rejected_Topic,
		messaging.Order_State_Approved_Topic,

		messaging.Warehouse_Reservation_Failed_Topic,
		messaging.Shipping_Shipment_Allocation_Topic,
		messaging.Payment_Processing_Topic,
	)
	if err != nil {
		log.Warn().Err(err).Msg("kafka: raised attempting to ensure svc topics")
	}

	// Add the consumption topics
	cl.AddConsumeTopics(
		messaging.Product_PriceQuotation_Topic,
		messaging.Warehouse_Reservation_Failed_Topic,
		messaging.Shipping_Shipment_Allocation_Topic,
		messaging.Payment_Processing_Topic,
	)

	return &kafkaController{cl: cl, ctx: ctx, ctxCancel: ctxCancel}
}

func (c *kafkaController) Attach(svc pb.OrderServiceServer) {
	c.svc = svc
}

func (c *kafkaController) Start() {
	if c.svc == nil {
		log.Panic().Msg("consumer: no service interface attached")
	}

	for {
		fetches := c.cl.PollFetches(c.ctx)
		if errs := fetches.Errors(); len(errs) > 0 {
			log.Panic().Any("kafka-errs", errs).Msg("consumer: unrecoverable kafka errors")
		}

		fetches.EachTopic(func(ft kgo.FetchTopic) {
			switch ft.Topic {
			case messaging.Product_PriceQuotation_Topic:
				c.consumeProductPriceQuoteEventTopic(ft)
			case messaging.Warehouse_Reservation_Failed_Topic:
				c.consumeStockReservationEventTopic(ft)
			case messaging.Shipping_Shipment_Allocation_Topic:
				c.consumeShipmentAllocationEventTopic(ft)
			case messaging.Payment_Processing_Topic:
				c.consumePaymentProcessedEventTopic(ft)
			default:
				log.Warn().Str("topic", ft.Topic).Msg("consumer: recieved records from unexpected topic")
			}
		})
	}
}

func (c *kafkaController) Stop() {
	// Cancel the consumer context
	c.ctxCancel()
}

func (c *kafkaController) consumeProductPriceQuoteEventTopic(ft kgo.FetchTopic) {
	log.Info().Str("topic", ft.Topic).Msg("consumer: recieved records from topic")

	// Process each message from the topic
	ft.EachRecord(func(record *kgo.Record) {
		// Unmarshal the event
		var event eventpb.ProductPriceQuoteEvent
		err := proto.Unmarshal(record.Value, &event)
		if err != nil {
			log.Panic().Err(err).Msg("consumer: failed to unmarshal event")
		}

		// Process the event
		ctx := context.Background()
		c.svc.ProcessProductPriceQuoteEvent(ctx, &event)
	})
}

func (c *kafkaController) consumeStockReservationEventTopic(ft kgo.FetchTopic) {
	log.Info().Str("topic", ft.Topic).Msg("consumer: recieved records from topic")

	// Process each message from the topic
	ft.EachRecord(func(record *kgo.Record) {
		// Unmarshal the event
		var event eventpb.StockReservationEvent
		err := proto.Unmarshal(record.Value, &event)
		if err != nil {
			log.Panic().Err(err).Msg("consumer: failed to unmarshal event")
		}

		// Process the event
		ctx := context.Background()
		c.svc.ProcessStockReservationEvent(ctx, &event)
	})
}

func (c *kafkaController) consumeShipmentAllocationEventTopic(ft kgo.FetchTopic) {
	log.Info().Str("topic", ft.Topic).Msg("consumer: recieved records from topic")

	// Process each message from the topic
	ft.EachRecord(func(record *kgo.Record) {
		// Unmarshal the event
		var event eventpb.ShipmentAllocationEvent
		err := proto.Unmarshal(record.Value, &event)
		if err != nil {
			log.Panic().Err(err).Msg("consumer: failed to unmarshal event")
		}

		// Process the event
		ctx := context.Background()
		c.svc.ProcessShipmentAllocationEvent(ctx, &event)
	})
}

func (c *kafkaController) consumePaymentProcessedEventTopic(ft kgo.FetchTopic) {
	log.Info().Str("topic", ft.Topic).Msg("consumer: recieved records from topic")

	// Process each message from the topic
	ft.EachRecord(func(record *kgo.Record) {
		// Unmarshal the event
		var event eventpb.PaymentProcessedEvent
		err := proto.Unmarshal(record.Value, &event)
		if err != nil {
			log.Panic().Err(err).Msg("consumer: failed to unmarshal event")
		}

		// Process the event
		ctx := context.Background()
		c.svc.ProcessPaymentProcessedEvent(ctx, &event)
	})
}
