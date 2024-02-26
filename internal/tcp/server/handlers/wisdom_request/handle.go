package wisdom_request

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/dgrijalva/jwt-go"

	"github.com/Vladimir77715/my-tcp-pow/internal/errror"
	"github.com/Vladimir77715/my-tcp-pow/internal/models"
	"github.com/Vladimir77715/my-tcp-pow/internal/tcp/server"
)

var quotes = []string{
	"We are not what we know but what we are willing to learn.",

	"Good people are good because they've come to wisdom through failure.",

	"The first problem for all of us, men and women, is not to learn, but to unlearn.",

	"Just as treasures are uncovered from the earth, so virtue appears " +
		"from good deeds, and wisdom appears from a pure and peaceful mind. " +
		"To walk safely through the maze of human life, " +
		"one needs the light of wisdom and the guidance of virtue.",
}

var ErrTokenExpire = errors.New("token expired")

type Handler struct{}

func New() *Handler {
	return &Handler{}
}

func (h *Handler) Handle(ctx context.Context, data *models.RawData) (*models.RawData, error) {
	const methodName = "wisdom_request.Handle"

	at := models.AccessToken{}

	if err := json.Unmarshal(data.Payload, &at); err != nil {
		return nil, errror.WrapError(ErrTokenExpire, methodName)
	}

	_, err := validateToken(ctx, at.Token)
	if err != nil {
		return nil, errror.WrapError(ErrTokenExpire, methodName)
	}

	return &models.RawData{Command: server.WisdomRequest, Payload: []byte(quotes[rand.Intn(4)])}, nil
}

func validateToken(ctx context.Context, tokenString string) (*jwt.Token, error) {
	// Парсинг и валидация токена
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		const methodName = "jwt.Parse"

		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errror.WrapError(fmt.Errorf("неожиданный метод подписи: %v",
				token.Header["alg"]), methodName)
		}

		val, ok := ctx.Value(server.SecretKeyVal).([]byte)
		if !ok {
			return nil, errror.WrapError(server.ErrSecretKeyNotFound, methodName)
		}
		return val, nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		fmt.Println("Токен валиден")
		fmt.Println("Admin:", claims["admin"])
		fmt.Println("Expires:", time.Unix(int64(claims["exp"].(float64)), 0))
	} else {
		fmt.Println("Ошибка валидации токена:", err)
	}

	return token, err
}
