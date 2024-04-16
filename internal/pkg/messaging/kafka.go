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

package messaging

import (
	"context"

	"github.com/twmb/franz-go/pkg/kadm"
	"github.com/twmb/franz-go/pkg/kgo"

	"github.com/hexolan/stocklet/internal/pkg/config"
	"github.com/hexolan/stocklet/internal/pkg/errors"
)

func NewKafkaConn(conf *config.KafkaConfig, opts ...kgo.Opt) (*kgo.Client, error) {
	opts = append(opts, kgo.SeedBrokers(conf.Brokers...))
	cl, err := kgo.NewClient(opts...)
	if err != nil {
		return nil, errors.WrapServiceError(errors.ErrCodeExtService, "failed to connect to Kafka", err)
	}

	return cl, nil
}

func EnsureKafkaTopics(cl *kgo.Client, topics ...string) error {
	ctx := context.Background()
	kadmCl := kadm.NewClient(cl)

	_, err := kadmCl.CreateTopics(ctx, -1, -1, nil, topics...)
	if err != nil {
		return errors.WrapServiceError(errors.ErrCodeExtService, "failed to create Kafka topics", err)
	}

	return nil
}
