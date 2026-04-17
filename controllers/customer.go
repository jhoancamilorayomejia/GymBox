package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jhoancamilorayomejia/GymBox/db"
	"github.com/jhoancamilorayomejia/GymBox/models"
)

func GetCustomer(c *gin.Context) {
	rows, err := db.DB.Query(`
		SELECT idcustomer, cedula, name, lastname, phone, password FROM customers`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var customers []models.Customer

	for rows.Next() {
		var r models.Customer

		if err := rows.Scan(
			&r.IDCustomer,
			&r.Cedula,
			&r.Name,
			&r.LastName,
			&r.Phone,
			&r.Password,
		); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		customers = append(customers, r)
	}

	if err := rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, customers)
}
