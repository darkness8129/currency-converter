package main

import (
	"context"
	"darkness8129/currency-converter/app/api/currency"
	"darkness8129/currency-converter/app/api/email"
	httpcontroller "darkness8129/currency-converter/app/controller/http"
	"darkness8129/currency-converter/app/entity"
	"darkness8129/currency-converter/app/service"
	"darkness8129/currency-converter/app/storage"
	"darkness8129/currency-converter/config"
	"darkness8129/currency-converter/packages/database"
	"darkness8129/currency-converter/packages/httpserver"
	"darkness8129/currency-converter/packages/logging"

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

	sql, err := database.NewPostgreSQLDatabase(database.Options{
		User:     cfg.PostgreSQL.User,
		Password: cfg.PostgreSQL.Password,
		Database: cfg.PostgreSQL.Database,
		Port:     cfg.PostgreSQL.Port,
		Host:     cfg.PostgreSQL.Host,
		Logger:   logger,
	})
	if err != nil {
		logger.Fatal("failed to init postgresql db", "err", err)
	}

	db := sql.DB()

	err = db.AutoMigrate(&entity.Subscriber{})
	if err != nil {
		logger.Fatal("automigration failed", "err", err)
	}

	storages := service.Storages{
		Subscriber: storage.NewSubscriberStorage(logger, db),
	}

	apis := service.APIs{
		Currency: currency.NewOxrAPI(logger, cfg),
		Email:    email.NewGomailAPIAPI(logger, cfg),
	}

	currencyService := service.NewCurrencyService(logger, apis)
	services := service.Services{
		Currency:     currencyService,
		Subscription: service.NewSubscriptionService(logger, cfg, storages, apis, currencyService),
	}

	go services.Subscription.StartMailing(context.Background())

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
