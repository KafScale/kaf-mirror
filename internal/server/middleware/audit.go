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


package middleware

import (
	"encoding/json"
	"fmt"
	"kaf-mirror/internal/database"
	"regexp"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
)

// AuditLog is a middleware that logs state-modifying actions.
func AuditLog(db *sqlx.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if c.Method() == "GET" || c.Method() == "HEAD" || c.Method() == "OPTIONS" {
			return c.Next()
		}

		err := c.Next()

		if err != nil {
			return err
		}

		if c.Response().StatusCode() < 200 || c.Response().StatusCode() >= 300 {
			return nil
		}

		// Safe user extraction
		var initiator string
		if user := c.Locals("user"); user != nil {
			if u, ok := user.(*database.User); ok && u != nil {
				initiator = u.Username
			} else {
				initiator = "System"
			}
		} else {
			initiator = "Anonymous"
		}

		// Create readable details
		details := formatAuditDetails(c)

		event := &database.OperationalEvent{
			EventType: strings.ToUpper(c.Method()) + " " + c.Path(),
			Initiator: initiator,
			Details:   details,
		}

		if err := database.CreateOperationalEvent(db, event); err != nil {
			// Log to stderr in production
		}

		return nil
	}
}

// formatAuditDetails creates human-readable audit details while masking sensitive data
func formatAuditDetails(c *fiber.Ctx) string {
	path := c.Path()
	method := c.Method()
	
	// Create summary based on endpoint
	switch {
	case strings.Contains(path, "/jobs") && method == "POST":
		return formatJobCreation(c)
	case strings.Contains(path, "/jobs") && strings.Contains(path, "/start"):
		return fmt.Sprintf("Started replication job: %s", c.Params("id"))
	case strings.Contains(path, "/jobs") && strings.Contains(path, "/stop"):
		return fmt.Sprintf("Stopped replication job: %s", c.Params("id"))
	case strings.Contains(path, "/clusters") && method == "POST":
		return formatClusterCreation(c)
	case strings.Contains(path, "/config") && method == "PUT":
		return formatConfigUpdate(c)
	case strings.Contains(path, "/users") && method == "POST":
		return formatUserCreation(c)
	default:
		// Generic format with sensitive data masking
		body := string(c.Request().Body())
		if body == "" {
			return "No additional details"
		}
		return maskSensitiveData(body)
	}
}

// formatJobCreation creates readable job creation details
func formatJobCreation(c *fiber.Ctx) string {
	body := string(c.Request().Body())
	if body == "" {
		return "Created new replication job"
	}

	var jobData map[string]interface{}
	if err := json.Unmarshal([]byte(body), &jobData); err != nil {
		return "Created new replication job"
	}

	name := getStringValue(jobData, "name")
	source := getStringValue(jobData, "source_cluster_name")
	target := getStringValue(jobData, "target_cluster_name")

	return fmt.Sprintf("Created job '%s' from %s to %s", name, source, target)
}

// formatClusterCreation creates readable cluster creation details
func formatClusterCreation(c *fiber.Ctx) string {
	body := string(c.Request().Body())
	if body == "" {
		return "Created new cluster"
	}

	var clusterData map[string]interface{}
	if err := json.Unmarshal([]byte(body), &clusterData); err != nil {
		return "Created new cluster"
	}

	name := getStringValue(clusterData, "name")
	provider := getStringValue(clusterData, "provider")

	return fmt.Sprintf("Created cluster '%s' (provider: %s)", name, provider)
}

// formatConfigUpdate creates readable config update details
func formatConfigUpdate(c *fiber.Ctx) string {
	body := string(c.Request().Body())
	if body == "" {
		return "Updated system configuration"
	}

	var configData map[string]interface{}
	if err := json.Unmarshal([]byte(body), &configData); err != nil {
		return "Updated system configuration"
	}

	details := []string{}
	
	if ai, ok := configData["AI"].(map[string]interface{}); ok {
		if provider := getStringValue(ai, "Provider"); provider != "" {
			details = append(details, fmt.Sprintf("AI Provider: %s", provider))
		}
		if model := getStringValue(ai, "Model"); model != "" {
			details = append(details, fmt.Sprintf("AI Model: %s", model))
		}
	}

	if server, ok := configData["Server"].(map[string]interface{}); ok {
		if port := getIntValue(server, "Port"); port > 0 {
			details = append(details, fmt.Sprintf("Server Port: %d", port))
		}
	}

	if len(details) > 0 {
		return fmt.Sprintf("Updated configuration: %s", strings.Join(details, ", "))
	}

	return "Updated system configuration"
}

// formatUserCreation creates readable user creation details
func formatUserCreation(c *fiber.Ctx) string {
	body := string(c.Request().Body())
	if body == "" {
		return "Created new user"
	}

	var userData map[string]interface{}
	if err := json.Unmarshal([]byte(body), &userData); err != nil {
		return "Created new user"
	}

	username := getStringValue(userData, "username")
	role := getStringValue(userData, "role")

	return fmt.Sprintf("Created user '%s' with role '%s'", username, role)
}

// maskSensitiveData removes or masks sensitive information
func maskSensitiveData(data string) string {
	// Mask API keys and tokens
	apiKeyRegex := regexp.MustCompile(`"(api_key|token|password|secret)":\s*"[^"]*"`)
	data = apiKeyRegex.ReplaceAllString(data, `"$1":"***MASKED***"`)

	// Mask long tokens that look like API keys
	tokenRegex := regexp.MustCompile(`"([^"]*)":\s*"(sk-|Bearer |[A-Za-z0-9+/]{20,})"`)
	data = tokenRegex.ReplaceAllString(data, `"$1":"***MASKED***"`)

	// Truncate very long data
	if len(data) > 500 {
		data = data[:497] + "..."
	}

	return data
}

// Helper functions
func getStringValue(m map[string]interface{}, key string) string {
	if val, ok := m[key]; ok {
		if str, ok := val.(string); ok {
			return str
		}
	}
	return ""
}

func getIntValue(m map[string]interface{}, key string) int {
	if val, ok := m[key]; ok {
		if num, ok := val.(float64); ok {
			return int(num)
		}
	}
	return 0
}
