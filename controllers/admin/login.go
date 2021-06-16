package admin

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/cache"
	"github.com/astaxie/beego/utils/captcha"
	"zbpblog_go/request"
	"zbpblog_go/service"
	"zbpblog_go/utils"
)

var cpt *captcha.Captcha
type LoginController struct {
	beego.Controller
}

func init(){
	store := cache.NewMemoryCache()
	cpt = captcha.NewWithFilter("/captcha/", store)			// 生成验证码的方法要放在init才能生效
	cpt.FieldCaptchaName = "cap"		// 表单中验证码input的name属性值
}

func (c *LoginController) Prepare(){
	beegoController = &c.Controller
}

func (c *LoginController) Login(){
	c.TplName = "admin/login/index.html"
}

func (c *LoginController) DoLogin(){
	redirect := "/zbpadmin/login"
	loginRequest := new(request.LoginRequest)
	loginRequest.BuildReq(&c.Controller)
	is_valid, err := loginRequest.Valid(&c.Controller, cpt)	// 校验参数是否正确

	if !is_valid {
		utils.SetErrInfoForTemplate(beegoController, err.Error(), redirect)
		return
	}

	res, err := service.ExecuteLogin(beegoController, loginRequest)
	if !res {
		utils.SetErrInfoForTemplate(beegoController, err.Error(), redirect)
		return
	}else{
		c.Ctx.Redirect(302, "/zbpadmin")
	}
}

func (c *LoginController) Logout(){
	c.Ctx.SetCookie("user_name", "", -1)
	c.Ctx.SetCookie("sign", "", -1)
	c.Ctx.Redirect(302, "/zbpadmin/login")
}