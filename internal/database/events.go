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
	"github.com/jmoiron/sqlx"
)

// CreateOperationalEvent inserts a new operational event into the database.
func CreateOperationalEvent(db *sqlx.DB, event *OperationalEvent) error {
	query := `INSERT INTO operational_events (event_type, initiator, details)
              VALUES (?, ?, ?)`
	_, err := db.Exec(query, event.EventType, event.Initiator, event.Details)
	return err
}

// GetOperationalEvent retrieves a single operational event by its ID.
func GetOperationalEvent(db *sqlx.DB, eventID int) (*OperationalEvent, error) {
	var event OperationalEvent
	err := db.Get(&event, "SELECT * FROM operational_events WHERE id = ?", eventID)
	return &event, err
}

// ListOperationalEvents retrieves all operational events from the database.
func ListOperationalEvents(db *sqlx.DB) ([]OperationalEvent, error) {
	var events []OperationalEvent
	err := db.Select(&events, "SELECT * FROM operational_events ORDER BY timestamp DESC LIMIT 100")
	return events, err
}
