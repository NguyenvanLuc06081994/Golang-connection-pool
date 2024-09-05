package routes

import (
	"demo_connection_pool/models"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"log"
	"net/http"
	"time"
)

var (
	Conn      *sqlx.DB
	PoolConn  *sqlx.DB
	query     = "SELECT id, name, price, description FROM products limit 1000"
	allCount  int64
	allTime   int64
	poolCount int64
	poolTime  int64
	newCount  int64
	newTime   int64
	dsn       = "postgres://postgres:password1@localhost:5433/postgres?sslmode=disable"
)

func InitRoutes(router *gin.Engine) {
	router.GET("/products/normal", func(c *gin.Context) {
		startTime := time.Now()

		// Query the database for all products
		rows, err := Conn.Queryx(query)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		products, err := scanProducts(rows)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		elapsed := time.Since(startTime).Microseconds()
		allCount++
		allTime += elapsed
		c.JSON(http.StatusOK, models.Response{Elapsed: elapsed, Average: float64(allTime / allCount), Products: products})
	})

	router.GET("/products/pooled", func(c *gin.Context) {
		startTime := time.Now()
		// Query the database for all products
		rows, err := PoolConn.Queryx(query)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		products, err := scanProducts(rows)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		elapsed := time.Since(startTime).Microseconds()
		poolCount++
		poolTime += elapsed
		c.JSON(http.StatusOK, models.Response{Elapsed: elapsed, Average: float64(poolTime / poolCount), Products: products})
	})

	router.GET("/products/new", func(c *gin.Context) {
		startTime := time.Now()
		conn, err := sqlx.Open("postgres", dsn)
		if err != nil {
			log.Fatalf("Unable to connect to database: %v\n", err)
		}

		rows, err := conn.Queryx(query)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		products, err := scanProducts(rows)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		elapsed := time.Since(startTime).Microseconds()
		newCount++
		newTime += elapsed
		c.JSON(http.StatusOK, models.Response{Elapsed: elapsed, Average: float64(newTime / newCount), Products: products})
	})
}

func scanProducts(rows *sqlx.Rows) ([]*models.Product, error) {
	var products []*models.Product
	for rows.Next() {
		var product models.Product
		if err := rows.StructScan(&product); err != nil {
			return nil, err
		}
		products = append(products, &product)
	}
	return products, nil
}
