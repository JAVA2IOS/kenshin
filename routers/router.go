package routers

import (
	"kenshin/controllers"

	beego "github.com/beego/beego/v2/server/web"
)

func init() {

	/* 开启websocket */

	beego.Router("/ws", &controllers.KenShinSocketController{})

	beego.Router("/", &controllers.BaseController{})

	beego.Router("/url/:url", &controllers.DirecotryController{}, "get:Url")

	beego.Router("/auth/login", &controllers.BaseController{}, "get:Auth")

	beego.Router("/auth/login", &controllers.BaseController{}, "post:ValidateAuth")
	// 一个路径配置一个controller

	//文件
	beego.Router("/file/upload/xlsx", &controllers.FileController{}, "post:Upload")

	// go kenshinUtil.CreatNewExcelFile()
}
