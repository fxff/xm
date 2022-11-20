package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"

	"xm/internal/api"
	"xm/internal/auth"
	"xm/internal/company"
	"xm/internal/recovermw"
	"xm/internal/storage"
)

func main() {
	config, err := Parse(os.Getenv("CONFIG_PATH"))
	if err != nil {
		panic(err)
	}

	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}

	db, err := pgx.Connect(context.Background(), config.DB.DSN())
	if err != nil {
		logger.Fatal("connect to db", zap.Error(err))
	}
	defer func() {
		if err := db.Close(context.Background()); err != nil {
			logger.Error("close db", zap.Error(err))
		}
	}()

	r := mux.NewRouter()

	authSubrouter := r.PathPrefix(config.Auth.Subroute).Subrouter()
	authHandler := auth.NewHandler(logger,
		auth.NewMockClientService(),
		config.Auth.Auth)
	authHandler.RegisterRoutes(authSubrouter)

	service := company.NewService(logger, storage.New(db))

	companySubrouter := r.PathPrefix(config.Company.Subroute).Subrouter()
	companySubrouter.Use(authHandler.AuthMiddleware)
	api.NewHandler(logger, service).RegisterRoutes(companySubrouter)

	addr := fmt.Sprintf(":%d", config.Server.Port)
	server := http.Server{
		Addr:              addr,
		Handler:           recovermw.New(logger).Middleware(r),
		ReadHeaderTimeout: 2 * time.Second,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Fatal("serve", zap.Error(err))
		}
	}()

	logger.Info("application started", zap.String("addr", addr))

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)
	<-stop

	if err := server.Shutdown(context.Background()); err != nil {
		logger.Fatal("shutdown server", zap.Error(err))
	}
}
