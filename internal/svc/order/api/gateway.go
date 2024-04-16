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

package api

import (
	"context"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/rs/zerolog/log"

	pb "github.com/hexolan/stocklet/internal/pkg/protogen/order/v1"
	"github.com/hexolan/stocklet/internal/pkg/serve"
	"github.com/hexolan/stocklet/internal/svc/order"
)

func PrepareGateway(cfg *order.ServiceConfig) *runtime.ServeMux {
	mux, clientOpts := serve.NewGatewayServeBase(&cfg.Shared)

	ctx := context.Background()
	err := pb.RegisterOrderServiceHandlerFromEndpoint(ctx, mux, serve.GetAddrToGrpc("localhost"), clientOpts)
	if err != nil {
		log.Panic().Err(err).Msg("failed to register endpoint for gateway")
	}

	return mux
}
