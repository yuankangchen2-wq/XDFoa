package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/oa-portal/crm-scrm-service/internal/config"
	"github.com/oa-portal/crm-scrm-service/internal/handler"
	"github.com/oa-portal/crm-scrm-service/internal/model"
	"github.com/oa-portal/crm-scrm-service/internal/repo"
	"github.com/oa-portal/crm-scrm-service/internal/service"
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
		&model.Activity{},
		&model.Lead{},
		&model.ImportTask{},
		&model.DataRecord{},
		&model.ExceptionLead{},
		&model.HandoverRecord{},
		&model.Schedule{},
		&model.ConsultantPerf{},
		&model.TodoSeed{},
		&model.FollowLog{},
		&model.Notification{},
		&model.Opportunity{},
		&model.OppAssignRecord{},
		&model.EmployeeTask{},
		&model.Appeal{},
		&model.DealApproval{},
		&model.DealLog{},
		&model.BizTicket{},
		&model.MallOrder{},
		&model.GridConfig{},
	); err != nil {
		log.Fatalf("migrate: %v", err)
	}

	if err := service.SeedDefaults(db); err != nil {
		log.Printf("seed defaults skipped: %v", err)
	}

	// 兼容旧数据：线索状态「已转化」更名为「已转商」（幂等）
	if err := db.Model(&model.Lead{}).Where("status = ?", "已转化").Update("status", model.StatusConverted).Error; err != nil {
		log.Printf("migrate lead status: %v", err)
	}

	// 顶部栏：通知/通话记录种子（各自幂等，老库也会补）
	if err := service.SeedNotifications(db); err != nil {
		log.Printf("seed notifications: %v", err)
	}
	if err := service.SeedCallLogs(db); err != nil {
		log.Printf("seed call logs: %v", err)
	}

	// 二三期：商机 / 分配流水 / 员工任务 / 申诉 / 审批 / 成单 / 工单 / 订单（各自幂等）
	if err := service.SeedP2P3(db); err != nil {
		log.Printf("seed p2p3: %v", err)
	}

	svc := service.NewScrmService(
		db,
		repo.NewActivityRepo(db),
		repo.NewLeadRepo(db),
		repo.NewImportTaskRepo(db),
		repo.NewDataRecordRepo(db),
		repo.NewExceptionRepo(db),
		repo.NewScheduleRepo(db),
		repo.NewPerfRepo(db),
		repo.NewTodoSeedRepo(db),
		repo.NewHandoverRepo(db),
		repo.NewFollowLogRepo(db),
		repo.NewNotificationRepo(db),
		repo.NewOppRepo(db),
		repo.NewAssignRecordRepo(db),
		repo.NewEmployeeTaskRepo(db),
		repo.NewAppealRepo(db),
		repo.NewDealApprovalRepo(db),
		repo.NewDealRepo(db),
		repo.NewTicketRepo(db),
		repo.NewOrderRepo(db),
		repo.NewGridConfigRepo(db),
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
	log.Printf("crm-scrm-service listening on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("run server: %v", err)
	}
}
