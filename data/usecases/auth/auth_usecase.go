package authusecase

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/mazurco066/playliter-api-go/domain/models/account"
)

type AuthUseCase interface {
	IssueToken(u account.Account) (string, error)
	ParseToken(token string) (*Claims, error)
}

type authUseCase struct {
	jwtSecret string
}

func NewAuthUseCase(jwtSecret string) AuthUseCase {
	return &authUseCase{
		jwtSecret: jwtSecret,
	}
}

type Claims struct {
	Account string `json:"account"`
	Role    string `json:"role"`
	jwt.RegisteredClaims
}

func (auth *authUseCase) IssueToken(a account.Account) (string, error) {
	nowTime := time.Now()
	expireTime := nowTime.Add(2160 * time.Hour) // 3 Months

	claims := Claims{
		a.Email,
		a.Role,
		jwt.RegisteredClaims{
			Issuer:    "Playliter API",
			ExpiresAt: jwt.NewNumericDate(expireTime),
		},
	}

	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return tokenClaims.SignedString([]byte(auth.jwtSecret))
}

func (auth *authUseCase) ParseToken(token string) (*Claims, error) {
	tokenClaims, err := jwt.ParseWithClaims(
		token,
		&Claims{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(auth.jwtSecret), nil
		},
	)

	if tokenClaims != nil {
		claims, ok := tokenClaims.Claims.(*Claims)
		if ok && tokenClaims.Valid {
			return claims, nil
		}
	}

	return nil, err
}
