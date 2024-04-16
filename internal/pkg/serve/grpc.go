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

package serve

import (
	"net"

	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"

	"github.com/hexolan/stocklet/internal/pkg/config"
)

func NewGrpcServeBase(cfg *config.SharedConfig) *grpc.Server {
	// Attach OTEL metrics middleware
	svr := grpc.NewServer(
		grpc.StatsHandler(
			otelgrpc.NewServerHandler(),
		),
	)

	// Attach the health service
	svc := health.NewServer()
	grpc_health_v1.RegisterHealthServer(svr, svc)

	// Enable reflection in dev mode
	// Eases usage of tools like grpcurl and grpcui
	if cfg.DevMode {
		reflection.Register(svr)
	}

	return svr
}

func Grpc(svr *grpc.Server) {
	lis, err := net.Listen("tcp", GetAddrToGrpc("0.0.0.0"))
	if err != nil {
		log.Panic().Err(err).Str("port", grpcPort).Msg("failed to listen on gRPC port")
	}

	err = svr.Serve(lis)
	if err != nil {
		log.Panic().Err(err).Msg("failed to serve gRPC server")
	}
}
