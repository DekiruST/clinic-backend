// auth/middleware.go
package auth

import (
	"clinic-backend/database"
	"clinic-backend/utils"
	"fmt"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	jwt "github.com/golang-jwt/jwt/v5"
)

func ValidateJWT(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Token no proporcionado", nil)
	}

	tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("método de firma inesperado")
		}
		return jwtSecret, nil
	})
	if err != nil || !token.Valid {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Token inválido", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Claims inválidos", nil)
	}

	c.Locals("user_id", int(claims["user_id"].(float64)))
	c.Locals("rol", claims["rol"].(string))
	if v, ok := claims["id_paciente"]; ok && v != nil {
		switch val := v.(type) {
		case float64:
			c.Locals("id_paciente", int(val))
		case int:
			c.Locals("id_paciente", val)
		case string:
			tmp, err := strconv.Atoi(val)
			if err == nil {
				c.Locals("id_paciente", tmp)
			} else {
				fmt.Printf("DEBUG JWT: id_paciente como string pero no convertible: %v\n", val)
			}
		default:
			fmt.Printf("DEBUG JWT: id_paciente en claims tiene tipo inesperado: %T, valor: %#v\n", val, val)
		}
		fmt.Printf("DEBUG JWT: asignando c.Locals(\"id_paciente\") = %v (tipo: %T)\n", c.Locals("id_paciente"), c.Locals("id_paciente"))
	} else {
		fmt.Println("DEBUG JWT: id_paciente no está presente o es nil en claims")
	}

	return c.Next()

}

func RequirePermission(permiso string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		rol, ok := c.Locals("rol").(string)
		if !ok {
			fmt.Println("DEBUG RequirePermission: Rol no encontrado en token")
			return utils.ErrorResponse(c, fiber.StatusForbidden, "Rol no encontrado en token", nil)
		}

		// Consultar los permisos del rol
		rows, err := database.DB.Query(`
			SELECT p.nombre
			FROM permisos p
			INNER JOIN rol_permisos rp ON p.id = rp.permiso_id
			WHERE rp.rol = $1
		`, rol)
		if err != nil {
			fmt.Println("DEBUG RequirePermission: Error consultando permisos")
			return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Error consultando permisos", err)
		}
		defer rows.Close()

		allowed := false
		for rows.Next() {
			var nombre string
			if err := rows.Scan(&nombre); err != nil {
				continue
			}
			if nombre == permiso || nombre == "*" {
				allowed = true
				break
			}
		}

		if !allowed {
			fmt.Printf("DEBUG RequirePermission: Acceso DENEGADO a rol %s para permiso %s\n", rol, permiso)
			return utils.ErrorResponse(c, fiber.StatusForbidden, fmt.Sprintf("No tienes permiso: %s", permiso), nil)
		}

		fmt.Printf("DEBUG RequirePermission: Acceso PERMITIDO a rol %s para permiso %s\n", rol, permiso)
		return c.Next()
	}
}

func OnlyPaciente(c *fiber.Ctx) error {
	rol, ok := c.Locals("rol").(string)
	if !ok || rol != "paciente" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"success": false,
			"message": "Solo pacientes pueden acceder a este recurso",
		})
	}
	return c.Next()
}
