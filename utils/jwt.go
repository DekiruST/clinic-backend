// utils/jwt.go
package utils

import (
	"clinic-backend/database"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte("supersecret")

func GenerateJWT(userID int, rol string, idPaciente *int, durationSecs int64) (string, error) {
	rows, err := database.DB.Query(`
	SELECT p.nombre 
	FROM permisos p
	JOIN rol_permisos rp ON rp.permiso_id = p.id
	WHERE rp.rol = $1
`, rol)
	if err != nil {
		return "", err
	}
	defer rows.Close()

	var permisos []string
	for rows.Next() {
		var permiso string
		if err := rows.Scan(&permiso); err == nil {
			permisos = append(permisos, permiso)
		}
	}

	claims := jwt.MapClaims{
		"user_id":  userID,
		"rol":      rol,
		"exp":      time.Now().Add(time.Second * time.Duration(durationSecs)).Unix(),
		"permisos": permisos,
	}

	if idPaciente != nil {
		claims["id_paciente"] = *idPaciente
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}
