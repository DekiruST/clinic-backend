package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func Connect() error {
	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		return fmt.Errorf("DATABASE_URL no está configurada")
	}

	log.Println("Conectando a Supabase con:", connStr[:30]+"...") // Log parcial por seguridad

	var err error
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		return fmt.Errorf("error abriendo conexión: %w", err)
	}

	if err = DB.Ping(); err != nil {
		return fmt.Errorf("error conectando a la base de datos: %w", err)
	}

	log.Println("✅ Conectado exitosamente a Supabase")
	return nil
}
