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

package database

import (
	"time"

	"github.com/jmoiron/sqlx"
)

// InsertMetrics stores a new metrics data point in the database.
func InsertMetrics(db *sqlx.DB, metric *ReplicationMetric) error {
	// Get the last metric for this job
	lastMetric, err := GetLatestMetrics(db, metric.JobID)
	if err != nil {
		return err
	}

	// Calculate deltas
	messagesDelta := metric.MessagesReplicated - lastMetric.MessagesReplicated
	bytesDelta := metric.BytesTransferred - lastMetric.BytesTransferred
	consumedMessagesDelta := metric.MessagesConsumed - lastMetric.MessagesConsumed
	consumedBytesDelta := metric.BytesConsumed - lastMetric.BytesConsumed
	errorsDelta := metric.ErrorCount - lastMetric.ErrorCount

	if messagesDelta < 0 {
		messagesDelta = metric.MessagesReplicated
	}
	if bytesDelta < 0 {
		bytesDelta = metric.BytesTransferred
	}
	if consumedMessagesDelta < 0 {
		consumedMessagesDelta = metric.MessagesConsumed
	}
	if consumedBytesDelta < 0 {
		consumedBytesDelta = metric.BytesConsumed
	}
	if errorsDelta < 0 {
		errorsDelta = metric.ErrorCount
	}

	// Insert into the aggregated table
	query := `INSERT INTO aggregated_metrics (
			  job_id,
			  timestamp,
			  messages_replicated_delta,
			  bytes_transferred_delta,
			  messages_consumed_delta,
			  bytes_consumed_delta,
			  avg_lag,
			  error_count_delta
			) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`
	_, err = db.Exec(query, metric.JobID, time.Now(), messagesDelta, bytesDelta, consumedMessagesDelta, consumedBytesDelta, metric.CurrentLag, errorsDelta)
	return err
}

// GetLatestMetrics retrieves the latest metrics for a given job.
func GetLatestMetrics(db *sqlx.DB, jobID string) (*ReplicationMetric, error) {
	var totals struct {
		MessagesReplicated int `db:"messages_replicated"`
		BytesTransferred   int `db:"bytes_transferred"`
		MessagesConsumed   int `db:"messages_consumed"`
		BytesConsumed      int `db:"bytes_consumed"`
		ErrorCount         int `db:"error_count"`
	}
	totalsQuery := `
        SELECT
            COALESCE(SUM(messages_replicated_delta), 0) as messages_replicated,
            COALESCE(SUM(bytes_transferred_delta), 0) as bytes_transferred,
            COALESCE(SUM(messages_consumed_delta), 0) as messages_consumed,
            COALESCE(SUM(bytes_consumed_delta), 0) as bytes_consumed,
            COALESCE(SUM(error_count_delta), 0) as error_count
        FROM aggregated_metrics
        WHERE job_id = ?
    `
	err := db.Get(&totals, totalsQuery, jobID)
	if err != nil {
		return nil, err
	}

	var lastMetric struct {
		CurrentLag int       `db:"avg_lag"`
		Timestamp  time.Time `db:"timestamp"`
	}
	lastMetricQuery := "SELECT avg_lag, timestamp FROM aggregated_metrics WHERE job_id = ? ORDER BY timestamp DESC LIMIT 1"
	err = db.Get(&lastMetric, lastMetricQuery, jobID)
	// Ignore error if no rows, lag and timestamp will be zero.

	return &ReplicationMetric{
		JobID:              jobID,
		MessagesReplicated: totals.MessagesReplicated,
		BytesTransferred:   totals.BytesTransferred,
		MessagesConsumed:   totals.MessagesConsumed,
		BytesConsumed:      totals.BytesConsumed,
		ErrorCount:         totals.ErrorCount,
		CurrentLag:         lastMetric.CurrentLag,
		Timestamp:          lastMetric.Timestamp,
	}, nil
}

// GetHistoricalMetrics retrieves historical metrics for a given job within a time range.
func GetHistoricalMetrics(db *sqlx.DB, jobID string, start, end time.Time) ([]AggregatedMetric, error) {
	var metrics []AggregatedMetric
	err := db.Select(&metrics, "SELECT * FROM aggregated_metrics WHERE job_id = ? AND timestamp BETWEEN ? AND ? ORDER BY timestamp", jobID, start, end)
	return metrics, err
}

// GetAggregatedHistoricalMetrics retrieves aggregated historical metrics for a given job.
func GetAggregatedHistoricalMetrics(db *sqlx.DB, jobID string, periodDays int, granularity string) ([]AggregatedMetric, error) {
	var metrics []AggregatedMetric
	var groupBy string

	switch granularity {
	case "hourly":
		groupBy = "strftime('%Y-%m-%d %H:00:00', timestamp)"
	case "daily":
		groupBy = "strftime('%Y-%m-%d', timestamp)"
	default:
		groupBy = "strftime('%Y-%m-%d', timestamp)" // Default to daily
	}

	query := `
        SELECT
            ` + groupBy + ` as period,
            SUM(messages_replicated_delta) as avg_throughput,
            AVG(avg_lag) as avg_lag,
            SUM(error_count_delta) as total_errors
        FROM aggregated_metrics
        WHERE job_id = ? AND timestamp >= datetime('now', '-' || ? || ' days')
        GROUP BY period
        ORDER BY period ASC
    `

	err := db.Select(&metrics, query, jobID, periodDays)
	return metrics, err
}
