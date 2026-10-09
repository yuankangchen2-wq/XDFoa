package main

import (
	"fmt"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/oa-portal/portal-bff/internal/config"
	"github.com/oa-portal/portal-bff/internal/handler"
	"github.com/oa-portal/portal-bff/internal/model"
	"github.com/oa-portal/portal-bff/internal/repo"
	"github.com/oa-portal/portal-bff/internal/service"
)

func main() {
	cfg, err := config.Load("configs/config.yaml")
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	// 门户库
	portalDB, err := gorm.Open(postgres.Open(cfg.Database.PortalDSN), &gorm.Config{})
	if err != nil {
		log.Fatalf("connect portal db: %v", err)
	}
	if err := portalDB.AutoMigrate(
		&model.PortalApproval{},
		&model.PortalSchedule{},
		&model.PortalFavorite{},
	); err != nil {
		log.Fatalf("migrate portal db: %v", err)
	}

	// iam 只读库（员工搜索）
	iamDB, err := gorm.Open(postgres.Open(cfg.Database.IamDSN), &gorm.Config{})
	if err != nil {
		log.Fatalf("connect iam db: %v", err)
	}

	seedDefaults(portalDB)

	r := gin.Default()
	bffSvc := service.NewBFFService(cfg.Services)
	portalSvc := service.NewPortalService(
		repo.NewPortalRepo(portalDB),
		repo.NewEmployeeRepo(iamDB),
	)
	h := handler.NewHandler(bffSvc, portalSvc)
	h.Register(r)

	addr := fmt.Sprintf(":%d", cfg.Server.HTTPPort)
	log.Printf("portal-bff listening on %s", addr)
	r.Run(addr)
}

// seedDefaults 首次启动播种默认常用入口与示例日程
func seedDefaults(db *gorm.DB) {
	var count int64
	db.Model(&model.PortalFavorite{}).Where("user_id = ?", service.CurrentUserID).Count(&count)
	if count == 0 {
		defaults := []model.PortalFavorite{
			{Name: "我的考勤", Icon: "Clock", Path: "", Sort: 1},
			{Name: "我的报销", Icon: "Wallet", Path: "", Sort: 2},
			{Name: "会议室预定", Icon: "Calendar", Path: "", Sort: 3},
			{Name: "400小助手", Icon: "Service", Path: "", Sort: 4},
			{Name: "OA系统", Icon: "Monitor", Path: "/iam", Sort: 5},
			{Name: "ERP", Icon: "Box", Path: "/erp", Sort: 6},
			{Name: "SCRM", Icon: "ChatDotRound", Path: "/crm", Sort: 7},
		}
		for i := range defaults {
			defaults[i].UserID = service.CurrentUserID
		}
		if err := db.Create(&defaults).Error; err != nil {
			log.Printf("seed favorites: %v", err)
		}
	}

	db.Model(&model.PortalSchedule{}).Where("user_id = ?", service.CurrentUserID).Count(&count)
	if count == 0 {
		now := time.Now()
		loc := time.Local
		day := func(d int) time.Time {
			return time.Date(now.Year(), now.Month(), d, 0, 0, 0, 0, loc)
		}
		schedules := []model.PortalSchedule{
			{UserID: service.CurrentUserID, Title: "部门周例会", EventDate: day(12)},
			{UserID: service.CurrentUserID, Title: "季度经营复盘", EventDate: day(20)},
		}
		if err := db.Create(&schedules).Error; err != nil {
			log.Printf("seed schedules: %v", err)
		}
	}
}
