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


package database_test

import (
	"kaf-mirror/internal/database"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestDatabase(t *testing.T) {
	db, err := database.InitDB(":memory:")
	assert.NoError(t, err)
	defer db.Close()

	t.Run("Jobs", func(t *testing.T) {
		jobID := uuid.NewString()
		job := &database.ReplicationJob{
			ID:                jobID,
			Name:              "test-job",
			SourceClusterName: "src",
			TargetClusterName: "tgt",
			Status:            "paused",
		}
		err := database.CreateJob(db, job)
		assert.NoError(t, err)

		fetchedJob, err := database.GetJob(db, jobID)
		assert.NoError(t, err)
		assert.Equal(t, "test-job", fetchedJob.Name)

		jobs, err := database.ListJobs(db)
		assert.NoError(t, err)
		assert.Len(t, jobs, 1)

		err = database.DeleteJob(db, jobID)
		assert.NoError(t, err)

		_, err = database.GetJob(db, jobID)
		assert.Error(t, err)
	})

	t.Run("Mappings", func(t *testing.T) {
		jobID := uuid.NewString()
		job := &database.ReplicationJob{ID: jobID, Name: "mapping-job", SourceClusterName: "a", TargetClusterName: "b", Status: "paused"}
		database.CreateJob(db, job)

		mappings := []database.TopicMapping{
			{JobID: jobID, SourceTopicPattern: "a", TargetTopicPattern: "b", Enabled: true},
		}
		err := database.UpdateMappingsForJob(db, jobID, mappings)
		assert.NoError(t, err)

		fetchedMappings, err := database.GetMappingsForJob(db, jobID)
		assert.NoError(t, err)
		assert.Len(t, fetchedMappings, 1)
		assert.Equal(t, "a", fetchedMappings[0].SourceTopicPattern)
	})

	t.Run("Metrics", func(t *testing.T) {
		jobID := uuid.NewString()
		job := &database.ReplicationJob{ID: jobID, Name: "metrics-job", SourceClusterName: "a", TargetClusterName: "b", Status: "paused"}
		database.CreateJob(db, job)

		metric1 := &database.ReplicationMetric{
			JobID:              jobID,
			MessagesReplicated: 100,
			BytesTransferred:   1000,
			CurrentLag:         10,
			ErrorCount:         0,
			Timestamp:          time.Now().Add(-10 * time.Second),
		}
		err := database.InsertMetrics(db, metric1)
		assert.NoError(t, err)

		metric2 := &database.ReplicationMetric{
			JobID:              jobID,
			MessagesReplicated: 123,
			BytesTransferred:   4560,
			CurrentLag:         12,
			ErrorCount:         1,
			Timestamp:          time.Now(),
		}
		err = database.InsertMetrics(db, metric2)
		assert.NoError(t, err)

		metrics, err := database.GetHistoricalMetrics(db, jobID, time.Now().Add(-1*time.Hour), time.Now())
		assert.NoError(t, err)
		assert.Len(t, metrics, 2)
		assert.Equal(t, 23, metrics[1].MessagesReplicatedDelta)
		assert.Equal(t, 3560, metrics[1].BytesTransferredDelta)
		assert.Equal(t, 1, metrics[1].ErrorCountDelta)
	})

	t.Run("TestConfluentClusterUniqueness", func(t *testing.T) {
		// Create a confluent cluster
		cluster1 := &database.KafkaCluster{
			Name:      "confluent-1",
			Provider:  "confluent",
			ClusterID: "lkc-12345",
			Brokers:   "localhost:9092",
		}
		err := database.CreateCluster(db, cluster1)
		assert.NoError(t, err)

		// Attempt to create another confluent cluster with the same cluster_id (should fail)
		cluster2 := &database.KafkaCluster{
			Name:      "confluent-2",
			Provider:  "confluent",
			ClusterID: "lkc-12345",
			Brokers:   "localhost:9092",
		}
		err = database.CreateCluster(db, cluster2)
		assert.Error(t, err)

		// Attempt to create another confluent cluster with the same brokers but different cluster_id (should succeed)
		cluster3 := &database.KafkaCluster{
			Name:      "confluent-3",
			Provider:  "confluent",
			ClusterID: "lkc-67890",
			Brokers:   "localhost:9092",
		}
		err = database.CreateCluster(db, cluster3)
		assert.NoError(t, err)

		// Attempt to create a plain cluster with the same brokers (should succeed)
		cluster4 := &database.KafkaCluster{
			Name:     "plain-1",
			Provider: "plain",
			Brokers:  "localhost:9092",
		}
		err = database.CreateCluster(db, cluster4)
		assert.NoError(t, err)
	})
}
