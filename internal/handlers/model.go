package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
	"github.com/smakkking/avito_test/internal/handlers/utils"
	"github.com/smakkking/avito_test/internal/models"
	"github.com/smakkking/avito_test/internal/services"
)

type Handler struct {
	bannerService *services.Service
}

func NewHandler(service *services.Service) *Handler {
	return &Handler{
		bannerService: service,
	}
}

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

func (h *Handler) UpdateBanner(w http.ResponseWriter, r *http.Request) {
	const op = "handler.UpdateBanner"

	data := chi.URLParam(r, "id")

	bannerID, err := strconv.Atoi(data)
	if err != nil {
		logrus.Error(fmt.Errorf("%s: %w", op, err))

		utils.JSONwithCode(w, r, http.StatusBadRequest, utils.ErrMessage("некорректные данные"))
		return
	}

	banner := new(models.BasicBannnerInfo)
	err = json.NewDecoder(r.Body).Decode(banner)
	if err != nil {
		logrus.Error(fmt.Errorf("%s: %w", op, err))

		utils.JSONwithCode(w, r, http.StatusBadRequest, utils.ErrMessage("некорректные данные"))
		return
	}

	err = h.bannerService.UpdateBanner(r.Context(), bannerID, banner)
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

	w.WriteHeader(http.StatusOK)
}

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
