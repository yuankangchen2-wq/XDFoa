# 启动全部后端服务（二进制，无窗口，工作目录为服务根目录）
# 前置：MySQL 3306、PostgreSQL 5432 已运行
$base = "c:\Users\86158\Documents\trae_projects\NewAgent\oa-portal\services"
$svcs = @(
  "mdm-service","iam-service","erp-inventory-service","erp-purchase-service",
  "erp-production-service","order-service","crm-customer-service",
  "crm-opportunity-service","crm-sales-service","crm-service-ticket-service",
  "erp-finance-service","erp-cost-service","crm-analytics-service",
  "crm-scrm-service",
  "audit-service","portal-bff"
)
foreach ($s in $svcs) {
  $exe = "$base\$s\bin\$s.exe"
  if (Test-Path $exe) {
    Start-Process -FilePath $exe -WorkingDirectory "$base\$s" -WindowStyle Hidden
    Write-Output "started $s"
  } else {
    Write-Output "MISSING $s (run: go build -o bin\$s.exe ./cmd/server)"
  }
}
