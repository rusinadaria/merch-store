package handlers

import (
	"merch-store/models"
	"merch-store/internal/common"
	"encoding/json"
	"net/http"
	"github.com/go-chi/chi"
)

func (h *Handler) SendHandler (w http.ResponseWriter, r *http.Request) { // Отправить монетки другому сотруднику.
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	var req models.SendCoinRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		common.WriteErrorResponse(w, http.StatusBadRequest, "Неверный запрос")
		return
	}

	cookie, err := r.Cookie("auth_token")
	if err != nil {
		common.WriteErrorResponse(w, http.StatusUnauthorized, "Ошибка при попытке получить cookies")
		return
	}
	
	id, err := h.services.ParseToken(cookie.Value)
	if err != nil {
		common.WriteErrorResponse(w, http.StatusBadRequest, "Ошибка при попытке распарсить cookies")
		return
	}

	err = h.services.SendCoin(id, req.ToUser, req.Amount)
	if err != nil {
		common.WriteErrorResponse(w, http.StatusInternalServerError, "Ошибка при попытке отправить коины")
		return 
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) BuyItemHandler (w http.ResponseWriter, r *http.Request) { // Покупка товаров
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	name := chi.URLParam(r, "item")

	cookie, err := r.Cookie("auth_token")
	if err != nil {
		common.WriteErrorResponse(w, http.StatusBadRequest, "Ошибка при попытке получить cookies")
		return
	}

	id, err := h.services.ParseToken(cookie.Value)
	if err != nil {
		common.WriteErrorResponse(w, http.StatusBadRequest, "Ошибка при попытке распарсить cookies")
		return
	}

	err = h.services.BuyItem(id, name)
	if err != nil {
		common.WriteErrorResponse(w, http.StatusInternalServerError, "Не возможно приобрести товар")
		return
	}

	w.WriteHeader(http.StatusOK)
}