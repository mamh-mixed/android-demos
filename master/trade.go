package master

import (
	"fmt"
	"net/http"

	"github.com/tealeg/xlsx"
)

// tradeReport 处理查找所有商户的请求
func tradeReport(w http.ResponseWriter, r *http.Request) {
	filename := r.URL.Query().Get("filename")

	var file *xlsx.File
	var sheet *xlsx.Sheet
	var row *xlsx.Row
	var cell *xlsx.Cell

	file = xlsx.NewFile()
	sheet = file.AddSheet("测试")
	row = sheet.AddRow()
	cell = row.AddCell()
	cell.Value = "I am a cell!"
	cell = row.AddCell()
	cell.Value = "中文"

	row = sheet.AddRow()
	cell = row.AddCell()
	cell.SetInt(123)
	cell = row.AddCell()
	/*
	 * The following are samples of format samples.
	 * "0.00e+00"
	 * "0", "#,##0"
	 * "0.00", "#,##0.00", "@"
	 * "#,##0 ;(#,##0)", "#,##0 ;[red](#,##0)"
	 * "#,##0.00;(#,##0.00)", "#,##0.00;[red](#,##0.00)"
	 * "0%", "0.00%"
	 * "0.00e+00", "##0.0e+0"
	 */
	cell.SetFloatWithFormat(-24245.137555, "#,##0.00")

	row = sheet.AddRow()
	row.WriteStruct(&struct {
		A int
		B float64
		C bool
	}{
		1, 3.2, true,
	}, 3)

	w.Header().Set(`Content-Type`, `application/vnd.openxmlformats-officedocument.spreadsheetml.sheet`)
	w.Header().Set(`Content-Disposition`,
		fmt.Sprintf(`attachment; filename="%s";  filename*=utf-8''%s`, filename, filename))
	file.Write(w)
}
