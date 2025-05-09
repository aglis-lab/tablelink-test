package utils

import (
	"errors"
	"tablelink/src/entity"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type accessToken struct {
	jwt.RegisteredClaims
	IsAccess bool             `json:"is_access"`
	Data     entity.UserToken `json:"data"`
}

// HashPassword generates a bcrypt hash for the given password.
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// VerifyPassword verifies if the given password matches the stored hash.
func VerifyPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func SignAccessToken(data entity.UserToken) (string, error) {
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS512, &accessToken{
		Data:     data,
		IsAccess: true,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt: &jwt.NumericDate{
				Time: time.Now(),
			},
			ExpiresAt: &jwt.NumericDate{
				Time: time.Now().AddDate(1, 0, 0),
			},
		},
	})

	token, err := accessToken.SignedString([]byte("accesskey"))
	if err != nil {
		return "", err
	}

	return token, nil
}

func ParseAccessToken(tokenStr string) (*entity.UserToken, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &accessToken{}, func(t *jwt.Token) (interface{}, error) {
		return []byte("accesskey"), nil
	})
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("token is not valid")
	}

	claims, ok := token.Claims.(*accessToken)
	if !ok {
		return nil, errors.New("miss configure token data")
	}

	exp, err := claims.GetExpirationTime()
	if err != nil {
		return nil, err
	}

	if exp.Unix() < time.Now().Unix() {
		return nil, errors.New("access token already expired")
	}

	return &claims.Data, nil
}

func DistinctSlice[T comparable](inputSlice []T) []T {
	uniqueSlice := make([]T, 0, len(inputSlice))
	seen := make(map[T]bool, len(inputSlice))
	for _, element := range inputSlice {
		if _, ok := seen[element]; !ok {
			uniqueSlice = append(uniqueSlice, element)
		}
	}

	return uniqueSlice
}

func SliceToMap[T any, E comparable](data []T, callback func(T) E) map[E]T {
	mapping := map[E]T{}
	for _, val := range data {
		mapping[callback(val)] = val
	}

	return mapping
}
