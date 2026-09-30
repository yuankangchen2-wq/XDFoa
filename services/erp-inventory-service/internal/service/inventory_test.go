package service

import (
	"context"
	"fmt"
	"sync"
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/glebarez/sqlite"
	"github.com/oa-portal/inventory-service/internal/cache"
	"github.com/oa-portal/inventory-service/internal/event"
	"github.com/oa-portal/inventory-service/internal/model"
	"github.com/oa-portal/inventory-service/internal/repo"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func setupTest(t *testing.T) (*InventoryService, *gorm.DB, *miniredis.Miniredis) {
	t.Helper()
	dbFile := fmt.Sprintf("file:%s?mode=memory&cache=shared", t.Name())
	db, err := gorm.Open(sqlite.Open(dbFile), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	db.AutoMigrate(&model.Warehouse{}, &model.StockAccount{}, &model.StockJournal{},
		&model.StockTransfer{}, &model.Stocktake{}, &model.StocktakeItem{})

	mr, _ := miniredis.Run()
	rdb := redis.NewClient(&redis.Options{Addr: mr.Addr()})

	stockAcc := repo.NewStockAccountRepo(db)
	svc := NewInventoryService(
		db, stockAcc, repo.NewWarehouseRepo(db), repo.NewStockJournalRepo(db),
		repo.NewStockTransferRepo(db), repo.NewStocktakeRepo(db),
		cache.NewStockRedis(rdb), event.NoopPublisher{},
	)
	return svc, db, mr
}

// 测试入库
func TestStockIn(t *testing.T) {
	svc, _, _ := setupTest(t)
	ctx := context.Background()

	if err := svc.StockIn(ctx, "t1", "SKU001", 1, 100, "PO001", "B001", 0); err != nil {
		t.Fatal(err)
	}

	// 验证库存
	res, _ := svc.ListStocks(ctx, "t1", "SKU001", 1, repo.Page{Page: 1, PageSize: 10})
	if len(res.Items) != 1 || res.Items[0].Quantity != 100 {
		t.Fatalf("expected qty=100, got %d", res.Items[0].Quantity)
	}

	// 验证流水
	journals, _ := svc.ListJournals(ctx, "t1", "SKU001", 1, repo.Page{Page: 1, PageSize: 10})
	if len(journals.Items) != 1 {
		t.Fatalf("expected 1 journal, got %d", len(journals.Items))
	}
	if journals.Items[0].ChangeType != "in" || journals.Items[0].Quantity != 100 {
		t.Fatalf("unexpected journal: %+v", journals.Items[0])
	}
}

// 测试预扣 + 真实扣减
func TestAllocateAndDeduct(t *testing.T) {
	svc, _, _ := setupTest(t)
	ctx := context.Background()

	// 入库 100
	svc.StockIn(ctx, "t1", "SKU001", 1, 100, "PO001", "", 0)

	// 预扣 30
	if err := svc.AllocateStock(ctx, "t1", "SKU001", 1, 30, "ORDER001"); err != nil {
		t.Fatal(err)
	}

	// 再预扣 80 应该失败（可用 70 < 80）
	if err := svc.AllocateStock(ctx, "t1", "SKU001", 1, 80, "ORDER002"); err == nil {
		t.Fatal("expected insufficient stock error")
	}

	// 真实扣减 30
	if err := svc.DeductStock(ctx, "t1", "SKU001", 1, 30, "ORDER001"); err != nil {
		t.Fatal(err)
	}

	// 验证库存 = 70
	res, _ := svc.ListStocks(ctx, "t1", "SKU001", 1, repo.Page{Page: 1, PageSize: 10})
	if res.Items[0].Quantity != 70 {
		t.Fatalf("expected qty=70, got %d", res.Items[0].Quantity)
	}
}

// 测试释放预扣
func TestReleaseStock(t *testing.T) {
	svc, _, _ := setupTest(t)
	ctx := context.Background()

	svc.StockIn(ctx, "t1", "SKU001", 1, 100, "PO001", "", 0)
	svc.AllocateStock(ctx, "t1", "SKU001", 1, 30, "ORDER001")

	// 释放预扣 30
	if err := svc.ReleaseStock(ctx, "t1", "SKU001", 1, 30, "ORDER001"); err != nil {
		t.Fatal(err)
	}

	// 再预扣 100 应该成功（释放后可用=100）
	if err := svc.AllocateStock(ctx, "t1", "SKU001", 1, 100, "ORDER002"); err != nil {
		t.Fatalf("expected success after release, got %v", err)
	}
}

// 测试并发预扣不超卖
func TestConcurrentAllocate(t *testing.T) {
	svc, _, _ := setupTest(t)
	ctx := context.Background()

	// 入库 1000
	svc.StockIn(ctx, "t1", "SKU001", 1, 1000, "PO001", "", 0)

	// 100 个 goroutine 各预扣 15，总共请求 1500，应该只有 66 个成功
	var wg sync.WaitGroup
	var success int64
	var mu sync.Mutex
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := svc.AllocateStock(ctx, "t1", "SKU001", 1, 15, "ORDER"); err == nil {
				mu.Lock()
				success++
				mu.Unlock()
			}
		}()
	}
	wg.Wait()

	if success != 66 {
		t.Fatalf("expected 66 successful allocations, got %d", success)
	}
}

// 测试调拨
func TestTransfer(t *testing.T) {
	svc, _, _ := setupTest(t)
	ctx := context.Background()

	// 仓1 入库 100
	svc.StockIn(ctx, "t1", "SKU001", 1, 100, "PO001", "", 0)

	// 创建调拨 30 从仓1到仓2
	tr, err := svc.CreateTransfer(ctx, "t1", 1, 2, "SKU001", 30)
	if err != nil {
		t.Fatal(err)
	}

	// 确认出库
	if err := svc.ConfirmTransferOut(ctx, tr.ID); err != nil {
		t.Fatal(err)
	}
	// 仓1 = 70
	res1, _ := svc.ListStocks(ctx, "t1", "SKU001", 1, repo.Page{Page: 1, PageSize: 10})
	if res1.Items[0].Quantity != 70 {
		t.Fatalf("expected wh1 qty=70, got %d", res1.Items[0].Quantity)
	}

	// 确认入库
	if err := svc.ConfirmTransferIn(ctx, tr.ID); err != nil {
		t.Fatal(err)
	}
	// 仓2 = 30
	res2, _ := svc.ListStocks(ctx, "t1", "SKU001", 2, repo.Page{Page: 1, PageSize: 10})
	if len(res2.Items) != 1 || res2.Items[0].Quantity != 30 {
		t.Fatalf("expected wh2 qty=30, got %+v", res2.Items)
	}
}
