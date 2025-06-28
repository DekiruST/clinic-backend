package handlers

import (
	"clinic-backend/database"
	"clinic-backend/models"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func CreateConsulta(c *fiber.Ctx) error {
	var consulta models.Consulta
	if err := c.BodyParser(&consulta); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Datos inválidos"})
	}

	err := database.DB.QueryRow(`
        INSERT INTO consulta (tipo, horario, diagnostico, costo, id_consultorio, id_paciente, id_medico) 
        VALUES ($1, $2, $3, $4, $5, $6, $7) 
        RETURNING id_consulta`,
		consulta.Tipo, consulta.Horario, consulta.Diagnostico, consulta.Costo,
		consulta.IDConsultorio, consulta.IDPaciente, consulta.IDMedico,
	).Scan(&consulta.ID)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error al crear consulta"})
	}

	return c.Status(fiber.StatusCreated).JSON(consulta)
}

func GetConsultas(c *fiber.Ctx) error {
	rows, err := database.DB.Query(`
        SELECT * FROM consulta
    `)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error en consulta"})
	}
	defer rows.Close()

	var consultas []models.Consulta
	for rows.Next() {
		var c models.Consulta
		if err := rows.Scan(
			&c.ID, &c.Tipo, &c.Horario, &c.Diagnostico, &c.Costo,
			&c.IDConsultorio, &c.IDPaciente, &c.IDMedico,
		); err != nil {
			continue
		}
		consultas = append(consultas, c)
	}

	return c.JSON(consultas)
}

func GetConsultasByPaciente(c *fiber.Ctx) error {
	idPaciente, err := strconv.Atoi(c.Params("id_paciente"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "ID inválido"})
	}

	rows, err := database.DB.Query(`
        SELECT * FROM consulta WHERE id_paciente = $1
    `, idPaciente)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error en consulta"})
	}
	defer rows.Close()

	var consultas []models.Consulta
	for rows.Next() {
		var c models.Consulta
		if err := rows.Scan(
			&c.ID, &c.Tipo, &c.Horario, &c.Diagnostico, &c.Costo,
			&c.IDConsultorio, &c.IDPaciente, &c.IDMedico,
		); err != nil {
			continue
		}
		consultas = append(consultas, c)
	}

	return c.JSON(consultas)
}

func UpdateConsulta(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "ID inválido"})
	}

	var consulta models.Consulta
	if err := c.BodyParser(&consulta); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Datos inválidos"})
	}

	_, err = database.DB.Exec(`
        UPDATE consulta 
        SET tipo = $1, horario = $2, diagnostico = $3, costo = $4, 
            id_consultorio = $5, id_paciente = $6, id_medico = $7 
        WHERE id_consulta = $8`,
		consulta.Tipo, consulta.Horario, consulta.Diagnostico, consulta.Costo,
		consulta.IDConsultorio, consulta.IDPaciente, consulta.IDMedico, id,
	)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error al actualizar consulta"})
	}

	consulta.ID = id
	return c.JSON(consulta)
}

func DeleteConsulta(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "ID inválido"})
	}

	_, err = database.DB.Exec("DELETE FROM consulta WHERE id_consulta = $1", id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error al eliminar consulta"})
	}

	return c.SendStatus(fiber.StatusNoContent)
}
