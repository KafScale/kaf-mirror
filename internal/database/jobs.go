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


package database

import (
	"errors"
	"time"

	"github.com/jmoiron/sqlx"
)

// ListJobs retrieves all replication jobs from the database.
func ListJobs(db *sqlx.DB) ([]ReplicationJob, error) {
	var jobs []ReplicationJob
	err := db.Select(&jobs, "SELECT * FROM replication_jobs ORDER BY created_at DESC")
	return jobs, err
}

// GetJob retrieves a single replication job by its ID.
func GetJob(db *sqlx.DB, id string) (*ReplicationJob, error) {
	var job ReplicationJob
	err := db.Get(&job, "SELECT * FROM replication_jobs WHERE id = ?", id)
	return &job, err
}

// CreateJob inserts a new replication job into the database.
func CreateJob(db *sqlx.DB, job *ReplicationJob) error {
	// Check for duplicate name
	var count int
	err := db.Get(&count, "SELECT COUNT(*) FROM replication_jobs WHERE name = ?", job.Name)
	if err != nil {
		return err
	}
	if count > 0 {
		return errors.New("a job with this name already exists")
	}

	query := `INSERT INTO replication_jobs (id, name, source_cluster_name, target_cluster_name, status, batch_size, parallelism, compression, preserve_partitions, created_at, updated_at)
              VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	_, err = db.Exec(query, job.ID, job.Name, job.SourceClusterName, job.TargetClusterName, job.Status, job.BatchSize, job.Parallelism, job.Compression, job.PreservePartitions, time.Now(), time.Now())
	return err
}

// UpdateJob updates an existing replication job in the database.
func UpdateJob(db *sqlx.DB, job *ReplicationJob) error {
	// Check for duplicate name
	var count int
	err := db.Get(&count, "SELECT COUNT(*) FROM replication_jobs WHERE name = ? AND id != ?", job.Name, job.ID)
	if err != nil {
		return err
	}
	if count > 0 {
		return errors.New("a job with this name already exists")
	}

	query := `UPDATE replication_jobs 
              SET name = ?, source_cluster_name = ?, target_cluster_name = ?, status = ?, failed_reason = ?, updated_at = ?
              WHERE id = ?`
	_, err = db.Exec(query, job.Name, job.SourceClusterName, job.TargetClusterName, job.Status, job.FailedReason, time.Now(), job.ID)
	return err
}

// DeleteJob removes a replication job from the database.
func DeleteJob(db *sqlx.DB, id string) error {
	_, err := db.Exec("DELETE FROM replication_jobs WHERE id = ?", id)
	return err
}
