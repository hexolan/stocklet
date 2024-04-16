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

type OtelConfig struct {
	// Env Var: "OTEL_COLLECTOR_GRPC"
	CollectorGrpc string
}

func (cfg *OtelConfig) Load() error {
	// Load configurations from env
	if collectorGrpc, err := RequireFromEnv("OTEL_COLLECTOR_GRPC"); err != nil {
		return err
	} else {
		cfg.CollectorGrpc = collectorGrpc
	}

	// Succesfully loaded config properties
	return nil
}
