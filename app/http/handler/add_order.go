package handler

import "net/http"

func (h *Handler) AddOrder(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
}
