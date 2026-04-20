package main

import (
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jhoancamilorayomejia/GymBox/controllers"
	"github.com/jhoancamilorayomejia/GymBox/db"
)

// Middleware JWT (MEJORADO)
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		// ❌ No hay token
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token requerido"})
			c.Abort()
			return
		}

		// ❌ Formato incorrecto
		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Formato inválido. Use Bearer <token>"})
			c.Abort()
			return
		}

		// ✅ Extraer token correctamente
		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

		// 🔎 Validar token
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

	// Crear router
	r := gin.Default()

	// CORS (permite conexión con Vue)
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// ✅ Ruta de login Publica
	r.POST("/api/login", controllers.Login)

	//ruta para admin, con permisos token
	r.GET("/api/users", AuthMiddleware(), controllers.GetUser)
	r.GET("/api/customers", AuthMiddleware(), controllers.GetCustomer)
	r.GET("/api/plans", AuthMiddleware(), controllers.GetPlan)
	r.POST("/api/plans", AuthMiddleware(), controllers.CreatePlan)

	r.POST("/api/new-users", AuthMiddleware(), controllers.CreateUser)
	r.DELETE("/api/delete-users/:id", controllers.DeleteUser)
	r.PUT("/api/update-users/:id", controllers.UpdateUser)

	r.POST("/api/new-customers", controllers.CreateCustomer)
	r.PUT("/api/update-customers/:id", controllers.UpdateCustomer)
	r.DELETE("/api/delete-customers/:id", controllers.DeleteCustomer)

	//
	r.GET("/api/customers/by-cedula/:cedula", AuthMiddleware(), controllers.GetCustomerByCedula)

	log.Println("🚀 Servidor corriendo en http://localhost:8080")

	// Levantar servidor
	if err := r.Run(":8080"); err != nil {
		log.Fatal("❌ Error iniciando servidor:", err)
	}
}
