package routers

import (
	"github.com/astaxie/beego"
	"zbpblog_go/controllers/api"
	"zbpblog_go/controllers/index"
	"zbpblog_go/controllers/admin"
)

func init() {
	/***  前台界面路由   ***/
	beego.Router("/", &index.IndexController{}, "get:Index")                                         // 首页
	beego.Router("/cate-:cate_name(\\w+)", &index.IndexController{}, "get:Category")                 // 栏目页
	beego.Router("/cate-:cate_name(\\w+)-page:page(\\d+)", &index.IndexController{}, "get:Category") // 栏目页
	beego.Router("/blog-:blog_id.html", &index.IndexController{}, "get:Blog")                        // 详情页
	beego.Router("/tag:tag_id(\\d+)-:tag_name(.*?)-page:page(\\d+)", &index.IndexController{}, "get:TagBlogs")	// 标签下文章页(该路由必须放在下面这个路由之前，否则含page的url会匹配到下面的路由)
	beego.Router("/tag:tag_id(\\d+)-:tag_name(.*?)", &index.IndexController{}, "get:TagBlogs")	// 标签下文章页
	beego.Router("/more_blogs", &api.ApiController{}, "get:MoreBlogs")	// 首页加载更多文章
	beego.Router("/tags", &index.IndexController{}, "get:Tags")		// 标签页
	beego.Router("/search", &index.IndexController{}, "get:Search")		// 搜索页
	beego.Router("/about", &index.IndexController{}, "get:About")		// 关于作者页


	/***   后台界面路由   ***/
	beego.Router("/zbpadmin", &admin.IndexController{}, "get:Index")			// 后台首页
	beego.Router("/zbpadmin/login", &admin.LoginController{}, "get:Login")		// 登录页
	beego.Router("/zbpadmin/dologin", &admin.LoginController{}, "post:DoLogin")	// 执行登录
	beego.Router("/zbpadmin/logout", &admin.LoginController{}, "get:Logout")	// 执行退出

	// 文章管理
	beego.Router("/zbpadmin/blogs/lists", &admin.BlogController{}, "get:List") // 后台文章列表
	beego.Router("/zbpadmin/blogs/add", &admin.BlogController{}, "get:Add")	// 添加文章页面
	beego.Router("/zbpadmin/blogs/edit", &admin.BlogController{}, "get:Edit")	// 修改文章页面
	beego.Router("/zbpadmin/blogs/getTagsPic", &api.ApiController{}, "post:GetTagsPic")	// 加载标签图片接口
	beego.Router("/zbpadmin/series/search", &api.ApiController{}, "get:SearchSeries")  // 搜索系列接口
	beego.Router("/zbpadmin/blogs/doadd", &api.ApiController{}, "post:BlogDoAdd")		// 执行添加文章
	beego.Router("/zbpadmin/blogs/doedit", &admin.BlogController{}, "post:DoEdit")		// 执行编辑文章

	// 栏目管理
	beego.Router("/zbpadmin/index/cate_list", &admin.CategoryController{}, "get:List")	// 栏目列表页
	beego.Router("/zbpadmin/category/add", &admin.CategoryController{}, "get:Add")  // 添加栏目页
	beego.Router("/zbpadmin/category/edit", &admin.CategoryController{}, "get:Detail")  // 栏目详情
	beego.Router("/zbpadmin/category/do_add",&api.ApiController{}, "post:CategoryDoAdd")  // 栏目新增
	beego.Router("/zbpadmin/category/do_edit",&api.ApiController{}, "post:CategoryDoEdit")  // 栏目修改

	// 修改状态（通用）
	beego.Router("/zbpadmin/common/tran_status", &api.ApiController{}, "*:ApiTranStatus")

	// 修改排序（通用）
	beego.Router("/zbpadmin/common/sort", &api.ApiController{}, "*:ApiSort")

	// 删除(通用)
	beego.Router("/zbpadmin/common/del",&api.ApiController{}, "*:ApiDel")

	// 查询博客使用的所有图片（通用）
	beego.Router("/zbpadmin/common/getBlogThumbs",&api.ApiController{}, "*:GetBlogThumbs")

	// 站点设置
	beego.Router("/zbpadmin/site", &admin.SiteController{}, "get:Show")		// 展示站点信息
	beego.Router("/zbpadmin/site/doedit", &api.ApiController{}, "post:SiteDoEdit")		// 修改站点信息

	// 标签页
	beego.Router("/zbpadmin/tags/lists", &admin.TagsController{}, "get:Lists")	// 标签展示
	beego.Router("/zbpadmin/tags/add", &admin.TagsController{}, "get:Add")	// 标签添加页
	beego.Router("/zbpadmin/tags/edit", &admin.TagsController{}, "get:Detail")	// 标签展详情示
	beego.Router("/zbpadmin/tags/doadd", &admin.TagsController{}, "post:DoAdd")	// 执行添加标签
	beego.Router("/zbpadmin/tags/doedit", &admin.TagsController{}, "post:DoEdit")	// 执行修改标签

	// 系列页
	beego.Router("/zbpadmin/series/lists", &admin.SeriesController{}, "get:Lists")	// 系列展示
	beego.Router("/zbpadmin/series/add", &admin.SeriesController{}, "get:Add")	// 标签添加页
	beego.Router("/zbpadmin/series/edit", &admin.SeriesController{}, "get:Detail")	// 标签展详情示
	beego.Router("/zbpadmin/series/doadd", &admin.SeriesController{}, "post:DoAdd")	// 执行添加标签
	beego.Router("/zbpadmin/series/doedit", &admin.SeriesController{}, "post:DoEdit")	// 执行修改标签

	// 管理员页面
	beego.Router("/zbpadmin/user/admin_list", &admin.AdminController{}, "get:Lists")	// 管理员展示
	beego.Router("/zbpadmin/user/admin_add", &admin.AdminController{}, "get:Add")	// 管理员新增页面
	beego.Router("/zbpadmin/user/doadd", &api.ApiController{}, "post:UserDoAdd")	// 管理员新增页面

	// 百度推送接口
	beego.Router("/zbpadmin/api/post_baidu", &api.ApiController{}, "post:PostBaidu")
	//beego.Router("/zbpadmin/api/test_send_post", &api.ApiController{}, "get:TestSendPost")
	//beego.Router("/zbpadmin/api/test_receive_post", &api.ApiController{}, "post:TestReceivePost")
}

