package admin

import (
	"zbpblog_go/request"
	"zbpblog_go/service"
)

type SeriesController struct{
	BaseController
}

func (c *SeriesController) Lists(){
	req := new(request.PageRequest)
	req.BuildReq(beegoController)
	service.GetSeriesListsAdmin(c.Data, req)
	c.TplName = "admin/series/lists.html"
}


func (c *SeriesController) Add(){
	c.TplName = "admin/series/add.html"
}

func (c *SeriesController) Detail(){
	req := new(request.EditByIdRequest)
	req.BuildReq(beegoController)
	service.GetAdminSeriesDetail(c.Data, req)
	c.TplName = "admin/series/edit.html"
}

func (c *SeriesController) DoAdd(){
	req := new(request.SeriesDoAddEditRequest)
	req.BuildReq(beegoController)
	service.AddSeries(req)
	c.Ctx.Redirect(302, "/zbpadmin/series/lists")
}

func (c *SeriesController) DoEdit(){
	req := new(request.SeriesDoAddEditRequest)
	req.BuildReq(beegoController)
	service.EditSeries(req)
	c.Ctx.Redirect(302, "/zbpadmin/series/edit?id=" + req.Id)
}