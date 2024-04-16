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
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/rs/zerolog/log"

	"github.com/hexolan/stocklet/internal/pkg/config"
)

func applyPostgresOutbox(cfg *InitConfig, conf *config.PostgresConfig) {
	payloadB, err := json.Marshal(map[string]string{
		"connector.class":        "io.debezium.connector.postgresql.PostgresConnector",
		"plugin.name":            "pgoutput",
		"tasks.max":              "1",
		"table.include.list":     "public.event_outbox",
		"transforms":             "outbox",
		"transforms.outbox.type": "io.debezium.transforms.outbox.EventRouter",
		"transforms.outbox.route.topic.replacement": "${routedByValue}",
		"value.converter": "io.debezium.converters.BinaryDataConverter",

		"topic.prefix":      cfg.ServiceName,
		"database.hostname": conf.Host,
		"database.port":     conf.Port,
		"database.user":     conf.Username,
		"database.password": conf.Password,
		"database.dbname":   conf.Database,
	})
	if err != nil {
		log.Panic().Err(err).Msg("debezium connect: failed to marshal debezium cfg")
	}

	url := cfg.DebeziumHost + "/connectors/" + cfg.ServiceName + "-outbox/config"
	log.Info().Str("url", url).Msg("debezium url")
	req, err := http.NewRequest(
		"PUT",
		url,
		bytes.NewReader(payloadB),
	)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		log.Panic().Err(err).Msg("debezium connect: failed to perform debezium request")
	}

	log.Info().Str("status", res.Status).Msg("debezium connect: applied outbox config")
}
