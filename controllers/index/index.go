package index

import (
	"fmt"
	"github.com/astaxie/beego/logs"
	"strconv"
	"time"
	"zbpblog_go/models"
	"zbpblog_go/request"
	"zbpblog_go/service"
	"zbpblog_go/utils"
	//"github.com/davecgh/go-spew/spew"
)

type IndexController struct {
	BaseController			// 这里不要写成 *BaseController 因为这会是一个空指针然后报错。除非在全局环境下为IndexController的*BaseController成员赋初值
}

// 首页
func (c *IndexController) Index() {
	// 加载首页数据
	service.GetIndexData(c.Data)

	// 加载全局tags
	tags := service.GetTags(15)
	utils.AddMapInner(c.Data, "global", "blogs_tags", tags)
	c.TplName = "index/index.html"
}

// 栏目页
func (c *IndexController) Category() {
	categoryRequest := &request.IndexCategoryRequest{}
	categoryRequest.BuildReq(beegoController) // 初始话前端参数

	// 加载详情页数据
	service.GetCategoryData(c.Data, categoryRequest)

	// 加载全局tags
	tags := service.GetTags(15)
	utils.AddMapInner(c.Data, "global", "blogs_tags", tags)
	c.TplName = fmt.Sprintf("index/%s.html", service.CurrentCate["list_template"].(string))
}

// 文章详情页
func (c *IndexController) Blog() {
	// 初始话前端参数
	blogRequest := &request.IndexBlogRequest{}
	blogRequest.BuildReq(beegoController)

	// 加载详情页数据
	service.GetBlogDetail(c.Data, blogRequest)

	// 加载全局tags
	tags := service.GetTags(15)
	utils.AddMapInner(c.Data, "global", "blogs_tags", tags)
	c.TplName = fmt.Sprintf("index/%s.html", service.CurrentCate["detail_template"].(string))

	// 更改文章点击数（阅读数）
	service.UpdateViewTimeForBlog(blogRequest.BlogId, c.Ctx.Request)
}

// 标签文章页
func (c *IndexController) TagBlogs(){
	// 初始话前端参数
	tagBlogRequest := &request.IndexTagBlogsRequest{}
	tagBlogRequest.BuildReq(beegoController)

	// 加载标签文章页数据
	service.GetTagBlogsDetail(c.Data, tagBlogRequest)

	// 加载全局tags
	tags := service.GetTags(15)
	utils.AddMapInner(c.Data, "global", "blogs_tags", tags)
	c.TplName = "index/tag_articles.html"
}

// 标签页
func (c *IndexController) Tags(){
	m_tags := new(models.Tags)
	tags, err := m_tags.GetAllTags()
	if err != nil{
		logs.Error("GetAllTags failed: " + err.Error())
	}
	c.Data["tags"] = tags
	c.TplName = "index/tags.html"
}

// 搜索页
func (c *IndexController) Search(){
	// 初始话前端参数
	searchRequest := &request.IndexSearchRequest{}
	searchRequest.BuildReq(beegoController)

	// 加载搜索数据
	service.GetSearchData(c.Data, searchRequest)

	// 加载全局tags
	tags := service.GetTags(15)
	utils.AddMapInner(c.Data, "global", "blogs_tags", tags)
	c.TplName = "index/search.html"
}

// 作者介绍页
func (c *IndexController) About(){
	// 获取当前年份
	date := time.Now()
	year_diff := date.Year() - 2018
	cont := `<h3>个人信息</h3><b>名字:张柏沛</b><br/>
<b>简称:沛沛</b><br/>
<b>职业:PHP程序员</b><br/>
<b>爱好:跳舞机,敲代码</b><br/>
<b>所在地：深圳</b><br/>
<b>语言:PHP,Python,Go</b><br/>

<h3>关于沛沛</h3>我是一个毕业`+ strconv.Itoa(year_diff) +`年的php程序员,但是我的大学专业是国际贸易,大二的时候才发现自己喜欢和适合编程,于是开始了自学,所以我其实是一个半路出家的程序员。

毕业之后，我开始从事Web开发，涉及到的一些PHP，Mysql，js，Linux，Go，vue等前后端的技术，写过几个网站类的小程序，以及会使用Python写爬虫。

当知识积累到了一定的程度之后，终于下决心做一个自己的网站，于是就有了这个个人技术博客。一个是为了记录下自己所学的知识然后和大家分享；一个是因为是自己喜欢做的事情。我始终认为做自己喜欢的事情能更有动力和热情，进步的更快，而这也是我当初没有选择经济类的工作而选择去做一个程序员的原因，即使整个过程我遇到了很多难题。所以希望大家能坚持所爱，勿忘初心。
 
最后，也希望这个博客网站能给大家带来帮助，当然我写的可能不一定对，错了的大家请忽略或者请大家指正，你们的肯定也是我最大的动力。
</h3>
`
	c.Data["about"] = cont
	c.TplName = "index/about.html"
}