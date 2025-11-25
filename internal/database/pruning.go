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
	"time"

	"github.com/jmoiron/sqlx"
)

func PruneOldData(db *sqlx.DB, retentionDays int) error {
	if retentionDays != 30 {
		retentionDays = 30
	}
	cutoff := time.Now().AddDate(0, 0, -retentionDays)

	_, err := db.Exec(`DELETE FROM aggregated_metrics WHERE timestamp < ?`, cutoff)
	if err != nil {
		return err
	}

	_, err = db.Exec(`DELETE FROM operational_events WHERE timestamp < ?`, cutoff)
	if err != nil {
		return err
	}

	mirrorStateCutoff := time.Now().AddDate(0, 0, -7) // 7 days ago
	
	_, err = db.Exec(`DELETE FROM mirror_progress WHERE last_updated < ?`, mirrorStateCutoff)
	if err != nil {
		return err
	}
	
	_, err = db.Exec(`DELETE FROM resume_points WHERE calculated_at < ?`, mirrorStateCutoff)
	if err != nil {
		return err
	}
	
	_, err = db.Exec(`DELETE FROM mirror_gaps WHERE detected_at < ?`, mirrorStateCutoff)
	if err != nil {
		return err
	}
	
	_, err = db.Exec(`DELETE FROM mirror_state_analysis WHERE analyzed_at < ?`, mirrorStateCutoff)
	if err != nil {
		return err
	}

	return nil
}
