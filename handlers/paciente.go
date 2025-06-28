package handlers

import (
	"clinic-backend/database"
	"clinic-backend/models"
	"log"
	"strconv"

	fiber "github.com/gofiber/fiber/v2"
)

func CreatePaciente(c *fiber.Ctx) error {
	var paciente models.Paciente
	if err := c.BodyParser(&paciente); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Datos inválidos"})
	}

	err := database.DB.QueryRow(`
        INSERT INTO paciente (nombre, seguro, contacto) 
        VALUES ($1, $2, $3) 
        RETURNING id_paciente`,
		paciente.Nombre, paciente.Seguro, paciente.Contacto,
	).Scan(&paciente.ID)

	if err != nil {
		log.Println("Error en BD:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error al crear paciente"})
	}

	return c.Status(fiber.StatusCreated).JSON(paciente)
}

func GetPacientes(c *fiber.Ctx) error {
	rows, err := database.DB.Query("SELECT * FROM paciente")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error en consulta"})
	}
	defer rows.Close()

	var pacientes []models.Paciente
	for rows.Next() {
		var p models.Paciente
		if err := rows.Scan(&p.ID, &p.Nombre, &p.Seguro, &p.Contacto); err != nil {
			continue
		}
		pacientes = append(pacientes, p)
	}

	return c.JSON(pacientes)
}

func GetPaciente(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "ID inválido"})
	}

	var paciente models.Paciente
	err = database.DB.QueryRow(`
        SELECT * FROM paciente WHERE id_paciente = $1`, id,
	).Scan(&paciente.ID, &paciente.Nombre, &paciente.Seguro, &paciente.Contacto)

	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Paciente no encontrado"})
	}

	return c.JSON(paciente)
}

func UpdatePaciente(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "ID inválido"})
	}

	var paciente models.Paciente
	if err := c.BodyParser(&paciente); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Datos inválidos"})
	}

	_, err = database.DB.Exec(`
        UPDATE paciente 
        SET nombre = $1, seguro = $2, contacto = $3 
        WHERE id_paciente = $4`,
		paciente.Nombre, paciente.Seguro, paciente.Contacto, id,
	)

	if err != nil {
		log.Println("Error en BD:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error al actualizar paciente"})
	}

	paciente.ID = id
	return c.JSON(paciente)
}

func DeletePaciente(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "ID inválido"})
	}

	_, err = database.DB.Exec("DELETE FROM paciente WHERE id_paciente = $1", id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error al eliminar paciente"})
	}

	return c.SendStatus(fiber.StatusNoContent)
}
