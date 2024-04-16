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

package metrics

import (
	"context"

	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	sdkresource "go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"

	"github.com/hexolan/stocklet/internal/pkg/config"
)

// Initiate the OpenTelemetry tracer provider
func InitTracerProvider(cfg *config.OtelConfig, svcName string) *sdktrace.TracerProvider {
	// Create resource and trace exporter (to otel-collector)
	resource := initTracerResource(svcName)
	exporter := initTracerExporter(cfg.CollectorGrpc)

	// Create the trace provider
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),

		sdktrace.WithResource(resource),
	)

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.TraceContext{})

	return tp
}

// Establishes a connection to otel-collector over gRPC
func initTracerExporter(collectorEndpoint string) sdktrace.SpanExporter {
	ctx := context.Background()

	exporter, err := otlptracegrpc.New(
		ctx,
		otlptracegrpc.WithEndpoint(collectorEndpoint),
		otlptracegrpc.WithInsecure(),
	)
	if err != nil {
		log.Panic().Err(err).Msg("otel: failed to start otlp gRPC exporter")
	}

	return exporter
}

// Prepare a tracer resource to use with the tracing provider
func initTracerResource(svcName string) *sdkresource.Resource {
	ctx := context.Background()

	resource, err := sdkresource.New(
		ctx,
		sdkresource.WithAttributes(
			semconv.ServiceName(svcName),
		),
	)
	if err != nil {
		log.Panic().Err(err).Msg("otel: failed to create tracer resource")
	}

	return resource
}
