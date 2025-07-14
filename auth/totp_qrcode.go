package auth

import (
	"clinic-backend/database"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/skip2/go-qrcode"
)

func GetTOTPQRCode(c *fiber.Ctx) error {
	userIDRaw := c.Locals("user_id")
	if userIDRaw == nil {
		return c.Status(fiber.StatusUnauthorized).SendString("Usuario no autenticado")
	}
	userID := userIDRaw.(int)

	var secret string
	err := database.DB.QueryRow(
		"SELECT totp_secret FROM usuario WHERE id_usuario = $1",
		userID,
	).Scan(&secret)
	if err != nil {
		return c.Status(500).SendString("Error obteniendo el secreto: " + err.Error())
	}
	if secret == "" {
		return c.Status(400).SendString("Este usuario no tiene TOTP habilitado")
	}

	otpURL := fmt.Sprintf(
		"otpauth://totp/ClinicApp:user%d@clinic?secret=%s&issuer=ClinicApp",
		userID,
		secret,
	)

	png, err := qrcode.Encode(otpURL, qrcode.Medium, 256)
	if err != nil {
		return c.Status(500).SendString("Error generando QR: " + err.Error())
	}

	return c.Type("image/png").Send(png)
}
