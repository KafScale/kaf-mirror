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

	"github.com/stretchr/testify/assert"
	"github.com/twmb/franz-go/pkg/kgo"
)

func TestProducer_Produce(t *testing.T) {
	record := &kgo.Record{Value: []byte("hello")}
	var producedRecord *kgo.Record

	mock := &mocks.MockKgoClient{
		ProduceFunc: func(ctx context.Context, r *kgo.Record, f func(*kgo.Record, error)) {
			producedRecord = r
			f(r, nil)
		},
	}

	producer := &kafka.Producer{Client: mock}
	producer.Produce(context.Background(), record, func(r *kgo.Record, err error) {
		assert.NoError(t, err)
		assert.Equal(t, record.Value, r.Value)
	})

	assert.Equal(t, record.Value, producedRecord.Value)
}

func TestNewProducer_Compression_Gzip(t *testing.T) {
	cfg := config.ClusterConfig{
		Brokers: "localhost:9092",
	}

	replicationCfg := config.ReplicationConfig{
		BatchSize:   1000,
		Parallelism: 4,
		Compression: "gzip",
	}

	producer, err := kafka.NewProducer(cfg, replicationCfg, "test-job")
	assert.NoError(t, err)
	assert.NotNil(t, producer)
}

func TestNewProducer_Compression_Snappy(t *testing.T) {
	cfg := config.ClusterConfig{
		Brokers: "localhost:9092",
	}

	replicationCfg := config.ReplicationConfig{
		BatchSize:   2000,
		Parallelism: 8,
		Compression: "snappy",
	}

	producer, err := kafka.NewProducer(cfg, replicationCfg, "test-job")
	assert.NoError(t, err)
	assert.NotNil(t, producer)
}

func TestNewProducer_Compression_Lz4(t *testing.T) {
	cfg := config.ClusterConfig{
		Brokers: "localhost:9092",
	}

	replicationCfg := config.ReplicationConfig{
		BatchSize:   500,
		Parallelism: 2,
		Compression: "lz4",
	}

	producer, err := kafka.NewProducer(cfg, replicationCfg, "test-job")
	assert.NoError(t, err)
	assert.NotNil(t, producer)
}

func TestNewProducer_Compression_Zstd(t *testing.T) {
	cfg := config.ClusterConfig{
		Brokers: "localhost:9092",
	}

	replicationCfg := config.ReplicationConfig{
		BatchSize:   1500,
		Parallelism: 6,
		Compression: "zstd",
	}

	producer, err := kafka.NewProducer(cfg, replicationCfg, "test-job")
	assert.NoError(t, err)
	assert.NotNil(t, producer)
}

func TestNewProducer_NoCompression(t *testing.T) {
	cfg := config.ClusterConfig{
		Brokers: "localhost:9092",
	}

	replicationCfg := config.ReplicationConfig{
		BatchSize:   1000,
		Parallelism: 4,
		Compression: "none",
	}

	producer, err := kafka.NewProducer(cfg, replicationCfg, "test-job")
	assert.NoError(t, err)
	assert.NotNil(t, producer)
}

func TestNewProducer_SASL_PLAIN(t *testing.T) {
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

	producer, err := kafka.NewProducer(cfg, replicationCfg, "test-job")
	assert.NoError(t, err)
	assert.NotNil(t, producer)
}

func TestNewProducer_TLS_SSL(t *testing.T) {
	cfg := config.ClusterConfig{
		Brokers: "localhost:9092",
		Security: config.SecurityConfig{
			Enabled:  true,
			Protocol: "SSL",
		},
	}

	replicationCfg := config.ReplicationConfig{
		BatchSize:   1000,
		Parallelism: 4,
		Compression: "snappy",
	}

	producer, err := kafka.NewProducer(cfg, replicationCfg, "test-job")
	assert.NoError(t, err)
	assert.NotNil(t, producer)
}

func TestNewProducer_Kerberos(t *testing.T) {
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
		BatchSize:   1000,
		Parallelism: 4,
		Compression: "lz4",
	}

	producer, err := kafka.NewProducer(cfg, replicationCfg, "test-job")
	assert.NoError(t, err)
	assert.NotNil(t, producer)
}

func TestNewProducer_InvalidSASL(t *testing.T) {
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

	producer, err := kafka.NewProducer(cfg, replicationCfg, "test-job")
	assert.Error(t, err)
	assert.Nil(t, producer)
	assert.Contains(t, err.Error(), "unsupported SASL mechanism")
}

func TestNewProducer_MissingKerberosServiceName(t *testing.T) {
	cfg := config.ClusterConfig{
		Brokers: "localhost:9092",
		Security: config.SecurityConfig{
			Enabled:       true,
			SASLMechanism: "GSSAPI",
		},
	}

	replicationCfg := config.ReplicationConfig{
		BatchSize:   1000,
		Parallelism: 4,
		Compression: "gzip",
	}

	producer, err := kafka.NewProducer(cfg, replicationCfg, "test-job")
	assert.Error(t, err)
	assert.Nil(t, producer)
	assert.Contains(t, err.Error(), "Kerberos service name is required")
}
