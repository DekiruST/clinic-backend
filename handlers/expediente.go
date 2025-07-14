package handlers

import (
	"clinic-backend/database"
	"clinic-backend/models"

	"github.com/gofiber/fiber/v2"
)

// handlers/expediente.go
func GetAllExpedientes(c *fiber.Ctx) error {
	rows, err := database.DB.Query(`SELECT id_expediente, antecedentes, historial_clinico, id_paciente FROM expediente`)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Error al obtener expedientes"})
	}
	defer rows.Close()

	var expedientes []models.Expediente
	for rows.Next() {
		var exp models.Expediente
		if err := rows.Scan(&exp.ID, &exp.Antecedentes, &exp.HistorialClinico, &exp.IDPaciente); err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Error al leer expediente"})
		}
		expedientes = append(expedientes, exp)
	}
	return c.JSON(expedientes)
}
