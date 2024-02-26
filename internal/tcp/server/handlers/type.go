package handlers

import (
	"context"

	"github.com/Vladimir77715/my-tcp-pow/internal/models"
)

type Handler interface {
	Handle(ctx context.Context, in *models.RawData) (*models.RawData, error)
}
