package handler

import (
	"context"
	"net/http"
)

func (h *Handler) auth(w http.ResponseWriter, r *http.Request, do func(ctx context.Context, login string)) {
	ctx, closeFn := h.Context(r.Context())
	defer closeFn()

	login, err := h.AuthService(ctx).Auth(getToken(r))
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	do(ctx, login)
}
