package db

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func ConnectDB() (*sql.DB, error) {
	// ✅ Railway provee DATABASE_URL automáticamente
	// En desarrollo puedes setearla en un archivo .env o directamente
	dsn := os.Getenv("DATABASE_URL")

	if dsn == "" {
		// Fallback para desarrollo local
		dsn = "host=localhost user=postgres password=tu_password dbname=barbershop sslmode=disable"
		log.Println("⚠️  DATABASE_URL no encontrada, usando conexión local")
	}

	var err error
	DB, err = sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	if err = DB.Ping(); err != nil {
		return nil, err
	}

	log.Println("✅ Conectado a PostgreSQL")
	return DB, nil
}
