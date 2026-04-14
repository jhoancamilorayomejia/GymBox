package controllers

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/jhoancamilorayomejia/barbershop/db"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// 🔐 Clave secreta para JWT
var jwtSecret = []byte("mi_clave_secreta_super_segura")

// =========================
// 📥 REQUEST LOGIN
// =========================
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// =========================
// 👤 MODELO USER
// =========================
type User struct {
	Username string
	Password string // HASH (bcrypt)
	Rol      string
}

// =========================
// 🔐 LOGIN
// =========================
func Login(c *gin.Context) {
	var req LoginRequest

	// Validar JSON
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
		return
	}

	// Buscar usuario en DB
	var user User
	err := db.DB.QueryRow(
		`SELECT username, password, rol FROM users WHERE username=$1`,
		req.Username,
	).Scan(&user.Username, &user.Password, &user.Rol)

	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuario o contraseña incorrectos"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al consultar la base de datos"})
		}
		return
	}

	// 🔑 Comparar contraseña (bcrypt)
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuario o contraseña incorrectos"})
		return
	}

	// =========================
	// 🎯 CREAR TOKEN (2 HORAS)
	// =========================
	expirationTime := time.Now().Add(time.Hour * 2)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": user.Username,
		"rol":      user.Rol,
		"exp":      expirationTime.Unix(),
		"iat":      time.Now().Unix(), // opcional pero recomendado
	})

	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo generar el token"})
		return
	}

	// Respuesta
	c.JSON(http.StatusOK, gin.H{
		"message": "Login exitoso",
		"user":    user.Username,
		"rol":     user.Rol,
		"token":   tokenString,
		"expires": expirationTime, // 👈 útil para frontend
	})
}

// =========================
// 🔎 VALIDAR TOKEN
// =========================
func ValidateToken(tokenString string) (*jwt.Token, error) {

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {

		// Validar método de firma
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}

		return jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	return token, nil
}
