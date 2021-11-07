package controllers

import (
	util "kenshin/util"
	"time"
)

// 所有的controller都继承当前controller
type FileController struct {
	BaseController
}

func (c *FileController) Get() {
}

func (c *FileController) Upload() {
	// f, h, err := c.GetFile("file")

	data := c.GetString("deviceName")

	if len(data) == 0 {
		println("参数获取失败")
		data = time.Now().Format("2021_11_03")
	}

	println("额外参数: ", data)

	go util.ReadExcelDataStream(c.Ctx.ResponseWriter, c.Ctx.Request)
	/*
		file, h, err := c.GetFile("file")
		if err != nil {
			// http.Error(w, err.Error(), http.StatusNoContent)
			println("获取文件失败")
			c.Ctx.WriteString("file fetched failured")
			return
		}

		defer file.Close()


		fileDirectory := "file/tmp/" + data

		_, fileError := util.CreateFileDirectory(fileDirectory)

		if fileError != nil {
			println("文件创建失败", fileError.Error())
			c.Ctx.WriteString("file created failured")
			return
		}

		saveErr := c.SaveToFile("file", fileDirectory+"/"+h.Filename)

		if saveErr != nil {
			println("保存文件失败", saveErr.Error())
			c.Ctx.WriteString("file saved failured")
			return
		}

		println("file: ", h.Filename)
		c.Data["json"] = "file uploaded successfully!"
		c.Ctx.WriteString("upload successfully!") */
}
