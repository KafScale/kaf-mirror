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
	"kaf-mirror/internal/database"
	"kaf-mirror/internal/kafka"
)

// MockKafMirror is a mock implementation of the KafMirror interface.
type MockKafMirror struct {
	StartFunc func(jobID string, metricsCallback func(database.ReplicationMetric), onPanic func(jobID string, reason string))
	StopFunc  func()
}

// Start calls the mock StartFunc.
func (m *MockKafMirror) Start(jobID string, metricsCallback func(database.ReplicationMetric), onPanic func(jobID string, reason string)) {
	if m.StartFunc != nil {
		m.StartFunc(jobID, metricsCallback, onPanic)
	}
}

// Stop calls the mock StopFunc.
func (m *MockKafMirror) Stop() {
	if m.StopFunc != nil {
		m.StopFunc()
	}
}

// GetConsumer is a mock implementation.
func (m *MockKafMirror) GetConsumer() *kafka.Consumer {
	return nil
}

// GetProducer is a mock implementation.
func (m *MockKafMirror) GetProducer() *kafka.Producer {
	return nil
}
