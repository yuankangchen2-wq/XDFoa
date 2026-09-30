package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type OrderClient struct {
	baseURL string
	client  *http.Client
}

func NewOrderClient(baseURL string) *OrderClient {
	return &OrderClient{baseURL: baseURL, client: &http.Client{}}
}

// CreateOrder 商机转订单：调用 order-service 创建订单
func (c *OrderClient) CreateOrder(tenantID string, userID, warehouseID int64, items []map[string]interface{}) error {
	req := map[string]interface{}{
		"tenant_id":    tenantID,
		"user_id":      userID,
		"warehouse_id": warehouseID,
		"items":        items,
	}
	body, _ := json.Marshal(req)
	httpReq, _ := http.NewRequest("POST", c.baseURL+"/api/order/orders", bytes.NewReader(body))
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("X-Tenant-Id", tenantID)
	resp, err := c.client.Do(httpReq)
	if err != nil {
		return fmt.Errorf("call order-service create: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("order-service create failed: status=%d", resp.StatusCode)
	}
	return nil
}
