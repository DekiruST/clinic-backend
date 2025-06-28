package handlers

import (
	"clinic-backend/database"
	"clinic-backend/models"
	"strconv"

	fiber "github.com/gofiber/fiber/v2"
)

func CreateReceta(c *fiber.Ctx) error {
	var receta models.Receta
	if err := c.BodyParser(&receta); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Datos inválidos"})
	}

	err := database.DB.QueryRow(`
        INSERT INTO receta (fecha, id_consulta) 
        VALUES (CURRENT_DATE, $1) 
        RETURNING id_receta`,
		receta.IDConsulta,
	).Scan(&receta.ID)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error al crear receta"})
	}

	return c.Status(fiber.StatusCreated).JSON(receta)
}

func GetRecetasByConsulta(c *fiber.Ctx) error {
	idConsulta, err := strconv.Atoi(c.Params("id_consulta"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "ID inválido"})
	}

	var receta models.Receta
	err = database.DB.QueryRow(`
        SELECT * FROM receta WHERE id_consulta = $1`, idConsulta,
	).Scan(&receta.ID, &receta.Fecha, &receta.IDConsulta)

	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Receta no encontrada"})
	}

	return c.JSON(receta)
}
