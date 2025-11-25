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


package server

import (
	"fmt"
	"log"
	"kaf-mirror/internal/ai"
	"kaf-mirror/internal/config"
	"kaf-mirror/internal/database"
	"kaf-mirror/internal/manager"
	"kaf-mirror/internal/server/middleware"
	"kaf-mirror/pkg/logger"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"

	_ "kaf-mirror/web/docu/swagger"
)

// Server wraps the Fiber application.
type Server struct {
	App       *fiber.App
	cfg       *config.Config
	Db        *sqlx.DB
	manager   *manager.JobManager
	aiClient  *ai.Client
	hub       *Hub
	startTime time.Time
	Version   string
}

// New creates a new server instance.
func New(cfg *config.Config, db *sqlx.DB, manager *manager.JobManager, hub *Hub, version string) *Server {
	app := fiber.New()

	app.Use(middleware.Cors(cfg.Server.CORS.AllowedOrigins))

	var aiConfig config.AIConfig
	if dbConfig, err := database.LoadConfig(db); err == nil {
		aiConfig = dbConfig.AI
	} else {
		aiConfig = cfg.AI
	}

	s := &Server{
		App:       app,
		cfg:       cfg,
		Db:        db,
		manager:   manager,
		aiClient:  ai.NewClient(aiConfig),
		hub:       hub,
		startTime: time.Now(),
		Version:   version,
	}

	go hub.Run()
	s.setupRoutes()

	go func() {
		time.Sleep(2 * time.Second)
		if err := manager.RestartAllJobs(); err != nil {
			log.Printf("Error restarting jobs on server startup: %v", err)
		}
	}()

	return s
}

// Start runs the Fiber server.
func (s *Server) Start() error {
	addr := fmt.Sprintf("%s:%d", s.cfg.Server.Host, s.cfg.Server.Port)
	
	// Check if TLS is enabled
	if s.cfg.Server.TLS.Enabled {
		certFile := s.cfg.Server.TLS.CertFile
		keyFile := s.cfg.Server.TLS.KeyFile
		
		if certFile == "" || keyFile == "" {
			return fmt.Errorf("TLS is enabled but certificate or key file path is missing")
		}
		
		logger.Info("Starting HTTPS server on %s with TLS certificate: %s", addr, certFile)
		return s.App.ListenTLS(addr, certFile, keyFile)
	}
	
	logger.Info("Starting HTTP server on %s", addr)
	return s.App.Listen(addr)
}

// getWebPath returns the configured web path with fallback to development default
func (s *Server) getWebPath() string {
	if s.cfg.Server.WebPath != "" {
		return s.cfg.Server.WebPath
	}
	return "./web"
}

// Shutdown gracefully shuts down the server.
func (s *Server) Shutdown() error {
	return s.App.Shutdown()
}
