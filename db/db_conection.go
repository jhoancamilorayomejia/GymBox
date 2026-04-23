package db

import (
	"database/sql"
	"log"
	"os"
	"strings"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func ConnectDB() (*sql.DB, error) {
	dsn := os.Getenv("DATABASE_URL")

	if dsn == "" {
		dsn = "host=localhost user=postgres password=tu_password dbname=rayobox sslmode=disable"
		log.Println("⚠️  DATABASE_URL no encontrada, usando conexión local")
	} else {
		// ✅ Railway Postgres sin SSL
		if !strings.Contains(dsn, "sslmode") {
			dsn += "?sslmode=disable"
		}
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
