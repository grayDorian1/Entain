package app

import (
	"context"
	"log"
	"net/http"

	"github.com/grayDorian1/Entain/internal/config"
	"github.com/grayDorian1/Entain/internal/db"
	"github.com/grayDorian1/Entain/internal/handler"
	"github.com/grayDorian1/Entain/internal/repository"
	"github.com/grayDorian1/Entain/internal/service"
	"github.com/grayDorian1/Entain/internal/logger"
)

func Run() {
	logger.Init()
	logger.Info("starting application")
	cfg := config.Load()
	ctx := context.Background()

	pool, err := db.NewPool(ctx, cfg)
	if err != nil {
		log.Fatalf("init db: %v", err)
	}
	defer pool.Close()

	repo := repository.New(pool)
	svc := service.New(repo)
	h := handler.New(svc)

	addr := ":" + cfg.ServerPort
	log.Printf("server listening on %s", addr)
	if err := http.ListenAndServe(addr, h.SetupRoutes()); err != nil {
		log.Fatalf("server error: %v", err)
	}
}