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
	"bytes"
	"encoding/json"
	"fmt"
	"kaf-mirror/internal/config"
	"kaf-mirror/internal/database"
	"net/http"
	"time"
)

// LokiSink sends metrics to a Grafana Loki instance.
type LokiSink struct {
	client *http.Client
	cfg    config.LokiConfig
}

// NewLokiSink creates a new Loki sink.
func NewLokiSink(cfg config.LokiConfig) (*LokiSink, error) {
	return &LokiSink{
		client: &http.Client{},
		cfg:    cfg,
	}, nil
}

// Send sends a metric to Loki.
func (s *LokiSink) Send(metric database.ReplicationMetric) error {
	logEntry := map[string]interface{}{
		"streams": []map[string]interface{}{
			{
				"stream": map[string]string{
					"job_id": metric.JobID,
				},
				"values": [][]string{
					{
						fmt.Sprintf("%d", time.Now().UnixNano()),
						fmt.Sprintf("messages_replicated=%d bytes_transferred=%d current_lag=%d error_count=%d",
							metric.MessagesReplicated, metric.BytesTransferred, metric.CurrentLag, metric.ErrorCount),
					},
				},
			},
		},
	}

	body, err := json.Marshal(logEntry)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", s.cfg.Endpoint, bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("failed to send metric to Loki: %s", resp.Status)
	}

	return nil
}
