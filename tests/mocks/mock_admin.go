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
	"kaf-mirror/internal/kafka"

	"github.com/stretchr/testify/mock"
)

// MockAdminClient provides a mock implementation of AdminClient for testing
type MockAdminClient struct {
	mock.Mock
}

func (m *MockAdminClient) GetClusterInfo(ctx context.Context) (*kafka.ClusterInfo, error) {
	args := m.Called(ctx)
	return args.Get(0).(*kafka.ClusterInfo), args.Error(1)
}

func (m *MockAdminClient) GetTopicDetails(ctx context.Context, topicNames []string) (map[string]kafka.TopicDetails, error) {
	args := m.Called(ctx, topicNames)
	return args.Get(0).(map[string]kafka.TopicDetails), args.Error(1)
}

func (m *MockAdminClient) GetConsumerGroupOffsets(ctx context.Context, groupID string, topics []string) (map[string][]kafka.OffsetInfo, error) {
	args := m.Called(ctx, groupID, topics)
	return args.Get(0).(map[string][]kafka.OffsetInfo), args.Error(1)
}

func (m *MockAdminClient) GetTopicHighWaterMarks(ctx context.Context, topics []string) (map[string][]kafka.OffsetInfo, error) {
	args := m.Called(ctx, topics)
	return args.Get(0).(map[string][]kafka.OffsetInfo), args.Error(1)
}

func (m *MockAdminClient) ValidateTopicCompatibility(ctx context.Context, sourceInfo, targetInfo kafka.TopicInfo) error {
	args := m.Called(ctx, sourceInfo, targetInfo)
	return args.Error(0)
}

func (m *MockAdminClient) EnsureTopicExists(ctx context.Context, topicName string, partitions int32, replicationFactor int16) error {
	args := m.Called(ctx, topicName, partitions, replicationFactor)
	return args.Error(0)
}

func (m *MockAdminClient) ListTopics(ctx context.Context, topics ...string) (map[string]kafka.TopicInfo, error) {
	args := m.Called(ctx, topics)
	return args.Get(0).(map[string]kafka.TopicInfo), args.Error(1)
}

func (m *MockAdminClient) GetTopicLag(ctx context.Context, groupID, topic string) (int64, error) {
	args := m.Called(ctx, groupID, topic)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockAdminClient) Close() {
	m.Called()
}
