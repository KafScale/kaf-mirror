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

// ListClusters retrieves all Kafka clusters from the database.
func ListClusters(db *sqlx.DB) ([]KafkaCluster, error) {
	var clusters []KafkaCluster
	err := db.Select(&clusters, "SELECT * FROM kafka_clusters ORDER BY name")
	return clusters, err
}

// GetCluster retrieves a single Kafka cluster by its name.
func GetCluster(db *sqlx.DB, name string) (*KafkaCluster, error) {
	var cluster KafkaCluster
	err := db.Get(&cluster, "SELECT * FROM kafka_clusters WHERE name = ?", name)
	return &cluster, err
}

// CreateCluster inserts a new Kafka cluster into the database.
func CreateCluster(db *sqlx.DB, cluster *KafkaCluster) error {
	// Check for duplicate name
	var count int
	err := db.Get(&count, "SELECT COUNT(*) FROM kafka_clusters WHERE name = ?", cluster.Name)
	if err != nil {
		return err
	}
	if count > 0 {
		return errors.New("a cluster with this name already exists")
	}

	// Provider-aware uniqueness check
	if cluster.Provider == "confluent" && cluster.ClusterID != "" {
		err = db.Get(&count, "SELECT COUNT(*) FROM kafka_clusters WHERE cluster_id = ?", cluster.ClusterID)
		if err != nil {
			return err
		}
		if count > 0 {
			return errors.New("a confluent cluster with this cluster_id already exists")
		}
	}

	query := `INSERT INTO kafka_clusters (name, provider, cluster_id, brokers, security_config, api_key, api_secret, connection_string)
              VALUES (?, ?, ?, ?, ?, ?, ?, ?)`
	_, err = db.Exec(query, cluster.Name, cluster.Provider, cluster.ClusterID, cluster.Brokers, cluster.SecurityConfig, cluster.APIKey, cluster.APISecret, cluster.ConnectionString)
	return err
}

// UpdateCluster updates an existing Kafka cluster in the database.
func UpdateCluster(db *sqlx.DB, cluster *KafkaCluster) error {
	// Check for duplicate name
	var count int
	err := db.Get(&count, "SELECT COUNT(*) FROM kafka_clusters WHERE name = ? AND name != ?", cluster.Name, cluster.Name)
	if err != nil {
		return err
	}
	if count > 0 {
		return errors.New("a cluster with this name already exists")
	}

	// Provider-aware uniqueness check
	if cluster.Provider == "confluent" && cluster.ClusterID != "" {
		err = db.Get(&count, "SELECT COUNT(*) FROM kafka_clusters WHERE cluster_id = ? AND name != ?", cluster.ClusterID, cluster.Name)
		if err != nil {
			return err
		}
		if count > 0 {
			return errors.New("a confluent cluster with this cluster_id already exists")
		}
	}

	query := `UPDATE kafka_clusters 
              SET provider = ?, cluster_id = ?, brokers = ?, security_config = ?, api_key = ?, api_secret = ?, connection_string = ?
              WHERE name = ?`
	_, err = db.Exec(query, cluster.Provider, cluster.ClusterID, cluster.Brokers, cluster.SecurityConfig, cluster.APIKey, cluster.APISecret, cluster.ConnectionString, cluster.Name)
	return err
}

// DeleteCluster removes a Kafka cluster from the database.
func DeleteCluster(db *sqlx.DB, name string) error {
	_, err := db.Exec("DELETE FROM kafka_clusters WHERE name = ?", name)
	return err
}

// SetClusterStatus updates the status of a Kafka cluster.
func SetClusterStatus(db *sqlx.DB, name, status string) error {
	query := `UPDATE kafka_clusters SET status = ?, updated_at = CURRENT_TIMESTAMP WHERE name = ?`
	_, err := db.Exec(query, status, name)
	return err
}

// PurgeArchivedClusters permanently deletes all archived Kafka clusters.
func PurgeArchivedClusters(db *sqlx.DB) error {
	_, err := db.Exec("DELETE FROM kafka_clusters WHERE status = 'archived'")
	return err
}

// ArchiveInactiveClusters moves clusters from inactive to archived after a certain duration.
func ArchiveInactiveClusters(db *sqlx.DB, duration time.Duration) error {
	cutoff := time.Now().Add(-duration)
	query := `UPDATE kafka_clusters SET status = 'archived' WHERE status = 'inactive' AND updated_at < ?`
	_, err := db.Exec(query, cutoff)
	return err
}
