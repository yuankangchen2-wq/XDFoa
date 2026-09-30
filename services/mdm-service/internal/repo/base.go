package repo

import (
	"context"

	"gorm.io/gorm"
)

// 通用分页参数
type PageQuery struct {
	Page     int
	PageSize int
}

type PageResult[T any] struct {
	Items    []T
	Total    int64
	Page     int
	PageSize int
}

// 分页查询辅助
func Paginate[T any](db *gorm.DB, query *gorm.DB, pq PageQuery) (*PageResult[T], error) {
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}
	var items []T
	offset := (pq.Page - 1) * pq.PageSize
	if err := query.Offset(offset).Limit(pq.PageSize).Find(&items).Error; err != nil {
		return nil, err
	}
	return &PageResult[T]{Items: items, Total: total, Page: pq.Page, PageSize: pq.PageSize}, nil
}

// ctxKey 用于在 context 中传递 db（事务）
type ctxKey struct{}

func WithTx(ctx context.Context, tx *gorm.DB) context.Context {
	return context.WithValue(ctx, ctxKey{}, tx)
}

func GetDB(ctx context.Context, defaultDB *gorm.DB) *gorm.DB {
	if tx, ok := ctx.Value(ctxKey{}).(*gorm.DB); ok {
		return tx
	}
	return defaultDB
}
