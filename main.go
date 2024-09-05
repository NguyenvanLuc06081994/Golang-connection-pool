package main

import (
	"demo_connection_pool/routes"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // Import the PostgreSQL driver
	"log"
	"time"
)

var dsn = "postgres://postgres:password1@localhost:5433/postgres?sslmode=disable"

func main() {
	// Postgres allows 100 connections in default
	// Set the maximum number of idle connections in the pool
	idleConn := 50
	// Set the maximum number of connections in the pool
	maxConnections := 100
	// Set the maximum amount of time a connection can be reused
	maxConnLifetime := 2 * time.Minute
	poolConn, err := sqlx.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer poolConn.Close()
	poolConn.SetMaxOpenConns(maxConnections)
	poolConn.SetMaxIdleConns(idleConn)
	poolConn.SetConnMaxLifetime(maxConnLifetime)

	// normal connection
	conn, err := sqlx.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	// default will be 2 idle connections
	// so set it to 1 to simulate
	conn.SetMaxIdleConns(1)

	// Initialize the HTTP router
	router := gin.Default()
	router.StaticFile("/", "./index.html")

	// Initialize routes with the database connections
	routes.Conn = conn
	routes.PoolConn = poolConn
	routes.InitRoutes(router)

	// Start the HTTP server
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Unable to start HTTP server: %v\n", err)
	}
}
