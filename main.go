package main

import (
	_ "kenshin/routers"

	beego "github.com/beego/beego/v2/server/web"
)

func main() {

	// 文件路径下是静态文件 url: ***/sugar/***
	go beego.SetStaticPath("sugar", "file")

	beego.Run()
}
