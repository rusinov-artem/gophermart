package handler

import "net/http"

func (h *Handler) GetBalance(w http.ResponseWriter, r *http.Request) {
	_, cancelFN := h.Context(r.Context())
	defer cancelFN()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}
