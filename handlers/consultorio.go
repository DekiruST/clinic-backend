package handlers

import (
	"clinic-backend/database"
	"clinic-backend/models"
	"clinic-backend/utils"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

// CreateConsultorio crea un nuevo consultorio
func CreateConsultorio(c *fiber.Ctx) error {
	var consultorio models.Consultorio
	if err := c.BodyParser(&consultorio); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Datos inválidos", err)
	}

	// Validar datos de entrada
	if consultorio.Nombre == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Nombre del consultorio requerido", nil)
	}
	if consultorio.Ubicacion == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Ubicación requerida", nil)
	}
	if consultorio.Especialidad == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Especialidad requerida", nil)
	}
	if consultorio.Categoria == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Categoría requerida", nil)
	}

	// Verificar si ya existe un consultorio con el mismo nombre
	var exists bool
	err := database.DB.QueryRow(`
		SELECT EXISTS(SELECT 1 FROM consultorio WHERE nombre = $1)
	`, consultorio.Nombre).Scan(&exists)

	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Error verificando consultorio", err)
	}
	if exists {
		return utils.ErrorResponse(c, fiber.StatusConflict, "Ya existe un consultorio con ese nombre", nil)
	}

	// Crear el consultorio
	err = database.DB.QueryRow(`
        INSERT INTO consultorio (nombre, ubicacion, especialidad, categoria) 
        VALUES ($1, $2, $3, $4) 
        RETURNING id_consultorio`,
		consultorio.Nombre, consultorio.Ubicacion, consultorio.Especialidad, consultorio.Categoria,
	).Scan(&consultorio.ID)

	if err != nil {
		log.Printf("Error al crear consultorio: %v", err)
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Error al crear consultorio", err)
	}

	return c.Status(fiber.StatusCreated).JSON(consultorio)
}

// GetConsultorios obtiene todos los consultorios con paginación y filtros
func GetConsultorios(c *fiber.Ctx) error {
	// Parámetros de paginación
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	offset := (page - 1) * limit

	// Parámetros de filtrado
	especialidad := c.Query("especialidad")
	categoria := c.Query("categoria")
	nombre := c.Query("nombre")

	// Construir consulta SQL con filtros
	query := "SELECT * FROM consultorio WHERE 1=1"
	args := []interface{}{}
	argCounter := 1

	if especialidad != "" {
		query += " AND especialidad = $" + strconv.Itoa(argCounter)
		args = append(args, especialidad)
		argCounter++
	}

	if categoria != "" {
		query += " AND categoria = $" + strconv.Itoa(argCounter)
		args = append(args, categoria)
		argCounter++
	}

	if nombre != "" {
		query += " AND nombre ILIKE $" + strconv.Itoa(argCounter)
		args = append(args, "%"+nombre+"%")
		argCounter++
	}

	// Ordenar por nombre
	query += " ORDER BY nombre ASC"

	// Paginación
	query += " LIMIT $" + strconv.Itoa(argCounter)
	args = append(args, limit)
	argCounter++

	query += " OFFSET $" + strconv.Itoa(argCounter)
	args = append(args, offset)

	// Ejecutar consulta
	rows, err := database.DB.Query(query, args...)
	if err != nil {
		log.Printf("Error en consulta: %v", err)
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Error al obtener consultorios", err)
	}
	defer rows.Close()

	var consultorios []models.Consultorio
	for rows.Next() {
		var c models.Consultorio
		if err := rows.Scan(&c.ID, &c.Nombre, &c.Ubicacion, &c.Especialidad, &c.Categoria); err != nil {
			log.Printf("Error escaneando fila: %v", err)
			continue
		}
		consultorios = append(consultorios, c)
	}

	// Obtener total de registros para paginación
	var total int
	countQuery := strings.Split(query, "ORDER BY")[0]
	countQuery = "SELECT COUNT(*) FROM (" + countQuery + ") AS subquery"
	err = database.DB.QueryRow(countQuery, args[:len(args)-2]...).Scan(&total)
	if err != nil {
		log.Printf("Error contando registros: %v", err)
	}

	return c.JSON(fiber.Map{
		"data":       consultorios,
		"pagination": fiber.Map{"page": page, "limit": limit, "total": total},
	})
}

// GetConsultorio obtiene un consultorio específico por ID
func GetConsultorio(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "ID de consultorio inválido", err)
	}

	var consultorio models.Consultorio
	err = database.DB.QueryRow(`
        SELECT * FROM consultorio WHERE id_consultorio = $1
    `, id).Scan(
		&consultorio.ID, &consultorio.Nombre, &consultorio.Ubicacion,
		&consultorio.Especialidad, &consultorio.Categoria,
	)

	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "Consultorio no encontrado", err)
	}

	return c.JSON(consultorio)
}

// UpdateConsultorio actualiza un consultorio existente
func UpdateConsultorio(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "ID de consultorio inválido", err)
	}

	var consultorio models.Consultorio
	if err := c.BodyParser(&consultorio); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Datos inválidos", err)
	}

	// Validar datos de entrada
	if consultorio.Nombre == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Nombre del consultorio requerido", nil)
	}
	if consultorio.Ubicacion == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Ubicación requerida", nil)
	}
	if consultorio.Especialidad == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Especialidad requerida", nil)
	}
	if consultorio.Categoria == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Categoría requerida", nil)
	}

	// Verificar si el consultorio existe
	var consultorioExists bool
	err = database.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM consultorio WHERE id_consultorio = $1)", id).Scan(&consultorioExists)
	if err != nil || !consultorioExists {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "Consultorio no encontrado", err)
	}

	// Verificar si el nuevo nombre ya existe en otro consultorio
	var nameExists bool
	err = database.DB.QueryRow(`
		SELECT EXISTS(SELECT 1 FROM consultorio WHERE nombre = $1 AND id_consultorio != $2)
	`, consultorio.Nombre, id).Scan(&nameExists)

	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Error verificando nombre", err)
	}
	if nameExists {
		return utils.ErrorResponse(c, fiber.StatusConflict, "Ya existe otro consultorio con ese nombre", nil)
	}

	// Actualizar consultorio
	_, err = database.DB.Exec(`
        UPDATE consultorio 
        SET nombre = $1, ubicacion = $2, especialidad = $3, categoria = $4 
        WHERE id_consultorio = $5`,
		consultorio.Nombre, consultorio.Ubicacion, consultorio.Especialidad,
		consultorio.Categoria, id,
	)

	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Error al actualizar consultorio", err)
	}

	consultorio.ID = id
	return c.JSON(consultorio)
}

// DeleteConsultorio elimina un consultorio
func DeleteConsultorio(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "ID de consultorio inválido", err)
	}

	// Verificar si el consultorio existe
	var consultorioExists bool
	err = database.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM consultorio WHERE id_consultorio = $1)", id).Scan(&consultorioExists)
	if err != nil || !consultorioExists {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "Consultorio no encontrado", err)
	}

	// Verificar si hay consultas programadas para este consultorio
	var hasConsultas bool
	err = database.DB.QueryRow(`
		SELECT EXISTS(SELECT 1 FROM consulta WHERE id_consultorio = $1 AND horario > NOW())
	`, id).Scan(&hasConsultas)

	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Error verificando consultas", err)
	}
	if hasConsultas {
		return utils.ErrorResponse(c, fiber.StatusConflict, "No se puede eliminar, hay consultas programadas", nil)
	}

	// Eliminar consultorio
	_, err = database.DB.Exec("DELETE FROM consultorio WHERE id_consultorio = $1", id)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Error al eliminar consultorio", err)
	}

	return c.SendStatus(fiber.StatusNoContent)
}

// GetConsultoriosDisponibles obtiene los consultorios disponibles en un horario específico
func GetConsultoriosDisponibles(c *fiber.Ctx) error {
	// Obtener parámetros de consulta
	fecha := c.Query("fecha")
	hora := c.Query("hora")
	especialidad := c.Query("especialidad")

	if fecha == "" || hora == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Fecha y hora requeridas", nil)
	}

	// Formatear fecha y hora
	horarioStr := fecha + " " + hora
	horario, err := time.Parse("2006-01-02 15:04", horarioStr)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Formato de fecha/hora inválido", err)
	}

	// Construir consulta SQL
	query := `
		SELECT c.* 
		FROM consultorio c
		WHERE NOT EXISTS (
			SELECT 1 
			FROM consulta con 
			WHERE con.id_consultorio = c.id_consultorio 
			AND con.horario = $1
		)
	`
	args := []interface{}{horario}
	argCounter := 2

	if especialidad != "" {
		query += " AND c.especialidad = $" + strconv.Itoa(argCounter)
		args = append(args, especialidad)
		argCounter++
	}

	// Ejecutar consulta
	rows, err := database.DB.Query(query, args...)
	if err != nil {
		log.Printf("Error en consulta: %v", err)
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Error al obtener consultorios disponibles", err)
	}
	defer rows.Close()

	var consultorios []models.Consultorio
	for rows.Next() {
		var c models.Consultorio
		if err := rows.Scan(&c.ID, &c.Nombre, &c.Ubicacion, &c.Especialidad, &c.Categoria); err != nil {
			log.Printf("Error escaneando fila: %v", err)
			continue
		}
		consultorios = append(consultorios, c)
	}

	return c.JSON(consultorios)
}

// GetConsultoriosPorEspecialidad obtiene consultorios por especialidad
func GetConsultoriosPorEspecialidad(c *fiber.Ctx) error {
	especialidad := c.Params("especialidad")
	if especialidad == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Especialidad requerida", nil)
	}

	rows, err := database.DB.Query(`
        SELECT * FROM consultorio 
        WHERE especialidad = $1
        ORDER BY nombre ASC
    `, especialidad)

	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Error al obtener consultorios", err)
	}
	defer rows.Close()

	var consultorios []models.Consultorio
	for rows.Next() {
		var c models.Consultorio
		if err := rows.Scan(&c.ID, &c.Nombre, &c.Ubicacion, &c.Especialidad, &c.Categoria); err != nil {
			log.Printf("Error escaneando fila: %v", err)
			continue
		}
		consultorios = append(consultorios, c)
	}

	return c.JSON(consultorios)
}
