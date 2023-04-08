package auth

import (
	"errors"

	"github.com/golang-jwt/jwt/v5"
)

type Service interface {
	GenerateToken(userUD int) (string, error)
	ValidateToken(token string) (*jwt.Token, error)
}

type jwtService struct{}

var SECRET_KEY = []byte("startup_secret_k3y")

func NewService() *jwtService {
	return &jwtService{}
}

func (s *jwtService) GenerateToken(userID int) (string, error) {
	claim := jwt.MapClaims{}
	claim["user_id"] = userID

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)

	signedToken, err := token.SignedString(SECRET_KEY)

	if err != nil {
		return signedToken, err
	}

	return signedToken, nil

}

func (s *jwtService) ValidateToken(encodedToken string) (*jwt.Token, error) {
	token, err := jwt.Parse(encodedToken, func(t *jwt.Token) (interface{}, error) {
		_, ok := t.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, errors.New("Invalid Token")
		}

		return []byte(SECRET_KEY), nil

	})

	if err != nil {
		return token, nil
	}

	return token, nil
}
