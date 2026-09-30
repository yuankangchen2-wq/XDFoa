package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type InventoryClient struct {
	baseURL string
	client  *http.Client
}

func NewInventoryClient(baseURL string) *InventoryClient {
	return &InventoryClient{baseURL: baseURL, client: &http.Client{}}
}

func (c *InventoryClient) call(path string, tenantID string, body map[string]interface{}) error {
	b, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", c.baseURL+path, bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Tenant-Id", tenantID)
	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("call inventory %s: %w", path, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("inventory %s failed: status=%d", path, resp.StatusCode)
	}
	return nil
}

// AllocateStock 下单预扣库存
func (c *InventoryClient) AllocateStock(tenantID, skuCode string, warehouseID, qty int64, referenceID string) error {
	return c.call("/api/inventory/stocks/allocate", tenantID, map[string]interface{}{
		"sku_code": skuCode, "warehouse_id": warehouseID,
		"quantity": qty, "reference_id": referenceID,
	})
}

// ReleaseStock 取消订单释放预扣
func (c *InventoryClient) ReleaseStock(tenantID, skuCode string, warehouseID, qty int64, referenceID string) error {
	return c.call("/api/inventory/stocks/release", tenantID, map[string]interface{}{
		"sku_code": skuCode, "warehouse_id": warehouseID,
		"quantity": qty, "reference_id": referenceID,
	})
}

// DeductStock 支付后真实扣减
func (c *InventoryClient) DeductStock(tenantID, skuCode string, warehouseID, qty int64, referenceID string) error {
	return c.call("/api/inventory/stocks/deduct", tenantID, map[string]interface{}{
		"sku_code": skuCode, "warehouse_id": warehouseID,
		"quantity": qty, "reference_id": referenceID,
	})
}
