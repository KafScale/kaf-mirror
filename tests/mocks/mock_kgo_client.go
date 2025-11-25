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

package mocks

import (
	"context"

	"github.com/twmb/franz-go/pkg/kgo"
	"github.com/twmb/franz-go/pkg/kmsg"
)

// MockKgoClient is a mock implementation of the KgoClient interface.
type MockKgoClient struct {
	RequestFunc     func(context.Context, kmsg.Request) (kmsg.Response, error)
	PollFetchesFunc func(context.Context) kgo.Fetches
	ProduceFunc     func(context.Context, *kgo.Record, func(*kgo.Record, error))
	CloseFunc       func()
}

func (m *MockKgoClient) Request(ctx context.Context, req kmsg.Request) (kmsg.Response, error) {
	if m.RequestFunc != nil {
		return m.RequestFunc(ctx, req)
	}
	return nil, nil
}

func (m *MockKgoClient) PollFetches(ctx context.Context) kgo.Fetches {
	if m.PollFetchesFunc != nil {
		return m.PollFetchesFunc(ctx)
	}
	return kgo.Fetches{}
}

func (m *MockKgoClient) Produce(ctx context.Context, r *kgo.Record, f func(*kgo.Record, error)) {
	if m.ProduceFunc != nil {
		m.ProduceFunc(ctx, r, f)
	}
}

func (m *MockKgoClient) Close() {
	if m.CloseFunc != nil {
		m.CloseFunc()
	}
}
