package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/oa-portal/mdm-service/internal/config"
	"github.com/oa-portal/mdm-service/internal/event"
	"github.com/oa-portal/mdm-service/internal/handler"
	"github.com/oa-portal/mdm-service/internal/model"
	"github.com/oa-portal/mdm-service/internal/repo"
	"github.com/oa-portal/mdm-service/internal/service"
)

func main() {
	// 加载配置
	configPath := "configs/config.yaml"
	if p := os.Getenv("CONFIG_PATH"); p != "" {
		configPath = p
	}
	cfg, err := config.Load(configPath)
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	// 初始化日志
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	// 连接数据库
	db, err := gorm.Open(postgres.Open(cfg.Database.DSN()), &gorm.Config{})
	if err != nil {
		log.Fatalf("connect db: %v", err)
	}
	// 自动迁移（开发环境，生产用 migrations）
	db.AutoMigrate(&model.Customer{}, &model.Product{}, &model.Supplier{}, &model.Organization{})

	// 事件发布器（Noop，后续接 Kafka）
	publisher := event.NoopPublisher{}

	// 初始化 repo/service/handler
	customerRepo := repo.NewCustomerRepo(db)
	productRepo := repo.NewProductRepo(db)
	supplierRepo := repo.NewSupplierRepo(db)
	orgRepo := repo.NewOrganizationRepo(db)

	customerSvc := service.NewCustomerService(customerRepo, publisher)
	productSvc := service.NewProductService(productRepo, publisher)
	supplierSvc := service.NewSupplierService(supplierRepo, publisher)
	orgSvc := service.NewOrganizationService(orgRepo, publisher)

	customerHandler := handler.NewCustomerHandler(customerSvc)
	productHandler := handler.NewProductHandler(productSvc)
	supplierHandler := handler.NewSupplierHandler(supplierSvc)
	orgHandler := handler.NewOrganizationHandler(orgSvc)

	// 路由
	r := gin.Default()
	r.Use(corsMiddleware())
	api := r.Group("/api/mdm")
	customerHandler.RegisterRoutes(api)
	productHandler.RegisterRoutes(api)
	supplierHandler.RegisterRoutes(api)
	orgHandler.RegisterRoutes(api)

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// 启动 HTTP
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Server.HTTPPort),
		Handler: r,
	}

	go func() {
		logger.Info("mdm-service starting", zap.Int("port", cfg.Server.HTTPPort))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %v", err)
		}
	}()

	// 优雅关闭
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info("shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	srv.Shutdown(ctx)
	logger.Info("mdm-service stopped")
}

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, X-Tenant-Id, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}
