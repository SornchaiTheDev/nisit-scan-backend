package libs

import (
	"os"
	"time"

	"github.com/SornchaiTheDev/nisit-scan-backend/domain/nerrors"
	"github.com/golang-jwt/jwt/v5"
)

func GenerateAccessToken(email string, name string, picture string, role string) (*string, *time.Time, error) {
	accessExp := time.Now().Add(time.Hour * 1)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email":   email,
		"name":    name,
		"picture": picture,
		"role":    role,
		"exp":     accessExp.Unix(),
	})

	signedToken, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return nil, nil, nerrors.ErrSomethingWentWrong
	}

	return &signedToken, &accessExp, nil
}

func GenerateRefreshToken(email string) (*string, *time.Time, error) {
	refreshExp := time.Now().Add(time.Hour * 24 * 10)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": email,
		"exp":   refreshExp.Unix(),
	})

	signedToken, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return nil, nil, nerrors.ErrSomethingWentWrong
	}

	return &signedToken, &refreshExp, nil
}

func ParseJwt(tok string) (map[string]interface{}, error) {

	claims := jwt.MapClaims{}

	token, err := jwt.ParseWithClaims(tok, &claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, nerrors.ErrTokenNotValid
	}

	return claims, nil
}
