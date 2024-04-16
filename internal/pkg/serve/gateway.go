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
	"context"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"

	"github.com/hexolan/stocklet/internal/pkg/config"
	"github.com/hexolan/stocklet/internal/pkg/gwauth"
)

func withGatewayErrorHandler() runtime.ServeMuxOption {
	return runtime.WithErrorHandler(
		func(ctx context.Context, mux *runtime.ServeMux, marshaler runtime.Marshaler, w http.ResponseWriter, r *http.Request, err error) {
			runtime.DefaultHTTPErrorHandler(ctx, mux, marshaler, w, r, err)
			log.Error().Err(err).Stack().Str("path", r.URL.Path).Str("reqURI", r.RequestURI).Str("remoteAddr", r.RemoteAddr).Msg("")
		},
	)
}

func withGatewayMetadataOpt() runtime.ServeMuxOption {
	return runtime.WithMetadata(
		func(ctx context.Context, req *http.Request) metadata.MD {
			return metadata.MD{"from-gateway": {"true"}}
		},
	)
}

func withGatewayHeaderOpt() runtime.ServeMuxOption {
	return runtime.WithIncomingHeaderMatcher(
		func(key string) (string, bool) {
			switch key {
			case gwauth.JWTPayloadHeader:
				// Envoy will validate JWT tokens and provide a payload header
				// containing a base64 string of the token claims.
				return "jwt-payload", true
			default:
				return key, false
			}
		},
	)
}

func withGatewayLogger(h http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			h.ServeHTTP(w, r)
			log.Info().Str("path", r.URL.Path).Str("reqURI", r.RequestURI).Str("remoteAddr", r.RemoteAddr).Msg("")
		},
	)
}

func NewGatewayServeBase(cfg *config.SharedConfig) (*runtime.ServeMux, []grpc.DialOption) {
	// Create the base runtime ServeMux
	mux := runtime.NewServeMux(
		withGatewayErrorHandler(),
		withGatewayMetadataOpt(),
		withGatewayHeaderOpt(),
	)

	// Attach open telemetry instrumentation through the gRPC client options
	clientOpts := []grpc.DialOption{
		grpc.WithStatsHandler(
			otelgrpc.NewClientHandler(),
		),

		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	return mux, clientOpts
}

func Gateway(mux *runtime.ServeMux) error {
	// Create OTEL instrumentation handler
	handler := otelhttp.NewHandler(
		mux,
		"grpc-gateway",
		otelhttp.WithSpanNameFormatter(
			func(operation string, r *http.Request) string {
				return operation + ": " + r.RequestURI
			},
		),
	)

	// Create gateway HTTP server
	svr := &http.Server{
		Addr:    GetAddrToGateway("0.0.0.0"),
		Handler: withGatewayLogger(handler),
	}

	return svr.ListenAndServe()
}
