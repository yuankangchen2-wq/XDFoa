package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/oa-portal/portal-bff/internal/config"
	"github.com/oa-portal/portal-bff/internal/handler"
	"github.com/oa-portal/portal-bff/internal/service"
)

func main() {
	cfg, err := config.Load("configs/config.yaml")
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	r := gin.Default()
	svc := service.NewBFFService(cfg.Services)
	h := handler.NewHandler(svc)
	h.Register(r)

	addr := fmt.Sprintf(":%d", cfg.Server.HTTPPort)
	log.Printf("portal-bff listening on %s", addr)
	r.Run(addr)
}
