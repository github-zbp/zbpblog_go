package admin

import (
	"zbpblog_go/request"
	"zbpblog_go/service"
)

type BlogController struct {
	BaseController
}

func (c *BlogController) List(){
	req := new(request.BlogListRequest)
	req.BuildReq(beegoController)

	service.GetAdminBlogs(c.Data, req)
	c.TplName = "admin/blogs/blog_list.html"
}

func (c *BlogController) Add(){
	service.GetAddBlogData(c.Data)
	c.TplName = "admin/blogs/add_blogs.html"
}

func (c *BlogController) Edit(){
	req := new(request.EditByIdRequest)
	req.BuildReq(beegoController)

	service.GetEditBlogData(c.Data, req)
	c.TplName = "admin/blogs/edit_blogs.html"
}

func (c *BlogController) DoEdit(){
	req := new(request.BlogDoAddEditRequest)
	req.BuildReq(beegoController)

	service.EditBlog(req)
	c.Ctx.Redirect(302, "/zbpadmin/blogs/edit?id=" + req.Id)
}