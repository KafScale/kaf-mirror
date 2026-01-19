// Copyright 2025 Scalytics, Inc. and Scalytics Europe, LTD
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//     http://www.apache.org/licenses/LICENSE-2.0
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
	messagesConsumed   prometheus.Gauge
	bytesConsumed      prometheus.Gauge
	currentLag         prometheus.Gauge
	errorCount         prometheus.Gauge
	sourceStalled      prometheus.Gauge
	targetStalled      prometheus.Gauge
	criticalLag        prometheus.Gauge
	highErrorRate      prometheus.Gauge
	errorSpike         prometheus.Gauge
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
	messagesConsumed := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "kaf_mirror_messages_consumed",
		Help: "Number of messages consumed.",
	})
	bytesConsumed := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "kaf_mirror_bytes_consumed",
		Help: "Number of bytes consumed.",
	})
	currentLag := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "kaf_mirror_current_lag",
		Help: "Current consumer lag.",
	})
	errorCount := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "kaf_mirror_error_count",
		Help: "Number of errors.",
	})
	sourceStalled := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "kaf_mirror_incident_source_stalled",
		Help: "Source consumption stalled (1=true).",
	})
	targetStalled := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "kaf_mirror_incident_target_stalled",
		Help: "Target production stalled (1=true).",
	})
	criticalLag := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "kaf_mirror_incident_critical_lag",
		Help: "Critical lag detected (1=true).",
	})
	highErrorRate := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "kaf_mirror_incident_high_error_rate",
		Help: "High error rate detected (1=true).",
	})
	errorSpike := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "kaf_mirror_incident_error_spike",
		Help: "Error spike detected (1=true).",
	})

	registry := prometheus.NewRegistry()
	registry.MustRegister(
		messagesReplicated,
		bytesTransferred,
		messagesConsumed,
		bytesConsumed,
		currentLag,
		errorCount,
		sourceStalled,
		targetStalled,
		criticalLag,
		highErrorRate,
		errorSpike,
	)

	pusher := push.New(cfg.PushGateway, "kaf-mirror").Gatherer(registry)

	return &PrometheusSink{
		pusher:             pusher,
		messagesReplicated: messagesReplicated,
		bytesTransferred:   bytesTransferred,
		messagesConsumed:   messagesConsumed,
		bytesConsumed:      bytesConsumed,
		currentLag:         currentLag,
		errorCount:         errorCount,
		sourceStalled:      sourceStalled,
		targetStalled:      targetStalled,
		criticalLag:        criticalLag,
		highErrorRate:      highErrorRate,
		errorSpike:         errorSpike,
	}, nil
}

// Send sends a metric to Prometheus.
func (s *PrometheusSink) Send(metric database.ReplicationMetric) error {
	s.messagesReplicated.Set(float64(metric.MessagesReplicated))
	s.bytesTransferred.Set(float64(metric.BytesTransferred))
	s.messagesConsumed.Set(float64(metric.MessagesConsumed))
	s.bytesConsumed.Set(float64(metric.BytesConsumed))
	s.currentLag.Set(float64(metric.CurrentLag))
	s.errorCount.Set(float64(metric.ErrorCount))
	s.sourceStalled.Set(boolToFloat(metric.SourceStalled))
	s.targetStalled.Set(boolToFloat(metric.TargetStalled))
	s.criticalLag.Set(boolToFloat(metric.CriticalLag))
	s.highErrorRate.Set(boolToFloat(metric.HighErrorRate))
	s.errorSpike.Set(boolToFloat(metric.ErrorSpike))

	return s.pusher.Push()
}

func boolToFloat(val bool) float64 {
	if val {
		return 1
	}
	return 0
}
