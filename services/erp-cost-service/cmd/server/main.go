package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/oa-portal/erp-cost-service/internal/config"
	"github.com/oa-portal/erp-cost-service/internal/event"
	"github.com/oa-portal/erp-cost-service/internal/handler"
	"github.com/oa-portal/erp-cost-service/internal/model"
	"github.com/oa-portal/erp-cost-service/internal/repo"
	"github.com/oa-portal/erp-cost-service/internal/service"
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
	db.AutoMigrate(&model.CostCenter{}, &model.ProductCost{}, &model.CostRecord{})

	r := gin.Default()
	ccRepo := repo.NewCostCenterRepo(db)
	pcRepo := repo.NewProductCostRepo(db)
	recordRepo := repo.NewCostRecordRepo(db)
	svc := service.NewCostService(db, ccRepo, pcRepo, recordRepo, event.NoopPublisher{})
	h := handler.NewHandler(svc)
	h.Register(r)

	addr := fmt.Sprintf(":%d", cfg.Server.HTTPPort)
	log.Printf("erp-cost-service listening on %s", addr)
	r.Run(addr)
}
