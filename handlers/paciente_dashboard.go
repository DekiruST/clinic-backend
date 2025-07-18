// handlers/paciente_dashboard.go
package handlers

import (
	"clinic-backend/database"
	"clinic-backend/models"

	"github.com/gofiber/fiber/v2"
)

// Próximas citas del paciente
// Próximas citas del paciente
func GetMisCitas(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int)
	rows, err := database.DB.Query(`
		SELECT id_consulta, tipo, horario, diagnostico, costo, id_consultorio, id_paciente, id_medico
		FROM consulta
		WHERE id_paciente = $1
		ORDER BY horario ASC`, userID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Error al obtener citas"})
	}
	defer rows.Close()

	var consultas []models.Consulta
	for rows.Next() {
		var con models.Consulta
		err := rows.Scan(&con.ID, &con.Tipo, &con.Horario, &con.Diagnostico, &con.Costo, &con.IDConsultorio, &con.IDPaciente, &con.IDMedico)
		if err == nil {
			consultas = append(consultas, con)
		}
	}
	return c.JSON(consultas)
}

// Histórico de citas
func GetHistoricoCitas(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int)
	rows, err := database.DB.Query(`
		SELECT id_consulta, tipo, horario, diagnostico, costo, id_consultorio, id_paciente, id_medico
		FROM consulta
		WHERE id_paciente = $1 AND horario < NOW()
		ORDER BY horario DESC`, userID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Error al obtener histórico"})
	}
	defer rows.Close()

	var consultas []models.Consulta
	for rows.Next() {
		var con models.Consulta
		err := rows.Scan(&con.ID, &con.Tipo, &con.Horario, &con.Diagnostico, &con.Costo, &con.IDConsultorio, &con.IDPaciente, &con.IDMedico)
		if err == nil {
			consultas = append(consultas, con)
		}
	}
	return c.JSON(consultas)
}

// Recetas del paciente (esto sí tiene fecha en la tabla receta, está bien)
// Handlers
func GetMisRecetas(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int)
	rows, err := database.DB.Query(`
		SELECT r.id_receta, r.fecha, r.id_consulta
		FROM receta r
		INNER JOIN consulta c ON r.id_consulta = c.id_consulta
		WHERE c.id_paciente = $1
		ORDER BY r.fecha DESC`, userID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Error al obtener recetas"})
	}
	defer rows.Close()

	recetas := make([]models.Receta, 0)
	for rows.Next() {
		var receta models.Receta
		err := rows.Scan(&receta.ID, &receta.Fecha, &receta.IDConsulta)
		if err == nil {
			recetas = append(recetas, receta)
		}
	}
	return c.JSON(recetas)
}

type SolicitarCitaInput struct {
	Horario       string `json:"horario"`
	Tipo          string `json:"tipo"`
	IDConsultorio int    `json:"id_consultorio"`
}

func SolicitarCita(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int)
	var input SolicitarCitaInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Datos inválidos"})
	}

	// Puedes validar la fecha y los campos aquí
	_, err := database.DB.Exec(`
	INSERT INTO consulta (tipo, horario, costo, id_consultorio, id_paciente)
	VALUES ($1, $2, $3, $4, $5)
`, input.Tipo, input.Horario, 0.0, input.IDConsultorio, userID)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "No se pudo agendar cita"})
	}
	return c.Status(201).JSON(fiber.Map{"success": true, "message": "Cita solicitada"})
}
