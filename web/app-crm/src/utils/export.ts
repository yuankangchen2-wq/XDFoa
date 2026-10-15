// 前端 CSV 导出（列表页按当前筛选快速导出，后端另有 xlsx 全量导出）
function csvCell(v: any): string {
  const s = String(v ?? '')
  if (/[",\n]/.test(s)) {
    return '"' + s.replace(/"/g, '""') + '"'
  }
  return s
}

export function downloadCsv(filename: string, headers: string[], rows: any[][]) {
  const lines = [headers, ...rows].map(r => r.map(csvCell).join(','))
  // BOM 保证 Excel 打开中文不乱码
  const blob = new Blob(['\uFEFF' + lines.join('\r\n')], { type: 'text/csv;charset=utf-8;' })
  const url = URL.createObjectURL(blob)
  const a = document.createElement('a')
  a.href = url
  a.download = filename
  a.click()
  URL.revokeObjectURL(url)
}
