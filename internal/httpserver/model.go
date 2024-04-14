package httpserver

import (
	"context"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/sirupsen/logrus"
	"github.com/smakkking/avito_test/internal/app"
	"github.com/smakkking/avito_test/internal/handlers"
	myMiddleware "github.com/smakkking/avito_test/internal/httpserver/middleware"
)

type HTTPService struct {
	srv http.Server
	mux *chi.Mux
}

func NewServer(cfg app.Config) *HTTPService {
	service := &HTTPService{
		mux: chi.NewRouter(),
	}

	service.srv = http.Server{
		Addr:         cfg.HTTPAddress,
		Handler:      service.mux,
		ReadTimeout:  cfg.HTTPReadTimeout,
		WriteTimeout: cfg.HTTPWriteTimeout,
		IdleTimeout:  cfg.HTTPIdleTimeout,
	}

	return service
}

func (h *HTTPService) SetupHandlers(bannerHandler *handlers.Handler) {
	// setup middleware
	h.mux.Use(middleware.RequestID)
	h.mux.Use(middleware.Recoverer)
	h.mux.Use(myMiddleware.Authorization)

	h.mux.Get("/user_banner", bannerHandler.GetUserBanner)
	h.mux.Get("/banner", bannerHandler.GetAllBannersFiltered)
	h.mux.Post("/banner", bannerHandler.CreateBanner)
	h.mux.Patch("/banner/{id}", bannerHandler.UpdateBanner)
	h.mux.Delete("/banner/{id}", bannerHandler.DeleteBanner)
}

func (h *HTTPService) Run() {
	err := h.srv.ListenAndServe()
	if err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			logrus.Errorf("cannot start server: %s", err.Error())
		}
	}
}

func (h *HTTPService) Shutdown(ctx context.Context) {
	if err := h.srv.Shutdown(ctx); err != nil {
		logrus.Errorf("failed to stop server %s", err.Error())
		return
	}
}
