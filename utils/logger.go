package utils

import (
	"clinic-backend/database"
	"time"
)

func LogOperacion(userID *int, operacion, entidad, detalle string, exito bool) {
	_, err := database.DB.Exec(`
		INSERT INTO logs_auditoria (usuario_id, operacion, entidad, detalle, exito, timestamp)
		VALUES ($1, $2, $3, $4, $5, $6)
	`, userID, operacion, entidad, detalle, exito, time.Now())

	if err != nil {
		println("Error guardando log:", err.Error())
	}
}
