package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/oa-portal/audit-service/internal/config"
	"github.com/oa-portal/audit-service/internal/event"
	"github.com/oa-portal/audit-service/internal/handler"
	"github.com/oa-portal/audit-service/internal/model"
	"github.com/oa-portal/audit-service/internal/repo"
	"github.com/oa-portal/audit-service/internal/service"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	cfg, err := config.Load("configs/config.yaml")
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	db, err := gorm.Open(postgres.Open(cfg.Database.DSN()), &gorm.Config{})
	if err != nil {
		log.Fatalf("connect db: %v", err)
	}
	db.AutoMigrate(&model.AuditLog{})

	r := gin.Default()
	auditRepo := repo.NewAuditLogRepo(db)
	svc := service.NewAuditService(auditRepo, event.NoopPublisher{})
	h := handler.NewHandler(svc)
	h.Register(r)

	addr := fmt.Sprintf(":%d", cfg.Server.HTTPPort)
	log.Printf("audit-service listening on %s", addr)
	r.Run(addr)
}
