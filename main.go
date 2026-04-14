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

	// ✅ Ruta correcta
	r.POST("/api/reservations", controllers.CreateReservation)

	// ✅ Ruta de login
	r.POST("/api/login", controllers.Login)

	//ruta para admin
	r.GET("/api/proyectos", AuthMiddleware(), controllers.GetReservation)

	log.Println("🚀 Servidor corriendo en http://localhost:8080")

	// Levantar servidor
	if err := r.Run(":8080"); err != nil {
		log.Fatal("❌ Error iniciando servidor:", err)
	}
}

/*package main

import (
	"fmt"
	"log"
	"net/http"

	"supplier-jwt/controllers"
	"supplier-jwt/db"

	"github.com/gin-gonic/gin"
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

		var tokenStr string
		_, err := fmt.Sscanf(authHeader, "Bearer %s", &tokenStr)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Formato inválido. Use Bearer <token>"})
			c.Abort()
			return
		}

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
	database, err := db.ConnectDB()
	if err != nil {
		log.Fatal("No se pudo conectar a la base de datos:", err)
	}
	defer database.Close()

	r := gin.Default()

	// CORS CORRECTO
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(200)
			return
		}
		c.Next()
	})

	// Endpoint publico
	r.POST("/api/login", controllers.Login)

	// Endpoints protegidos
	r.GET("/api/suppliers", AuthMiddleware(), controllers.GetSuppliers)
	r.POST("/api/suppliers", AuthMiddleware(), controllers.CreateSupplier)
	r.PUT("/api/suppliers/:id", AuthMiddleware(), controllers.UpdateSupplier)
	r.DELETE("/api/suppliers/:id", AuthMiddleware(), controllers.DeleteSupplier)

	//para proyectos
	r.GET("/api/proyectos", AuthMiddleware(), controllers.GetSuppliers2)

	log.Println("Servidor corriendo en http://localhost:8080")
	r.Run(":8080")
}
*/
