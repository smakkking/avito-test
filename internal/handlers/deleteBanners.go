package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
	"github.com/smakkking/avito_test/internal/handlers/utils"
	"github.com/smakkking/avito_test/internal/services"
)

func (h *Handler) DeleteBanner(w http.ResponseWriter, r *http.Request) {
	const op = "handler.DeleteBanner"

	data := chi.URLParam(r, "id")

	bannerID, err := strconv.Atoi(data)
	if err != nil {
		logrus.Error(fmt.Errorf("%s: %w", op, err))

		utils.JSONwithCode(w, r, http.StatusBadRequest, utils.ErrMessage("некорректные данные"))
		return
	}

	err = h.bannerService.DeleteBanner(r.Context(), bannerID)
	if err != nil {
		logrus.Error(fmt.Errorf("%s: %w", op, err))

		if errors.Is(err, services.ErrNotAllowed) {
			w.WriteHeader(http.StatusForbidden)
			return
		}

		if errors.Is(err, services.ErrNotFound) {
			// это значит, что не нашли строк, куда вставлять данные
			w.WriteHeader(http.StatusNotFound)
			return
		}

		utils.JSONwithCode(w, r, http.StatusInternalServerError, utils.ErrMessage("Внутренняя ошибка сервера"))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
