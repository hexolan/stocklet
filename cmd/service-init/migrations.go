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
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/rs/zerolog/log"

	"github.com/hexolan/stocklet/internal/pkg/config"
)

func applyPostgresMigrations(conf *config.PostgresConfig) {
	m, err := migrate.New("file:///migrations", conf.GetDSN())
	if err != nil {
		log.Panic().Err(err).Msg("migrate: failed to open client")
	}

	err = m.Up()
	if err != nil {
		if err.Error() == "no change" {
			log.Info().Err(err).Msg("migrate: migrations up to date")
		} else {
			log.Panic().Err(err).Msg("migrate: raised when performing db migration")
		}
	}

	log.Info().Msg("migrate: succesfully performed postgres migrations")
}
