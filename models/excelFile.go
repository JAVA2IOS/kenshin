package models

import (
	"errors"
	"strconv"

	kenshinUtil "kenshin/util"

	"github.com/beego/beego/v2/core/logs"
	"github.com/xuri/excelize/v2"
)

type ExcelFileInfo struct {
	Error      error  // 错误信息
	File       string // 保存后的文件绝对路径
	StaticFile string // 保存后的文件静态路径
}

// 京东excel文件处理
type JDExcelFile struct {
	ExcelFileInfo

	ERP string // erp表

	CostFile string // 成本表
}

// 京东文件处理
func (file *JDExcelFile) FileAccess(filePath string) {

	// 订单明细数据
	erpRows, err := openXlsxFile(file.ERP)

	if err != nil {
		file.Error = err
		return
	}

	mapSheets := mergeJDErpFile(erpRows[1:])

	if len(mapSheets) == 0 {
		file.Error = errors.New("获取数据为空")
		return
	}

	// 成本表数据
	cosRows, err := openXlsxFile(file.CostFile)

	if err != nil {
		file.Error = err
		return
	}

	// 创建excel
	createExcelFile := excelize.NewFile()

	sheetName := "表单1"

	sheet := createExcelFile.NewSheet(sheetName)

	// 行索引
	rowIndex := 1

	for _, value := range mapSheets {

		for columnIndex := 0; columnIndex < len(value); columnIndex++ {
			columnName, err := excelize.ColumnNumberToName(columnIndex + 1)

			if err != nil {
				logs.Info("获取列名失败:", err.Error())
				continue
			}

			cellIndex, err := excelize.JoinCellName(columnName, rowIndex)

			if err != nil {
				logs.Info("列名索引创建失败:", err.Error())
				continue
			}

			setErr := createExcelFile.SetCellValue(sheetName, cellIndex, value[columnIndex])

			if setErr != nil {
				logs.Info("设置单元格数据失败:", setErr.Error())
			}
		}

		rowIndex += 1
	}

	createExcelFile.SetActiveSheet(sheet)

	kenshinUtil.CreateFileDirectory(filePath)

	if kenshinUtil.PathExists(filePath) {
		kenshinUtil.RemoveFile(filePath)
	}

	saveErr := createExcelFile.SaveAs(filePath)

	if saveErr != nil {
		println("保存失败:", saveErr.Error())
		file.Error = saveErr
	} else {
		println("xlsx文件保存成功!")
		file.File = filePath
	}
}

/* 打开xlsx文件，读取表单数据
@file 文件路径
*/
func openXlsxFile(file string) ([][]string, error) {

	excelFile, err := excelize.OpenFile(file)
	if err != nil {

		return nil, err
	}

	// defer excelFile.Close()

	sheetlist := excelFile.GetSheetList()

	if len(sheetlist) == 0 {

		return nil, errors.New("文件为空")
	}

	rows, err := excelFile.GetRows(sheetlist[0])

	if err != nil {

		return nil, err
	}

	return rows, nil
}

/* 京东ERP文件数据获取合并买家支付金额
@rows 表单数据
*/
func mergeJDErpFile(rows [][]string) map[string][]interface{} {
	/* 店铺名称[B], 商品代码[E], 商品名称[F], 发货数[H], 买家支付金额[I + J], 平台规格名称[N] */

	mergeRows := make(map[string][]interface{})

	for _, row := range rows {

		if row[11] == "退款" || row[12] != "不是" {
			continue
		}

		// 索引下标
		mergeValue := row[4]

		savedRow := mergeRows[mergeValue]

		// 发货数
		sentGoodsStringValue := row[7]

		sentGoodsIntValue, err := strconv.ParseInt(sentGoodsStringValue, 10, 64)

		if err != nil {
			sentGoodsIntValue = 0
		}

		// 买家支付金额
		buyerPayedValue := row[8]

		buyerPayedFloatValue, err := strconv.ParseFloat(buyerPayedValue, 64)

		if err != nil {
			buyerPayedFloatValue = 0.00
		}

		// 平台支付金额
		platformPayValue := row[9]

		platformPayedFloatValue, err := strconv.ParseFloat(platformPayValue, 64)

		if err != nil {
			platformPayedFloatValue = 0.00
		}

		mergedPayedFloatValue := platformPayedFloatValue + buyerPayedFloatValue

		// 不存在直接添加
		if len(savedRow) == 0 || savedRow == nil {

			newRow := []interface{}{row[1], row[4], row[5], sentGoodsIntValue, mergedPayedFloatValue, row[13]}

			mergeRows[mergeValue] = newRow

		} else {
			// 合并数据

			// 发货数量
			savedRow[3] = savedRow[3].(int64) + sentGoodsIntValue

			// 买家金额
			savedRow[4] = savedRow[4].(float64) + mergedPayedFloatValue
		}
	}

	return mergeRows
}

/* 合并成本表数据 */
func excuteMergedCostFile(rows [][]string) map[string][]interface{} {

	mergeRows := make(map[string][]interface{})

	for _, row := range rows {

		// 索引下标
		mergeValue := row[1]

		savedRow := mergeRows[mergeValue]

		// 发货数
		sentGoodsStringValue := row[4]

		sentGoodsIntValue, err := strconv.ParseInt(sentGoodsStringValue, 10, 64)

		if err != nil {
			sentGoodsIntValue = 0
		}

		// 不存在直接添加
		if len(savedRow) == 0 || savedRow == nil {

			newRow := []interface{}
			

			mergeRows[mergeValue] = newRow

		} else {
			// 合并数据

			// 发货数量
			savedRow[3] = savedRow[3].(int64) + sentGoodsIntValue

			// 买家金额
			savedRow[4] = savedRow[4].(float64) + mergedPayedFloatValue
		}
	}

	return mergeRows
}
