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

	MoneyFile string // 快车表
}

// 京东文件处理
func (file *JDExcelFile) FileAccess(filePath string) {

	// 订单明细数据
	erpRows, err := openXlsxFile(file.ERP)

	if err != nil {
		file.Error = err
		return
	}

	mapSheets, invalidatedSheets := mergeJDErpFile(erpRows[1:])

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

	cosMapSheets := excuteMergedCostFile(cosRows[1:])

	// 创建excel
	createExcelFile := excelize.NewFile()

	sheetName := "表单1"

	sheet := 0

	if createExcelFile.SheetCount > 0 {

		sheetName = createExcelFile.GetSheetName(sheet)

	} else {

		sheet = createExcelFile.NewSheet(sheetName)

		createExcelFile.SetActiveSheet(sheet)
	}

	createExcelFile.SetSheetPrOptions(sheetName, excelize.EnableFormatConditionsCalculation(true), excelize.FitToPage(true))

	/* 商品代码, 商品名称, 买家支付金额, 发货数/入仓费, 税费, 平台扣点, [成本], [京准通]*/

	rate := 0.0834

	platFormPayedRate := 0.05

	// 行索引
	rowIndex := 1

	// 插入标题
	columnTitleName, err := excelize.ColumnNumberToName(1)

	if err == nil {
		cellTitleIndex, err := excelize.JoinCellName(columnTitleName, rowIndex)

		if err == nil {
			// 插入row
			createExcelFile.SetSheetRow(sheetName, cellTitleIndex, &[]string{"商品代码", "商品名称", "买家支付金额", "发货数/入仓费", "税费", "平台扣点", "成本", "京准通"})

			rowIndex += 1
		} else {
			logs.Info("列名索引创建失败:", err.Error())
		}
	} else {
		logs.Info("标题索引创建失败", err.Error())
	}

	for _, value := range mapSheets {

		columnValueIndex := 0

		for _, cell := range value {

			columnValueIndex += 1

			columnName, err := excelize.ColumnNumberToName(columnValueIndex)

			if err != nil {
				logs.Info("获取列名失败:", err.Error())
				continue
			}

			cellIndex, err := excelize.JoinCellName(columnName, rowIndex)

			if err != nil {
				logs.Info("列名索引创建失败:", err.Error())
				continue
			}

			createExcelFile.SetCellValue(sheetName, cellIndex, cell)
		}

		// `买家金额` 坐标 `*3`
		cellNameIndex, success := getCellDescriptionName(rowIndex, 3)

		if success {
			// 税费 = 买家金额 * rate

			columnValueIndex += 1

			strScore := strconv.FormatFloat(rate, 'f', 10, 64)

			setSheetCellFormula(createExcelFile, sheetName, rowIndex, columnValueIndex, "PRODUCT("+cellNameIndex+","+strScore+")")

			// 平台扣点 = 买家金额 * platFormPayedRate
			columnValueIndex += 1

			platFormRateString := strconv.FormatFloat(platFormPayedRate, 'f', 10, 64)

			setSheetCellFormula(createExcelFile, sheetName, rowIndex, columnValueIndex, "PRODUCT("+cellNameIndex+","+platFormRateString+")")
		}

		if len(value) > 1 {
			// 添加`成本`数据
			costGoosFloatValue := mergeCostColumns(cosMapSheets, value[0].(string))

			// 成本 = 单价 * 发货数量
			costValue := costGoosFloatValue * float64(value[len(value)-1].(int64))

			columnValueIndex += 1

			insertSheetFileCell(createExcelFile, sheetName, rowIndex, columnValueIndex, costValue)
		}

		rowIndex += 1
	}

	// 保留未符合筛选条件的数据
	if len(invalidatedSheets) > 0 {

		invalidatedSheetName := "未符合条件订单明细"

		createExcelFile.NewSheet(invalidatedSheetName)

		invalidateRowIndex := 1

		createExcelFile.SetSheetPrOptions(invalidatedSheetName, excelize.EnableFormatConditionsCalculation(true), excelize.FitToPage(true))

		// 插入表头
		columnTitleName, err := excelize.ColumnNumberToName(1)

		if err == nil {
			cellTitleIndex, err := excelize.JoinCellName(columnTitleName, invalidateRowIndex)

			if err == nil {
				createExcelFile.SetSheetRow(invalidatedSheetName, cellTitleIndex, &erpRows[0])
			}
		}

		invalidateRowIndex += 1

		for _, value := range invalidatedSheets {

			columnName, err := excelize.ColumnNumberToName(1)

			if err != nil {
				logs.Info("获取列名失败:", err.Error())
				continue
			}

			cellIndex, err := excelize.JoinCellName(columnName, invalidateRowIndex)

			if err != nil {
				logs.Info("列名索引创建失败:", err.Error())
				continue
			}

			createExcelFile.SetSheetRow(invalidatedSheetName, cellIndex, &value)

			invalidateRowIndex += 1
		}
	}

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

	if len(rows) == 0 {
		return nil, errors.New("文件数据获取为空")
	}

	return rows, nil
}

func getCellDescriptionName(rowIndex int, columnIndex int) (string, bool) {
	columnName, err := excelize.ColumnNumberToName(columnIndex)

	if err != nil {
		logs.Info("获取列名失败:", err.Error())
		return "", false
	}

	cellIndex, err := excelize.JoinCellName(columnName, rowIndex)

	if err != nil {
		logs.Info("列名索引创建失败:", err.Error())
		return "", false
	}

	return cellIndex, true
}

func setSheetCellFormula(file *excelize.File, sheetName string, rowIndex int, columnIndex int, formatString string) {

	cellIndex, success := getCellDescriptionName(rowIndex, columnIndex)

	if !success {
		return
	}

	file.SetCellFormula(sheetName, cellIndex, formatString)
}

func insertSheetFileCell(file *excelize.File, sheetName string, rowIndex int, columnIndex int, value interface{}) {

	cellIndex, success := getCellDescriptionName(rowIndex, columnIndex)

	if !success {
		return
	}

	file.SetCellValue(sheetName, cellIndex, value)
}

/// 订单明细表

/* 京东ERP文件数据获取合并买家支付金额
@rows 表单数据
*/
func mergeJDErpFile(rows [][]string) (map[string][]interface{}, map[string][]string) {
	/* 商品代码[E], 商品名称[F], 买家支付金额[I + J], 发货数[H]*/

	mergeRows := make(map[string][]interface{})

	invalidatedMergeRows := make(map[string][]string)

	for rowIndex, row := range rows {

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
		if len(savedRow) == 0 || savedRow == nil || len(mergeValue) == 0 {

			newRow := []interface{}{row[4], row[5], mergedPayedFloatValue, sentGoodsIntValue}

			if len(mergeValue) == 0 {

				invalidatedMergeRows[string(rowIndex)] = row

			} else {
				mergeRows[mergeValue] = newRow
			}

		} else {
			// 合并数据

			// 发货数量
			savedRow[3] = savedRow[3].(int64) + sentGoodsIntValue

			// 买家金额
			savedRow[2] = savedRow[2].(float64) + mergedPayedFloatValue
		}
	}

	return mergeRows, invalidatedMergeRows
}

///  成本表数据筛选

/* 从成本数据内筛选获取具体的成本数据 */
func mergeCostColumns(rows map[string][]interface{}, rowIndexValue string) float64 {

	costGoodsNumber := 0.0

	for _, value := range rows {

		if len(value) > 4 {
			if rowIndexValue == value[1] {
				costGoodsNumber = value[4].(float64)
				break
			}
		}
	}

	return costGoodsNumber
}

/* 合并成本表数据 */
func excuteMergedCostFile(rows [][]string) map[string][]interface{} {

	mergeRows := make(map[string][]interface{})

	for rowIndex, row := range rows {

		if len(row) > 4 {
			// 索引下标
			mergeValue := row[1]

			savedRow := mergeRows[mergeValue]

			// 成本
			costGoodsStringValue := row[4]

			costGoodsIntValue, err := strconv.ParseFloat(costGoodsStringValue, 64)

			if err != nil {
				costGoodsIntValue = 0
			}

			// 不存在直接添加
			if len(savedRow) == 0 || savedRow == nil || len(mergeValue) == 0 {

				newRow := make([]interface{}, len(row))

				for key, value := range row {

					if key == 4 {
						newRow[key] = costGoodsIntValue
					} else {
						newRow[key] = value
					}

				}

				if len(mergeValue) == 0 {
					mergeRows[string(rowIndex)] = newRow[:]
				} else {
					mergeRows[mergeValue] = newRow[:]
				}

			} else {
				// 合并数据
				// 成本
				savedRow[4] = savedRow[4].(float64) + costGoodsIntValue
			}
		} else {
			logs.Info("当前行数据:", row)
		}
	}

	return mergeRows
}
