package main

import (
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jhoancamilorayomejia/barbershop/controllers"
	"github.com/jhoancamilorayomejia/barbershop/db"
)

// Middleware JWT
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token requerido"})
			c.Abort()
			return
		}

		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Formato inválido. Use Bearer <token>"})
			c.Abort()
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

		token, err := controllers.ValidateToken(tokenStr)
		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token inválido o expirado"})
			c.Abort()
			return
		}

		c.Next()
	}
}

func main() {
	_, err := db.ConnectDB()
	if err != nil {
		log.Fatal("❌ Error conectando a la DB:", err)
	}

	r := gin.Default()

	// ✅ CORS — en producción usa tu dominio de Railway
	allowedOrigin := os.Getenv("ALLOWED_ORIGIN")
	if allowedOrigin == "" {
		allowedOrigin = "*" // fallback para desarrollo local
	}

	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", allowedOrigin)
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// ── Rutas API ──────────────────────────────────────────────────
	r.POST("/api/login", controllers.Login)
	r.POST("/api/reservations", controllers.CreateReservation)
	r.GET("/api/proyectos", AuthMiddleware(), controllers.GetReservation)
	r.DELETE("/api/reservations/:id", AuthMiddleware(), controllers.DeleteReservation)
	r.PUT("/api/users/change-password", AuthMiddleware(), controllers.ChangePassword)

	// ── Servir frontend Vue (dist/) ────────────────────────────────
	r.Static("/assets", "./dist/assets")
	r.StaticFile("/favicon.ico", "./dist/favicon.ico")

	// Todo lo que no sea una ruta real del backend → index.html (Vue Router)
	r.NoRoute(func(c *gin.Context) {
		c.File("./dist/index.html")
	})

	// ✅ Puerto dinámico — Railway asigna el puerto via variable PORT
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // fallback para desarrollo local
	}

	log.Println("🚀 Servidor corriendo en puerto", port)

	if err := r.Run(":" + port); err != nil {
		log.Fatal("❌ Error iniciando servidor:", err)
	}
}
