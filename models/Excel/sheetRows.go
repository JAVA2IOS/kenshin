package models

type ExcelCell struct {
	ColumnIndex int         // 列坐标
	Value       interface{} // 值
}

type ExcelRow struct {
	Cells []ExcelCell // 数组
}

// 长度
func (row *ExcelRow) Len() int {
	return len(row.Cells)
}

func (row *ExcelRow) Each(block func(columnIndex int, cellValue interface{})) {
	for index, value := range row.Cells {
		block(index, value)
	}
}

type SheetRow struct {
	Rows      []ExcelRow // 数组
	Titles    []string   // 标题
	SheetName string
}

func (sheet *SheetRow) Len() int {
	return len(sheet.Rows)
}

func (sheet *SheetRow) Each(block func(rowIndex int, title string, row ExcelRow, columnIndex int, cellValue interface{})) {

	for index, value := range sheet.Rows {

		value.Each(func(columnIndex int, cellValue interface{}) {

			title := ""

			if len(sheet.Titles) > index {
				title = sheet.Titles[index]
			}

			block(index, title, value, columnIndex, cellValue)
		})
	}
}
