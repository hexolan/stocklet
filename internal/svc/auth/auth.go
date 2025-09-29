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

package auth

import (
	"context"

	"buf.build/go/protovalidate"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/hexolan/stocklet/internal/pkg/errors"
	"github.com/hexolan/stocklet/internal/pkg/gwauth"
	"github.com/hexolan/stocklet/internal/pkg/messaging"
	pb "github.com/hexolan/stocklet/internal/pkg/protogen/auth/v1"
	commonpb "github.com/hexolan/stocklet/internal/pkg/protogen/common/v1"
	eventpb "github.com/hexolan/stocklet/internal/pkg/protogen/events/v1"
)

// Interface for the service
type AuthService struct {
	pb.UnimplementedAuthServiceServer

	cfg *ServiceConfig

	store StorageController
	pbVal *protovalidate.Validator
}

// Interface for database methods
// Allows implementing separate controllers for different databases (e.g. Postgres, MongoDB, etc)
type StorageController interface {
	SetPassword(ctx context.Context, userId string, password string) error
	VerifyPassword(ctx context.Context, userId string, password string) (bool, error)

	DeleteAuthMethods(ctx context.Context, userId string) error
}

// Interface for event consumption
// Flexibility for separate controllers for different messaging systems (e.g. Kafka, NATS, etc)
type ConsumerController interface {
	messaging.ConsumerController

	Attach(svc pb.AuthServiceServer)
}

// Create the auth service
func NewAuthService(cfg *ServiceConfig, store StorageController) *AuthService {
	// Initialise the protobuf validator
	pbVal, err := protovalidate.New()
	if err != nil {
		log.Panic().Err(err).Msg("failed to initialise protobuf validator")
	}

	svc := &AuthService{
		cfg:   cfg,
		store: store,
		pbVal: &pbVal,
	}

	return svc
}

func (svc AuthService) ServiceInfo(ctx context.Context, req *commonpb.ServiceInfoRequest) (*commonpb.ServiceInfoResponse, error) {
	return &commonpb.ServiceInfoResponse{
		Name:          "auth",
		Source:        "https://github.com/hexolan/stocklet",
		SourceLicense: "AGPL-3.0",
	}, nil
}

func (svc AuthService) GetOpenIDProviderConfig(ctx context.Context, req *pb.GetOpenIDProviderConfigRequest) (*pb.GetOpenIDProviderConfigResponse, error) {
	return &pb.GetOpenIDProviderConfigResponse{
		Issuer:                            svc.cfg.ServiceOpts.ApiPublicUrl + "/v1/auth",
		AuthorizationEndpoint:             svc.cfg.ServiceOpts.ApiPublicUrl + "/v1/auth/oauth/authorize",
		TokenEndpoint:                     svc.cfg.ServiceOpts.ApiPublicUrl + "/v1/auth/oauth/token",
		IntrospectionEndpoint:             svc.cfg.ServiceOpts.ApiPublicUrl + "/v1/auth/oauth/introspect",
		RevocationEndpoint:                svc.cfg.ServiceOpts.ApiPublicUrl + "/v1/auth/oauth/revoke",
		UserinfoEndpoint:                  svc.cfg.ServiceOpts.ApiPublicUrl + "/v1/auth/oidc/userinfo",
		JwksUri:                           svc.cfg.ServiceOpts.ApiPublicUrl + "/v1/auth/.well-known/jwks",
		TokenEndpointAuthMethodsSupported: []string{},
		ScopesSupported:                   []string{},
		ClaimsSupported:                   []string{},
	}, nil
}

func (svc AuthService) OAuthAuthorize(ctx context.Context, req *pb.OAuthAuthorizeRequest) (*pb.OAuthAuthorizeResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method OAuthAuthorize not implemented")
}

func (svc AuthService) OAuthToken(ctx context.Context, req *pb.OAuthTokenRequest) (*pb.OAuthTokenResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method OAuthToken not implemented")
}

func (svc AuthService) OAuthTokenRevocation(ctx context.Context, req *pb.OAuthTokenRevocationRequest) (*pb.OAuthTokenRevocationResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method OAuthTokenRevocation not implemented")
}

func (svc AuthService) OAuthTokenIntrospection(ctx context.Context, req *pb.OAuthTokenIntrospectionRequest) (*pb.OAuthTokenIntrospectionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method OAuthTokenIntrospection not implemented")
}

func (svc AuthService) OpenIDUserInfo(ctx context.Context, req *pb.OpenIDUserInfoRequest) (*pb.OpenIDUserInfoResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method OpenIDUserInfo not implemented")
}

func (svc AuthService) LoginPassword(ctx context.Context, req *pb.LoginPasswordRequest) (*pb.LoginPasswordResponse, error) {
	// Validate the request args
	if err := (*svc.pbVal).Validate(req); err != nil {
		// provide validation err context to user
		return nil, errors.NewServiceError(errors.ErrCodeInvalidArgument, "invalid request: "+err.Error())
	}

	// Verify password
	match, err := svc.store.VerifyPassword(ctx, req.UserId, req.Password)
	if err != nil || !match {
		return nil, errors.WrapServiceError(errors.ErrCodeForbidden, "invalid user id or password", err)
	}

	// Issue token for the user
	token, err := issueToken(svc.cfg, req.UserId)
	if err != nil {
		return nil, errors.WrapServiceError(errors.ErrCodeService, "error issuing token", err)
	}

	return &pb.LoginPasswordResponse{Detail: "Success", Data: token}, nil
}

func (svc AuthService) SetPassword(ctx context.Context, req *pb.SetPasswordRequest) (*pb.SetPasswordResponse, error) {
	// If the request is through the gateway,
	// then perform permission checking
	gatewayRequest, gwMd := gwauth.IsGatewayRequest(ctx)
	if gatewayRequest {
		log.Info().Msg("is a gateway request")
		// Ensure user is authenticated
		claims, err := gwauth.GetGatewayUser(gwMd)
		if err != nil {
			return nil, err
		}

		// Only allow changing of own password
		req.UserId = claims.Subject
	}

	// Validate the request args
	if err := (*svc.pbVal).Validate(req); err != nil {
		// provide validation err context to user
		return nil, errors.NewServiceError(errors.ErrCodeInvalidArgument, "invalid request: "+err.Error())
	}

	// Set the password
	err := svc.store.SetPassword(ctx, req.UserId, req.Password)
	if err != nil {
		return nil, err
	}

	return &pb.SetPasswordResponse{Detail: "Successfully updated password"}, nil
}

func (svc AuthService) ProcessUserDeletedEvent(ctx context.Context, req *eventpb.UserDeletedEvent) (*emptypb.Empty, error) {
	// Validate the request args
	if err := (*svc.pbVal).Validate(req); err != nil {
		// provide validation err context to user
		return nil, errors.NewServiceError(errors.ErrCodeInvalidArgument, "invalid request: "+err.Error())
	}

	err := svc.store.DeleteAuthMethods(ctx, req.UserId)
	if err != nil {
		return nil, errors.WrapServiceError(errors.ErrCodeService, "failed to process event", err)
	}

	return &emptypb.Empty{}, nil
}

// Provide the JWK ECDSA public key as part of a JSON Web Key set.
// This method is called by the API gateway for usage when validating inbound JWT tokens.
func (svc AuthService) GetJwks(ctx context.Context, req *pb.GetJwksRequest) (*pb.GetJwksResponse, error) {
	return &pb.GetJwksResponse{Keys: []*pb.PublicEcJWK{svc.cfg.ServiceOpts.PublicJwk}}, nil
}
