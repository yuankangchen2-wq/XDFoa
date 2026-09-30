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

	"github.com/oa-portal/iam-service/internal/auth"
	"github.com/oa-portal/iam-service/internal/config"
	"github.com/oa-portal/iam-service/internal/handler"
	"github.com/oa-portal/iam-service/internal/model"
	"github.com/oa-portal/iam-service/internal/repo"
	"github.com/oa-portal/iam-service/internal/service"
)

func main() {
	configPath := "configs/config.yaml"
	if p := os.Getenv("CONFIG_PATH"); p != "" {
		configPath = p
	}
	cfg, err := config.Load(configPath)
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	logger, _ := zap.NewProduction()
	defer logger.Sync()

	// 数据库
	db, err := gorm.Open(postgres.Open(cfg.Database.DSN()), &gorm.Config{})
	if err != nil {
		log.Fatalf("connect db: %v", err)
	}
	db.AutoMigrate(
		&model.User{}, &model.Role{}, &model.UserRole{},
		&model.Permission{}, &model.CasbinRule{}, &model.AuditLog{},
	)

	// JWT + Casbin
	jwtMgr := auth.NewJWTManager(cfg.JWT.Secret, cfg.JWT.AccessTokenTTL, cfg.JWT.RefreshTokenTTL)
	casbinMgr, err := auth.NewCasbinManager(db, cfg.Casbin.ModelPath)
	if err != nil {
		log.Fatalf("init casbin: %v", err)
	}

	// Repo
	userRepo := repo.NewUserRepo(db)
	roleRepo := repo.NewRoleRepo(db)
	urRepo := repo.NewUserRoleRepo(db)
	permRepo := repo.NewPermissionRepo(db)
	auditRepo := repo.NewAuditRepo(db)

	// Service
	authSvc := service.NewAuthService(userRepo, roleRepo, urRepo, auditRepo, jwtMgr, casbinMgr)
	userSvc := service.NewUserService(userRepo, roleRepo, urRepo, auditRepo, casbinMgr)
	roleSvc := service.NewRoleService(roleRepo, auditRepo)
	permSvc := service.NewPermissionService(permRepo, casbinMgr)

	// Handler
	authHandler := handler.NewAuthHandler(authSvc)
	userHandler := handler.NewUserHandler(userSvc)
	roleHandler := handler.NewRoleHandler(roleSvc)
	permHandler := handler.NewPermissionHandler(permSvc)

	// 路由
	r := gin.Default()
	r.Use(corsMiddleware())

	api := r.Group("/api/iam")
	authHandler.RegisterRoutes(api)

	// 需鉴权的接口
	authed := api.Group("")
	authed.Use(handler.AuthMiddleware(jwtMgr))
	userHandler.RegisterRoutes(authed)
	roleHandler.RegisterRoutes(authed)
	permHandler.RegisterRoutes(authed)

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Server.HTTPPort),
		Handler: r,
	}

	go func() {
		logger.Info("iam-service starting", zap.Int("port", cfg.Server.HTTPPort))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info("shutting down...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	srv.Shutdown(ctx)
	logger.Info("iam-service stopped")
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
