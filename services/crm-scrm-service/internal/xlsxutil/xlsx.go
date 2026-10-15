// Package xlsxutil 提供最小可用的 xlsx 读取与写出能力（仅依赖标准库）。
// xlsx 本质是 zip 包 + XML，这里只处理单行表头 + 内联字符串/共享字符串两种最常见形态。
package xlsxutil

import (
	"archive/zip"
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"strings"
)

// ReadAll 读取 xlsx 中第一个工作表，返回所有行（含表头行）
func ReadAll(data []byte) ([][]string, error) {
	zr, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		return nil, fmt.Errorf("打开xlsx失败: %w", err)
	}

	shared := readSharedStrings(zr)

	var sheetFile *zip.File
	for _, f := range zr.File {
		if f.Name == "xl/worksheets/sheet1.xml" {
			sheetFile = f
			break
		}
	}
	if sheetFile == nil {
		for _, f := range zr.File {
			if strings.HasPrefix(f.Name, "xl/worksheets/sheet") && strings.HasSuffix(f.Name, ".xml") {
				sheetFile = f
				break
			}
		}
	}
	if sheetFile == nil {
		return nil, fmt.Errorf("xlsx中未找到工作表")
	}

	rc, err := sheetFile.Open()
	if err != nil {
		return nil, err
	}
	defer rc.Close()
	raw, err := io.ReadAll(rc)
	if err != nil {
		return nil, err
	}

	dec := xml.NewDecoder(bytes.NewReader(raw))
	var rows [][]string
	var row []string
	colIndex := 0
	inCell := false
	cellType := ""
	inValue := false
	inInlineT := false
	var valBuf strings.Builder

	flushCell := func() {
		for len(row) <= colIndex {
			row = append(row, "")
		}
		s := valBuf.String()
		if cellType == "s" {
			idx := 0
			fmt.Sscanf(strings.TrimSpace(s), "%d", &idx)
			if idx >= 0 && idx < len(shared) {
				s = shared[idx]
			}
		}
		row[colIndex] = strings.TrimSpace(s)
		valBuf.Reset()
	}

	for {
		tok, err := dec.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		switch t := tok.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "row":
				row = nil
				colIndex = 0
			case "c":
				inCell = true
				cellType = ""
				valBuf.Reset()
				ref := ""
				for _, a := range t.Attr {
					if a.Name.Local == "t" {
						cellType = a.Value
					}
					if a.Name.Local == "r" {
						ref = a.Value
					}
				}
				colIndex = colLettersToIndex(ref)
				if colIndex < 0 {
					colIndex = len(row)
				}
			case "v":
				inValue = true
			case "t":
				if inCell && cellType == "inlineStr" {
					inInlineT = true
				}
			}
		case xml.EndElement:
			switch t.Name.Local {
			case "c":
				flushCell()
				inCell = false
				colIndex++
			case "v":
				inValue = false
			case "t":
				inInlineT = false
			case "row":
				rows = append(rows, row)
			}
		case xml.CharData:
			if inValue || inInlineT {
				valBuf.Write(t)
			}
		}
	}
	return rows, nil
}

func readSharedStrings(zr *zip.Reader) []string {
	var f *zip.File
	for _, cand := range zr.File {
		if cand.Name == "xl/sharedStrings.xml" {
			f = cand
			break
		}
	}
	if f == nil {
		return nil
	}
	rc, err := f.Open()
	if err != nil {
		return nil
	}
	defer rc.Close()
	raw, err := io.ReadAll(rc)
	if err != nil {
		return nil
	}
	type si struct {
		T string   `xml:"t"`
		R []struct {
			T string `xml:"t"`
		} `xml:"r"`
	}
	var doc struct {
		SI []si `xml:"si"`
	}
	if err := xml.Unmarshal(raw, &doc); err != nil {
		return nil
	}
	out := make([]string, 0, len(doc.SI))
	for _, s := range doc.SI {
		if len(s.R) > 0 {
			var b strings.Builder
			for _, r := range s.R {
				b.WriteString(r.T)
			}
			out = append(out, b.String())
		} else {
			out = append(out, s.T)
		}
	}
	return out
}

// colLettersToIndex 把单元格引用（如 A1、AB12）转为从 0 开始的列号
func colLettersToIndex(ref string) int {
	i := 0
	col := 0
	for i < len(ref) && ref[i] >= 'A' && ref[i] <= 'Z' {
		col = col*26 + int(ref[i]-'A'+1)
		i++
	}
	if i == 0 {
		return -1
	}
	return col - 1
}

// ---- 写出 ----

const contentTypes = `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Types xmlns="http://schemas.openxmlformats.org/package/2006/content-types"><Default Extension="rels" ContentType="application/vnd.openxmlformats-package.relationships+xml"/><Default Extension="xml" ContentType="application/xml"/><Override PartName="/xl/workbook.xml" ContentType="application/vnd.openxmlformats-officedocument.spreadsheetml.sheet.main+xml"/><Override PartName="/xl/worksheets/sheet1.xml" ContentType="application/vnd.openxmlformats-officedocument.spreadsheetml.worksheet+xml"/></Types>`

const rootRels = `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships"><Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/officeDocument" Target="xl/workbook.xml"/></Relationships>`

const workbook = `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<workbook xmlns="http://schemas.openxmlformats.org/spreadsheetml/2006/main" xmlns:r="http://schemas.openxmlformats.org/officeDocument/2006/relationships"><sheets><sheet name="Sheet1" sheetId="1" r:id="rId1"/></sheets></workbook>`

const workbookRels = `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships"><Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/worksheet" Target="worksheets/sheet1.xml"/></Relationships>`

func xmlEscape(s string) string {
	var b strings.Builder
	xml.EscapeText(&b, []byte(s))
	return b.String()
}

func colName(idx int) string {
	name := ""
	for idx >= 0 {
		name = string(rune('A'+idx%26)) + name
		idx = idx/26 - 1
	}
	return name
}

// Write 生成单工作表 xlsx，第一行由调用方自行放入表头
func Write(rows [][]string) ([]byte, error) {
	var sheet strings.Builder
	sheet.WriteString(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>`)
	sheet.WriteString(`<worksheet xmlns="http://schemas.openxmlformats.org/spreadsheetml/2006/main"><sheetData>`)
	for r, row := range rows {
		sheet.WriteString(fmt.Sprintf(`<row r="%d">`, r+1))
		for c, v := range row {
			ref := fmt.Sprintf("%s%d", colName(c), r+1)
			sheet.WriteString(fmt.Sprintf(`<c r="%s" t="inlineStr"><is><t xml:space="preserve">%s</t></is></c>`, ref, xmlEscape(v)))
		}
		sheet.WriteString(`</row>`)
	}
	sheet.WriteString(`</sheetData></worksheet>`)

	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)
	files := map[string]string{
		"[Content_Types].xml":              contentTypes,
		"_rels/.rels":                      rootRels,
		"xl/workbook.xml":                  workbook,
		"xl/_rels/workbook.xml.rels":       workbookRels,
		"xl/worksheets/sheet1.xml":         sheet.String(),
	}
	for name, content := range files {
		w, err := zw.Create(name)
		if err != nil {
			return nil, err
		}
		if _, err := io.WriteString(w, content); err != nil {
			return nil, err
		}
	}
	if err := zw.Close(); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
