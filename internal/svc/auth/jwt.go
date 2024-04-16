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
	"time"

	"github.com/lestrrat-go/jwx/v2/jwa"
	"github.com/lestrrat-go/jwx/v2/jwt"

	"github.com/hexolan/stocklet/internal/pkg/errors"
	"github.com/hexolan/stocklet/internal/pkg/gwauth"
	pb "github.com/hexolan/stocklet/internal/pkg/protogen/auth/v1"
)

// Issues a JWT token
func issueToken(cfg *ServiceConfig, sub string) (*pb.AuthToken, error) {
	const expiryTime = 86400 // 1 day

	// Set token claims
	token := jwt.New()
	claims := &gwauth.JWTClaims{
		Subject:  sub,
		IssuedAt: time.Now().Unix(),
		Expiry:   time.Now().Unix() + expiryTime,
	}

	token.Set("sub", claims.Subject)
	token.Set("iat", claims.IssuedAt)
	token.Set("exp", claims.Expiry)

	// Sign token
	accessToken, err := jwt.Sign(token, jwt.WithKey(jwa.KeyAlgorithmFrom(cfg.ServiceOpts.PrivateKey), cfg.ServiceOpts.PrivateKey))
	if err != nil {
		return nil, errors.WrapServiceError(errors.ErrCodeService, "failed to sign JWT token", err)
	}

	return &pb.AuthToken{
		TokenType:   "Bearer",
		AccessToken: string(accessToken),
		ExpiresIn:   expiryTime,
	}, nil
}
