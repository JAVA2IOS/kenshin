package controllers

import (
	"kenshin/models"
	kenshinUtil "kenshin/util"
	"net/http"
	"strings"

	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
)

type TemplateName string

const (
	HeaderTemplate       TemplateName = "header.tpl"      // 顶部模板
	SideBarTemplate      TemplateName = "slider.tpl"      // 左边导航模板
	MainTemplate         TemplateName = "main.html"       // 主题内容默认模板
	JDUploadFileTemplate TemplateName = "fileUpload.html" // 京东xlsx文件上传模板
)

// 请求回调返回参数
type SweetyResponse struct {
	Code    int
	Message string
	Data    interface{}
}

/* 所有的controller都继承当前controller */
type BaseController struct {
	beego.Controller
}

// 渲染内容模板
func (baseController *BaseController) RenderContentHtml(tmplate TemplateName) {
	baseController.TplName = string(tmplate)
	baseController.Layout = "layoutContent.html"
	baseController.LayoutSections = make(map[string]string)
	baseController.LayoutSections["Header"] = string(HeaderTemplate)
	baseController.LayoutSections["Slider"] = string(SideBarTemplate)
	baseController.Render()
}

// 初始化默认参数配置
func (baseController *BaseController) RenderDefaultConfig() {
	baseController.Data["SiteName"] = "0NE DATE"
	baseController.Data["AppName"] = "0NEDATE	FILE"

	baseController.Data["Uid"] = baseController.GetSession("uid")

	navMap := models.NaviActions()

	navInnerHtml := ""
	for _, value := range navMap {
		navInnerHtml += "<dd><a href=\"/url/" + value.Action + "\">" + value.Name + "</a></dd>"
	}

	baseController.Data["Nav"] = navInnerHtml
}

func (baseController *BaseController) Prepare() {
	baseController.RenderDefaultConfig()
}

func (c *BaseController) Get() {

	c.RenderContentHtml(MainTemplate)
}

// 成功回调参数
func (c *BaseController) Success(data interface{}) {

	sweetyResponse := SweetyResponse{
		Code:    http.StatusOK,
		Message: "",
		Data:    data,
	}

	c.Data["json"] = sweetyResponse

	c.ServeJSON()
}

// 失败回调
func (c *BaseController) Failure(code int, message string) {

	sweetyResponse := SweetyResponse{
		Code:    code,
		Message: message,
		Data:    nil,
	}

	c.Data["json"] = sweetyResponse

	c.ServeJSON()
}

// 登录认证
func (c *BaseController) Auth() {

	c.TplName = "index.html"
	c.Render()
}

// 登录跳转
func (c *BaseController) ValidateAuth() {

	uid := c.GetString("uid")

	uid = strings.TrimSpace(uid)

	if len(uid) == 0 {
		c.Failure(500, "请输入用户名")
		return
	}

	// 创建根目录

	go func() {
		fileDirectory := "/file/tmp/" + strings.ReplaceAll(uid, " ", "_")

		_, err := kenshinUtil.CreateFileDirectory(fileDirectory)

		if err != nil {
			logs.Warn("目录创建失败:", err.Error())
		} else {
			logs.Info("目录创建成功:", fileDirectory)
		}
	}()

	c.SetSession("uid", uid)

	c.Ctx.Redirect(302, "/")
}
