package routers

import (
	"kenshin/controllers"

	beego "github.com/beego/beego/v2/server/web"
)

func init() {

	/* 开启websocket */

	beego.Router("/ws", &controllers.KenShinSocketController{})

	beego.Router("/", &controllers.BaseController{})

	// 一个路径配置一个controller

	//文件
	beego.Router("/file/upload/xlsx", &controllers.FileController{}, "post:Upload")

	// go kenshinUtil.CreatNewExcelFile()
}
