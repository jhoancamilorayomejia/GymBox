package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jhoancamilorayomejia/GymBox/db"
	"github.com/jhoancamilorayomejia/GymBox/models"
)

func GetPlan(c *gin.Context) {
	idcustomer := c.Query("idcustomer")

	// Validación básica
	if idcustomer == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "idcustomer es requerido",
		})
		return
	}

	rows, err := db.DB.Query(`
		SELECT idplan, idcustomer, typeplan, price, datepay, datestart, datefinish, state
		FROM plan
		WHERE idcustomer = $1
	`, idcustomer)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var plans []models.Plan

	for rows.Next() {
		var r models.Plan

		if err := rows.Scan(
			&r.IDPlan,
			&r.IDCustomer,
			&r.TypePlan,
			&r.Price,
			&r.DatePay,
			&r.DateStart,
			&r.DateFinish,
			&r.State,
		); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		plans = append(plans, r)
	}

	if err := rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, plans)
}

func CreatePlan(c *gin.Context) {
	var plan models.Plan

	// ✅ Bind JSON
	if err := c.ShouldBindJSON(&plan); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// ✅ Validación básica
	if plan.IDCustomer == 0 || plan.TypePlan == "" || plan.Price == 0 ||
		plan.DatePay == "" || plan.DateStart == "" || plan.DateFinish == "" || plan.State == "" {

		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Campos obligatorios faltantes",
		})
		return
	}

	// ✅ Insertar plan
	err := db.DB.QueryRow(
		`INSERT INTO plan (idcustomer, typeplan, price, datepay, datestart, datefinish, state)
		VALUES ($1,$2,$3,$4,$5,$6,$7)
		RETURNING idplan`,
		plan.IDCustomer,
		plan.TypePlan,
		plan.Price,
		plan.DatePay,
		plan.DateStart,
		plan.DateFinish,
		plan.State,
	).Scan(&plan.IDPlan)

	if err != nil {
		println("ERROR INSERT PLAN:", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	// ✅ Respuesta
	c.JSON(http.StatusCreated, plan)
}

func DeletePlan(c *gin.Context) {
	id := c.Param("id")

	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID es requerido"})
		return
	}

	query := `DELETE FROM plan WHERE idplan=$1`

	result, err := db.DB.Exec(query, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al eliminar el plan"})
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error verificando eliminación"})
		return
	}

	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Plan no encontrado"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Plan eliminado correctamente"})
}

func UpdatePlan(c *gin.Context) {
	id := c.Param("id")

	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID es requerido"})
		return
	}

	var plan models.Plan

	if err := c.ShouldBindJSON(&plan); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Validación básica
	if plan.TypePlan == "" || plan.Price == 0 ||
		plan.DatePay == "" || plan.DateStart == "" ||
		plan.DateFinish == "" || plan.State == "" {

		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Campos obligatorios faltantes",
		})
		return
	}

	query := `
		UPDATE plan
		SET typeplan=$1, price=$2, datepay=$3, datestart=$4, datefinish=$5, state=$6
		WHERE idplan=$7
	`

	result, err := db.DB.Exec(
		query,
		plan.TypePlan,
		plan.Price,
		plan.DatePay,
		plan.DateStart,
		plan.DateFinish,
		plan.State,
		id,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al actualizar el plan"})
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error verificando actualización"})
		return
	}

	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Plan no encontrado"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Plan actualizado correctamente"})
}
