package services

import (
	"errors"
	"fmt"
	"os"
	"strings"
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

func CreateToken(userID string) (*TokenDetails, error) {
	td := TokenDetails{
		AtExpires:   time.Now().Add(time.Minute * 15),
		AccessUUID:  helpers.GetUUID(),
		RtExpires:   time.Now().Add(time.Hour * 24 * 7),
		RefreshUuid: helpers.GetUUID(),
	}

	var err error

	at := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		ID:        td.AccessUUID,
		Issuer:    userID,
		ExpiresAt: jwt.NewNumericDate(td.AtExpires),
	})

	td.AccessToken, err = at.SignedString([]byte(os.Getenv("JWT_ACCESS_SECRET")))
	if err != nil {
		return nil, err
	}

	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		ID:        td.RefreshUuid,
		Issuer:    userID,
		ExpiresAt: jwt.NewNumericDate(td.RtExpires),
	})

	td.RefreshToken, err = rt.SignedString([]byte(os.Getenv("JWT_REFRESH_SECRET")))
	if err != nil {
		return nil, err
	}

	return &td, nil
}

func CreateAuth(userID string, td *TokenDetails) error {
	now := time.Now()

	err := redisConn.redisClient.Set(redisConn.ctx, td.AccessUUID, userID, td.AtExpires.Sub(now)).Err()
	if err != nil {
		return err
	}

	err = redisConn.redisClient.Set(redisConn.ctx, td.RefreshUuid, userID, td.RtExpires.Sub(now)).Err()
	if err != nil {
		return err
	}

	// AuthDetails.AccessUUID = td.AccessUUID
	// AuthDetails.UserID = userID

	return nil
}

func FetchAuth(accessUUID string) (string, error) {
	userID, err := redisConn.redisClient.Get(redisConn.ctx, accessUUID).Result()
	if err != nil {
		return "", err
	}

	return userID, nil
}

func DeleteAuth(givenUuid string) (int64, error) {
	deleted, err := redisConn.redisClient.Del(redisConn.ctx, givenUuid).Result()
	if err != nil {
		return 0, err
	}
	return deleted, nil
}

func ExtractFromToken(bearToken string) (jwt.MapClaims, error) {
	strArr := strings.Split(bearToken, " ")
	var tokenString string
	if len(strArr) == 2 {
		tokenString = strArr[1]
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("JWT_ACCESS_SECRET")), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok && !token.Valid {
		return nil, errors.New("invalid token")
	}
	return claims, nil
}
