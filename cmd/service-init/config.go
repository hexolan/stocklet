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
	"github.com/hexolan/stocklet/internal/pkg/config"
)

// Init Container Configuration
type InitConfig struct {
	// Name of the service (e.g. 'auth' or 'order')
	//
	// Env Var: "INIT_SVC_NAME"
	ServiceName string

	// Env Var: "INIT_MIGRATIONS" (optional. accepts 'false')
	// Defaults to true
	ApplyMigrations bool

	// 'ApplyDebezium' will default to false unless
	// the debezium host is provided.
	//
	// Env Var: "INIT_DEBEZIUM_HOST" (optional)
	// e.g. "http://debezium:8083"
	ApplyDebezium bool
	DebeziumHost  string
}

func (opts *InitConfig) Load() error {
	// ServiceName
	opt, err := config.RequireFromEnv("INIT_SVC_NAME")
	if err != nil {
		return err
	}
	opts.ServiceName = opt

	// ApplyMigrations
	opts.ApplyMigrations = true
	if opt, _ := config.RequireFromEnv("INIT_MIGRATIONS"); opt == "false" {
		opts.ApplyMigrations = false
	}

	// ApplyDebezium and DebeziumHost
	if opt, err := config.RequireFromEnv("INIT_DEBEZIUM_HOST"); err == nil {
		opts.ApplyDebezium = true
		opts.DebeziumHost = opt
	}

	return nil
}
