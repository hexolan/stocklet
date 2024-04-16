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

package main

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
	"github.com/twmb/franz-go/pkg/kgo"

	"github.com/hexolan/stocklet/internal/pkg/messaging"
	"github.com/hexolan/stocklet/internal/pkg/metrics"
	"github.com/hexolan/stocklet/internal/pkg/serve"
	"github.com/hexolan/stocklet/internal/pkg/storage"
	"github.com/hexolan/stocklet/internal/svc/shipping"
	"github.com/hexolan/stocklet/internal/svc/shipping/api"
	"github.com/hexolan/stocklet/internal/svc/shipping/controller"
)

func loadConfig() *shipping.ServiceConfig {
	// Load the core service configuration
	cfg, err := shipping.NewServiceConfig()
	if err != nil {
		log.Panic().Err(err).Msg("")
	}

	// Configure metrics (logging and OTEL)
	metrics.ConfigureLogger()
	metrics.InitTracerProvider(&cfg.Shared.Otel, "shipping")

	return cfg
}

func usePostgresController(cfg *shipping.ServiceConfig) (shipping.StorageController, *pgxpool.Pool) {
	// load the Postgres configuration
	if err := cfg.Postgres.Load(); err != nil {
		log.Panic().Err(err).Msg("")
	}

	// open a Postgres connection
	client, err := storage.NewPostgresConn(&cfg.Postgres)
	if err != nil {
		log.Panic().Err(err).Msg("")
	}

	controller := controller.NewPostgresController(client)
	return controller, client
}

func useKafkaController(cfg *shipping.ServiceConfig) (shipping.ConsumerController, *kgo.Client) {
	// load the Kafka configuration
	if err := cfg.Kafka.Load(); err != nil {
		log.Panic().Err(err).Msg("")
	}

	// open a Kafka connection
	client, err := messaging.NewKafkaConn(&cfg.Kafka, kgo.ConsumerGroup("shipping-service"))
	if err != nil {
		log.Panic().Err(err).Msg("")
	}

	controller := controller.NewKafkaController(client)
	return controller, client
}

func main() {
	cfg := loadConfig()

	// Create the storage controller
	store, storeCl := usePostgresController(cfg)
	defer storeCl.Close()

	// Create the service (& API interfaces)
	svc := shipping.NewShippingService(cfg, store)
	grpcSvr := api.PrepareGrpc(cfg, svc)
	gatewayMux := api.PrepareGateway(cfg)

	// Create the consumer
	consumer, consCl := useKafkaController(cfg)
	defer consCl.Close()
	consumer.Attach(svc)

	// Serve/start the interfaces
	go consumer.Start()
	go serve.Gateway(gatewayMux)
	serve.Grpc(grpcSvr)
}
