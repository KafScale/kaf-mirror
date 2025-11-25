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
)

// SplunkSink sends metrics to a Splunk HTTP Event Collector (HEC).
type SplunkSink struct {
	client *http.Client
	cfg    config.SplunkConfig
}

// NewSplunkSink creates a new Splunk sink.
func NewSplunkSink(cfg config.SplunkConfig) (*SplunkSink, error) {
	return &SplunkSink{
		client: &http.Client{},
		cfg:    cfg,
	}, nil
}

// Send sends a metric to Splunk.
func (s *SplunkSink) Send(metric database.ReplicationMetric) error {
	payload := map[string]interface{}{
		"event":       metric,
		"source":      "kaf-mirror",
		"sourcetype":  "_json",
		"index":       s.cfg.Index,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", s.cfg.HECEndpoint, bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Splunk "+s.cfg.HECToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to send metric to Splunk: %s", resp.Status)
	}

	return nil
}
