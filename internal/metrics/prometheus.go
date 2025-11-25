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


package metrics

import (
	"kaf-mirror/internal/config"
	"kaf-mirror/internal/database"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/push"
)

// PrometheusSink sends metrics to a Prometheus Pushgateway.
type PrometheusSink struct {
	pusher *push.Pusher

	messagesReplicated prometheus.Gauge
	bytesTransferred   prometheus.Gauge
	currentLag         prometheus.Gauge
	errorCount         prometheus.Gauge
}

// NewPrometheusSink creates a new Prometheus sink.
func NewPrometheusSink(cfg config.PrometheusConfig) (*PrometheusSink, error) {
	messagesReplicated := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "kaf_mirror_messages_replicated",
		Help: "Number of messages replicated.",
	})
	bytesTransferred := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "kaf_mirror_bytes_transferred",
		Help: "Number of bytes transferred.",
	})
	currentLag := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "kaf_mirror_current_lag",
		Help: "Current consumer lag.",
	})
	errorCount := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "kaf_mirror_error_count",
		Help: "Number of errors.",
	})

	registry := prometheus.NewRegistry()
	registry.MustRegister(messagesReplicated, bytesTransferred, currentLag, errorCount)

	pusher := push.New(cfg.PushGateway, "kaf-mirror").Gatherer(registry)

	return &PrometheusSink{
		pusher:             pusher,
		messagesReplicated: messagesReplicated,
		bytesTransferred:   bytesTransferred,
		currentLag:         currentLag,
		errorCount:         errorCount,
	}, nil
}

// Send sends a metric to Prometheus.
func (s *PrometheusSink) Send(metric database.ReplicationMetric) error {
	s.messagesReplicated.Set(float64(metric.MessagesReplicated))
	s.bytesTransferred.Set(float64(metric.BytesTransferred))
	s.currentLag.Set(float64(metric.CurrentLag))
	s.errorCount.Set(float64(metric.ErrorCount))

	return s.pusher.Push()
}
