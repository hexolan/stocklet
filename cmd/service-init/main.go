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
	"github.com/rs/zerolog/log"

	"github.com/hexolan/stocklet/internal/pkg/config"
)

func main() {
	// Load the init container cfg options
	cfg := InitConfig{}
	if err := cfg.Load(); err != nil {
		log.Panic().Err(err).Msg("missing required configuration")
	}

	// If migrations or debezium are enabled,
	// then a database configuration will be required.
	if cfg.ApplyMigrations || cfg.ApplyDebezium {
		// Support for dynamic loading of configuration
		// (e.g. mongo config instead of postgres config)
		pgConf := config.PostgresConfig{}
		if err := pgConf.Load(); err == nil {
			// Using postgres as a database.
			if cfg.ApplyMigrations {
				applyPostgresMigrations(&pgConf)
			}

			if cfg.ApplyDebezium {
				applyPostgresOutbox(&cfg, &pgConf)
			}
		} else {
			log.Panic().Msg("unable to load any db configs (unable to perform migrations or apply connector cfgs)")
		}
	}

	log.Info().Str("svc", cfg.ServiceName).Msg("completed init for service")
}
