package admin

import (
	"zbpblog_go/request"
	"zbpblog_go/service"
)

type TagsController struct{
	BaseController
}

func (c *TagsController) Lists(){
	req := new(request.PageRequest)
	req.BuildReq(beegoController)
	service.GetTagsListsAdmin(c.Data, req)
	c.TplName = "admin/tags/lists.html"
}

func (c *TagsController) Add(){
	c.TplName = "admin/tags/add.html"
}

func (c *TagsController) Detail(){
	req := new(request.EditByIdRequest)
	req.BuildReq(beegoController)
	service.GetAdminTagDetail(c.Data, req)
	c.TplName = "admin/tags/edit.html"
}

func (c *TagsController) DoAdd(){
	req := new(request.TagsDoAddEditRequest)
	req.BuildReq(beegoController)
	service.AddTag(req)
	c.Ctx.Redirect(302, "/zbpadmin/tags/lists")
}

func (c *TagsController) DoEdit(){
	req := new(request.TagsDoAddEditRequest)
	req.BuildReq(beegoController)
	service.EditTag(req)
	c.Ctx.Redirect(302, "/zbpadmin/tags/edit?id=" + req.Id)
}