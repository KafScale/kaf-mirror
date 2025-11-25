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
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"log"
	"time"

	"github.com/jmoiron/sqlx"
)

// CreateApiToken creates a new API token for a user.
func CreateApiToken(db *sqlx.DB, userID int, description string, expiresAt time.Time) (string, *ApiToken, error) {
	tokenBytes := make([]byte, 32)
	if _, err := rand.Read(tokenBytes); err != nil {
		return "", nil, err
	}
	token := hex.EncodeToString(tokenBytes)

	hash := sha256.Sum256([]byte(token))
	tokenHash := hex.EncodeToString(hash[:])

	query := `INSERT INTO api_tokens (user_id, token_hash, description, expires_at) VALUES (?, ?, ?, ?)`
	result, err := db.Exec(query, userID, tokenHash, description, expiresAt)
	if err != nil {
		return "", nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return "", nil, err
	}

	var apiToken ApiToken
	err = db.Get(&apiToken, "SELECT * FROM api_tokens WHERE id = ?", id)
	if err != nil {
		return "", nil, err
	}

	return token, &apiToken, nil
}

// ValidateApiToken checks if a token is valid and returns the user ID.
func ValidateApiToken(db *sqlx.DB, token string) (int, error) {
	hash := sha256.Sum256([]byte(token))
	tokenHash := hex.EncodeToString(hash[:])

	var tokenInfo ApiToken
	err := db.Get(&tokenInfo, "SELECT * FROM api_tokens WHERE token_hash = ?", tokenHash)
	if err != nil {
		log.Printf("ERROR: DB: Token not found in database: %v", err)
		return 0, err
	}

	if !tokenInfo.ExpiresAt.IsZero() && tokenInfo.ExpiresAt.Before(time.Now()) {
		log.Printf("ERROR: DB: Token expired for user %d", tokenInfo.UserID)
		return 0, sql.ErrNoRows
	}

	return tokenInfo.UserID, nil
}

// RevokeAllUserTokens revokes all API tokens for a user.
func RevokeAllUserTokens(db *sqlx.DB, userID int) error {
	query := `DELETE FROM api_tokens WHERE user_id = ?`
	_, err := db.Exec(query, userID)
	return err
}
