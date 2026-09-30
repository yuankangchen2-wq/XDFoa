package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/oa-portal/inventory-service/internal/cache"
	"github.com/oa-portal/inventory-service/internal/config"
	"github.com/oa-portal/inventory-service/internal/event"
	"github.com/oa-portal/inventory-service/internal/handler"
	"github.com/oa-portal/inventory-service/internal/model"
	"github.com/oa-portal/inventory-service/internal/repo"
	"github.com/oa-portal/inventory-service/internal/service"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	cfg, err := config.Load("configs/config.yaml")
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	// DB
	db, err := gorm.Open(postgres.Open(cfg.Database.DSN()), &gorm.Config{})
	if err != nil {
		log.Fatalf("connect db: %v", err)
	}
	if err := db.AutoMigrate(
		&model.Warehouse{},
		&model.StockAccount{},
		&model.StockJournal{},
		&model.StockTransfer{},
		&model.Stocktake{},
		&model.StocktakeItem{},
	); err != nil {
		log.Fatalf("migrate: %v", err)
	}

	// Redis
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Addr,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})

	// Repos
	stockAccRepo := repo.NewStockAccountRepo(db)
	warehouseRepo := repo.NewWarehouseRepo(db)
	journalRepo := repo.NewStockJournalRepo(db)
	transferRepo := repo.NewStockTransferRepo(db)
	stocktakeRepo := repo.NewStocktakeRepo(db)

	// Service
	svc := service.NewInventoryService(db, stockAccRepo, warehouseRepo, journalRepo, transferRepo, stocktakeRepo, cache.NewStockRedis(rdb), event.NoopPublisher{})

	// HTTP
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
	log.Printf("inventory-service listening on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("run server: %v", err)
	}
}
