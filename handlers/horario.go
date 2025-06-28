package handlers

import (
	"clinic-backend/database"
	"clinic-backend/models"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func CreateHorario(c *fiber.Ctx) error {
	var horario models.Horario
	if err := c.BodyParser(&horario); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Datos inválidos"})
	}

	err := database.DB.QueryRow(`
        INSERT INTO horario (turno, id_consultorio, id_medico) 
        VALUES ($1, $2, $3) 
        RETURNING id_horario`,
		horario.Turno, horario.IDConsultorio, horario.IDMedico,
	).Scan(&horario.ID)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error al crear horario"})
	}

	return c.Status(fiber.StatusCreated).JSON(horario)
}

func GetHorariosByMedico(c *fiber.Ctx) error {
	idMedico, err := strconv.Atoi(c.Params("id_medico"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "ID inválido"})
	}

	rows, err := database.DB.Query(`
        SELECT * FROM horario WHERE id_medico = $1`, idMedico)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error en consulta"})
	}
	defer rows.Close()

	var horarios []models.Horario
	for rows.Next() {
		var h models.Horario
		if err := rows.Scan(&h.ID, &h.Turno, &h.IDConsultorio, &h.IDMedico); err != nil {
			continue
		}
		horarios = append(horarios, h)
	}

	return c.JSON(horarios)
}

func DeleteHorario(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "ID inválido"})
	}

	_, err = database.DB.Exec("DELETE FROM horario WHERE id_horario = $1", id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error al eliminar horario"})
	}

	return c.SendStatus(fiber.StatusNoContent)
}
