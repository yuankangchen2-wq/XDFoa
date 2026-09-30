package service

import (
	"context"
	"net/http"
	"time"
)

type BFFService struct {
	services map[string]string
	client   *http.Client
}

func NewBFFService(services map[string]string) *BFFService {
	return &BFFService{
		services: services,
		client:   &http.Client{Timeout: 3 * time.Second},
	}
}

type HealthResult struct {
	Service string `json:"service"`
	URL     string `json:"url"`
	Status  string `json:"status"` // up / down
	Latency int64  `json:"latency_ms"`
}

// CheckAllHealth 聚合各服务健康状态
func (s *BFFService) CheckAllHealth(ctx context.Context) []HealthResult {
	results := make([]HealthResult, 0, len(s.services))
	for name, url := range s.services {
		start := time.Now()
		req, err := http.NewRequestWithContext(ctx, "GET", url+"/health", nil)
		status := "down"
		if err == nil {
			resp, err := s.client.Do(req)
			if err == nil {
				resp.Body.Close()
				if resp.StatusCode < 500 {
					status = "up"
				}
			}
		}
		results = append(results, HealthResult{
			Service: name,
			URL:     url,
			Status:  status,
			Latency: time.Since(start).Milliseconds(),
		})
	}
	return results
}

type MenuItem struct {
	Key   string `json:"key"`
	Title string `json:"title"`
	Path  string `json:"path"`
	Icon  string `json:"icon"`
}

// GetMenus 返回主应用菜单
func (s *BFFService) GetMenus() []MenuItem {
	return []MenuItem{
		{Key: "portal", Title: "工作台", Path: "/portal", Icon: "HomeFilled"},
		{Key: "iam", Title: "权限中心", Path: "/iam", Icon: "Lock"},
		{Key: "mdm", Title: "主数据", Path: "/mdm", Icon: "Files"},
		{Key: "erp", Title: "ERP", Path: "/erp", Icon: "Box"},
		{Key: "crm", Title: "CRM", Path: "/crm", Icon: "UserFilled"},
	}
}
