package main

import (
	"os"

	"github.com/sirupsen/logrus"
	"github.com/smakkking/avito_test/internal/app"
	"github.com/smakkking/avito_test/internal/handlers"
	"github.com/smakkking/avito_test/internal/httpserver"
	"github.com/smakkking/avito_test/internal/infrastructure/cache"
	"github.com/smakkking/avito_test/internal/infrastructure/postgres"
	"github.com/smakkking/avito_test/internal/services"
)

const (
	configPath = "./config/config.yaml"
)

var (
	modeDev = os.Getenv("ENV") == ""
)

func main() {
	// init логгер
	setupLogger()

	// загрузка конфига
	config, _ := app.MustLoadConfig(configPath)

	logrus.Infoln("service started...")
	logrus.Debugln("debug messages are available")

	// init репозитории
	bannerStorage, err := postgres.NewStorage(config)
	if err != nil {
		panic(err)
	}

	bannerCache := cache.NewCache(config.CacheExpirationTime)

	// init сервисы
	bannerService := services.NewService(bannerStorage, bannerCache)

	// init хендлеры
	bannerHandler := handlers.NewHandler(bannerService)

	// запуск сервера HTTP
	srv := httpserver.NewServer(config)
	srv.SetupHandlers(bannerHandler)

	srv.Run()
}

func setupLogger() {
	logrus.SetFormatter(
		&logrus.JSONFormatter{
			PrettyPrint:     true,
			TimestampFormat: "2006-01-02 15:04:05",
		},
	)
	logrus.SetOutput(os.Stdout)

	if modeDev {
		logrus.SetLevel(logrus.DebugLevel)
	} else {
		logrus.SetLevel(logrus.InfoLevel)
	}

}
