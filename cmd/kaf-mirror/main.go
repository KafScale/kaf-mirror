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

// @title kaf-mirror API
// @version 1.0
// @description This is the API for kaf-mirror, a high-performance Kafka replication tool.
// @host localhost:8080
// @BasePath /api/v1
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
package main

import (
	"fmt"
	"kaf-mirror/internal/config"
	"kaf-mirror/internal/database"
	"kaf-mirror/internal/manager"
	"kaf-mirror/internal/server"
	"kaf-mirror/pkg/logger"
	"kaf-mirror/pkg/utils"
	"log"
	"os"
	"os/signal"
	"syscall"
)

var (
	Version string
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}
	fmt.Println("Configuration loaded successfully.")

	if err := logger.InitializeFromConfig(cfg.Logging.File, cfg.Logging.Level, cfg.Logging.Console); err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	logger.Info("Logger initialized with level %s, console=%t", cfg.Logging.Level, cfg.Logging.Console)

	// Initialize database
	db, err := database.InitDB(cfg.Database.Path)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	// Check if there are any users in the database
	users, err := database.ListUsers(db)
	if err != nil {
		log.Fatalf("Failed to check for users: %v", err)
	}
	if len(users) == 0 {
		fmt.Println("No users found in the database. Seeding default roles and creating initial admin user...")
		if err := database.SeedDefaultRolesAndPermissions(db); err != nil {
			log.Fatalf("Failed to seed default roles and permissions: %v", err)
		}

		password, err := utils.GenerateRandomPassword(16)
		if err != nil {
			log.Fatalf("Failed to generate password for initial admin: %v", err)
		}

		user, err := database.CreateUser(db, "admin@localhost", password, true)
		if err != nil {
			log.Fatalf("Failed to create initial admin user: %v", err)
		}

		var adminRoleID int
		if err := db.Get(&adminRoleID, "SELECT id FROM roles WHERE name = 'admin'"); err != nil {
			log.Fatalf("Failed to find admin role: %v", err)
		}
		if err := database.AssignRoleToUser(db, user.ID, adminRoleID); err != nil {
			log.Fatalf("Failed to assign admin role: %v", err)
		}

		fmt.Println("=================================================================")
		fmt.Println("  INITIAL ADMIN USER CREATED")
		fmt.Println("=================================================================")
		fmt.Printf("  Username: %s\n", user.Username)
		fmt.Printf("  Password: %s\n", password)
		fmt.Println("=================================================================")
		fmt.Println("  Please store this password in a secure location.")
		fmt.Println("=================================================================")
	}
	fmt.Println("Database initialized successfully.")

	// Initialize the Hub and JobManager
	hub := server.NewHub()
	jobManager := manager.New(db, cfg, hub)

	// Start all jobs
	if err := jobManager.RestartAllJobs(); err != nil {
		logger.Error("Failed to restart jobs on startup: %v", err)
	}

	// Initialize and start the API server
	srv := server.New(cfg, db, jobManager, hub, Version)
	go func() {
		addr := fmt.Sprintf("0.0.0.0:%d", cfg.Server.Port)
		fmt.Println("Starting API server on", addr)
		if cfg.Server.TLS.Enabled {
			if err := srv.App.ListenTLS(addr, cfg.Server.TLS.CertFile, cfg.Server.TLS.KeyFile); err != nil {
				log.Printf("API server error: %v", err)
			}
		} else {
			if err := srv.App.Listen(addr); err != nil {
				log.Printf("API server error: %v", err)
			}
		}
	}()
	fmt.Println("API server started successfully.")

	// Wait for a shutdown signal
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)
	<-shutdown

	fmt.Println("Shutting down...")
	// The kaf-mirror is now managed by the JobManager, so we don't need to stop it here.
	// In a real implementation, the JobManager would have a StopAll method.
	if err := srv.Shutdown(); err != nil {
		log.Printf("API server shutdown error: %v", err)
	}
	fmt.Println("Shutdown complete.")
}
