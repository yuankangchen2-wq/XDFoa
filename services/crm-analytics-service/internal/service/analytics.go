package service

import (
	"context"

	"github.com/oa-portal/crm-analytics-service/internal/event"
	"github.com/oa-portal/crm-analytics-service/internal/model"
	"github.com/oa-portal/crm-analytics-service/internal/repo"
)

type Dashboard struct {
	TotalSales     float64 `json:"total_sales"`
	TotalOrders    int32   `json:"total_orders"`
	NewCustomers   int32   `json:"new_customers"`
	OpenTickets    int32   `json:"open_tickets"`
	ConversionRate float64 `json:"conversion_rate"`
	AvgTicket      float64 `json:"avg_ticket"`
}

type AnalyticsService struct {
	salesRepo  *repo.SalesSnapshotRepo
	funnelRepo *repo.FunnelStageRepo
	custRepo   *repo.CustomerStatRepo
	publisher  event.Publisher
}

func NewAnalyticsService(
	salesRepo *repo.SalesSnapshotRepo,
	funnelRepo *repo.FunnelStageRepo,
	custRepo *repo.CustomerStatRepo,
	pub event.Publisher,
) *AnalyticsService {
	return &AnalyticsService{
		salesRepo: salesRepo, funnelRepo: funnelRepo,
		custRepo: custRepo, publisher: pub,
	}
}

// ---------- 销售快照 ----------

func (s *AnalyticsService) RecordSalesSnapshot(ctx context.Context, o *model.SalesSnapshot) error {
	if o.OrderCount > 0 {
		o.AvgAmount = o.TotalAmount / float64(o.OrderCount)
	}
	return s.salesRepo.Create(ctx, o)
}

func (s *AnalyticsService) ListSalesSnapshots(ctx context.Context, tenantID, period string, p repo.Page) (*repo.PageResult[model.SalesSnapshot], error) {
	return s.salesRepo.List(ctx, tenantID, period, p)
}

func (s *AnalyticsService) SalesSummaryByPeriod(ctx context.Context, tenantID string) ([]repo.SalesSummary, error) {
	return s.salesRepo.SummaryByPeriod(ctx, tenantID)
}

// ---------- 漏斗 ----------

func (s *AnalyticsService) RecordFunnelStage(ctx context.Context, o *model.FunnelStage) error {
	return s.funnelRepo.Create(ctx, o)
}

func (s *AnalyticsService) ListFunnel(ctx context.Context, tenantID, period string) ([]model.FunnelStage, error) {
	return s.funnelRepo.List(ctx, tenantID, period)
}

// ---------- 客户统计 ----------

func (s *AnalyticsService) RecordCustomerStat(ctx context.Context, o *model.CustomerStat) error {
	return s.custRepo.Create(ctx, o)
}

func (s *AnalyticsService) ListCustomerStats(ctx context.Context, tenantID string, p repo.Page) (*repo.PageResult[model.CustomerStat], error) {
	return s.custRepo.List(ctx, tenantID, p)
}

// ---------- 看板 ----------

// GetDashboard 综合看板：销售总额、订单数、新增客户、待处理工单、转化率、客单价
func (s *AnalyticsService) GetDashboard(ctx context.Context, tenantID string) (*Dashboard, error) {
	totalSales, totalOrders, err := s.salesRepo.TotalSales(ctx, tenantID)
	if err != nil {
		return nil, err
	}

	avgTicket := 0.0
	if totalOrders > 0 {
		avgTicket = totalSales / float64(totalOrders)
	}

	var newCustomers int32
	latest, err := s.custRepo.Latest(ctx, tenantID)
	if err == nil && latest != nil {
		newCustomers = latest.NewCustomers
	}

	// 转化率：取漏斗第一阶段到最后阶段的转化
	funnel, _ := s.funnelRepo.List(ctx, tenantID, "")
	conversionRate := 0.0
	if len(funnel) >= 2 && funnel[0].OpportunityCount > 0 {
		last := funnel[len(funnel)-1]
		conversionRate = float64(last.OpportunityCount) / float64(funnel[0].OpportunityCount) * 100
	}

	return &Dashboard{
		TotalSales:     totalSales,
		TotalOrders:    totalOrders,
		NewCustomers:   newCustomers,
		OpenTickets:    0, // 工单数据通过事件同步，此处暂为 0
		ConversionRate: conversionRate,
		AvgTicket:      avgTicket,
	}, nil
}
