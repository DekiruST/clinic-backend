package auth

import (
	"clinic-backend/database"
	"clinic-backend/models"
	"clinic-backend/utils"
	"database/sql"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var jwtSecret = []byte("supersecret")

// --------------------- REGISTRO ---------------------

func Register(c *fiber.Ctx) error {
	var input struct {
		Email        string  `json:"email"`
		Rol          string  `json:"rol"`
		Especialidad *string `json:"especialidad"`
		Password     string  `json:"password"`
		Nombre       string  `json:"nombre"`
		Contacto     string  `json:"contacto"`
		Seguro       *string `json:"seguro"`
	}

	if err := c.BodyParser(&input); err != nil {
		utils.LogOperacion(nil, "register", "usuario", "JSON inválido", false)
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Datos inválidos", err)
	}

	if input.Email == "" || !utils.IsValidEmail(input.Email) {
		utils.LogOperacion(nil, "register", "usuario", "Email inválido", false)
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Correo electrónico inválido", nil)
	}

	if !utils.ValidPassword(input.Password) {
		utils.LogOperacion(nil, "register", "usuario", "Contraseña débil", false)
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Contraseña débil", nil)
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		utils.LogOperacion(nil, "register", "usuario", "Error al encriptar", false)
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Error al encriptar", err)
	}

	var idPaciente *int = nil

	if input.Rol == "paciente" {
		var id int
		err := database.DB.QueryRow(`
			INSERT INTO paciente (nombre, contacto, seguro)
			VALUES ($1, $2, $3)
			RETURNING id_paciente
		`, input.Nombre, input.Contacto, input.Seguro).Scan(&id)
		if err != nil {
			utils.LogOperacion(nil, "register", "paciente", "Error al crear paciente", false)
			return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Error al registrar paciente", err)
		}
		idPaciente = &id
	}

	var userID int
	err = database.DB.QueryRow(`
	INSERT INTO usuario (email, rol, especialidad, id_paciente, password_hash, nombre)
	VALUES ($1, $2, $3, $4, $5, $6)
	RETURNING id_usuario
`, input.Email, input.Rol, input.Especialidad, idPaciente, string(hashed), input.Nombre).Scan(&userID)
	if err != nil {
		utils.LogOperacion(nil, "register", "usuario", "Error en BD", false)
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Error al registrar usuario", err)
	}

	utils.LogOperacion(&userID, "register", "usuario", "Registro exitoso", true)
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"id_usuario": userID})
}

// --------------------- LOGIN ---------------------

func Login(c *fiber.Ctx) error {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		TOTPCode string `json:"totp_code"`
	}

	if err := c.BodyParser(&input); err != nil {
		utils.LogOperacion(nil, "login", "usuario", "JSON inválido", false)
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Datos inválidos", err)
	}

	var user models.Usuario
	var hash string
	var totpSecret sql.NullString
	var idPaciente sql.NullInt64

	err := database.DB.QueryRow(`
		SELECT id_usuario, rol, especialidad, id_paciente, password_hash, totp_secret, email
		FROM usuario WHERE email = $1`, input.Email,
	).Scan(
		&user.ID,
		&user.Rol,
		&user.Especialidad,
		&idPaciente,
		&hash,
		&totpSecret,
		&user.Email,
	)
	if err == sql.ErrNoRows {
		utils.LogOperacion(nil, "login", "usuario", "Usuario no encontrado", false)
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Credenciales inválidas", nil)
	} else if err != nil {
		utils.LogOperacion(nil, "login", "usuario", "Error en BD", false)
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Error en BD", err)
	}
	if idPaciente.Valid {
		tmp := int(idPaciente.Int64)
		user.IDPaciente = &tmp
	} else {
		user.IDPaciente = nil
	}

	if bcrypt.CompareHashAndPassword([]byte(hash), []byte(input.Password)) != nil {
		utils.LogOperacion(&user.ID, "login", "usuario", "Contraseña incorrecta", false)
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Contraseña incorrecta", nil)
	}

	if totpSecret.Valid {
		if input.TOTPCode == "" {
			utils.LogOperacion(&user.ID, "login", "usuario", "TOTP requerido pero no proporcionado", false)
			return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Este usuario requiere un código TOTP", nil)
		}
		if !utils.ValidateTOTP(totpSecret.String, input.TOTPCode) {
			utils.LogOperacion(&user.ID, "login", "usuario", "TOTP inválido", false)
			return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Código TOTP inválido", nil)
		}
	}

	token, err := generateJWT(user)
	if err != nil {
		utils.LogOperacion(&user.ID, "login", "usuario", "Error creando token", false)
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Error creando token", err)
	}
	refreshToken, err := generateRefreshToken(user.ID)
	if err != nil {
		utils.LogOperacion(&user.ID, "login", "usuario", "Error creando refresh", false)
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Error creando refresh token", err)
	}
	utils.LogOperacion(&user.ID, "login", "usuario", "Login exitoso", true)

	return c.JSON(fiber.Map{
		"access_token":  token,
		"refresh_token": refreshToken,
		"user": fiber.Map{
			"id":           user.ID,
			"rol":          user.Rol,
			"email":        user.Email,
			"especialidad": user.Especialidad,
			"id_paciente":  user.IDPaciente,
		},
	})
}

func generateRefreshToken(i int) (any, error) {
	var token string
	expiresAt := time.Now().Add(time.Hour * 24 * 7) // 7 días de expiración
	err := database.DB.QueryRow(`
		INSERT INTO refresh_tokens (usuario_id, token, expires_at)
		VALUES ($1, gen_random_uuid(), $2)
		RETURNING token`, i, expiresAt).Scan(&token)
	if err != nil {
		return nil, err
	}
	return fiber.Map{"token": token, "expires_at": expiresAt}, nil
}

// --------------------- REFRESH TOKEN ---------------------

func Refresh(c *fiber.Ctx) error {
	var input struct {
		RefreshToken string `json:"refresh_token"`
	}
	if err := c.BodyParser(&input); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Token inválido", err)
	}

	var userID int
	err := database.DB.QueryRow(`
		SELECT usuario_id FROM refresh_tokens WHERE token = $1 AND expires_at > NOW()`,
		input.RefreshToken,
	).Scan(&userID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Token inválido o expirado", err)
	}

	var user models.Usuario
	var idPaciente sql.NullInt64
	err = database.DB.QueryRow(`
		SELECT id_usuario, rol, especialidad, id_paciente, email
		FROM usuario WHERE id_usuario = $1`,
		userID,
	).Scan(&user.ID, &user.Rol, &user.Especialidad, &idPaciente, &user.Email)
	if idPaciente.Valid {
		tmp := int(idPaciente.Int64)
		user.IDPaciente = &tmp
	} else {
		user.IDPaciente = nil
	}

	token, err := generateJWT(user)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Error creando nuevo token", err)
	}

	return c.JSON(fiber.Map{"access_token": token})
}

// --------------------- JWT Y REFRESH TOKEN HELPERS ---------------------

func generateJWT(user models.Usuario) (string, error) {
	permisos := []string{}
	rows, err := database.DB.Query(`
		SELECT p.nombre
		FROM permisos p
		INNER JOIN rol_permisos rp ON rp.permiso_id = p.id
		WHERE rp.rol = $1
	`, user.Rol)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var permiso string
			if err := rows.Scan(&permiso); err == nil {
				permisos = append(permisos, permiso)
			}
		}
	}

	claims := jwt.MapClaims{
		"user_id":  user.ID,
		"rol":      user.Rol,
		"permisos": permisos,
		"exp":      time.Now().Add(time.Minute * 15).Unix(),
	}

	if user.IDPaciente != nil {
		claims["id_paciente"] = *user.IDPaciente
	}

	if user.Rol == "medico" {
		var idMedico int
		err := database.DB.QueryRow(`
			SELECT id_medico FROM medico WHERE id_usuario = $1
		`, user.ID).Scan(&idMedico)
		if err == nil && idMedico != 0 {
			claims["id_medico"] = idMedico
		}
	}

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return t.SignedString(jwtSecret)
}
