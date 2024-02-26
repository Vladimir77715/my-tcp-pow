package send_task

import (
	"context"
	"encoding/json"
	"math/rand"

	"github.com/Vladimir77715/my-tcp-pow/internal/errror"
	"github.com/Vladimir77715/my-tcp-pow/internal/models"
	"github.com/Vladimir77715/my-tcp-pow/internal/tcp/server"

	"github.com/google/uuid"
)

type Handler struct {
	maxN int
	minN int
}

func New(minN, maxN int) *Handler {
	return &Handler{maxN: maxN, minN: minN}
}

func (h *Handler) Handle(_ context.Context, _ *models.RawData) (*models.RawData, error) {
	const methodName = "send_task.Handle"

	n := rand.Intn(h.maxN-h.minN+1) + h.minN
	b, err := json.Marshal(models.HashData{ZeroCount: n, Data: uuid.New().String()})
	if err != nil {
		return nil, errror.WrapError(err, methodName)
	}

	return &models.RawData{Command: server.RequestChallenge, Payload: b}, nil
}
