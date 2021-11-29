package controllers

import (
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
	logs.Info("当前类型: ", url)
	if url == "0" {
		c.excuteJDFile()
	}

	c.SuccessMessage("文件上传成功")
}

func (c *FileController) excuteJDFile() bool {

	erpFile := c.GetString("erp")

	if len(erpFile) == 0 {

		c.Failure(201, "找不到当前文件")
		return false
	}

	logs.Info("当前文件:", erpFile)

	return true
}
