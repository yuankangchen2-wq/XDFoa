package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/oa-portal/order-service/internal/client"
	"github.com/oa-portal/order-service/internal/config"
	"github.com/oa-portal/order-service/internal/event"
	"github.com/oa-portal/order-service/internal/handler"
	"github.com/oa-portal/order-service/internal/model"
	"github.com/oa-portal/order-service/internal/repo"
	"github.com/oa-portal/order-service/internal/service"
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
	if err := db.AutoMigrate(
		&model.Order{},
		&model.OrderItem{},
		&model.OrderLog{},
	); err != nil {
		log.Fatalf("migrate: %v", err)
	}

	invClient := client.NewInventoryClient(cfg.InventoryService.BaseURL)
	svc := service.NewOrderService(
		db,
		repo.NewOrderRepo(db),
		repo.NewOrderLogRepo(db),
		invClient,
		event.NoopPublisher{},
	)

	r := gin.Default()
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Tenant-Id")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})
	handler.NewHandler(svc).Register(r)

	addr := fmt.Sprintf(":%d", cfg.Server.HTTPPort)
	log.Printf("order-service listening on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("run server: %v", err)
	}
}
