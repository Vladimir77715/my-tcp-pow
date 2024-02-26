package receive_solution

import (
	"context"
	"encoding/json"
	"time"

	"github.com/Vladimir77715/my-tcp-pow/internal/errror"
	"github.com/Vladimir77715/my-tcp-pow/internal/models"
	"github.com/Vladimir77715/my-tcp-pow/internal/tcp/server"
	"github.com/dgrijalva/jwt-go"
)

type Handler struct{}

func New() *Handler {
	return &Handler{}
}

func (h *Handler) Handle(ctx context.Context, _ *models.RawData) (*models.RawData, error) {
	const methodName = "receive_solution.Handle"

	t, err := h.generateToken(ctx)
	if err != nil {
		return nil, errror.WrapError(err, methodName)
	}

	at := models.AccessToken{Token: t}
	b, err := json.Marshal(&at)
	if err != nil {
		return nil, errror.WrapError(err, methodName)
	}

	return &models.RawData{Command: server.SendSolution, Payload: b}, nil
}

func (h *Handler) generateToken(ctx context.Context) (string, error) {
	const methodName = "generateToken"
	token := jwt.New(jwt.SigningMethodHS256)

	expire := time.Now().Add(time.Minute * 5)
	// Создание токена с Claims (утверждениями)
	claims := token.Claims.(jwt.MapClaims)
	claims["admin"] = true
	claims["exp"] = expire.Unix()

	val, ok := ctx.Value(server.SecretKeyVal).([]byte)
	if !ok {
		return "", errror.WrapError(server.ErrSecretKeyNotFound, methodName)
	}

	// Подписание токена нашим секретным ключом
	tokenString, err := token.SignedString(val)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
