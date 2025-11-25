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
	"log"

	"github.com/jmoiron/sqlx"
)

// AssignRoleToUser assigns a role to a user.
func AssignRoleToUser(db *sqlx.DB, userID, roleID int) error {
	query := `INSERT OR REPLACE INTO user_roles (user_id, role_id) VALUES (?, ?)`
	_, err := db.Exec(query, userID, roleID)
	return err
}

// GrantPermissionToRole grants a permission to a role.
func GrantPermissionToRole(db *sqlx.DB, roleID, permissionID int) error {
	query := `INSERT OR IGNORE INTO role_permissions (role_id, permission_id) VALUES (?, ?)`
	_, err := db.Exec(query, roleID, permissionID)
	return err
}

// UserHasPermission checks if a user has a specific permission.
func UserHasPermission(db *sqlx.DB, userID int, permissionName string) (bool, error) {
	query := `
        SELECT COUNT(*)
        FROM user_roles ur
        JOIN role_permissions rp ON ur.role_id = rp.role_id
        JOIN permissions p ON rp.permission_id = p.id
        WHERE ur.user_id = ? AND p.name = ?`

	var count int
	err := db.Get(&count, query, userID, permissionName)
	if err != nil {
		log.Printf("ERROR: RBAC: Error checking permission: %v", err)
		return false, err
	}

	hasPermission := count > 0
	return hasPermission, nil
}

// GetUserRole retrieves the role of a user.
func GetUserRole(db *sqlx.DB, userID int) (string, error) {
	var roleName string
	query := `
		SELECT r.name
		FROM roles r
		JOIN user_roles ur ON r.id = ur.role_id
		WHERE ur.user_id = ?`
	err := db.Get(&roleName, query, userID)
	return roleName, err
}

// SeedDefaultRolesAndPermissions creates the default roles and permissions.
func SeedDefaultRolesAndPermissions(db *sqlx.DB) error {
	roles := []string{"admin", "operator", "monitoring", "compliance"}
	permissions := []string{
		"jobs:view", "jobs:start", "jobs:stop", "jobs:pause",
		"jobs:create", "jobs:delete", "jobs:edit",
		"clusters:view", "clusters:create", "clusters:edit", "clusters:delete",
		"metrics:view", "ai:insights:view", "ai:analysis:trigger",
		"users:create", "users:delete", "users:list", "users:assign-roles",
		"roles:manage", "config:view", "config:edit",
		"compliance:generate", "compliance:view",
		"inventory:view", "inventory:create",
	}

	rolePermissions := map[string][]string{
		"admin":      permissions,
		"operator": {
			"jobs:view", "jobs:start", "jobs:stop", "jobs:pause", "jobs:edit",
			"clusters:view", "clusters:edit",
			"metrics:view", "ai:insights:view", "ai:analysis:trigger",
			"inventory:view", "inventory:create",
		},
		"monitoring": {"jobs:view", "clusters:view", "metrics:view", "ai:insights:view", "inventory:view"},
		"compliance": {"jobs:view", "clusters:view", "metrics:view", "compliance:generate", "compliance:view", "inventory:view"},
	}

	tx, err := db.Beginx()
	if err != nil {
		return err
	}

	for _, role := range roles {
		if _, err := tx.Exec("INSERT OR IGNORE INTO roles (name) VALUES (?)", role); err != nil {
			tx.Rollback()
			return err
		}
	}

	for _, permission := range permissions {
		if _, err := tx.Exec("INSERT OR IGNORE INTO permissions (name) VALUES (?)", permission); err != nil {
			tx.Rollback()
			return err
		}
	}

	for role, perms := range rolePermissions {
		var roleID int
		if err := tx.Get(&roleID, "SELECT id FROM roles WHERE name = ?", role); err != nil {
			tx.Rollback()
			return err
		}

		for _, perm := range perms {
			var permID int
			if err := tx.Get(&permID, "SELECT id FROM permissions WHERE name = ?", perm); err != nil {
				tx.Rollback()
				return err
			}
			if _, err := tx.Exec("INSERT OR IGNORE INTO role_permissions (role_id, permission_id) VALUES (?, ?)", roleID, permID); err != nil {
				tx.Rollback()
				return err
			}
		}
	}

	return tx.Commit()
}
