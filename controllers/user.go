package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jhoancamilorayomejia/GymBox/db"
	"github.com/jhoancamilorayomejia/GymBox/models"
)

func GetUser(c *gin.Context) {
	rows, err := db.DB.Query(`
		SELECT id, username, password, rol FROM users`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var users []models.User

	for rows.Next() {
		var r models.User

		if err := rows.Scan(
			&r.ID,
			&r.Username,
			&r.Password,
			&r.Rol,
		); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		users = append(users, r)
	}

	if err := rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, users)
}
