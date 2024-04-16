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
	pb "github.com/hexolan/stocklet/internal/pkg/protogen/product/v1"
	"github.com/hexolan/stocklet/internal/svc/product"
)

type kafkaController struct {
	cl *kgo.Client

	svc pb.ProductServiceServer

	ctx       context.Context
	ctxCancel context.CancelFunc
}

func NewKafkaController(cl *kgo.Client) product.ConsumerController {
	// Create a cancellable context for the consumer
	ctx, ctxCancel := context.WithCancel(context.Background())

	// Ensure the required Kafka topics exist
	err := messaging.EnsureKafkaTopics(
		cl,

		messaging.Product_State_Created_Topic,
		messaging.Product_State_Deleted_Topic,
		messaging.Product_Attribute_Price_Topic,
		messaging.Product_PriceQuotation_Topic,

		messaging.Order_State_Created_Topic,
	)
	if err != nil {
		log.Warn().Err(err).Msg("kafka: raised attempting to ensure svc topics")
	}

	// Add the consumption topics
	cl.AddConsumeTopics(
		messaging.Order_State_Created_Topic,
	)

	return &kafkaController{cl: cl, ctx: ctx, ctxCancel: ctxCancel}
}

func (c *kafkaController) Attach(svc pb.ProductServiceServer) {
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
			case messaging.Order_State_Created_Topic:
				c.consumeOrderCreatedEventTopic(ft)
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

func (c *kafkaController) consumeOrderCreatedEventTopic(ft kgo.FetchTopic) {
	log.Info().Str("topic", ft.Topic).Msg("consumer: recieved records from topic")

	// Process each message from the topic
	ft.EachRecord(func(record *kgo.Record) {
		// Unmarshal the event
		var event eventpb.OrderCreatedEvent
		err := proto.Unmarshal(record.Value, &event)
		if err != nil {
			log.Panic().Err(err).Msg("consumer: failed to unmarshal event")
		}

		// Process the event
		ctx := context.Background()
		c.svc.ProcessOrderCreatedEvent(ctx, &event)
	})
}
