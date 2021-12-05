package controllers

import (
	"kenshin/models"
	util "kenshin/util"
	"strconv"
	"strings"
	"time"

	"github.com/beego/beego/v2/core/logs"
)

// 所有的controller都继承当前controller
type FileController struct {
	BaseController
}

func (c *FileController) Get() {
}

func (c *FileController) Upload() {

	file, h, err := c.GetFile("file")
	if err != nil {
		// http.Error(w, err.Error(), http.StatusNoContent)
		logs.Info("获取文件失败")
		c.Failure(201, "读取上传文件失败")
		return
	}

	defer file.Close()

	dateDirectory := strconv.Itoa(time.Now().Local().Year()) + "_" + strconv.Itoa(int(time.Now().Local().Month())) + "_" + strconv.Itoa(time.Now().Local().Day())

	fileDirectory := "file/tmp/" + strings.ReplaceAll(c.GetSession("uid").(string), " ", "_") + "/" + dateDirectory + "/" + h.Filename

	_, fileError := util.CreateFileDirectory(fileDirectory)

	if fileError != nil {
		logs.Info("文件创建失败", fileError.Error())
		c.Failure(201, "文件保存失败")
		return
	}

	saveErr := c.SaveToFile("file", fileDirectory)

	if saveErr != nil {
		logs.Info("保存文件失败", saveErr.Error())
		c.Failure(201, "文件保存失败")
		return
	}

	logs.Info("文件名称: ", h.Filename)

	c.Data["XlsxFile"] = h.Filename

	c.Success(map[string]string{"xlsx": h.Filename, "url": fileDirectory})
}

func (c *FileController) AccessJDFile() {
	// /file/xlsx/:url
	url := c.Ctx.Input.Param(":url")

	if url == "0" {
		c.excuteJDFile()
		return
	}

	c.SuccessMessage("文件上传成功")
}

func (c *FileController) excuteJDFile() bool {

	erpFile := c.GetString("erp")

	if len(erpFile) == 0 {

		c.Failure(201, "找不到ERP订单明细文件")
		return false
	}

	cosFile := c.GetString("cos")

	if len(cosFile) == 0 {

		c.Failure(201, "找不到成本表文件")
		return false
	}

	dateDirectory := strconv.Itoa(time.Now().Local().Year()) + "_" + strconv.Itoa(int(time.Now().Local().Month())) + "_" + strconv.Itoa(time.Now().Local().Day())

	fileName := dateDirectory + "_" + strconv.Itoa(time.Now().Local().Hour()) + strconv.Itoa(time.Now().Local().Minute()) + strconv.Itoa(time.Now().Local().Second())

	childDirectory := "/tmp/" + strings.ReplaceAll(c.GetSession("uid").(string), " ", "_") + "/" + dateDirectory + "/" + "京东单品毛利润_" + fileName + ".xlsx"

	fileDirectory := "file" + childDirectory

	jdFile := new(models.JDExcelFile)

	jdFile.ERP = erpFile

	jdFile.CostFile = cosFile

	jdFile.FileAccess(fileDirectory)

	if jdFile.Error == nil {
		jdFile.StaticFile = "sugar" + childDirectory
	}

	if jdFile.Error != nil {

		c.Failure(201, jdFile.Error.Error())

		return false
	}

	c.Success(map[string]string{"file": jdFile.StaticFile})

	return true
}
