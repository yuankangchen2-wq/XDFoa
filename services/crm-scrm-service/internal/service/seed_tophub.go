package service

import (
	"time"

	"github.com/oa-portal/crm-scrm-service/internal/model"
	"gorm.io/gorm"
)

// SeedNotifications 顶部栏通知种子（幂等；老库启动时补齐）
func SeedNotifications(db *gorm.DB) error {
	var n int64
	db.Model(&model.Notification{}).Where("tenant_id = ?", CurrentTenant).Count(&n)
	if n > 0 {
		return nil
	}
	now := time.Now()
	ago := func(h int) time.Time { return now.Add(-time.Duration(h) * time.Hour) }
	rows := []model.Notification{
		{Type: "handover", Title: "您有新的待交接线索", Content: "活动「20261014-高中-迂回转学」新分配 16 条待交接线索，请在今日内完成交接", BizType: "lead", Route: "/handover-leads", IsRead: 0, CreatedAt: ago(1)},
		{Type: "todo", Title: "今日回访提醒", Content: "您有 5 条线索安排了今日回访，请及时跟进", BizType: "home", Route: "/all-leads?revisit=today", IsRead: 0, CreatedAt: ago(2)},
		{Type: "task", Title: "导入任务已完成", Content: "任务「20261014-高中-线下地推」导入完成：成功 581 条，失败 7 条，去重 396 条", BizType: "data", Route: "/data-record", IsRead: 0, CreatedAt: ago(3)},
		{Type: "todo", Title: "逾期未跟进线索", Content: "您有 7 条线索已超过下次回访时间仍未跟进", BizType: "lead", Route: "/all-leads?revisit=overdue", IsRead: 0, CreatedAt: ago(5)},
		{Type: "system", Title: "异常线索待处理", Content: "检测到 6 条异常线索（学员信息冲突/不完整），请前往核对处理", BizType: "exception", Route: "/exception-leads", IsRead: 0, CreatedAt: ago(6)},
		{Type: "handover", Title: "线索交接确认", Content: "王主管向您交接了 3 条重点线索，请确认接收", BizType: "lead", Route: "/handover-leads", IsRead: 0, CreatedAt: ago(8)},
		{Type: "task", Title: "活动即将开始", Content: "活动「初二期末冲刺讲座」将于 10 天后开始，请提前准备地推物料", BizType: "activity", Route: "/activity", IsRead: 0, CreatedAt: ago(20)},
		{Type: "system", Title: "云座席通话小结", Content: "昨日您共外呼 12 通，接通 8 通，接通率 67%，继续保持", BizType: "home", Route: "/home", IsRead: 0, CreatedAt: ago(24)},
		{Type: "todo", Title: "重复跟进提醒", Content: "系统检测到 3 条线索存在重复跟进记录，请核对归属", BizType: "lead", Route: "/all-leads?dup=follow", IsRead: 0, CreatedAt: ago(26)},
		{Type: "system", Title: "新版本上线通知", Content: "SCRM 顶部栏已上线全局搜索、云座席双呼与通知中心，欢迎体验", BizType: "home", Route: "/home", IsRead: 0, CreatedAt: ago(30)},
		{Type: "task", Title: "模板下载提醒", Content: "线索导入支持手机号/学员号两种模板，导入前请下载最新模板", BizType: "import", Route: "/lead-import", IsRead: 1, CreatedAt: ago(50)},
		{Type: "system", Title: "账号安全提示", Content: "您的账号于武汉登录，如非本人操作请及时修改密码", BizType: "home", Route: "/home", IsRead: 1, CreatedAt: ago(72)},
	}
	for i := range rows {
		rows[i].TenantID = CurrentTenant
		rows[i].UserID = CurrentUserID
	}
	return db.CreateInBatches(&rows, 50).Error
}

// SeedCallLogs 云座席历史通话种子（幂等），并回写线索通话维度
func SeedCallLogs(db *gorm.DB) error {
	var n int64
	db.Model(&model.FollowLog{}).Where("tenant_id = ?", CurrentTenant).Count(&n)
	if n > 0 {
		return nil
	}
	var leads []model.Lead
	if err := db.Where("tenant_id = ?", CurrentTenant).Order("id ASC").Limit(8).Find(&leads).Error; err != nil {
		return err
	}
	if len(leads) == 0 {
		return nil
	}
	now := time.Now()
	rows := []model.FollowLog{}
	data := []struct {
		connected bool
		duration  int32
		hoursAgo  int
	}{
		{true, 187, 2}, {true, 342, 6}, {false, 0, 9}, {true, 95, 26},
		{true, 421, 30}, {false, 0, 34}, {true, 63, 50}, {true, 258, 74},
	}
	for i, l := range leads {
		if i >= len(data) {
			break
		}
		d := data[i]
		content, result := "双呼接通，通话时长 "+fmtDur(d.duration), "已接通"
		if !d.connected {
			content, result = "双呼发起成功，客户未接通", "未接通"
		}
		rows = append(rows, model.FollowLog{
			TenantID: CurrentTenant, LeadID: l.ID, LeadName: l.Name, Type: "双呼",
			Content: content, Result: result, DurationSec: d.duration,
			CustomerPhone: l.Phone, Operator: CurrentUserName, OperatorID: CurrentUserID,
			CreatedAt: now.Add(-time.Duration(d.hoursAgo) * time.Hour),
		})
	}
	if err := db.CreateInBatches(&rows, 50).Error; err != nil {
		return err
	}
	// 回写线索通话计数（最后一通的结果作为接线状态）
	for i, l := range leads {
		if i >= len(data) {
			break
		}
		d := data[i]
		status := "未接通"
		if d.connected {
			status = "已接通"
		}
		callAt := now.Add(-time.Duration(d.hoursAgo) * time.Hour)
		db.Model(&model.Lead{}).Where("id = ?", l.ID).Updates(map[string]interface{}{
			"call_status": status, "last_call_at": callAt,
			"self_call_count": 1, "total_call_count": 1,
		})
	}
	return nil
}

func fmtDur(sec int32) string {
	return time.Unix(int64(sec), 0).UTC().Format("04:05")
}
