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
	r.GET("/api/login-price_plans", controllers.GetPricePlans)
	r.POST("/api/new-public-customers", controllers.CreateCustomer)

	//ruta para admin, con permisos token
	r.GET("/api/users", AuthMiddleware(), controllers.GetUser)
	r.GET("/api/customers", AuthMiddleware(), controllers.GetCustomer)

	//Users endpoint desde administrador
	r.POST("/api/new-users", AuthMiddleware(), controllers.CreateUser)
	r.DELETE("/api/delete-users/:id", AuthMiddleware(), controllers.DeleteUser)
	r.PUT("/api/update-users/:id", AuthMiddleware(), controllers.UpdateUser)

	//planes endpoint desde administrador
	r.GET("/api/plans", AuthMiddleware(), controllers.GetPlan)
	r.POST("/api/plans", AuthMiddleware(), controllers.CreatePlan)
	r.DELETE("/api/delete-plans/:id", AuthMiddleware(), controllers.DeletePlan)
	r.DELETE("/api/delete-plans-multiple", AuthMiddleware(), controllers.DeletePlan)
	r.PUT("/api/update-plans/:id", AuthMiddleware(), controllers.UpdatePlan)

	//para visualizar y actualizar precios de la membresia
	r.GET("/api/priceplan", AuthMiddleware(), controllers.GetPricePlans)
	r.PUT("/api/update-priceplan/:id", AuthMiddleware(), controllers.UpdatePricePlan)

	//mostrar precios de membresia para clientes customer
	r.GET("/api/customer-priceplans", AuthMiddleware(), controllers.GetPricePlans)
	// Pagos Mercado Pago
	r.POST("/api/payment/preference", AuthMiddleware(), controllers.CreatePreference)

	r.POST("/api/new-customers", AuthMiddleware(), controllers.CreateCustomer)
	r.PUT("/api/update-customers/:id", AuthMiddleware(), controllers.UpdateCustomer)
	r.DELETE("/api/delete-customers/:id", AuthMiddleware(), controllers.DeleteCustomer)

	// Con middleware de autenticación:
	//r.PUT("/api/users/update-password/:id", AuthMiddleware(), controllers.UpdatePassword)
	r.PUT("/api/users/update-password-by-username", AuthMiddleware(), controllers.UpdatePasswordByUsername)

	//
	r.GET("/api/customers/by-cedula/:cedula", AuthMiddleware(), controllers.GetCustomerByCedula)

	log.Println("🚀 Servidor corriendo en http://localhost:8080")

	// Levantar servidor
	if err := r.Run(":8080"); err != nil {
		log.Fatal("❌ Error iniciando servidor:", err)
	}
}
