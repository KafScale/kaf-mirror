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
	"encoding/json"
	"kaf-mirror/internal/config"
	"time"

	"github.com/jmoiron/sqlx"
)

// SaveConfig saves the entire configuration to the database.
// This is a simplified approach; a more granular approach might be better.
func SaveConfig(db *sqlx.DB, cfg *config.Config) error {
	// We'll store the entire config as a single JSON blob for simplicity.
	// A more robust solution would store individual keys.
	configJSON, err := json.Marshal(cfg)
	if err != nil {
		return err
	}

	query := `INSERT OR REPLACE INTO configuration (key, value, updated_at) VALUES (?, ?, ?)`
	_, err = db.Exec(query, "full_config", string(configJSON), time.Now())
	return err
}

// LoadConfig retrieves the configuration from the database.
func LoadConfig(db *sqlx.DB) (*config.Config, error) {
	var configJSON string
	err := db.Get(&configJSON, "SELECT value FROM configuration WHERE key = 'full_config'")
	if err != nil {
		return nil, err
	}

	var cfg config.Config
	if err := json.Unmarshal([]byte(configJSON), &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
