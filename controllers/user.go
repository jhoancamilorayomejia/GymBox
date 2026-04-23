package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jhoancamilorayomejia/GymBox/db"
	"github.com/jhoancamilorayomejia/GymBox/models"
	"golang.org/x/crypto/bcrypt"
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

func CreateUser(c *gin.Context) {
	var user models.User

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
		return
	}

	// ✅ Validar username duplicado
	var exists bool
	err := db.DB.QueryRow(
		"SELECT EXISTS(SELECT 1 FROM users WHERE username=$1)",
		user.Username,
	).Scan(&exists)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al validar username"})
		return
	}

	if exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "El username ya existe"})
		return
	}

	// 🔐 Encriptar contraseña
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al encriptar la contraseña"})
		return
	}

	// ✅ Insertar con rol fijo = admin
	query := `
		INSERT INTO users
		(username, password, rol)
		VALUES ($1,$2,'admin')
		RETURNING id
	`

	err = db.DB.QueryRow(
		query,
		user.Username,
		string(hashedPassword),
	).Scan(&user.ID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al crear un nuevo Administrador"})
		return
	}

	user.Rol = "admin"

	// ⚠️ Opcional: no devolver password
	user.Password = ""

	c.JSON(http.StatusCreated, user)
}

// funcion Delete
func DeleteUser(c *gin.Context) {
	id := c.Param("id")

	query := `DELETE FROM users WHERE id = $1`

	result, err := db.DB.Exec(query, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al eliminar usuario"})
		return
	}

	rowsAffected, _ := result.RowsAffected()

	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Usuario no encontrado"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Usuario eliminado correctamente"})
}

// funcion UPDATE
func UpdateUser(c *gin.Context) {
	id := c.Param("id")

	var user models.User

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
		return
	}

	// 🔍 Validar username duplicado (excepto el mismo usuario)
	var exists bool
	err := db.DB.QueryRow(
		"SELECT EXISTS(SELECT 1 FROM users WHERE username=$1 AND id != $2)",
		user.Username, id,
	).Scan(&exists)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al validar username"})
		return
	}

	if exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "El username ya existe"})
		return
	}

	// 🔐 Si envían password → encriptar
	var query string
	var args []interface{}

	if user.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al encriptar contraseña"})
			return
		}

		query = `
			UPDATE users
			SET username=$1, password=$2
			WHERE id=$3
		`
		args = []interface{}{user.Username, string(hashedPassword), id}

	} else {
		// 👉 Solo actualiza username
		query = `
			UPDATE users
			SET username=$1
			WHERE id=$2
		`
		args = []interface{}{user.Username, id}
	}

	result, err := db.DB.Exec(query, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al actualizar usuario"})
		return
	}

	rowsAffected, _ := result.RowsAffected()

	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Usuario no encontrado"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Usuario actualizado correctamente"})
}

// ✅ Simple: busca por username en tabla users y actualiza password
/*func UpdatePasswordByUsername(c *gin.Context) {
	var body struct {
		Username string `json:"username"`
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

	// 🔐 Encriptar directamente
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al encriptar la contraseña"})
		return
	}

	// ✅ UPDATE directo en tabla users donde username coincide
	result, err := db.DB.Exec(
		"UPDATE users SET password = $1 WHERE username = $2",
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
		c.JSON(http.StatusNotFound, gin.H{"error": "Usuario no encontrado"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Contraseña actualizada correctamente"})
}
*/
