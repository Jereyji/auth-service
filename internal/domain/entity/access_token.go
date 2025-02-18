package entity

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"time"
)

type TokenPayload struct {
	UserID      uuid.UUID `json:"user_id"`
	AccessLevel int       `json:"access_level"`
}

type TokenClaims struct {
	jwt.RegisteredClaims
	TokenPayload
}

type AccessToken struct {
	Token string
}

func NewAccessToken(userID uuid.UUID, accessLevel int, expiresIn time.Duration, secretKey string) (*AccessToken, error) {
	claims := TokenClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt: &jwt.NumericDate{
				Time: time.Now(),
			},
			ExpiresAt: &jwt.NumericDate{
				Time: time.Now().Add(expiresIn),
			},
		},
		TokenPayload: TokenPayload{
			UserID:      userID,
			AccessLevel: accessLevel,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	accessToken, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return nil, err
	}

	return &AccessToken{
		Token: accessToken,
	}, nil
}

func ValidateAccessToken(accessToken, secretKey string) (*TokenClaims, error) {
	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(secretKey), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(TokenClaims)
	if !ok {
		return nil, jwt.ErrTokenInvalidClaims
	}

	return &claims, nil
}

func (c TokenClaims) CheckAccessLevel(accessLevel int) error {
	if c.AccessLevel != accessLevel {
		return ErrInsufficientAccessLevel
	}

	return nil
}

