package service

import (
	"context"
	"fmt"

	"github.com/oa-portal/crm-opportunity-service/internal/event"
	"github.com/oa-portal/crm-opportunity-service/internal/model"
	"github.com/oa-portal/crm-opportunity-service/internal/repo"
)

// OrderClient 订单服务客户端接口
type OrderClient interface {
	CreateOrder(tenantID string, userID, warehouseID int64, items []map[string]interface{}) error
}

type OpportunityService struct {
	repo      *repo.OpportunityRepo
	orderClient OrderClient
	publisher event.Publisher
}

func NewOpportunityService(
	repo *repo.OpportunityRepo,
	orderClient OrderClient,
	pub event.Publisher,
) *OpportunityService {
	return &OpportunityService{
		repo: repo, orderClient: orderClient, publisher: pub,
	}
}

// 阶段定义
const (
	StageLead      = 1 // 线索
	StageIntent    = 2 // 意向
	StageQuote     = 3 // 报价
	StageNegotiate = 4 // 谈判
	StageWon       = 5 // 成交
	StageLost      = 6 // 输单
)

// 各阶段默认赢率
var stageWinRate = map[int32]int32{
	StageLead: 10, StageIntent: 30, StageQuote: 50,
	StageNegotiate: 70, StageWon: 100, StageLost: 0,
}

// CreateOpportunity 创建商机
func (s *OpportunityService) CreateOpportunity(ctx context.Context, o *model.Opportunity) error {
	if o.Stage == 0 {
		o.Stage = StageLead
	}
	if o.WinRate == 0 {
		o.WinRate = stageWinRate[o.Stage]
	}
	if err := s.repo.Create(ctx, o); err != nil {
		return err
	}
	s.publisher.Publish("crm-events", o.Name, event.OpportunityWon{
		TenantID: o.TenantID, OpportunityID: o.ID, CustomerID: o.CustomerID, Amount: o.Amount,
	})
	return nil
}

// ListOpportunities 查询商机
func (s *OpportunityService) ListOpportunities(ctx context.Context, tenantID, ownerID string, stage int32, p repo.Page) (*repo.PageResult[model.Opportunity], error) {
	return s.repo.List(ctx, tenantID, ownerID, stage, p)
}

// GetOpportunity 获取商机详情
func (s *OpportunityService) GetOpportunity(ctx context.Context, id int64) (*model.Opportunity, error) {
	return s.repo.Get(ctx, id)
}

// UpdateOpportunity 更新商机
func (s *OpportunityService) UpdateOpportunity(ctx context.Context, o *model.Opportunity) error {
	return s.repo.Update(ctx, o)
}

// DeleteOpportunity 删除商机
func (s *OpportunityService) DeleteOpportunity(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}

// TransitionStage 阶段流转
func (s *OpportunityService) TransitionStage(ctx context.Context, id, toStage int32) error {
	if toStage < StageLead || toStage > StageLost {
		return fmt.Errorf("invalid stage: %d", toStage)
	}
	o, err := s.repo.Get(ctx, int64(id))
	if err != nil {
		return err
	}
	// 已成交或输单不能再流转
	if o.Stage == StageWon || o.Stage == StageLost {
		return fmt.Errorf("opportunity already closed, stage=%d", o.Stage)
	}

	if err := s.repo.UpdateStage(ctx, int64(id), toStage); err != nil {
		return err
	}

	// 更新赢率
	if toStage == StageWon {
		s.publisher.Publish("crm-events", "", event.OpportunityWon{
			TenantID: o.TenantID, OpportunityID: o.ID, CustomerID: o.CustomerID, Amount: o.Amount,
		})
	}
	return nil
}

// ConvertToOrder 商机转订单（成交后创建销售订单）
func (s *OpportunityService) ConvertToOrder(ctx context.Context, id, userID, warehouseID int64) error {
	o, err := s.repo.Get(ctx, id)
	if err != nil {
		return err
	}
	if o.Stage != StageWon {
		return fmt.Errorf("only won opportunity can convert to order, stage=%d", o.Stage)
	}

	// 商机转订单：以商机金额作为订单金额，创建一个虚拟 SKU 明细
	items := []map[string]interface{}{
		{
			"sku_code":   fmt.Sprintf("OPP-%d", o.ID),
			"quantity":   1,
			"unit_price": o.Amount,
		},
	}
	if err := s.orderClient.CreateOrder(o.TenantID, userID, warehouseID, items); err != nil {
		return fmt.Errorf("create order failed: %w", err)
	}
	return nil
}
