package admin

import "zbpblog_go/service"

type SiteController struct{
	BaseController
}

func (c *SiteController) Show(){
	service.ShowSiteInfo(c.Data)
	c.TplName = "admin/index/edit_site.html"
}