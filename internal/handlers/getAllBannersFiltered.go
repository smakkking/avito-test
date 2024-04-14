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

func (h *Handler) GetAllBannersFiltered(w http.ResponseWriter, r *http.Request) {
	const op = "handler.GetAllBannersFiltered"

	var data string
	var err error

	var tagID int
	var tagSearch bool
	data = r.URL.Query().Get("tag_id")
	if data != "" {
		tagID, err = strconv.Atoi(r.URL.Query().Get("tag_id"))
		if err != nil {
			logrus.Error(fmt.Errorf("%s: %w", op, err))

			utils.JSONwithCode(w, r, http.StatusBadRequest, utils.ErrMessage("некорректные данные"))
			return
		}
		tagSearch = true
	}

	var featureID int
	var featureSearch bool
	data = r.URL.Query().Get("feature_id")
	if data != "" {
		featureID, err = strconv.Atoi(r.URL.Query().Get("feature_id"))
		if err != nil {
			logrus.Error(fmt.Errorf("%s: %w", op, err))

			utils.JSONwithCode(w, r, http.StatusBadRequest, utils.ErrMessage("некорректные данные"))
			return
		}
		featureSearch = true
	}

	limit := -1
	data = r.URL.Query().Get("limit")
	if data != "" {
		limit, err = strconv.Atoi(r.URL.Query().Get("limit"))
		if err != nil || limit < 0 {
			logrus.Error(fmt.Errorf("%s: %w", op, err))

			utils.JSONwithCode(w, r, http.StatusBadRequest, utils.ErrMessage("некорректные данные"))
			return
		}
	}

	offset := -1
	data = r.URL.Query().Get("offset")
	if data != "" {
		offset, err = strconv.Atoi(r.URL.Query().Get("offset"))
		if err != nil || offset < 0 {
			logrus.Error(fmt.Errorf("%s: %w", op, err))

			utils.JSONwithCode(w, r, http.StatusBadRequest, utils.ErrMessage("некорректные данные"))
			return
		}
	}

	banners, err := h.bannerService.GetAllBannersFiltered(
		r.Context(),
		tagID, tagSearch,
		featureID, featureSearch,
		limit, offset,
	)
	if err != nil {
		logrus.Error(fmt.Errorf("%s: %w", op, err))

		if errors.Is(err, services.ErrNotAllowed) {
			w.WriteHeader(http.StatusForbidden)
			return
		}

		utils.JSONwithCode(w, r, http.StatusInternalServerError, utils.ErrMessage("Внутренняя ошибка сервера"))
		return
	}

	logrus.Debug(banners)

	utils.JSONwithCode(w, r, http.StatusOK, banners)
}
