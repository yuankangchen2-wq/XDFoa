package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/oa-portal/crm-analytics-service/internal/config"
	"github.com/oa-portal/crm-analytics-service/internal/event"
	"github.com/oa-portal/crm-analytics-service/internal/handler"
	"github.com/oa-portal/crm-analytics-service/internal/model"
	"github.com/oa-portal/crm-analytics-service/internal/repo"
	"github.com/oa-portal/crm-analytics-service/internal/service"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	cfg, err := config.Load("configs/config.yaml")
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	db, err := gorm.Open(mysql.Open(cfg.Database.DSN()), &gorm.Config{})
	if err != nil {
		log.Fatalf("connect db: %v", err)
	}
	db.AutoMigrate(&model.SalesSnapshot{}, &model.FunnelStage{}, &model.CustomerStat{})

	r := gin.Default()
	salesRepo := repo.NewSalesSnapshotRepo(db)
	funnelRepo := repo.NewFunnelStageRepo(db)
	custRepo := repo.NewCustomerStatRepo(db)
	svc := service.NewAnalyticsService(salesRepo, funnelRepo, custRepo, event.NoopPublisher{})
	h := handler.NewHandler(svc)
	h.Register(r)

	addr := fmt.Sprintf(":%d", cfg.Server.HTTPPort)
	log.Printf("crm-analytics-service listening on %s", addr)
	r.Run(addr)
}
