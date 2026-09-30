package service

import (
	"context"
	"fmt"
	"time"

	"github.com/oa-portal/inventory-service/internal/cache"
	"github.com/oa-portal/inventory-service/internal/event"
	"github.com/oa-portal/inventory-service/internal/model"
	"github.com/oa-portal/inventory-service/internal/repo"
	"gorm.io/gorm"
)

type InventoryService struct {
	db          *gorm.DB
	stockAcc    *repo.StockAccountRepo
	warehouse   *repo.WarehouseRepo
	journal     *repo.StockJournalRepo
	transfer    *repo.StockTransferRepo
	stocktake   *repo.StocktakeRepo
	redis       *cache.StockRedis
	publisher   event.Publisher
}

func NewInventoryService(
	db *gorm.DB,
	stockAcc *repo.StockAccountRepo,
	warehouse *repo.WarehouseRepo,
	journal *repo.StockJournalRepo,
	transfer *repo.StockTransferRepo,
	stocktake *repo.StocktakeRepo,
	redis *cache.StockRedis,
	pub event.Publisher,
) *InventoryService {
	return &InventoryService{
		db: db, stockAcc: stockAcc, warehouse: warehouse, journal: journal,
		transfer: transfer, stocktake: stocktake, redis: redis, publisher: pub,
	}
}

// writeJournalAndUpdateAccount 在事务中写流水 + 更新库存账户
// 使用乐观锁（条件更新）保证并发安全，兼容 SQLite 和 PostgreSQL
func (s *InventoryService) writeJournalAndUpdateAccount(
	ctx context.Context, tx *gorm.DB,
	tenantID, skuCode string, warehouseID int64,
	changeType string, qty int64, referenceID, batchNo string,
) (int64, int64, error) {
	for attempt := 0; attempt < 10; attempt++ {
		var acc model.StockAccount
		if err := tx.WithContext(ctx).Where("tenant_id = ? AND sku_code = ? AND warehouse_id = ?", tenantID, skuCode, warehouseID).
			First(&acc).Error; err != nil {
			return 0, 0, err
		}

		beforeQty := acc.Quantity
		afterQty := beforeQty + qty
		if afterQty < 0 {
			return 0, 0, fmt.Errorf("insufficient stock: sku=%s warehouse=%d before=%d need=%d", skuCode, warehouseID, beforeQty, -qty)
		}

		// 乐观锁：条件更新，只有 quantity 未变时才成功
		result := tx.WithContext(ctx).Model(&model.StockAccount{}).
			Where("id = ? AND quantity = ?", acc.ID, beforeQty).
			Update("quantity", afterQty)
		if result.Error != nil {
			return 0, 0, result.Error
		}
		if result.RowsAffected == 0 {
			continue // 并发冲突，重试
		}

		// 写流水
		journal := model.StockJournal{
			TenantID: tenantID, SkuCode: skuCode, WarehouseID: warehouseID,
			ChangeType: changeType, Quantity: qty,
			BeforeQty: beforeQty, AfterQty: afterQty,
			ReferenceID: referenceID, BatchNo: batchNo,
		}
		if err := tx.WithContext(ctx).Create(&journal).Error; err != nil {
			return 0, 0, err
		}

		return beforeQty, afterQty, nil
	}
	return 0, 0, fmt.Errorf("concurrent update conflict after retries")
}

// AllocateStock 预扣库存（Redis 原子操作，不写 DB 流水）
func (s *InventoryService) AllocateStock(ctx context.Context, tenantID, skuCode string, warehouseID, qty int64, referenceID string) error {
	// 确保 Redis 有数据（懒加载）
	if err := s.ensureRedisWarm(ctx, tenantID, skuCode, warehouseID); err != nil {
		return err
	}

	ok, err := s.redis.Allocate(ctx, tenantID, skuCode, warehouseID, qty)
	if err != nil {
		return fmt.Errorf("redis allocate: %w", err)
	}
	if !ok {
		// 发布库存不足事件
		qty2, _, _ := s.redis.GetStock(ctx, tenantID, skuCode, warehouseID)
		s.publisher.Publish("inventory-events", skuCode, event.StockShortage{
			TenantID: tenantID, SkuCode: skuCode, WarehouseID: warehouseID,
			RequiredQty: qty, AvailableQty: qty2, ReferenceID: referenceID,
		})
		return fmt.Errorf("insufficient stock")
	}
	return nil
}

// ReleaseStock 释放预扣
func (s *InventoryService) ReleaseStock(ctx context.Context, tenantID, skuCode string, warehouseID, qty int64, referenceID string) error {
	return s.redis.Release(ctx, tenantID, skuCode, warehouseID, qty)
}

// DeductStock 真实扣减（订单支付后，从预扣转为真实扣减）
func (s *InventoryService) DeductStock(ctx context.Context, tenantID, skuCode string, warehouseID, qty int64, referenceID string) error {
	// Redis 扣减
	if err := s.redis.Deduct(ctx, tenantID, skuCode, warehouseID, qty); err != nil {
		return fmt.Errorf("redis deduct: %w", err)
	}

	// DB 事务：更新账户 + 写流水
	var beforeQty, afterQty int64
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var err error
		beforeQty, afterQty, err = s.writeJournalAndUpdateAccount(ctx, tx, tenantID, skuCode, warehouseID, "out", -qty, referenceID, "")
		return err
	})
	if err != nil {
		return err
	}

	// 发布库存变动事件
	s.publisher.Publish("inventory-events", skuCode, event.StockChanged{
		TenantID: tenantID, SkuCode: skuCode, WarehouseID: warehouseID,
		BeforeQty: beforeQty, AfterQty: afterQty, ChangeType: "out", ReferenceID: referenceID,
	})

	// 检查库存预警
	s.checkSafetyStock(ctx, tenantID, skuCode, warehouseID, afterQty)
	return nil
}

// StockIn 入库（采购到货、退货等）
func (s *InventoryService) StockIn(ctx context.Context, tenantID, skuCode string, warehouseID, qty int64, referenceID, batchNo string, expireAt int64) error {
	// 先确保账户存在（事务外执行）
	acc, err := s.stockAcc.GetOrCreate(ctx, tenantID, skuCode, warehouseID)
	if err != nil {
		return err
	}
	if batchNo != "" {
		s.db.Model(acc).Updates(map[string]interface{}{"batch_no": batchNo, "expire_at": expireAt})
	}

	// Redis 入库
	if err := s.redis.StockIn(ctx, tenantID, skuCode, warehouseID, qty); err != nil {
		return fmt.Errorf("redis stockin: %w", err)
	}

	// DB 事务
	var beforeQty, afterQty int64
	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var e error
		beforeQty, afterQty, e = s.writeJournalAndUpdateAccount(ctx, tx, tenantID, skuCode, warehouseID, "in", qty, referenceID, batchNo)
		return e
	})
	if err != nil {
		return err
	}

	s.publisher.Publish("inventory-events", skuCode, event.StockChanged{
		TenantID: tenantID, SkuCode: skuCode, WarehouseID: warehouseID,
		BeforeQty: beforeQty, AfterQty: afterQty, ChangeType: "in", ReferenceID: referenceID,
	})
	return nil
}

// CreateTransfer 创建调拨单（不扣库存，需要确认出库）
func (s *InventoryService) CreateTransfer(ctx context.Context, tenantID string, fromWH, toWH int64, skuCode string, qty int64) (*model.StockTransfer, error) {
	t := &model.StockTransfer{
		TenantID: tenantID, TransferNo: fmt.Sprintf("TR%d", time.Now().UnixNano()),
		FromWarehouseID: fromWH, ToWarehouseID: toWH, SkuCode: skuCode, Quantity: qty, Status: 1,
	}
	if err := s.transfer.Create(ctx, t); err != nil {
		return nil, err
	}
	return t, nil
}

// ConfirmTransferOut 调拨出库（从源仓扣减）
func (s *InventoryService) ConfirmTransferOut(ctx context.Context, transferID int64) error {
	t, err := s.transfer.Get(ctx, transferID)
	if err != nil {
		return err
	}
	if t.Status != 1 {
		return fmt.Errorf("transfer status invalid")
	}

	// Redis 扣减
	if err := s.redis.Deduct(ctx, t.TenantID, t.SkuCode, t.FromWarehouseID, t.Quantity); err != nil {
		return err
	}

	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		_, _, err := s.writeJournalAndUpdateAccount(ctx, tx, t.TenantID, t.SkuCode, t.FromWarehouseID, "transfer_out", -t.Quantity, t.TransferNo, "")
		return err
	})
	if err != nil {
		return err
	}
	return s.transfer.UpdateStatus(ctx, t.ID, 2)
}

// ConfirmTransferIn 调拨入库（到目标仓）
func (s *InventoryService) ConfirmTransferIn(ctx context.Context, transferID int64) error {
	t, err := s.transfer.Get(ctx, transferID)
	if err != nil {
		return err
	}
	if t.Status != 2 {
		return fmt.Errorf("transfer not out yet")
	}

	// 先确保目标仓账户存在（事务外执行）
	if _, err := s.stockAcc.GetOrCreate(ctx, t.TenantID, t.SkuCode, t.ToWarehouseID); err != nil {
		return err
	}

	// Redis 入库
	if err := s.redis.StockIn(ctx, t.TenantID, t.SkuCode, t.ToWarehouseID, t.Quantity); err != nil {
		return err
	}

	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		_, _, err := s.writeJournalAndUpdateAccount(ctx, tx, t.TenantID, t.SkuCode, t.ToWarehouseID, "transfer_in", t.Quantity, t.TransferNo, "")
		return err
	})
	if err != nil {
		return err
	}
	if err := s.transfer.UpdateStatus(ctx, t.ID, 3); err != nil {
		return err
	}

	s.publisher.Publish("inventory-events", t.SkuCode, event.StockTransferred{
		TenantID: t.TenantID, TransferNo: t.TransferNo,
		FromWarehouseID: t.FromWarehouseID, ToWarehouseID: t.ToWarehouseID,
		SkuCode: t.SkuCode, Quantity: t.Quantity,
	})
	return nil
}

// CreateStocktake 创建盘点单
func (s *InventoryService) CreateStocktake(ctx context.Context, tenantID string, warehouseID int64, items []model.StocktakeItem) (*model.Stocktake, error) {
	st := &model.Stocktake{
		TenantID: tenantID, StocktakeNo: fmt.Sprintf("ST%d", time.Now().UnixNano()),
		WarehouseID: warehouseID, Status: 1,
	}
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(st).Error; err != nil {
			return err
		}
		for i := range items {
			items[i].StocktakeID = st.ID
			items[i].DiffQty = items[i].ActualQty - items[i].SystemQty
		}
		return tx.Create(&items).Error
	})
	return st, err
}

// checkSafetyStock 检查库存预警
func (s *InventoryService) checkSafetyStock(ctx context.Context, tenantID, skuCode string, warehouseID int64, currentQty int64) {
	var acc model.StockAccount
	if err := s.db.WithContext(ctx).Where("tenant_id = ? AND sku_code = ? AND warehouse_id = ?", tenantID, skuCode, warehouseID).First(&acc).Error; err != nil {
		return
	}
	if acc.SafetyStock > 0 && currentQty < acc.SafetyStock {
		s.publisher.Publish("inventory-events", skuCode, event.StockShortage{
			TenantID: tenantID, SkuCode: skuCode, WarehouseID: warehouseID,
			RequiredQty: acc.SafetyStock, AvailableQty: currentQty, ReferenceID: "safety_check",
		})
	}
}

// ensureRedisWarm 懒加载：如果 Redis 没有该 sku 的库存数据，从 DB 加载
func (s *InventoryService) ensureRedisWarm(ctx context.Context, tenantID, skuCode string, warehouseID int64) error {
	_, _, err := s.redis.GetStock(ctx, tenantID, skuCode, warehouseID)
	if err == nil {
		return nil
	}
	// Redis 没有，从 DB 加载
	acc, err := s.stockAcc.GetOrCreate(ctx, tenantID, skuCode, warehouseID)
	if err != nil {
		return err
	}
	return s.redis.SetStock(ctx, tenantID, skuCode, warehouseID, acc.Quantity, acc.Allocated)
}

// ListStocks 查询库存列表
func (s *InventoryService) ListStocks(ctx context.Context, tenantID, skuCode string, warehouseID int64, p repo.Page) (*repo.PageResult[model.StockAccount], error) {
	return s.stockAcc.List(ctx, tenantID, skuCode, warehouseID, p)
}

// ListJournals 查询库存流水
func (s *InventoryService) ListJournals(ctx context.Context, tenantID, skuCode string, warehouseID int64, p repo.Page) (*repo.PageResult[model.StockJournal], error) {
	return s.journal.List(ctx, tenantID, skuCode, warehouseID, p)
}

// ListWarehouses 查询仓库
func (s *InventoryService) ListWarehouses(ctx context.Context, tenantID string, p repo.Page) (*repo.PageResult[model.Warehouse], error) {
	return s.warehouse.List(ctx, tenantID, p)
}

// CreateWarehouse 创建仓库
func (s *InventoryService) CreateWarehouse(ctx context.Context, w *model.Warehouse) error {
	return s.warehouse.Create(ctx, w)
}

// ListTransfers 查询调拨单
func (s *InventoryService) ListTransfers(ctx context.Context, tenantID string, p repo.Page) (*repo.PageResult[model.StockTransfer], error) {
	return s.transfer.List(ctx, tenantID, p)
}

// ListStocktakes 查询盘点单
func (s *InventoryService) ListStocktakes(ctx context.Context, tenantID string, p repo.Page) (*repo.PageResult[model.Stocktake], error) {
	return s.stocktake.List(ctx, tenantID, p)
}
