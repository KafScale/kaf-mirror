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

import "github.com/jmoiron/sqlx"

// GetMappingsForJob retrieves all topic mappings for a given job ID.
func GetMappingsForJob(db *sqlx.DB, jobID string) ([]TopicMapping, error) {
	var mappings []TopicMapping
	err := db.Select(&mappings, "SELECT * FROM topic_mappings WHERE job_id = ?", jobID)
	return mappings, err
}

// UpdateMappingsForJob replaces all topic mappings for a given job ID.
func UpdateMappingsForJob(db *sqlx.DB, jobID string, mappings []TopicMapping) error {
	tx, err := db.Beginx()
	if err != nil {
		return err
	}

	// It's often easiest to delete all existing mappings and re-insert them.
	// For very large sets, a more sophisticated diff-and-update would be better.
	_, err = tx.Exec("DELETE FROM topic_mappings WHERE job_id = ?", jobID)
	if err != nil {
		tx.Rollback()
		return err
	}

	for _, m := range mappings {
		query := `INSERT INTO topic_mappings (job_id, source_topic_pattern, target_topic_pattern, enabled)
				  VALUES (?, ?, ?, ?)`
		_, err = tx.Exec(query, jobID, m.SourceTopicPattern, m.TargetTopicPattern, m.Enabled)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}
