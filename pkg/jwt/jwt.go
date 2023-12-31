package jwt

import (
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
)

const JWT_SECRET = "wizh"

func GenerateToken(id int64, username string) (string, error) {
	expiresTimes := time.Now().Add(time.Hour * 24).Unix()
	claims := jwt.StandardClaims{
		Audience:  username,
		ExpiresAt: expiresTimes,
		Id:        strconv.FormatInt(id, 10),
		IssuedAt:  time.Now().Unix(),
		Issuer:    "wizh",
		NotBefore: time.Now().Unix(),
		Subject:   "token",
	}

	var jwtSercet = []byte(JWT_SECRET)
	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := tokenClaims.SignedString(jwtSercet)
	if err != nil {
		return "", err
	}
	return token, nil
}

func ParseToken(token string) (*jwt.StandardClaims, error) {
	jwtToken, err := jwt.ParseWithClaims(token, &jwt.StandardClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(JWT_SECRET), nil
	})
	if err == nil && jwtToken != nil {
		if cliam, ok := jwtToken.Claims.(*jwt.StandardClaims); ok && jwtToken.Valid {
			return cliam, nil
		}
	}
	return nil, err
}
