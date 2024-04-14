package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/sirupsen/logrus"
	"github.com/smakkking/avito_test/internal/handlers/utils"
	"github.com/smakkking/avito_test/internal/services"
)

func (h *Handler) GetUserBanner(w http.ResponseWriter, r *http.Request) {
	const op = "handler.GetUserBanner"

	tagID, err := strconv.Atoi(r.URL.Query().Get("tag_id"))
	if err != nil {
		logrus.Error(fmt.Errorf("%s: %w", op, err))

		utils.JSONwithCode(w, r, http.StatusBadRequest, utils.ErrMessage("некорректные данные"))
		return
	}

	featureID, err := strconv.Atoi(r.URL.Query().Get("feature_id"))
	if err != nil {
		logrus.Error(fmt.Errorf("%s: %w", op, err))

		utils.JSONwithCode(w, r, http.StatusBadRequest, utils.ErrMessage("некорректные данные"))
		return
	}

	var useLastRevision bool
	data := r.URL.Query().Get("use_last_revision")
	if data != "" {
		useLastRevision, err = strconv.ParseBool(data)
		if err != nil {
			logrus.Error(fmt.Errorf("%s: %w", op, err))

			utils.JSONwithCode(w, r, http.StatusBadRequest, utils.ErrMessage("некорректные данные"))
			return
		}
	} else {
		useLastRevision = false
	}

	banner, err := h.bannerService.GetUserBanner(r.Context(), tagID, featureID, useLastRevision)
	if err != nil {
		logrus.Error(fmt.Errorf("%s: %w", op, err))

		if errors.Is(err, services.ErrNotFound) {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		if errors.Is(err, services.ErrNotAllowed) {
			w.WriteHeader(http.StatusForbidden)
			return
		}

		utils.JSONwithCode(w, r, http.StatusInternalServerError, utils.ErrMessage("Внутренняя ошибка сервера"))
		return
	}

	logrus.Debug(banner)

	utils.JSONwithCode(w, r, http.StatusOK, banner)
}
