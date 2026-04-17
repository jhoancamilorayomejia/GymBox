package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jhoancamilorayomejia/GymBox/db"
)

type Reservation struct {
	ID              int    `json:"id"`
	Name            string `json:"name"`
	Phone           string `json:"phone"`
	ReservationDate string `json:"reservation_date"` // 👈 AHORA STRING
	Note            string `json:"note"`
}

func CreateReservation(c *gin.Context) {
	var reservation Reservation

	// ✅ Bind JSON
	if err := c.ShouldBindJSON(&reservation); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// ✅ Validación básica
	if reservation.Name == "" || reservation.Phone == "" || reservation.ReservationDate == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Campos obligatorios faltantes",
		})
		return
	}

	// ✅ Verificar si ya existe esa reserva (exacta)
	var exists bool
	err := db.DB.QueryRow(
		`SELECT EXISTS(
			SELECT 1 FROM reservations 
			WHERE reservation_date = $1
		)`,
		reservation.ReservationDate,
	).Scan(&exists)

	if err != nil {
		println("ERROR CHECK:", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	if exists {
		c.JSON(http.StatusConflict, gin.H{
			"error": "Ese horario ya está reservado",
		})
		return
	}

	// ✅ Insertar
	err = db.DB.QueryRow(
		`INSERT INTO reservations (name, phone, reservation_date, note)
		VALUES ($1,$2,$3,$4)
		RETURNING id`,
		reservation.Name,
		reservation.Phone,
		reservation.ReservationDate,
		reservation.Note,
	).Scan(&reservation.ID)

	if err != nil {
		println("ERROR INSERT:", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	// ✅ Respuesta
	c.JSON(http.StatusCreated, reservation)
}

func GetReservation(c *gin.Context) {
	rows, err := db.DB.Query(`
		SELECT id, name, phone, reservation_date, note FROM reservations
	`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	reservations := []Reservation{}

	for rows.Next() {
		var r Reservation

		if err := rows.Scan(
			&r.ID,
			&r.Name,
			&r.Phone,
			&r.ReservationDate,
			&r.Note,
		); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		reservations = append(reservations, r)
	}

	c.JSON(http.StatusOK, reservations)

	if err := rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
}
