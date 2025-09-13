package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jperdior/chatbot-kit/application/query"
)

// Config represents server configuration
type Config struct {
	Port         string
	Environment  string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

// NewConfig creates a new server configuration
func NewConfig() *Config {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	env := os.Getenv("ENVIRONMENT")
	if env == "" {
		env = "production"
	}

	return &Config{
		Port:         port,
		Environment:  env,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
}

// Server represents the HTTP server
type Server struct {
	config     *Config
	router     *gin.Engine
	httpServer *http.Server
	queryBus   query.Bus
}

// NewServer creates a new server instance
func NewServer(config *Config, queryBus query.Bus) *Server {
	// Set Gin mode based on environment
	if config.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	router := gin.Default()

	server := &Server{
		config:   config,
		router:   router,
		queryBus: queryBus,
	}

	// Register routes
	server.registerRoutes(queryBus)

	return server
}

// GetRouter returns the Gin router
func (s *Server) GetRouter() *gin.Engine {
	return s.router
}

// Start starts the HTTP server
func (s *Server) Start() error {
	s.httpServer = &http.Server{
		Addr:         ":" + s.config.Port,
		Handler:      s.router,
		ReadTimeout:  s.config.ReadTimeout,
		WriteTimeout: s.config.WriteTimeout,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("Starting WooCommerce MCP server on port %s", s.config.Port)
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	// Create a context with timeout for shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Shutdown the server
	if err := s.httpServer.Shutdown(ctx); err != nil {
		return fmt.Errorf("server forced to shutdown: %w", err)
	}

	log.Println("Server exited")
	return nil
}

// Stop stops the HTTP server
func (s *Server) Stop() error {
	if s.httpServer == nil {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	return s.httpServer.Shutdown(ctx)
}
