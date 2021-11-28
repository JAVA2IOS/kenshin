package controllers

// 所有的controller都继承当前controller
type DirecotryController struct {
	BaseController
}

// 请求路径 /url/:url
func (c *DirecotryController) Url() {
	action := c.Ctx.Input.Param(":url")

	if action == "jd_" {
		c.RenderContentHtml(JDUploadFileTemplate)
		return
	}
}
