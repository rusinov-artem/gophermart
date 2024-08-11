package handler

import (
	"context"
	"time"

	"github.com/rusinov-artem/gophermart/app/dto"
	appError "github.com/rusinov-artem/gophermart/app/error"
)

type RegisterAction interface {
	Validate(params dto.RegisterParams) *appError.ValidationError
	Register(params dto.RegisterParams) (string, *appError.InternalError)
}

type Handler struct {
	RegisterAction func(ctx context.Context) RegisterAction
}

func New() *Handler {
	return &Handler{}
}

func (h *Handler) Context(ctx context.Context) (context.Context, func()) {
	return context.WithTimeout(ctx, 5*time.Second)
}
