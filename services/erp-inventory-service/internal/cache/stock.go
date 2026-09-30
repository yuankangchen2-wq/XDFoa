package cache

import (
	"context"
	"errors"
	"fmt"

	"github.com/redis/go-redis/v9"
)

// Redis 库存预扣管理器
// 用 Hash 存储：key = stock:{tenant}:{sku}:{warehouse}
// field: quantity（可用）, allocated（已预扣）
type StockRedis struct {
	rdb *redis.Client
}

func NewStockRedis(rdb *redis.Client) *StockRedis {
	return &StockRedis{rdb: rdb}
}

// 预扣脚本：检查可用量 >= 预扣量，原子扣减
// KEYS[1] = stock key
// ARGV[1] = 预扣数量
// 返回：1=成功，0=库存不足
const allocateScript = `
local qty = tonumber(redis.call('HGET', KEYS[1], 'quantity') or '0')
local alloc = tonumber(redis.call('HGET', KEYS[1], 'allocated') or '0')
local need = tonumber(ARGV[1])
if qty - alloc >= need then
  redis.call('HINCRBY', KEYS[1], 'allocated', need)
  return 1
else
  return 0
end
`

// 释放预扣脚本
const releaseScript = `
redis.call('HINCRBY', KEYS[1], 'allocated', -tonumber(ARGV[1]))
return 1
`

// 真实扣减脚本：扣减可用量 + 释放预扣
const deductScript = `
redis.call('HINCRBY', KEYS[1], 'quantity', -tonumber(ARGV[1]))
redis.call('HINCRBY', KEYS[1], 'allocated', -tonumber(ARGV[1]))
return 1
`

// 入库脚本：增加可用量
const stockInScript = `
redis.call('HINCRBY', KEYS[1], 'quantity', tonumber(ARGV[1]))
return 1
`

func (s *StockRedis) key(tenantID, skuCode string, warehouseID int64) string {
	return fmt.Sprintf("stock:%s:%s:%d", tenantID, skuCode, warehouseID)
}

// Allocate 预扣库存
func (s *StockRedis) Allocate(ctx context.Context, tenantID, skuCode string, warehouseID int64, qty int64) (bool, error) {
	result, err := s.rdb.Eval(ctx, allocateScript, []string{s.key(tenantID, skuCode, warehouseID)}, qty).Result()
	if err != nil {
		return false, err
	}
	return result.(int64) == 1, nil
}

// Release 释放预扣
func (s *StockRedis) Release(ctx context.Context, tenantID, skuCode string, warehouseID int64, qty int64) error {
	_, err := s.rdb.Eval(ctx, releaseScript, []string{s.key(tenantID, skuCode, warehouseID)}, qty).Result()
	return err
}

// Deduct 真实扣减（可用量减少 + 预扣释放）
func (s *StockRedis) Deduct(ctx context.Context, tenantID, skuCode string, warehouseID int64, qty int64) error {
	_, err := s.rdb.Eval(ctx, deductScript, []string{s.key(tenantID, skuCode, warehouseID)}, qty).Result()
	return err
}

// StockIn 入库
func (s *StockRedis) StockIn(ctx context.Context, tenantID, skuCode string, warehouseID int64, qty int64) error {
	_, err := s.rdb.Eval(ctx, stockInScript, []string{s.key(tenantID, skuCode, warehouseID)}, qty).Result()
	return err
}

// SetStock 初始化库存（从 DB 同步到 Redis）
func (s *StockRedis) SetStock(ctx context.Context, tenantID, skuCode string, warehouseID int64, qty, allocated int64) error {
	key := s.key(tenantID, skuCode, warehouseID)
	pipe := s.rdb.Pipeline()
	pipe.HSet(ctx, key, "quantity", qty)
	pipe.HSet(ctx, key, "allocated", allocated)
	_, err := pipe.Exec(ctx)
	return err
}

// GetStock 获取 Redis 库存
func (s *StockRedis) GetStock(ctx context.Context, tenantID, skuCode string, warehouseID int64) (int64, int64, error) {
	key := s.key(tenantID, skuCode, warehouseID)
	vals, err := s.rdb.HMGet(ctx, key, "quantity", "allocated").Result()
	if err != nil {
		return 0, 0, err
	}
	if len(vals) < 2 || vals[0] == nil {
		return 0, 0, errors.New("stock not found in redis")
	}
	qty, _ := vals[0].(int64)
	alloc, _ := vals[1].(int64)
	return qty, alloc, nil
}
