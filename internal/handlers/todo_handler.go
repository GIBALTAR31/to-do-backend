package handlers

import (
	"net/http"
	"strconv"
	"todo_api/internal/repository"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CreateTodoInput struct {
	Title     string `json:"title" binding:"required"`
	Completed bool   `json:"completed"`
}

type UpdateTodoInput struct {
	Title     *string `json:"title"`
	Completed *bool   `json:"completed"`
}

// CreateTodoHandler godoc
// @Summary Create a new todo
// @Description Create a new todo item
// @Tags todos
// @Accept json
// @Produce json
// @Param todo body CreateTodoInput true "Todo input"
// @Success 201 {object} models.Todo
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /todos [post]
func CreateTodoHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		userIDInterface, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "user_id not found in context"})
			return
		}
		userID := userIDInterface.(string)

		var input CreateTodoInput

		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		todo, err := repository.CreateTodo(pool, input.Title, input.Completed, userID)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, todo)
	}
}

// GetAllTodosHandler godoc
// @Summary Get all todos
// @Description Retrieve all todo items
// @Tags todos
// @Produce json
// @Success 200 {array} models.Todo
// @Failure 500 {object} map[string]string
// @Router /todos [get]
func GetAllTodosHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		userIDInterface, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "user_id not found in context"})
			return
		}
		userID := userIDInterface.(string)

		todos, err := repository.GetAllTodos(pool, userID)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, todos)
	}
}

// GetTodoByIDHandler godoc
// @Summary Get a todo by ID
// @Description Retrieve a todo by its ID
// @Tags todos
// @Produce json
// @Param id path int true "Todo ID"
// @Success 200 {object} models.Todo
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /todos/{id} [get]
func GetTodoByIDHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		userIDInterface, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "user_id not found in context"})
			return
		}
		userID := userIDInterface.(string)

		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid todo ID"})
			return
		}
		todo, err := repository.GetTodoByID(pool, id, userID)

		if err != nil {
			if err == pgx.ErrNoRows {
				c.JSON(http.StatusNotFound, gin.H{"error": "Todo not found"})
				return
			}

			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		c.JSON(http.StatusOK, todo)
	}
}

// UpdateToDoHandler godoc
// @Summary Update a todo
// @Description Update an existing todo by ID (partial update)
// @Tags todos
// @Accept json
// @Produce json
// @Param id path int true "Todo ID"
// @Param todo body UpdateTodoInput true "Updated todo data"
// @Success 200 {object} models.Todo
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /todos/{id} [put]
func UpdateToDoHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		userIDInterface, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "user_id not found in context"})
			return
		}
		userID := userIDInterface.(string)
		
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid todo ID"})
		}

		var input UpdateTodoInput

		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if input.Title == nil && input.Completed == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "At least one field (title or completed) must be provided"})
			return
		}

		existing, err := repository.GetTodoByID(pool, id, userID)

		if err != nil {
			if err == pgx.ErrNoRows {
				c.JSON(http.StatusNotFound, gin.H{"error": "Todo not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		title := existing.Title
		if input.Title != nil {
			title = *input.Title
		}

		completed := existing.Completed
		if input.Completed != nil {
			completed = *input.Completed
		}

		todo, err := repository.UpdateTodoByID(pool, id, title, completed, userID)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, todo)
	}
}

// DeleteTodoHandler godoc
// @Summary Delete a todo
// @Description Delete a todo by ID
// @Tags todos
// @Produce json
// @Param id path int true "Todo ID"
// @Success 204 {string} string "No Content"
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /todos/{id} [delete]
func DeleteTodoHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		userIDInterface, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "user_id not found in context"})
			return
		}
		userID := userIDInterface.(string)
		
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid todo ID"})
			return
		}

		err = repository.DeleteTodoByID(pool, id, userID)

		if err != nil {
			if err == pgx.ErrNoRows {
				c.JSON(http.StatusNotFound, gin.H{"error": "Todo not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusNoContent, nil)
	}
}
