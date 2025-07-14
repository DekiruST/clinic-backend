package utils

import (
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte("supersecret")

func GenerateJWT(userID int, rol string, idPaciente *int, durationSeconds int64) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"rol":     rol,
		"exp":     time.Now().Add(time.Second * time.Duration(durationSeconds)).Unix(),
	}

	if idPaciente != nil {
		claims["id_paciente"] = *idPaciente
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}
