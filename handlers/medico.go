// handlers/medico.go

package handlers

import (
	"clinic-backend/database"
	"clinic-backend/models"
	"database/sql"
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func GetRecetasMedico(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(int)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "No autorizado"})
	}
	rows, err := database.DB.Query(`
		SELECT r.id_receta, r.fecha, r.id_consulta
		FROM receta r
		INNER JOIN consulta c ON r.id_consulta = c.id_consulta
		WHERE c.id_medico = $1
		ORDER BY r.fecha DESC
	`, userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error al obtener recetas"})
	}
	defer rows.Close()
	var recetas []models.Receta
	for rows.Next() {
		var receta models.Receta
		if err := rows.Scan(&receta.ID, &receta.Fecha, &receta.IDConsulta); err != nil {
			continue
		}
		recetas = append(recetas, receta)
	}
	return c.JSON(recetas)
}

func GetExpedientesMedico(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(int)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "No autorizado"})
	}
	rows, err := database.DB.Query(`
		SELECT e.id_expediente, e.antecedentes, e.historial_clinico, e.id_paciente
		FROM expediente e
		INNER JOIN consulta c ON e.id_paciente = c.id_paciente
		WHERE c.id_medico = $1
		GROUP BY e.id_expediente, e.antecedentes, e.historial_clinico, e.id_paciente
	`, userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error al obtener expedientes"})
	}
	defer rows.Close()
	var expedientes []models.Expediente
	for rows.Next() {
		var exp models.Expediente
		if err := rows.Scan(&exp.ID, &exp.Antecedentes, &exp.HistorialClinico, &exp.IDPaciente); err != nil {
			continue
		}
		expedientes = append(expedientes, exp)
	}
	return c.JSON(expedientes)
}

type ConsultaMedico struct {
	IDConsulta    int    `json:"id_consulta"`
	IDPaciente    int    `json:"id_paciente"`
	IDConsultorio int    `json:"id_consultorio"`
	Paciente      string `json:"paciente"`
	Motivo        string `json:"motivo"`
	Fecha         string `json:"fecha"`
	Tipo          string `json:"tipo"`
}

func GetConsultasMedico(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(int)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "No autorizado"})
	}

	rows, err := database.DB.Query(`
	SELECT 
		co.id_consulta,
		co.id_paciente,
		co.id_consultorio,
		p.nombre AS paciente,
		co.diagnostico AS motivo,
		co.horario AS fecha,
		co.tipo
	FROM consulta co
	JOIN paciente p ON co.id_paciente = p.id_paciente
	WHERE co.id_medico = $1
	ORDER BY co.horario DESC
`, userID)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error al obtener consultas"})
	}
	defer rows.Close()

	consultas := []ConsultaMedico{}
	for rows.Next() {
		var cMed ConsultaMedico
		var motivo sql.NullString
		err := rows.Scan(
			&cMed.IDConsulta,
			&cMed.IDPaciente,
			&cMed.IDConsultorio,
			&cMed.Paciente,
			&motivo,
			&cMed.Fecha,
			&cMed.Tipo,
		)
		if err != nil {
			continue
		}
		if motivo.Valid {
			cMed.Motivo = motivo.String
		} else {
			cMed.Motivo = ""
		}
		consultas = append(consultas, cMed)
	}

	return c.JSON(consultas)
}
func UpdateConsultaMedico(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "ID inválido"})
	}

	var consulta models.Consulta
	if err := c.BodyParser(&consulta); err != nil {
		fmt.Println("DEBUG BodyParser Error:", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Datos inválidos"})
	}

	fmt.Printf("DEBUG BODY (consulta): %+v\n", consulta)

	_, err = database.DB.Exec(`
        UPDATE consulta 
        SET tipo=$1, horario=$2, diagnostico=$3, costo=$4, id_consultorio=$5, id_paciente=$6, id_medico=$7 
        WHERE id_consulta=$8`,
		consulta.Tipo,
		consulta.Horario,
		consulta.Diagnostico,
		consulta.Costo,
		consulta.IDConsultorio,
		consulta.IDPaciente,
		consulta.IDMedico, // 👈 ahora sí, sin NullInt64
		id,
	)
	if err != nil {
		fmt.Println("DEBUG DB Exec Error:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error al actualizar consulta"})
	}

	consulta.ID = id
	return c.JSON(consulta)
}
