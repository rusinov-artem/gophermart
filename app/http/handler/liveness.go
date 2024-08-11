package handler

import "net/http"

func (h *Handler) Liveness(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
