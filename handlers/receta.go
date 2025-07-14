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

func GetRecetasByPaciente(c *fiber.Ctx) error {
	idPaciente, err := strconv.Atoi(c.Params("id_paciente"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "ID inválido"})
	}

	// Consulta todas las recetas del paciente
	rows, err := database.DB.Query(`
        SELECT r.id_receta, r.fecha, r.id_consulta
        FROM receta r
        INNER JOIN consulta c ON r.id_consulta = c.id_consulta
        WHERE c.id_paciente = $1
        ORDER BY r.fecha DESC
    `, idPaciente)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error en consulta de recetas"})
	}
	defer rows.Close()

	type DetalleReceta struct {
		ID          int    `json:"id_detalle"`
		Medicamento string `json:"medicamento"`
		Dosis       string `json:"dosis"`
	}
	type RecetaConDetalles struct {
		IDReceta   int             `json:"id_receta"`
		Fecha      string          `json:"fecha"`
		IDConsulta int             `json:"id_consulta"`
		Detalles   []DetalleReceta `json:"detalles"`
	}

	var recetas []RecetaConDetalles
	for rows.Next() {
		var receta RecetaConDetalles
		if err := rows.Scan(&receta.IDReceta, &receta.Fecha, &receta.IDConsulta); err != nil {
			continue
		}
		// Consulta los detalles para cada receta
		detallesRows, err := database.DB.Query(`
            SELECT id_detalle, medicamento, dosis
            FROM detalle_receta
            WHERE id_receta = $1
        `, receta.IDReceta)
		if err == nil {
			for detallesRows.Next() {
				var det DetalleReceta
				detallesRows.Scan(&det.ID, &det.Medicamento, &det.Dosis)
				receta.Detalles = append(receta.Detalles, det)
			}
			detallesRows.Close()
		}
		recetas = append(recetas, receta)
	}
	return c.JSON(recetas)
}
