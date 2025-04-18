package entity

import (
	"time"

	auth_errors "github.com/Jereyji/auth-service/internal/auth/domain/errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type TokenPayload struct {
	UserID uuid.UUID `json:"user_id"`
}

type TokenClaims struct {
	jwt.RegisteredClaims
	TokenPayload
}

type AccessToken struct {
	Token string
}

func NewAccessToken(userID uuid.UUID, expiresIn time.Duration, secretKey string) (AccessToken, error) {
	curTime := time.Now()
	claims := TokenClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt: &jwt.NumericDate{
				Time: curTime,
			},
			ExpiresAt: &jwt.NumericDate{
				Time: curTime.Add(expiresIn),
			},
		},
		TokenPayload: TokenPayload{
			UserID: userID,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	accessToken, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return AccessToken{}, err
	}

	return AccessToken{
		Token: accessToken,
	}, nil
}

func ValidateAccessToken(accessToken, secretKey string) (TokenClaims, error) {
	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, auth_errors.ErrInvalidSigningMethod
		}
		return []byte(secretKey), nil
	})
	if err != nil {
		return TokenClaims{}, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return TokenClaims{}, jwt.ErrTokenInvalidClaims
	}

	userIDStr, ok := claims["user_id"].(string)
	if !ok {
		return TokenClaims{}, jwt.ErrTokenInvalidClaims
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return TokenClaims{}, err
	}

	issuedAt, ok := claims["iat"].(float64)
	if !ok {
		return TokenClaims{}, jwt.ErrTokenInvalidClaims
	}

	expiresAt, ok := claims["exp"].(float64)
	if !ok {
		return TokenClaims{}, jwt.ErrTokenInvalidClaims
	}

	return TokenClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  &jwt.NumericDate{Time: time.Unix(int64(issuedAt), 0)},
			ExpiresAt: &jwt.NumericDate{Time: time.Unix(int64(expiresAt), 0)},
		},
		TokenPayload: TokenPayload{
			UserID: userID,
		},
	}, nil
}
