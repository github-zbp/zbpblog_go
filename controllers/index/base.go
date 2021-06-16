package index

import (
	"github.com/astaxie/beego"
	"zbpblog_go/service"
	//"github.com/davecgh/go-spew/spew"
	//"zbpblog_go/template_utils"
)

var beegoController *beego.Controller

type BaseController struct{
	beego.Controller
	//Request request.BaseRequest			// 嵌入 request.BaseRequest 接口，目的是为了避免import循环引用
	//Categories []*models.Category
}

func (c *BaseController) Prepare(){
	setGlobalDataForTemplate(c)
	c.Layout = "index/layout.html"		// 设置layout模板文件
	beegoController = &c.Controller
}

// 设置引入模板
func setLayoutSections() map[string]string{
	sections := make(map[string]string)
	sections["Right"] = "right/right.html"
	sections["RightHotBlogs"] = "right/hot_blogs.html"
	return sections
}

// 设置模板中的公共变量
func setGlobalDataForTemplate(c *BaseController){
	_, ok := c.Data["global"]
	var global map[string]interface{}
	if !ok || c.Data["global"] == nil{
		global = make(map[string]interface{})
	}else{
		global = c.Data["global"].(map[string]interface{})
	}
	global["seo"] = service.GetSeo()
	global["category"] = service.GetCategory()		// 加载栏目
	global["search_type"] = ""			// 设置默认值为空
	curCateName := c.GetString("cate_name")	// 当前目录名
	curCate := service.GetCurrentCate(curCateName)
	global["curCate"] = curCate
	global["hot_blogs"] = service.GetHotBlogs([]string{}, 15)		// 获取热门文章（不区分栏目）
	global["friend_links"] = service.GetAllFriendLinks()
	c.Data["global"] = global
}

