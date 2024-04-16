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

package config

import (
	"fmt"
)

type PostgresConfig struct {
	// Env Var: "PG_USER"
	Username string

	// Env Var: "PG_PASS"
	Password string

	// Env Var: "PG_HOST"
	Host string

	// Env Var: "PG_PORT" (optional)
	// Defaults to "5432"
	Port string

	// Env Var: "PG_DB"
	Database string
}

func (conf *PostgresConfig) GetDSN() string {
	return fmt.Sprintf(
		"postgresql://%s:%s@%s:%s/%s?sslmode=disable",
		conf.Username,
		conf.Password,
		conf.Host,
		conf.Port,
		conf.Database,
	)
}

func (cfg *PostgresConfig) Load() error {
	// Load configurations from env
	if opt, err := RequireFromEnv("PG_USER"); err != nil {
		return err
	} else {
		cfg.Username = opt
	}

	if opt, err := RequireFromEnv("PG_PASS"); err != nil {
		return err
	} else {
		cfg.Password = opt
	}

	if opt, err := RequireFromEnv("PG_HOST"); err != nil {
		return err
	} else {
		cfg.Host = opt
	}

	if opt, err := RequireFromEnv("PG_PORT"); err != nil {
		cfg.Port = "5432"
	} else {
		cfg.Port = opt
	}

	if opt, err := RequireFromEnv("PG_DB"); err != nil {
		return err
	} else {
		cfg.Database = opt
	}

	// Config properties succesfully loaded
	return nil
}
