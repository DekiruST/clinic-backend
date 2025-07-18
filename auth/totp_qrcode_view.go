// auth/totp_qrcode_view.go
package auth

import (
	"clinic-backend/database"
	"fmt"
	"net/url"

	"github.com/gofiber/fiber/v2"
)

func ServeTOTPQRViewByID(c *fiber.Ctx) error {
	idUsuario := c.Params("id_usuario")
	if idUsuario == "" {
		return c.Status(400).SendString("Falta id_usuario")
	}

	var secret string
	var rol string
	err := database.DB.QueryRow(`SELECT totp_secret, rol FROM usuario WHERE id_usuario=$1`, idUsuario).Scan(&secret, &rol)
	if err != nil || secret == "" {
		return c.Status(404).SendString("Usuario no encontrado o sin TOTP configurado")
	}

	// Construir otpauth URL
	label := fmt.Sprintf("ClinicApp:user%s", idUsuario)
	otpauth := fmt.Sprintf("otpauth://totp/%s?secret=%s&issuer=%s", url.PathEscape(label), secret, url.QueryEscape("ClinicApp"))

	// URL del QR
	qrURL := fmt.Sprintf("https://api.qrserver.com/v1/create-qr-code/?data=%s&size=250x250", url.QueryEscape(otpauth))

	// HTML simple con la imagen
	html := fmt.Sprintf(`
		<!DOCTYPE html>
		<html>
		<head><meta charset="UTF-8"><title>QR TOTP</title></head>
		<body style="text-align:center;font-family:sans-serif">
			<h2>Escanea este QR en Google Authenticator</h2>
			<p><strong>Usuario ID:</strong> %s</p>
			<img src="%s" alt="QR Code"><br>
			<p><code>%s</code></p>
		</body>
		</html>
	`, idUsuario, qrURL, otpauth)

	return c.Type("html").SendString(html)
}
