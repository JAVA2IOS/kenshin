package kenshinUtil

import (
	"fmt"
	"net/http"
	"os"
	"path"
	"strconv"
	"time"

	"github.com/xuri/excelize/v2"
)

// 创建文件目录
func CreateFileDirectory(filePath string) (bool, error) {

	fileDirectory := path.Dir(filePath)

	if _, fileError := os.Stat(fileDirectory); os.IsNotExist(fileError) {

		mkErr := os.MkdirAll(fileDirectory, 0777)

		if mkErr != nil {

			println("创建文件夹失败", mkErr.Error())

			return false, mkErr
		}

		cmdErr := os.Chmod(fileDirectory, 0777)

		if cmdErr != nil {
			println("权限授权失败", cmdErr.Error())

			return false, cmdErr
		}
	}

	return true, nil
}

// 是否是文件
func IsFile(path string) bool {

	file, err := os.Stat(path)

	if err != nil {
		return false
	}

	return !file.IsDir()
}

// 文件是否存在
// path 文件路径
// return bool 是否存在
func PathExists(path string) bool {

	_, statesError := os.Stat(path)

	return os.IsNotExist(statesError)
}

// 保存文件
func SaveFile(filePath string) {
	CreateFileDirectory(filePath)

}

// 删除文件
func RemoveFile(filePath string) {
	os.Remove(filePath)
}

// exlcel 列循环最大次数 A ~ Z, A ~ Z
const ExcelSheetColumnIndexCount = 26

// 创建excel文件
func CreatNewExcelFile() {

	f := excelize.NewFile()

	sheetName := "表单1"

	sheet := f.NewSheet(sheetName)

	var firstIndex int64 = 0

	for row := firstIndex; row < 1000; row++ {

		rowPrefixName := ExcelGetColumnName(row, 4)

		for column := 0; column < 10; column++ {

			f.SetCellValue(sheetName, ExcelGetRowName(rowPrefixName, int64(column)), row*int64(column))
		}
	}

	f.SetActiveSheet(sheet)

	filePath := "file/tmp/" + time.Now().Format("2021-11-04") + ".xlsx"

	CreateFileDirectory(filePath)

	if PathExists(filePath) {
		RemoveFile(filePath)
	}

	saveErr := f.SaveAs("file/tmp/" + time.Now().Format("2021-11-04") + ".xlsx")

	if saveErr != nil {
		println("保存失败:", saveErr.Error())
	} else {
		println("xlsx保存成功!")
	}
}

// getColumnName 生成列名
// Excel的列名规则是从A-Z往后排;超过Z以后用两个字母表示，比如AA,AB,AC;两个字母不够以后用三个字母表示，比如AAA,AAB,AAC
// 这里做数字到列名的映射：0 -> A, 1 -> B, 2 -> C
// maxColumnRowNameLen 表示名称框的最大长度，假设数据是10行，1000列，则最后一个名称框是J1000(如果有表头，则是J1001),是4位
// 这里根据 maxColumnRowNameLen 生成切片，后面生成名称框的时候可以复用这个切片，而无需扩容
func ExcelGetColumnName(column int64, maxColumnRowNameLen int) []byte {

	const A = 'A'
	if column < ExcelSheetColumnIndexCount {
		// 第一次就分配好切片的容量
		slice := make([]byte, 0, maxColumnRowNameLen)
		return append(slice, byte(A+column))
	} else {
		// 递归生成类似AA,AB,AAA,AAB这种形式的列名
		return append(ExcelGetColumnName(column/ExcelSheetColumnIndexCount-1, maxColumnRowNameLen), byte(A+column%ExcelSheetColumnIndexCount))
	}
}

// getColumnRowName 生成名称框
// Excel的名称框是用A1,A2,B1,B2来表示的，这里需要传入前一步生成的列名切片，然后直接加上行索引来生成名称框，就无需每次分配内存
func ExcelGetRowName(columnName []byte, rowIndex int64) (columnRowName string) {

	columnName = strconv.AppendInt(columnName, rowIndex+1, 10)
	rowName := string(columnName)
	println("当前列:", rowName)
	return rowName
}

// 打开excel文件
func OpenExcelFile(filePath string) {

}

// 获取excel数据流
func ReadExcelDataStream(w http.ResponseWriter, req *http.Request) {
	file, _, err := req.FormFile("file")
	if err != nil {
		// fmt.Fprintf(w, err.Error())
		return
	}
	defer file.Close()
	excelFile, err := excelize.OpenReader(file)
	if err != nil {
		// fmt.Fprintf(w, err.Error())
		return
	}

	sheetlist := excelFile.GetSheetList()

	if len(sheetlist) == 0 {

		println("找不到数据")
		return
	}

	rows, err := excelFile.GetRows(sheetlist[0])

	if err != nil {
		// fmt.Println(err)
		return
	}

	startTime := time.Now().UnixMilli()

	for _, row := range rows {
		for _, colCell := range row {
			print(colCell, "\t")
		}
		println()
	}

	endTime := time.Now().UnixMilli()

	totalTime := endTime - startTime

	// 7000行 7300ms左右
	println("读取耗时:", totalTime, "ms")

	fmt.Fprintf(w, "读取数据耗时%vms", totalTime)

	// 保存文件到本地
}
