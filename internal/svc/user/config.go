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

package user

import (
	"github.com/hexolan/stocklet/internal/pkg/config"
)

// Order Service Configuration
type ServiceConfig struct {
	// Core Configuration
	Shared      config.SharedConfig
	ServiceOpts ServiceConfigOpts

	// Dynamically loaded configuration
	Postgres config.PostgresConfig
}

// load the base service configuration
func NewServiceConfig() (*ServiceConfig, error) {
	cfg := ServiceConfig{}

	// Load the core configuration
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
	// Env Var: "AUTH_SERVICE_GRPC"
	AuthServiceGrpc string
}

// Load the ServiceConfigOpts
func (opts *ServiceConfigOpts) Load() error {
	// Load configurations from env
	if authServiceGrpc, err := config.RequireFromEnv("AUTH_SERVICE_GRPC"); err == nil {
		opts.AuthServiceGrpc = authServiceGrpc
	} else {
		return err
	}

	return nil
}
