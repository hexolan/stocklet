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
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"

	"github.com/lestrrat-go/jwx/v2/jwk"
	"github.com/rs/zerolog/log"
	"google.golang.org/protobuf/encoding/protojson"

	"github.com/hexolan/stocklet/internal/pkg/config"
	"github.com/hexolan/stocklet/internal/pkg/errors"
	pb "github.com/hexolan/stocklet/internal/pkg/protogen/auth/v1"
)

// Auth Service Configuration
type ServiceConfig struct {
	// Core configuration
	Shared      config.SharedConfig
	ServiceOpts ServiceConfigOpts

	// Dynamically loaded configuration
	Postgres config.PostgresConfig
	Kafka    config.KafkaConfig
}

// load the service configuration
func NewServiceConfig() (*ServiceConfig, error) {
	cfg := ServiceConfig{}

	// load the shared config options
	if err := cfg.Shared.Load(); err != nil {
		return nil, err
	}

	// load the service config opts
	if err := cfg.ServiceOpts.Load(); err != nil {
		return nil, err
	}

	return &cfg, nil
}

// Service specific config options
type ServiceConfigOpts struct {
	// Env Var: "AUTH_PRIVATE_KEY"
	// to be provided in base64 format
	PrivateKey *ecdsa.PrivateKey

	// Generated from PrivateKey
	PublicJwk *pb.PublicEcJWK
}

// Load the ServiceConfigOpts
//
// PrivateKey is loaded and decoded from the base64
// encoded PEM file exposed in the 'AUTH_PRIVATE_KEY'
// environment variable.
func (opts *ServiceConfigOpts) Load() error {
	// load the private key
	if err := opts.loadPrivateKey(); err != nil {
		return err
	}

	// prepare the JWK public key
	opts.PublicJwk = preparePublicJwk(opts.PrivateKey)

	return nil
}

// Load the ECDSA private key.
//
// Used for signing JWT tokens.
// The public key is also served in JWK format, from this service,
// for use when validating the tokens at the API ingress.
func (opts *ServiceConfigOpts) loadPrivateKey() error {
	// PEM private key file exposed as an environment variable encoded in base64
	opt, err := config.RequireFromEnv("AUTH_PRIVATE_KEY")
	if err != nil {
		return err
	}

	// Decode from base64
	pkBytes, err := base64.StdEncoding.DecodeString(opt)
	if err != nil {
		return errors.WrapServiceError(errors.ErrCodeService, "the provided 'AUTH_PRIVATE_KEY' is not valid base64", err)
	}

	// Decode the PEM key
	pkBlock, _ := pem.Decode(pkBytes)
	if pkBlock == nil {
		return errors.NewServiceError(errors.ErrCodeService, "the provided 'AUTH_PRIVATE_KEY' is not valid PEM format")
	}

	// Parse the block to a ecdsa.PrivateKey object
	privKey, err := x509.ParseECPrivateKey(pkBlock.Bytes)
	if err != nil {
		return errors.WrapServiceError(errors.ErrCodeService, "failed to parse the provided 'AUTH_PRIVATE_KEY' to an EC private key", err)
	}

	opts.PrivateKey = privKey
	return nil
}

// Converts the ECDSA key to a public JWK.
func preparePublicJwk(privateKey *ecdsa.PrivateKey) *pb.PublicEcJWK {
	// Assemble the public JWK
	jwk, err := jwk.FromRaw(privateKey.PublicKey)
	if err != nil {
		log.Panic().Err(err).Msg("something went wrong parsing public key from private key")
	}

	// denote use for signatures
	jwk.Set("use", "sig")

	// envoy includes support for ES256, ES384 and ES512
	alg := fmt.Sprintf("ES%v", privateKey.Curve.Params().BitSize)
	if alg != "ES256" && alg != "ES384" && alg != "ES512" {
		log.Panic().Err(err).Msg("unsupported bitsize for private key")
	}
	jwk.Set("alg", alg)

	// Convert the JWK to JSON
	jwkBytes, err := json.Marshal(jwk)
	if err != nil {
		log.Panic().Err(err).Msg("something went wrong preparing the public JWK (json marshal)")
	}

	// Unmarshal the JSON to Protobuf format
	publicJwkPB := pb.PublicEcJWK{}
	err = protojson.Unmarshal(jwkBytes, &publicJwkPB)
	if err != nil {
		log.Panic().Err(err).Msg("something went wrong preparing the public JWK (protonjson unmarshal)")
	}

	return &publicJwkPB
}
