package admin

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"zbpblog_go/service"
)

var beegoController *beego.Controller
var AdminUser orm.Params

type BaseController struct {
	beego.Controller
}

func (c *BaseController) Prepare(){
	beegoController = &c.Controller

	// 校验登录态
	var check_res bool
	AdminUser, check_res = service.CheckLoginStatus(beegoController)
	if !check_res{
		// 指定错误信息并输出错误页面
		c.Ctx.Redirect(302, "/zbpadmin/login")
		return
	}

	// 初始化公共参数
	setGlobalDataForTemplate(c)
}

func setGlobalDataForTemplate(c *BaseController){
	_, ok := c.Data["global"]
	var global map[string]interface{}
	if !ok || c.Data["global"] == nil {
		global = make(map[string]interface{})
	} else {
		global = c.Data["global"].(map[string]interface{})
	}

	global["user_name"] = AdminUser["user_name"]
}