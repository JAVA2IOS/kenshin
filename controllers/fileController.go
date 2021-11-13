package controllers

import (
	util "kenshin/util"
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

	data := c.GetString("deviceName")

	if len(data) == 0 {
		logs.Info("参数获取失败")
		data = time.Now().Format("2021_11_03")
	}

	logs.Info("额外参数: ", data)

	file, h, err := c.GetFile("file")
	if err != nil {
		// http.Error(w, err.Error(), http.StatusNoContent)
		logs.Info("获取文件失败")
		c.Failure(201, "读取上传文件失败")
		return
	}

	defer file.Close()

	fileDirectory := "file/tmp/" + strings.ReplaceAll(c.GetSession("uid").(string), " ", "_") + "/" + time.Now().Format("2021_11_03") + "/" + h.Filename

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

	c.Success(map[string]string{"xlsx": h.Filename})
}
