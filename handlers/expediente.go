package handlers

import (
	"clinic-backend/database"
	"clinic-backend/models"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func CreateExpediente(c *fiber.Ctx) error {
	var expediente models.Expediente
	if err := c.BodyParser(&expediente); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Datos inválidos"})
	}

	err := database.DB.QueryRow(`
        INSERT INTO expediente (antecedentes, historial_clinico, id_paciente) 
        VALUES ($1, $2, $3) 
        RETURNING id_expediente`,
		expediente.Antecedentes, expediente.HistorialClinico, expediente.IDPaciente,
	).Scan(&expediente.ID)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error al crear expediente"})
	}

	return c.Status(fiber.StatusCreated).JSON(expediente)
}

func GetExpedienteByPaciente(c *fiber.Ctx) error {
	idPaciente, err := strconv.Atoi(c.Params("id_paciente"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "ID inválido"})
	}

	var expediente models.Expediente
	err = database.DB.QueryRow(`
        SELECT * FROM expediente WHERE id_paciente = $1`, idPaciente,
	).Scan(&expediente.ID, &expediente.Antecedentes, &expediente.HistorialClinico, &expediente.IDPaciente)

	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Expediente no encontrado"})
	}

	return c.JSON(expediente)
}

func UpdateExpediente(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "ID inválido"})
	}

	var expediente models.Expediente
	if err := c.BodyParser(&expediente); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Datos inválidos"})
	}

	_, err = database.DB.Exec(`
        UPDATE expediente 
        SET antecedentes = $1, historial_clinico = $2 
        WHERE id_expediente = $3`,
		expediente.Antecedentes, expediente.HistorialClinico, id,
	)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error al actualizar expediente"})
	}

	expediente.ID = id
	return c.JSON(expediente)
}
