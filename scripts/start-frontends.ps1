# 启动全部前端 dev server（node 直启 vite，无窗口）
$web = "c:\Users\86158\Documents\trae_projects\NewAgent\oa-portal\web"
$apps = @("portal-main","app-mdm","app-iam","app-erp","app-portal","app-crm")
foreach ($a in $apps) {
  $vite = "$web\$a\node_modules\vite\bin\vite.js"
  if (Test-Path $vite) {
    Start-Process -FilePath "node" -ArgumentList $vite -WorkingDirectory "$web\$a" -WindowStyle Hidden
    Write-Output "started $a"
  } else {
    Write-Output "MISSING $a (run npm install first)"
  }
}
