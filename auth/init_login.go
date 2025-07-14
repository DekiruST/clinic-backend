// auth/init_login.go
package auth

import (
	"clinic-backend/database"
	"clinic-backend/utils"
	"database/sql"
	"strings"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

func InitLogin(c *fiber.Ctx) error {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.BodyParser(&input); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Datos inválidos", err)
	}

	input.Email = strings.TrimSpace(input.Email)
	if input.Email == "" || input.Password == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Correo o contraseña faltante", nil)
	}

	var userID int
	var passwordHash string
	var rol string
	var totpSecret sql.NullString

	err := database.DB.QueryRow(`
		SELECT id_usuario, password_hash, rol, totp_secret
		FROM usuario
		WHERE email = $1
	`, input.Email).Scan(&userID, &passwordHash, &rol, &totpSecret)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Credenciales inválidas", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(input.Password)); err != nil {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Contraseña incorrecta", nil)
	}

	if totpSecret.Valid {
		return utils.ErrorResponse(c, fiber.StatusForbidden, "MFA ya configurado. Usa /auth/login con TOTP", nil)
	}

	tokenStr, err := utils.GenerateJWT(userID, rol, nil, 5*60)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Error generando token", err)
	}

	return c.JSON(fiber.Map{
		"access_token": tokenStr,
		"mfa_required": true,
		"message":      "MFA no configurado. Continúa a configurar TOTP.",
	})
}
