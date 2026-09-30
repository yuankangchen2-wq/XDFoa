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

// DeductStock 领料：调用 inventory 扣减库存
func (c *InventoryClient) DeductStock(tenantID, skuCode string, warehouseID, qty int64, referenceID string) error {
	req := map[string]interface{}{
		"tenant_id": tenantID, "sku_code": skuCode,
		"warehouse_id": warehouseID, "quantity": qty, "reference_id": referenceID,
	}
	body, _ := json.Marshal(req)
	httpReq, _ := http.NewRequest("POST", c.baseURL+"/api/inventory/stocks/deduct", bytes.NewReader(body))
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("X-Tenant-Id", tenantID)
	resp, err := c.client.Do(httpReq)
	if err != nil {
		return fmt.Errorf("call inventory deduct: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("inventory deduct failed: status=%d", resp.StatusCode)
	}
	return nil
}

// StockIn 完工入库：调用 inventory 增加库存
func (c *InventoryClient) StockIn(tenantID, skuCode string, warehouseID, qty int64, referenceID string) error {
	req := map[string]interface{}{
		"tenant_id": tenantID, "sku_code": skuCode,
		"warehouse_id": warehouseID, "quantity": qty, "reference_id": referenceID,
	}
	body, _ := json.Marshal(req)
	httpReq, _ := http.NewRequest("POST", c.baseURL+"/api/inventory/stocks/in", bytes.NewReader(body))
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("X-Tenant-Id", tenantID)
	resp, err := c.client.Do(httpReq)
	if err != nil {
		return fmt.Errorf("call inventory stockin: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("inventory stockin failed: status=%d", resp.StatusCode)
	}
	return nil
}
