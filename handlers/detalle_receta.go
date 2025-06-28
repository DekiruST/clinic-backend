package handlers

import (
	"clinic-backend/database"
	"clinic-backend/models"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func CreateDetalleReceta(c *fiber.Ctx) error {
	var detalle models.DetalleReceta
	if err := c.BodyParser(&detalle); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Datos inválidos"})
	}

	err := database.DB.QueryRow(`
        INSERT INTO detalle_receta (medicamento, dosis, id_receta) 
        VALUES ($1, $2, $3) 
        RETURNING id_detalle`,
		detalle.Medicamento, detalle.Dosis, detalle.IDReceta,
	).Scan(&detalle.ID)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error al crear detalle"})
	}

	return c.Status(fiber.StatusCreated).JSON(detalle)
}

func GetDetallesByReceta(c *fiber.Ctx) error {
	idReceta, err := strconv.Atoi(c.Params("id_receta"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "ID inválido"})
	}

	rows, err := database.DB.Query(`
        SELECT * FROM detalle_receta WHERE id_receta = $1`, idReceta)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error en consulta"})
	}
	defer rows.Close()

	var detalles []models.DetalleReceta
	for rows.Next() {
		var d models.DetalleReceta
		if err := rows.Scan(&d.ID, &d.Medicamento, &d.Dosis, &d.IDReceta); err != nil {
			continue
		}
		detalles = append(detalles, d)
	}

	return c.JSON(detalles)
}
