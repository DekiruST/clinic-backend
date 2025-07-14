// paciente.go
package handlers

import (
	"clinic-backend/database"
	"clinic-backend/models"
	"database/sql"
	"fmt"
	"strconv"

	"clinic-backend/utils"

	fiber "github.com/gofiber/fiber/v2"
)

func CreatePaciente(c *fiber.Ctx) error {
	var paciente models.Paciente
	if err := c.BodyParser(&paciente); err != nil {
		userID := c.Locals("user_id").(int)
		utils.LogOperacion(&userID, "create", "paciente", "JSON inválido", false)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Datos inválidos"})
	}

	err := database.DB.QueryRow(`
        INSERT INTO paciente (nombre, seguro, contacto) 
        VALUES ($1, $2, $3) 
        RETURNING id_paciente`,
		paciente.Nombre, paciente.Seguro, paciente.Contacto,
	).Scan(&paciente.ID)
	if err != nil {
		userID := c.Locals("user_id").(int)
		utils.LogOperacion(&userID, "create", "paciente", "Error al crear", false)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error al crear paciente"})
	}

	userID := c.Locals("user_id").(int)
	utils.LogOperacion(&userID, "create", "paciente", "Paciente creado correctamente", true)
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
		userID := c.Locals("user_id").(int)
		utils.LogOperacion(&userID, "update", "paciente", "ID inválido", false)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "ID inválido"})
	}

	var paciente models.Paciente
	if err := c.BodyParser(&paciente); err != nil {
		userID := c.Locals("user_id").(int)
		utils.LogOperacion(&userID, "update", "paciente", "JSON inválido", false)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Datos inválidos"})
	}

	_, err = database.DB.Exec(`
        UPDATE paciente 
        SET nombre = $1, seguro = $2, contacto = $3 
        WHERE id_paciente = $4`,
		paciente.Nombre, paciente.Seguro, paciente.Contacto, id,
	)
	if err != nil {
		userID := c.Locals("user_id").(int)
		utils.LogOperacion(&userID, "update", "paciente", "Error al actualizar", false)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error al actualizar paciente"})
	}

	userID := c.Locals("user_id").(int)
	utils.LogOperacion(&userID, "update", "paciente", "Paciente actualizado correctamente", true)

	paciente.ID = id
	return c.JSON(paciente)
}

func DeletePaciente(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		userID := c.Locals("user_id").(int)
		utils.LogOperacion(&userID, "delete", "paciente", "ID inválido", false)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "ID inválido"})
	}

	_, err = database.DB.Exec("DELETE FROM paciente WHERE id_paciente = $1", id)
	if err != nil {
		userID := c.Locals("user_id").(int)
		utils.LogOperacion(&userID, "delete", "paciente", "Error al eliminar", false)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error al eliminar paciente"})
	}

	userID := c.Locals("user_id").(int)
	utils.LogOperacion(&userID, "delete", "paciente", "Paciente eliminado correctamente", true)

	return c.SendStatus(fiber.StatusNoContent)
}

func GetPacienteByToken(c *fiber.Ctx) error {
	fmt.Println("DEBUG Handler: ¡Entrando a GetPacienteByToken!")
	v := c.Locals("id_paciente")
	fmt.Printf("DEBUG Handler - c.Locals(\"id_paciente\"): %#v (type: %T)\n", v, v)

	var idPaciente int
	if v != nil {
		switch val := v.(type) {
		case int:
			idPaciente = val
		case float64:
			idPaciente = int(val)
		case string:
			tmp, err := strconv.Atoi(val)
			if err != nil {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "ID inválido"})
			}
			idPaciente = tmp
		default:
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "ID inválido"})
		}
	} else {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "ID inválido"})
	}

	var paciente models.Paciente
	err := database.DB.QueryRow(`
        SELECT id_paciente, nombre, seguro, contacto
        FROM paciente
        WHERE id_paciente = $1
    `, idPaciente).Scan(&paciente.ID, &paciente.Nombre, &paciente.Seguro, &paciente.Contacto)

	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Paciente no encontrado"})
	}

	return c.JSON(paciente)
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
        UPDATE consulta SET tipo=$1, horario=$2, diagnostico=$3, costo=$4, id_consultorio=$5, id_paciente=$6, id_medico=$7 WHERE id_consulta=$8`,
		consulta.Tipo, consulta.Horario, consulta.Diagnostico, consulta.Costo, consulta.IDConsultorio, consulta.IDPaciente, consulta.IDMedico, id,
	)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error al actualizar consulta"})
	}
	consulta.ID = id
	return c.JSON(consulta)
}

func GetConsultasPaciente(c *fiber.Ctx) error {
	v := c.Locals("id_paciente")
	var idPaciente int

	switch val := v.(type) {
	case int:
		idPaciente = val
	case float64:
		idPaciente = int(val)
	case string:
		tmp, err := strconv.Atoi(val)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "ID inválido"})
		}
		idPaciente = tmp
	default:
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "ID inválido"})
	}

	rows, err := database.DB.Query(`
		SELECT 
			c.id_consulta,
			c.tipo,
			c.horario,
			COALESCE(c.diagnostico, ''),
			COALESCE(c.costo, 0),
			c.id_consultorio,
			c.id_paciente,
			c.id_medico,
			COALESCE(m.nombre, '') AS nombre_medico
		FROM consulta c
		LEFT JOIN medico m ON c.id_medico = m.id_medico
		WHERE c.id_paciente = $1
		ORDER BY c.horario DESC
	`, idPaciente)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error al obtener consultas"})
	}
	defer rows.Close()

	var consultas []models.Consulta
	for rows.Next() {
		var consulta models.Consulta
		var nombreMedico sql.NullString

		err := rows.Scan(
			&consulta.ID,
			&consulta.Tipo,
			&consulta.Horario,
			&consulta.Diagnostico,
			&consulta.Costo,
			&consulta.IDConsultorio,
			&consulta.IDPaciente,
			&consulta.IDMedico,
			&nombreMedico,
		)
		if err != nil {
			continue
		}
		if nombreMedico.Valid {
			consulta.NombreMedico = &nombreMedico.String
		}
		consultas = append(consultas, consulta)
	}

	return c.JSON(consultas)
}
