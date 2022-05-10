package services

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/arnab333/golang-employee-management/helpers"
	"github.com/golang-jwt/jwt/v4"
)

type TokenDetails struct {
	AccessToken  string
	RefreshToken string
	AccessUUID   string
	RefreshUuid  string
	AtExpires    time.Time
	RtExpires    time.Time
}

type TokenClaims struct {
	Role   string
	UserID string
	jwt.RegisteredClaims
}

func createToken(userID string, role string) (*TokenDetails, error) {
	td := TokenDetails{
		AtExpires:  time.Now().Add(time.Minute * 15),
		AccessUUID: "ACCESS-" + helpers.GetUUID(),
		// RtExpires:   time.Now().Add(time.Hour * 24 * 7), // production
		RtExpires:   time.Now().Add(time.Hour * 8),
		RefreshUuid: "REFRESH-" + helpers.GetUUID(),
	}

	var err error

	at := jwt.NewWithClaims(jwt.SigningMethodHS256, TokenClaims{
		Role:   role,
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        td.AccessUUID,
			ExpiresAt: jwt.NewNumericDate(td.AtExpires),
		},
	})

	td.AccessToken, err = at.SignedString([]byte(os.Getenv(helpers.EnvKeys.JWT_ACCESS_SECRET)))
	if err != nil {
		return nil, err
	}

	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, TokenClaims{
		Role:   role,
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        td.RefreshUuid,
			ExpiresAt: jwt.NewNumericDate(td.RtExpires),
		},
	})

	td.RefreshToken, err = rt.SignedString([]byte(os.Getenv(helpers.EnvKeys.JWT_REFRESH_SECRET)))
	if err != nil {
		return nil, err
	}

	return &td, nil
}

func ExtractFromToken(tokenString string, envKey string) (*TokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv(envKey)), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*TokenClaims)
	if !ok && !token.Valid {
		return nil, errors.New("invalid token")
	}
	return claims, nil
}

func CreateAuth(ctx context.Context, userID string, role string) (*TokenDetails, error) {
	td, err := createToken(userID, role)
	if err != nil {
		return nil, err
	}

	now := time.Now()

	err = redisConn.redisClient.Set(ctx, td.AccessUUID, userID, td.AtExpires.Sub(now)).Err()
	if err != nil {
		return nil, err
	}

	err = redisConn.redisClient.Set(ctx, td.RefreshUuid, userID, td.RtExpires.Sub(now)).Err()
	if err != nil {
		return nil, err
	}

	return td, nil
}

func FetchAuth(ctx context.Context, accessUUID string) (string, error) {
	userID, err := redisConn.redisClient.Get(ctx, accessUUID).Result()
	if err != nil {
		return "", err
	}

	return userID, nil
}

func DeleteAuth(ctx context.Context, key string) (int64, error) {
	deleted, err := redisConn.redisClient.Del(ctx, key).Result()
	if err != nil {
		return 0, err
	}
	return deleted, nil
}
