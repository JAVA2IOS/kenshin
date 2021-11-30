package models

import (
	"math"
	"runtime/internal/math"

	"github.com/xuri/excelize/v2"
)

// 京东excel文件处理
type JDExcelFile struct {
	ERP string // erp表
}

// 京东文件处理
func JDFileAccess(file *JDExcelFile) (bool, string) {
	excelFile, err := excelize.OpenFile(file.ERP)
	if err != nil {
		// fmt.Println(err)
		return false, err.Error()
	}

	// defer excelFile.Close()

	sheetlist := excelFile.GetSheetList()

	if len(sheetlist) == 0 {
		return false, "当前文件为空"
	}

	rows, err := excelFile.GetRows(sheetlist[0])

	if err != nil {
		// fmt.Println(err)
		return false, err.Error()
	}

	mergeRows(rows, 1)

	return true, ""
}

/* 合并行
@rows 行数据
@mergeIndex 指定某列为合并的索引
*/
func mergeRows(rows [][]string, mergeIndex int) {

	mergeRows := make(map[string][]string)

	for _, row := range rows {

		// mergeIndex
		mergeValue := row[mergeIndex]

		savedRow := mergeRows[mergeValue]

		// 不存在直接添加
		if len(savedRow) == 0 || savedRow == nil {
			mergeRows[mergeValue] = row
		} else {
			// 合并数据
			maxCols := math.Max(float64(len(savedRow)), float64(len(row)))

			for columnIndex := 0; columnIndex < int(maxCols); columnIndex++ {

				mergeColValue := savedRow[columnIndex]

				currentColValue := row[columnIndex]

			}
		}

		for _, colCell := range row {
			print(colCell, "\t")
		}
	}
}
