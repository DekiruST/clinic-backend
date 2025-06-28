package utils

import (
	"log"

	"github.com/gofiber/fiber/v2"
)

// ErrorResponse devuelve una respuesta de error estandarizada
func ErrorResponse(c *fiber.Ctx, status int, message string, err error) error {
	// Registrar el error para depuración
	if err != nil {
		log.Printf("Error: %s - Detalle: %v", message, err)
	}

	// Construir respuesta
	response := fiber.Map{
		"success": false,
		"message": message,
	}

	// Incluir detalles de error solo en desarrollo
	if fiber.IsChild() {
		response["error"] = err.Error()
	}

	return c.Status(status).JSON(response)
}
