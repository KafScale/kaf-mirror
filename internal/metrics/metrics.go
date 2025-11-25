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
)

// Sink is an interface for sending metrics to a monitoring platform.
type Sink interface {
	Send(metric database.ReplicationMetric) error
}

// NewSink creates a new metrics sink based on the provided configuration.
func NewSink(cfg config.MonitoringConfig) (Sink, error) {
	if !cfg.Enabled {
		return nil, nil
	}

	switch cfg.Platform {
	case "splunk":
		return NewSplunkSink(cfg.Splunk)
	case "loki":
		return NewLokiSink(cfg.Loki)
	case "prometheus":
		return NewPrometheusSink(cfg.Prometheus)
	default:
		return nil, nil
	}
}
