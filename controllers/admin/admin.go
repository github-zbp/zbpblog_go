package admin

import (
	"zbpblog_go/service"
)

type AdminController struct{
	BaseController
}

func (c *AdminController) Lists(){
	service.GetUserAdmin(c.Data)
	c.TplName = "admin/user/admin_list.html"
}


func (c *AdminController) Add(){
	c.TplName = "admin/user/admin_add.html"
}