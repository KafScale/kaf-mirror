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

package kafka_test

import (
	"context"
	"errors"
	"kaf-mirror/internal/kafka"
	"kaf-mirror/tests/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAdminClient_GetTopicDetails(t *testing.T) {
	mockAdmin := &mocks.MockAdminClient{}
	topicDetails := map[string]kafka.TopicDetails{
		"test-topic": {
			Name:              "test-topic",
			Partitions:        3,
			ReplicationFactor: 2,
			CompressionType:   "gzip",
		},
	}

	mockAdmin.On("GetTopicDetails", mock.Anything, []string{"test-topic"}).Return(topicDetails, nil)

	details, err := mockAdmin.GetTopicDetails(context.Background(), []string{"test-topic"})
	assert.NoError(t, err)
	assert.Equal(t, topicDetails, details)
	mockAdmin.AssertExpectations(t)
}

func TestAdminClient_GetTopicDetails_Error(t *testing.T) {
	mockAdmin := &mocks.MockAdminClient{}
	mockAdmin.On("GetTopicDetails", mock.Anything, []string{"test-topic"}).Return(map[string]kafka.TopicDetails{}, errors.New("kafka error"))

	_, err := mockAdmin.GetTopicDetails(context.Background(), []string{"test-topic"})
	assert.Error(t, err)
	assert.Equal(t, "kafka error", err.Error())
	mockAdmin.AssertExpectations(t)
}

func TestAdminClient_GetClusterInfo(t *testing.T) {
	mockAdmin := &mocks.MockAdminClient{}
	clusterInfo := &kafka.ClusterInfo{
		ClusterID:   "test-cluster",
		Topics:      make(map[string]kafka.TopicInfo),
		BrokerCount: 3,
	}

	mockAdmin.On("GetClusterInfo", mock.Anything).Return(clusterInfo, nil)

	info, err := mockAdmin.GetClusterInfo(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, clusterInfo, info)
	mockAdmin.AssertExpectations(t)
}

func TestAdminClient_GetClusterInfo_Error(t *testing.T) {
	mockAdmin := &mocks.MockAdminClient{}
	mockAdmin.On("GetClusterInfo", mock.Anything).Return((*kafka.ClusterInfo)(nil), errors.New("kafka error"))

	_, err := mockAdmin.GetClusterInfo(context.Background())
	assert.Error(t, err)
	assert.Equal(t, "kafka error", err.Error())
	mockAdmin.AssertExpectations(t)
}

func TestAdminClient_GetConsumerGroupOffsets(t *testing.T) {
	mockAdmin := &mocks.MockAdminClient{}
	offsetInfo := map[string][]kafka.OffsetInfo{
		"test-topic": {
			{
				Topic:         "test-topic",
				Partition:     0,
				Offset:        100,
				HighWaterMark: 150,
				Lag:           50,
			},
		},
	}

	mockAdmin.On("GetConsumerGroupOffsets", mock.Anything, "test-group", []string{"test-topic"}).Return(offsetInfo, nil)

	offsets, err := mockAdmin.GetConsumerGroupOffsets(context.Background(), "test-group", []string{"test-topic"})
	assert.NoError(t, err)
	assert.Equal(t, offsetInfo, offsets)
	mockAdmin.AssertExpectations(t)
}

func TestAdminClient_GetConsumerGroupOffsets_Error(t *testing.T) {
	mockAdmin := &mocks.MockAdminClient{}
	mockAdmin.On("GetConsumerGroupOffsets", mock.Anything, "test-group", []string{"test-topic"}).Return((map[string][]kafka.OffsetInfo)(nil), errors.New("kafka error"))

	_, err := mockAdmin.GetConsumerGroupOffsets(context.Background(), "test-group", []string{"test-topic"})
	assert.Error(t, err)
	assert.Equal(t, "kafka error", err.Error())
	mockAdmin.AssertExpectations(t)
}
