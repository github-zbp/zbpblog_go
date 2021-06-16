package admin

import (
	"zbpblog_go/models"
	"zbpblog_go/request"
	"zbpblog_go/service"
)

type CategoryController struct{
	BaseController
}

func (c *CategoryController) List(){
	service.GetCategoryAdmin(c.Data)
	c.TplName = "admin/index/cate_list.html"
}

func (c *CategoryController) Detail(){
	req := new(request.EditByIdRequest)
	req.BuildReq(beegoController)
	service.GetSingleCateogryAdmin(c.Data, req)
	c.TplName = "admin/index/edit_category.html"
}

func (c *CategoryController) Add(){
	m_cate := new(models.Category)
	c.Data["cates"], _ = m_cate.GetAllCategory2()
	c.TplName = "admin/index/add_category.html"
}