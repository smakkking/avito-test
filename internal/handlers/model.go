package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/render"
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

	useLastRevision, err := strconv.ParseBool(r.URL.Query().Get("use_last_revision"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		render.JSON(w, r, ErrMessage("некорректные данные"))
		return
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

	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		render.JSON(w, r, ErrMessage("некорректные данные"))
		return
	}

	offset, err := strconv.Atoi(r.URL.Query().Get("offset"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		render.JSON(w, r, ErrMessage("некорректные данные"))
		return
	}

	banners, err := h.bannerService.GetAllBannersFiltered(r.Context(), tagID, featureID, limit, offset)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		render.JSON(w, r, ErrMessage("Внутренняя ошибка сервера"))
		return
	}

	w.WriteHeader(http.StatusOK)
	render.JSON(w, r, banners)
}

func ErrMessage(text string) struct {
	Err string `json:"error"`
} {
	return struct {
		Err string `json:"error"`
	}{Err: text}
}
