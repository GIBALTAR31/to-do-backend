package handlers

import (
	"net/http"
	"todo_api/internal/models"
	"todo_api/internal/repository"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

type RegisterRequest struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

// CreateUserHandler godoc
// @Summary Register a new user
// @Description Register a new user with email and password
// @Tags auth
// @Accept json
// @Produce json
// @Param user body RegisterRequest true "User registration input"
// @Success 201 {object} models.User
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /auth/register [post]
func CreateUserHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		var request RegisterRequest
		
		if err := c.BindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if len(request.Password) < 6 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "password must be at least 6 characters"})
			return
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to hash password" +  err.Error()})
			return
		}

		user := &models.User{
			Email:    request.Email,
			Password: string(hashedPassword),
		}

		createdUser, err := repository.CreateUser(pool, user)
		if err != nil {
			if strings.Contains(err.Error(), "duplicate") || strings.Contains(err.Error(), "unique") {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Email already registered"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create user" + " " + err.Error()})
			return
		}
		c.JSON(http.StatusCreated, createdUser)
	}
}
