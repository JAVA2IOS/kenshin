package main

import (
	_ "kenshin/routers"
	"strings"

	"github.com/beego/beego/v2/server/web/context"

	beego "github.com/beego/beego/v2/server/web"
)

func main() {

	beego.BConfig.WebConfig.AutoRender = false

	beego.InsertFilter("/*", beego.BeforeRouter, func(ctx *context.Context) {

		if strings.Contains(ctx.Request.URL.Path, "auth/login") {
			return
		}

		if ok := ctx.Input.Session("uid"); ok == nil {
			ctx.Redirect(302, "/auth/login")
		}
	})

	// 文件路径下是静态文件 url: ***/sugar/***
	go beego.SetStaticPath("sugar", "file")

	beego.Run()
}
