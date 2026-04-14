package controllers

/*
import (
	"encoding/json"
	"net/http"
	"supplier-jwt/db"

	"github.com/gin-gonic/gin"
)

// ---------- ESTRUCTURA ----------
type Supplier2 struct {
	ID             int             `json:"id"`
	NIT            string          `json:"nit"`
	Nombre         string          `json:"nombre"`
	Apellido       string          `json:"apellido"`
	Cedula         string          `json:"cedula"`
	TipoProveedor  string          `json:"tipo_proveedor"`
	TipoPersona    string          `json:"tipo_persona"`
	Beneficiarios  json.RawMessage `json:"beneficiarios"`
	DatosBancarios json.RawMessage `json:"datos_bancarios"`
	Estado         string          `json:"estado"`
}

// ================= LISTAR =================
func GetSuppliers2(c *gin.Context) {
	rows, err := db.DB.Query(`
		SELECT id, nit, nombre, apellido, cedula, tipo_proveedor, tipo_persona, beneficiarios, datos_bancarios, estado
		FROM supplier
	`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	suppliers := []Supplier{}
	for rows.Next() {
		var s Supplier
		var b, d []byte
		var estado string

		if err := rows.Scan(
			&s.ID,
			&s.NIT,
			&s.Nombre,
			&s.Apellido,
			&s.Cedula,
			&s.TipoProveedor,
			&s.TipoPersona,
			&b,
			&d,
			&estado,
		); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		s.Beneficiarios = b
		s.DatosBancarios = d
		s.Estado = estado
		suppliers = append(suppliers, s)
	}

	c.JSON(http.StatusOK, suppliers)
}
*/
