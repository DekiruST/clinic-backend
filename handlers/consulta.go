package handlers

import (
	"clinic-backend/database"
	"clinic-backend/models"
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func CreateConsulta(c *fiber.Ctx) error {
	var consulta models.Consulta
	if err := c.BodyParser(&consulta); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Datos inválidos"})
	}
	err := database.DB.QueryRow(`
    INSERT INTO consulta (tipo, horario, diagnostico, costo, id_consultorio, id_paciente)
    VALUES ($1,$2,$3,$4,$5,$6)
    RETURNING id_consulta`,
		consulta.Tipo, consulta.Horario, consulta.Diagnostico, consulta.Costo, consulta.IDConsultorio, consulta.IDPaciente,
	).Scan(&consulta.ID)

	if err != nil {
		fmt.Printf("DEBUG ERROR AL INSERTAR CONSULTA: %v\n", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error al crear consulta"})
	}
	return c.Status(fiber.StatusCreated).JSON(consulta)
}

func GetConsultas(c *fiber.Ctx) error {
	rows, err := database.DB.Query("SELECT id_consulta, tipo, horario, diagnostico, costo, id_consultorio, id_paciente, id_medico FROM consulta")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error en consulta"})
	}
	defer rows.Close()
	var consultas []models.Consulta
	for rows.Next() {
		var consulta models.Consulta
		if err := rows.Scan(&consulta.ID, &consulta.Tipo, &consulta.Horario, &consulta.Diagnostico, &consulta.Costo, &consulta.IDConsultorio, &consulta.IDPaciente, &consulta.IDMedico); err != nil {
			continue
		}
		consultas = append(consultas, consulta)
	}
	return c.JSON(consultas)
}

func GetConsultasByPaciente(c *fiber.Ctx) error {
	fmt.Printf("💡 ID PACIENTE del token: %v\n", c.Locals("id_paciente"))
	fmt.Println("🚨 Handler GetConsultasByPaciente ACTIVADO")
	idRaw := c.Locals("id_paciente")
	if idRaw == nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "id_paciente no disponible en el token",
		})
	}

	idPaciente, ok := idRaw.(int)
	if !ok {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "id_paciente en formato incorrecto",
		})
	}

	rows, err := database.DB.Query(`
		SELECT id_consulta, tipo, horario, diagnostico, costo, id_consultorio, id_paciente, id_medico 
		FROM consulta WHERE id_paciente = $1
	`, idPaciente)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Error al obtener las consultas",
			"error":   err.Error(),
		})
	}
	defer rows.Close()

	consultas := make([]models.Consulta, 0)

	for rows.Next() {
		var c models.Consulta
		if err := rows.Scan(
			&c.ID, &c.Tipo, &c.Horario, &c.Diagnostico, &c.Costo,
			&c.IDConsultorio, &c.IDPaciente, &c.IDMedico); err == nil {
			consultas = append(consultas, c)
		}
	}

	return c.Status(fiber.StatusOK).JSON(consultas)
}

func DeleteConsulta(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "ID inválido"})
	}
	_, err = database.DB.Exec("DELETE FROM consulta WHERE id_consulta=$1", id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error al eliminar consulta"})
	}
	return c.SendStatus(fiber.StatusNoContent)
}

func GetConsultasByPacienteIDParam(c *fiber.Ctx) error {
	idStr := c.Params("id_paciente")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "ID inválido",
		})
	}

	rows, err := database.DB.Query(`
		SELECT id_consulta, tipo, horario, diagnostico, costo, id_consultorio, id_paciente, id_medico
		FROM consulta
		WHERE id_paciente = $1
	`, id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error en la consulta",
			"error":   err.Error(),
		})
	}
	defer rows.Close()

	var consultas []models.Consulta
	for rows.Next() {
		var c models.Consulta
		if err := rows.Scan(&c.ID, &c.Tipo, &c.Horario, &c.Diagnostico, &c.Costo, &c.IDConsultorio, &c.IDPaciente, &c.IDMedico); err != nil {
			continue
		}
		consultas = append(consultas, c)
	}

	return c.JSON(consultas)
}
