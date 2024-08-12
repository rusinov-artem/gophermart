package handler

import (
	"context"
	"time"
)

type Handler struct {
	RegisterAction func(ctx context.Context) RegisterAction
	LoginAction    func(ctx context.Context) LoginAction
	AddOrderAction func(ctx context.Context) AddOrderAction

	AuthService      func(ctx context.Context) AuthService
	ListOrdersAction func(ctx context.Context) ListOrdersAction
}

func New() *Handler {
	return &Handler{}
}

func (h *Handler) Context(ctx context.Context) (context.Context, func()) {
	return context.WithTimeout(ctx, 5*time.Second)
}
