package controllers

import (
	//"fmt"
	"net/http"
	//"net/url"

	"github.com/gin-gonic/gin"
	"github.com/jhoancamilorayomejia/barbershop/db"
)

type Reservation struct {
	ID              int    `json:"id"`
	Name            string `json:"name"`
	Phone           string `json:"phone"`
	ReservationDate string `json:"reservation_date"` // 👈 AHORA STRING
	Note            string `json:"note"`
}

/*func sendWhatsApp(res Reservation) {
	phone := "573XXXXXXXXX" // 👈 TU número (con código país)
	apikey := "1234567"     // 👈 tu API KEY de CallMeBot

	msg := fmt.Sprintf(
		"💈 Nueva Reserva\n👤 %s\n📞 %s\n📅 %s\n📝 %s",
		res.Name,
		res.Phone,
		res.ReservationDate,
		res.Note,
	)

	apiURL := fmt.Sprintf(
		"https://api.callmebot.com/whatsapp.php?phone=%s&text=%s&apikey=%s",
		phone,
		url.QueryEscape(msg),
		apikey,
	)

	_, err := http.Get(apiURL)
	if err != nil {
		fmt.Println("Error enviando WhatsApp:", err)
	}
} */

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
	// 🚀 ENVÍO DE WHATSAPP (NO BLOQUEA)
	//go sendWhatsApp(reservation)

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
