package controllers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jhoancamilorayomejia/GymBox/db"
	"github.com/jhoancamilorayomejia/GymBox/models"
	"golang.org/x/crypto/bcrypt"
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

func GetCustomerByCedula(c *gin.Context) {
	cedula := c.Param("cedula")

	var customer models.Customer
	err := db.DB.QueryRow(`
        SELECT idcustomer, cedula, name, lastname, phone, password
        FROM customers
        WHERE cedula = $1
    `, cedula).Scan(
		&customer.IDCustomer,
		&customer.Cedula,
		&customer.Name,
		&customer.LastName,
		&customer.Phone,
		&customer.Password,
	)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Cliente no encontrado"})
		return
	}
	c.JSON(http.StatusOK, customer)
}

func CreateCustomer(c *gin.Context) {
	var customer models.Customer

	if err := c.ShouldBindJSON(&customer); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
		return
	}

	// ✅ Validar cédula duplicada
	var exists bool
	err := db.DB.QueryRow(
		"SELECT EXISTS(SELECT 1 FROM customers WHERE cedula=$1)",
		customer.Cedula,
	).Scan(&exists)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al validar cédula"})
		return
	}

	if exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "La cédula ya está registrada"})
		return
	}

	// 🔐 Encriptar contraseña
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(customer.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al encriptar contraseña"})
		return
	}

	query := `
		INSERT INTO customers
		(cedula, name, lastname, phone, password)
		VALUES ($1,$2,$3,$4,$5)
		RETURNING idcustomer
	`

	err = db.DB.QueryRow(
		query,
		customer.Cedula,
		customer.Name,
		customer.LastName,
		customer.Phone,
		string(hashedPassword),
	).Scan(&customer.IDCustomer)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al crear cliente"})
		return
	}

	customer.Password = ""

	c.JSON(http.StatusCreated, customer)
}

func DeleteCustomer(c *gin.Context) {
	id := c.Param("id")

	query := `DELETE FROM customers WHERE idcustomer=$1`

	result, err := db.DB.Exec(query, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al eliminar cliente"})
		return
	}

	rowsAffected, _ := result.RowsAffected()

	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Cliente no encontrado"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Cliente eliminado correctamente"})
}

func UpdateCustomer(c *gin.Context) {
	id := c.Param("id")

	var customer models.Customer

	if err := c.ShouldBindJSON(&customer); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
		return
	}

	// ✅ Validar cédula duplicada (excepto el mismo)
	var exists bool
	err := db.DB.QueryRow(
		"SELECT EXISTS(SELECT 1 FROM customers WHERE cedula=$1 AND idcustomer != $2)",
		customer.Cedula, id,
	).Scan(&exists)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al validar cédula"})
		return
	}

	if exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "La cédula ya está registrada"})
		return
	}

	var query string
	var args []interface{}

	// 🔐 Si envían password → encriptar
	if customer.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(customer.Password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al encriptar contraseña"})
			return
		}

		query = `
			UPDATE customers
			SET cedula=$1, name=$2, lastname=$3, phone=$4, password=$5
			WHERE idcustomer=$6
		`
		args = []interface{}{
			customer.Cedula,
			customer.Name,
			customer.LastName,
			customer.Phone,
			string(hashedPassword),
			id,
		}

	} else {
		// 👉 sin password
		query = `
			UPDATE customers
			SET cedula=$1, name=$2, lastname=$3, phone=$4
			WHERE idcustomer=$5
		`
		args = []interface{}{
			customer.Cedula,
			customer.Name,
			customer.LastName,
			customer.Phone,
			id,
		}
	}

	result, err := db.DB.Exec(query, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al actualizar cliente"})
		return
	}

	rowsAffected, _ := result.RowsAffected()

	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Cliente no encontrado"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Cliente actualizado correctamente"})
}

func UpdatePasswordByUsername(c *gin.Context) {
	var body struct {
		Username string `json:"username"` // aquí viene la cédula
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
		return
	}

	if body.Username == "" || body.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username y contraseña son requeridos"})
		return
	}

	// 🔐 Encriptar contraseña
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al encriptar la contraseña"})
		return
	}

	// ✅ UPDATE en customers usando cedula = username
	result, err := db.DB.Exec(
		"UPDATE customers SET password = $1 WHERE cedula = $2",
		string(hashedPassword),
		body.Username,
	)
	if err != nil {
		fmt.Println("❌ Error SQL:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Cliente no encontrado"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Contraseña actualizada correctamente"})
}
