package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// InventoryClient 调用 inventory-service 的 REST 接口
type InventoryClient struct {
	baseURL string
	client  *http.Client
}

func NewInventoryClient(baseURL string) *InventoryClient {
	return &InventoryClient{
		baseURL: baseURL,
		client:  &http.Client{},
	}
}

type stockInReq struct {
	TenantID    string `json:"tenant_id"`
	SkuCode     string `json:"sku_code"`
	WarehouseID int64  `json:"warehouse_id"`
	Quantity    int64  `json:"quantity"`
	ReferenceID string `json:"reference_id"`
	BatchNo     string `json:"batch_no"`
}

// StockIn 调用 inventory-service 入库
func (c *InventoryClient) StockIn(tenantID, skuCode string, warehouseID, qty int64, referenceID, batchNo string) error {
	req := stockInReq{
		TenantID: tenantID, SkuCode: skuCode, WarehouseID: warehouseID,
		Quantity: qty, ReferenceID: referenceID, BatchNo: batchNo,
	}
	body, _ := json.Marshal(req)

	httpReq, err := http.NewRequest("POST", c.baseURL+"/api/inventory/stocks/in", bytes.NewReader(body))
	if err != nil {
		return err
	}
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
