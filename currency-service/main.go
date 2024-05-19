package main

import (
	currapi "darkness8129/currency-converter/currency-service/app/api/currency"
	httpcontroller "darkness8129/currency-converter/currency-service/app/controller/http"
	"darkness8129/currency-converter/currency-service/app/service"
	"darkness8129/currency-converter/currency-service/config"
	"darkness8129/currency-converter/currency-service/packages/httpserver"
	"darkness8129/currency-converter/currency-service/packages/logging"

	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	logger, err := logging.NewZapLogger()
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	cfg, err := config.New()
	if err != nil {
		logger.Fatal("failed to get config", "err", err)
	}

	apis := service.APIs{
		Currency: currapi.NewOxrAPI(logger, cfg),
	}

	services := service.Services{
		Currency: service.NewCurrencyService(logger, apis),
	}

	// init http server and start it
	httpServer := httpserver.NewGinHTTPServer(httpserver.Options{
		Addr:         cfg.HTTP.Addr,
		WriteTimeout: cfg.HTTP.WriteTimeout,
		ReadTimeout:  cfg.HTTP.ReadTimeout,
	})

	router := httpServer.Router()

	httpcontroller.New(httpcontroller.Options{
		Router:   router,
		Logger:   logger,
		Services: services,
	})

	httpServer.Start()

	// graceful shutdown the http server with a timeout
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		logger.Info("app interrupt", "signal", s.String())
	case err := <-httpServer.Notify():
		logger.Error("err from notify ch", "err", err)
	}

	err = httpServer.Shutdown(cfg.ShutdownTimeout)
	if err != nil {
		logger.Error("failed to shutdown server", "err", err)
	}
}
