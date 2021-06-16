package request

import (
	"github.com/astaxie/beego"
)

// request包用于接收和处理请求参数，Request对象命名规范为 "{控制器名}{方法名}Reuqst" ，驼峰命名法

/***	文章详情页Request	***/
type IndexBlogRequest struct {
	BlogId string
}

func (blogRequest *IndexBlogRequest) BuildReq(c *beego.Controller) {
	blogRequest.BlogId = c.Ctx.Input.Param(":blog_id")
}

/***	栏目页Request	***/
type IndexCategoryRequest struct {
	CateName string
	Page     string
}

func (blogRequest *IndexCategoryRequest) BuildReq(c *beego.Controller) {
	blogRequest.CateName = c.Ctx.Input.Param(":cate_name")
	blogRequest.Page = c.GetString(":page", "1")
}

/***     标签文章页Request    ***/
type IndexTagBlogsRequest struct {
	TagId	string
	TagName	string
	Page 	string
	OrderType string
}

func (tagBlogRequest *IndexTagBlogsRequest) BuildReq(c *beego.Controller) {
	tagBlogRequest.TagId = c.Ctx.Input.Param(":tag_id")
	tagBlogRequest.TagName = c.Ctx.Input.Param(":tag_name")
	tagBlogRequest.Page = c.GetString(":page", "1")

	params := c.Ctx.Request.URL.Query()
	if params["order_type"] != nil && len(params["order_type"]) > 0 {
		tagBlogRequest.OrderType = params["order_type"][0]
	}else{
		tagBlogRequest.OrderType = "asc"
	}
}

/***     加载更多文章Request    ***/
type IndexMoreBlogsRequest struct {
	Page 	string
}

func (moreBlogRequest *IndexMoreBlogsRequest) BuildReq(c *beego.Controller) {
	moreBlogRequest.Page = c.GetString(":page", "1")

	params := c.Ctx.Request.URL.Query()
	if params["page"] != nil && len(params["page"]) > 0 {
		moreBlogRequest.Page = params["page"][0]
	}else{
		moreBlogRequest.Page = "1"
	}
}

/***     搜索页Request    ***/
type IndexSearchRequest struct {
	Kw string		// 搜索关键词
	SearchType string	// 搜索类型
	Page string		// 页数
}

func (searchRequest *IndexSearchRequest) BuildReq(c *beego.Controller) {
	params := c.Ctx.Request.URL.Query()
	if params["kw"] != nil && len(params["kw"]) > 0 {
		searchRequest.Kw = params["kw"][0]
	}else{
		searchRequest.Kw = ""
	}

	if params["search_type"] != nil && len(params["search_type"]) > 0 {
		searchRequest.SearchType = params["search_type"][0]
	}else{
		searchRequest.SearchType = ""
	}

	if params["page"] != nil && len(params["page"]) > 0 {
		searchRequest.Page = params["page"][0]
	}else{
		searchRequest.Page = "1"
	}
}