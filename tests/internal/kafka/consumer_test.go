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
	"kaf-mirror/internal/config"
	"kaf-mirror/internal/kafka"
	"kaf-mirror/tests/mocks"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/twmb/franz-go/pkg/kgo"
)

func TestConsumer_Consume(t *testing.T) {
	record := &kgo.Record{
		Topic:     "test-topic",
		Partition: 0,
		Value:     []byte("hello"),
	}

	fetches := kgo.Fetches{{
		Topics: []kgo.FetchTopic{{
			Topic: "test-topic",
			Partitions: []kgo.FetchPartition{{
				Partition: 0,
				Records:   []*kgo.Record{record},
			}},
		}},
	}}

	mock := &mocks.MockKgoClient{
		PollFetchesFunc: func(ctx context.Context) kgo.Fetches {
			return fetches
		},
	}

	consumer := &kafka.Consumer{
		Client: mock,
	}
	
	recordProcessed := false
	handler := func(r *kgo.Record) {
		assert.Equal(t, record.Value, r.Value)
		assert.Equal(t, "test-topic", r.Topic)
		assert.Equal(t, int32(0), r.Partition)
		recordProcessed = true
	}

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	consumer.Consume(ctx, handler)
	
	assert.True(t, recordProcessed, "Handler should have been called")
}

func TestNewConsumer_SASL_PLAIN(t *testing.T) {
	cfg := config.ClusterConfig{
		Brokers: "localhost:9092",
		Security: config.SecurityConfig{
			Enabled:       true,
			Protocol:      "SASL_PLAINTEXT",
			SASLMechanism: "PLAIN",
			Username:      "testuser",
			Password:      "testpass",
		},
	}

	replicationCfg := config.ReplicationConfig{
		BatchSize:   1000,
		Parallelism: 4,
		Compression: "gzip",
	}

	consumer, err := kafka.NewConsumer(cfg, "test-group", replicationCfg, "test-job", "test-topic")
	assert.NoError(t, err)
	assert.NotNil(t, consumer)
}

func TestNewConsumer_SASL_SCRAM(t *testing.T) {
	cfg := config.ClusterConfig{
		Brokers: "localhost:9092",
		Security: config.SecurityConfig{
			Enabled:       true,
			Protocol:      "SASL_SSL",
			SASLMechanism: "SCRAM-SHA-256",
			Username:      "testuser",
			Password:      "testpass",
		},
	}

	replicationCfg := config.ReplicationConfig{
		BatchSize:   2000,
		Parallelism: 8,
		Compression: "snappy",
	}

	consumer, err := kafka.NewConsumer(cfg, "test-group", replicationCfg, "test-job", "test-topic")
	assert.NoError(t, err)
	assert.NotNil(t, consumer)
}

func TestNewConsumer_Kerberos(t *testing.T) {
	cfg := config.ClusterConfig{
		Brokers: "localhost:9092",
		Security: config.SecurityConfig{
			Enabled:       true,
			Protocol:      "SASL_SSL",
			SASLMechanism: "GSSAPI",
			Kerberos: struct {
				ServiceName string `mapstructure:"service_name"`
			}{
				ServiceName: "kafka",
			},
		},
	}

	replicationCfg := config.ReplicationConfig{
		BatchSize:   500,
		Parallelism: 2,
		Compression: "lz4",
	}

	consumer, err := kafka.NewConsumer(cfg, "test-group", replicationCfg, "test-job", "test-topic")
	assert.NoError(t, err)
	assert.NotNil(t, consumer)
}

func TestNewConsumer_InvalidSASL(t *testing.T) {
	cfg := config.ClusterConfig{
		Brokers: "localhost:9092",
		Security: config.SecurityConfig{
			Enabled:       true,
			SASLMechanism: "INVALID_MECHANISM",
		},
	}

	replicationCfg := config.ReplicationConfig{
		BatchSize:   1000,
		Parallelism: 4,
		Compression: "gzip",
	}

	consumer, err := kafka.NewConsumer(cfg, "test-group", replicationCfg, "test-job", "test-topic")
	assert.Error(t, err)
	assert.Nil(t, consumer)
	assert.Contains(t, err.Error(), "unsupported SASL mechanism")
}

func TestNewConsumer_MissingCredentials(t *testing.T) {
	cfg := config.ClusterConfig{
		Brokers: "localhost:9092",
		Security: config.SecurityConfig{
			Enabled:       true,
			SASLMechanism: "PLAIN",
		},
	}

	replicationCfg := config.ReplicationConfig{
		BatchSize:   1000,
		Parallelism: 4,
		Compression: "gzip",
	}

	consumer, err := kafka.NewConsumer(cfg, "test-group", replicationCfg, "test-job", "test-topic")
	assert.Error(t, err)
	assert.Nil(t, consumer)
	assert.Contains(t, err.Error(), "username and password are required")
}
