package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/render"
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
	tagID, err := strconv.Atoi(r.URL.Query().Get("tag_id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		render.JSON(w, r, ErrMessage("некорректные данные"))
		return
	}

	featureID, err := strconv.Atoi(r.URL.Query().Get("feature_id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		render.JSON(w, r, ErrMessage("некорректные данные"))
		return
	}

	var useLastRevision bool
	data := r.URL.Query().Get("use_last_revision")
	if data != "" {
		useLastRevision, err = strconv.ParseBool(data)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, ErrMessage("некорректные данные"))
			return
		}
	} else {
		useLastRevision = false
	}

	banner, err := h.bannerService.GetUserBanner(r.Context(), tagID, featureID, useLastRevision)
	if err != nil {
		if errors.Is(err, services.ErrNotFound) {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		render.JSON(w, r, ErrMessage("Внутренняя ошибка сервера"))
		return
	}

	w.WriteHeader(http.StatusOK)
	render.JSON(w, r, banner)
}

func (h *Handler) GetAllBannersFiltered(w http.ResponseWriter, r *http.Request) {
	var data string
	var err error

	var tagID int
	var tagSearch bool
	data = r.URL.Query().Get("tag_id")
	if data != "" {
		tagID, err = strconv.Atoi(r.URL.Query().Get("tag_id"))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, ErrMessage("некорректные данные"))
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
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, ErrMessage("некорректные данные"))
			return
		}
		featureSearch = true
	}

	limit := -1
	data = r.URL.Query().Get("limit")
	if data != "" {
		limit, err = strconv.Atoi(r.URL.Query().Get("limit"))
		if err != nil || limit < 0 {
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, ErrMessage("некорректные данные"))
			return
		}
	}

	offset := -1
	data = r.URL.Query().Get("offset")
	if data != "" {
		offset, err = strconv.Atoi(r.URL.Query().Get("offset"))
		if err != nil || offset < 0 {
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, ErrMessage("некорректные данные"))
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
		w.WriteHeader(http.StatusInternalServerError)
		render.JSON(w, r, ErrMessage("Внутренняя ошибка сервера"))
		return
	}

	w.WriteHeader(http.StatusOK)
	render.JSON(w, r, banners)
}

func (h *Handler) CreateBanner(w http.ResponseWriter, r *http.Request) {
	banner := new(models.BasicBannnerInfo)
	err := json.NewDecoder(r.Body).Decode(banner)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		render.JSON(w, r, ErrMessage("некорректные данные"))
		return
	}

	bannerID, err := h.bannerService.CreateBanner(r.Context(), banner)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		render.JSON(w, r, ErrMessage("Внутренняя ошибка сервера"))
		return
	}

	w.WriteHeader(http.StatusCreated)
	render.JSON(w, r, struct {
		BannedID int `json:"banner_id"`
	}{BannedID: bannerID})
}

func ErrMessage(text string) struct {
	Err string `json:"error"`
} {
	return struct {
		Err string `json:"error"`
	}{Err: text}
}
