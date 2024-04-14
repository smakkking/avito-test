package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/sirupsen/logrus"
	"github.com/smakkking/avito_test/internal/handlers/utils"
	"github.com/smakkking/avito_test/internal/models"
	"github.com/smakkking/avito_test/internal/services"
)

func (h *Handler) CreateBanner(w http.ResponseWriter, r *http.Request) {
	const op = "handler.CreateBanner"

	banner := new(models.BasicBannnerInfo)
	err := json.NewDecoder(r.Body).Decode(banner)
	if err != nil {
		logrus.Error(fmt.Errorf("%s: %w", op, err))

		utils.JSONwithCode(w, r, http.StatusBadRequest, utils.ErrMessage("некорректные данные"))
		return
	}

	bannerID, err := h.bannerService.CreateBanner(r.Context(), banner)
	if err != nil {
		logrus.Error(fmt.Errorf("%s: %w", op, err))

		if errors.Is(err, services.ErrNotAllowed) {
			w.WriteHeader(http.StatusForbidden)
			return
		}

		utils.JSONwithCode(w, r, http.StatusInternalServerError, utils.ErrMessage("Внутренняя ошибка сервера"))
		return
	}

	logrus.Debug(bannerID)

	utils.JSONwithCode(w, r, http.StatusCreated, struct {
		BannedID int `json:"banner_id"`
	}{BannedID: bannerID})
}
