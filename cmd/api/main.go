package main

import (
	// "fmt"

	"log"
	"todo_api/internal/config"
	"todo_api/internal/database"
	"todo_api/internal/handlers"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	// Database entry point
	var cfg *config.Config
	var err error

	cfg, err = config.Load()
	if err != nil {
		log.Fatal("Failed to load configuration: ", err)
	}

	var pool *pgxpool.Pool
	pool, err = database.Connect(cfg.DatabaseURL)

	if err != nil {
		log.Fatal("Failed to connect to database: ", err)
	}

	defer pool.Close()

	var router *gin.Engine = gin.Default()
	router.SetTrustedProxies(nil)
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message":  "Todo API is running",
			"status":   "success",
			"database": "connected",
		})
	})

	router.POST("/todos", handlers.CreateTodoHandler(pool))

	router.GET("/todos", handlers.GetAllTodosHandler(pool))

	router.GET("/todos/:id", handlers.GetTodoByIDHandler(pool))

	router.PUT("/todos/:id", handlers.UpdateToDoHandler(pool))

	router.DELETE("/todos/:id", handlers.DeleteTodoHandler(pool))

	router.Run(":" + cfg.Port)
}
