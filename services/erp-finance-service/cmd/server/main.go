package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/oa-portal/erp-finance-service/internal/config"
	"github.com/oa-portal/erp-finance-service/internal/event"
	"github.com/oa-portal/erp-finance-service/internal/handler"
	"github.com/oa-portal/erp-finance-service/internal/model"
	"github.com/oa-portal/erp-finance-service/internal/repo"
	"github.com/oa-portal/erp-finance-service/internal/service"
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
	db.AutoMigrate(&model.Receivable{}, &model.Payable{}, &model.FinancePayment{})

	r := gin.Default()
	recvRepo := repo.NewReceivableRepo(db)
	payRepo := repo.NewPayableRepo(db)
	paymentRepo := repo.NewPaymentRepo(db)
	svc := service.NewFinanceService(db, recvRepo, payRepo, paymentRepo, event.NoopPublisher{})
	h := handler.NewHandler(svc)
	h.Register(r)

	addr := fmt.Sprintf(":%d", cfg.Server.HTTPPort)
	log.Printf("erp-finance-service listening on %s", addr)
	r.Run(addr)
}
