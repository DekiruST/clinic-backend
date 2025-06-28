package handlers

import (
	"clinic-backend/database"
	"clinic-backend/models"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func CreateUsuario(c *fiber.Ctx) error {
	var usuario models.Usuario
	if err := c.BodyParser(&usuario); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Datos inválidos"})
	}

	err := database.DB.QueryRow(`
        INSERT INTO usuario (rol, especialidad, id_paciente) 
        VALUES ($1, $2, $3) 
        RETURNING id_usuario`,
		usuario.Rol, usuario.Especialidad, usuario.IDPaciente,
	).Scan(&usuario.ID)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error al crear usuario"})
	}

	return c.Status(fiber.StatusCreated).JSON(usuario)
}

func GetUsuarios(c *fiber.Ctx) error {
	rows, err := database.DB.Query("SELECT * FROM usuario")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error en consulta"})
	}
	defer rows.Close()

	var usuarios []models.Usuario
	for rows.Next() {
		var u models.Usuario
		if err := rows.Scan(&u.ID, &u.Rol, &u.Especialidad, &u.IDPaciente); err != nil {
			continue
		}
		usuarios = append(usuarios, u)
	}

	return c.JSON(usuarios)
}

func GetUsuario(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "ID inválido"})
	}

	var usuario models.Usuario
	err = database.DB.QueryRow(`
        SELECT * FROM usuario WHERE id_usuario = $1`, id,
	).Scan(&usuario.ID, &usuario.Rol, &usuario.Especialidad, &usuario.IDPaciente)

	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Usuario no encontrado"})
	}

	return c.JSON(usuario)
}

func UpdateUsuario(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "ID inválido"})
	}

	var usuario models.Usuario
	if err := c.BodyParser(&usuario); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Datos inválidos"})
	}

	_, err = database.DB.Exec(`
        UPDATE usuario 
        SET rol = $1, especialidad = $2, id_paciente = $3 
        WHERE id_usuario = $4`,
		usuario.Rol, usuario.Especialidad, usuario.IDPaciente, id,
	)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error al actualizar usuario"})
	}

	usuario.ID = id
	return c.JSON(usuario)
}

func DeleteUsuario(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "ID inválido"})
	}

	_, err = database.DB.Exec("DELETE FROM usuario WHERE id_usuario = $1", id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error al eliminar usuario"})
	}

	return c.SendStatus(fiber.StatusNoContent)
}
