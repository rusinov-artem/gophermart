package handler

import (
	"context"
	"time"
)

type Handler struct {
	RegisterAction func(ctx context.Context) RegisterAction
	LoginAction    func(ctx context.Context) LoginAction
}

func New() *Handler {
	return &Handler{}
}

func (h *Handler) Context(ctx context.Context) (context.Context, func()) {
	return context.WithTimeout(ctx, 5*time.Second)
}
