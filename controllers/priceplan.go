package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jhoancamilorayomejia/GymBox/db"
	"github.com/jhoancamilorayomejia/GymBox/models"
)

func GetPricePlans(c *gin.Context) {
	rows, err := db.DB.Query(`
		SELECT idpriceplan, typeplan, price FROM price_plan`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var priceplans []models.PricePlan

	for rows.Next() {
		var r models.PricePlan

		if err := rows.Scan(
			&r.IDPricePlan,
			&r.TypePlan,
			&r.Price,
		); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		priceplans = append(priceplans, r)
	}

	if err := rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if len(priceplans) == 0 {
		c.JSON(http.StatusOK, []models.PricePlan{})
		return
	}

	c.JSON(http.StatusOK, priceplans)
}

func UpdatePricePlan(c *gin.Context) {
	id := c.Param("id")

	var body struct {
		Price float64 `json:"price"`
	}

	// Bind del JSON
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validación básica
	if body.Price <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "El precio debe ser mayor a 0"})
		return
	}

	// Ejecutar update
	result, err := db.DB.Exec(`
		UPDATE price_plan
		SET price = $1
		WHERE idpriceplan = $2
	`, body.Price, id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Verificar si realmente se actualizó algo
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Plan no encontrado"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Precio actualizado correctamente",
	})
}
