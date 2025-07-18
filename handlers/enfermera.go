package handlers

import (
	"clinic-backend/database"
	"clinic-backend/models"

	"github.com/gofiber/fiber/v2"
)

func GetPacientesParaEnfermera(c *fiber.Ctx) error {
	rows, err := database.DB.Query(`
		SELECT DISTINCT id_paciente, nombre, seguro, contacto
		FROM paciente
		ORDER BY id_paciente
	`)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error al obtener pacientes",
		})
	}
	defer rows.Close()

	var pacientes []models.Paciente
	for rows.Next() {
		var p models.Paciente
		if err := rows.Scan(&p.ID, &p.Nombre, &p.Seguro, &p.Contacto); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Error al leer paciente",
			})
		}
		pacientes = append(pacientes, p)
	}

	return c.JSON(pacientes)
}
