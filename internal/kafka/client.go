// kaf-mirror - A high-performance Kafka replication tool with AI-powered operational intelligence.
// Copyright (C) 2025 Scalytics
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program. If not, see <https://www.gnu.org/licenses/>.


package kafka

import (
	"context"

	"github.com/twmb/franz-go/pkg/kgo"
	"github.com/twmb/franz-go/pkg/kmsg"
)

// KgoClient is an interface that abstracts the methods we use from kgo.Client.
// This allows us to mock the client in tests.
type KgoClient interface {
	Request(context.Context, kmsg.Request) (kmsg.Response, error)
	PollFetches(context.Context) kgo.Fetches
	Produce(context.Context, *kgo.Record, func(*kgo.Record, error))
	Close()
}

// Ensure kgo.Client satisfies our interface.
var _ KgoClient = (*kgo.Client)(nil)
