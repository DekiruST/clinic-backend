// auth/totp.go
package auth

import (
	"clinic-backend/database"
	"clinic-backend/utils"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/pquerna/otp/totp"
)

func GenerateTOTPSecret(c *fiber.Ctx) error {
	userIDRaw := c.Locals("user_id")
	userID, ok := userIDRaw.(int)
	if !ok {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Usuario no autenticado", nil)
	}

	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "ClinicApp",
		AccountName: fmt.Sprintf("user%d@clinic", userID),
	})
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Error generando TOTP", err)
	}

	_, err = database.DB.Exec("UPDATE usuario SET totp_secret = $1 WHERE id_usuario = $2", key.Secret(), userID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Error guardando secreto", err)
	}

	return c.JSON(fiber.Map{
		"secret":  key.Secret(),
		"otp_url": key.URL(),
	})
}

func ValidateTOTPCode(userID int, code string) bool {
	var secret string
	err := database.DB.QueryRow("SELECT totp_secret FROM usuario WHERE id_usuario = $1", userID).Scan(&secret)
	if err != nil || secret == "" {
		return false
	}

	valid := totp.Validate(code, secret)
	return valid
}

func VerifyTOTP(c *fiber.Ctx) error {
	var input struct {
		Code string `json:"code"`
	}
	if err := c.BodyParser(&input); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Datos inválidos", err)
	}

	userIDRaw := c.Locals("user_id")
	userID, ok := userIDRaw.(int)
	if !ok {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Usuario no autenticado", nil)
	}

	if !ValidateTOTPCode(userID, input.Code) {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Código TOTP inválido", nil)
	}

	var rol string
	var idPaciente *int
	err := database.DB.QueryRow(`
		SELECT rol, id_paciente 
		FROM usuario 
		WHERE id_usuario = $1`, userID).Scan(&rol, &idPaciente)
	fmt.Printf("Usuario %d tiene rol '%s'\n", userID, rol)

	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Error cargando datos de usuario", err)
	}

	tokenStr, err := utils.GenerateJWT(userID, rol, idPaciente, 3600)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Error generando token", err)
	}
	fmt.Println("TOKEN ENVIADO TRAS TOTP:", tokenStr)

	return c.JSON(fiber.Map{
		"success":      true,
		"access_token": tokenStr,
	})
}

func ResetTOTPSecret(c *fiber.Ctx) error {
	userIDRaw := c.Locals("user_id")
	userID, ok := userIDRaw.(int)
	if !ok {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Usuario no autenticado", nil)
	}

	_, err := database.DB.Exec("UPDATE usuario SET totp_secret = NULL WHERE id_usuario = $1", userID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Error al reiniciar MFA", err)
	}

	return c.JSON(fiber.Map{"success": true, "message": "Secreto MFA reiniciado"})
}
