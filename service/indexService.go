package service

import (
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"github.com/garyburd/redigo/redis"
	"net/http"
	"strconv"
	"strings"
	"zbpblog_go/cache"
	"zbpblog_go/models"
	"zbpblog_go/request"
	"zbpblog_go/utils"
)

/*********************************** 控制器获取数据的方法 ************************************/

// 首页
func GetIndexData(data map[interface{}]interface{}) {
	// 获取博客文章
	var err error
	blogs := &models.Blogs{}
	cateNames := []string{}
	blog_data, err := blogs.GetBlogsByCateName(cateNames, 15, 0)
	data["blogs"] = blog_data

	if err != nil {
		logs.Error("GetBlogsByCateName failed: " + err.Error())
		return
	}

	blog_ids := utils.GetFieldSliceFromData(blog_data, "id")
	if len(blog_ids) > 0 {
		data["blogs_tags"] = GetBlogsTagsMap(blog_ids)
	}
}

// 文章详情页
func GetBlogDetail(data map[interface{}]interface{}, blogRequest *request.IndexBlogRequest) {
	m_series := new(models.Series)
	m_blog := new(models.Blogs)

	// 获取文章详情
	blog, err := m_blog.GetBlogById(blogRequest.BlogId)
	if err != nil {
		logs.Error("GetBlogById failed: " + err.Error())
	}

	// 获取文章相关联的栏目
	for _, category := range Category {
		if blog["tid"] == category["id"] {
			CurrentCate = category
			utils.SetSingleGlobalDataForTemplate(data, "curCate", category)
			break
		}
	}

	// 获取文章相关的系列下所有的文章
	series_blogs, err := m_blog.GetBlogsBySeriesId(blog["series_id"].(string))
	if err != nil {
		logs.Error("GetBlogsBySeriesId failed: " + err.Error())
	}

	// 获取文章关联的系列
	series, err := m_series.GetSeriesById(blog["series_id"].(string))
	if err != nil {
		logs.Error("GetSeriesById fail: " + err.Error())
	}

	// 获取文章的标签
	tags := GetTagsForSingleBlog(blog["id"])

	// 更改侧边文章
	hot_blogs := GetHotBlogs([]string{CurrentCate["id"].(string)}, 15)

	// 获取本文章的上一篇和下一篇文章
	data["next_blog"], _ = m_blog.GetSublingBlog(blog["id"].(string), []string{"id", "title"}, true)
	data["prev_blog"], _ = m_blog.GetSublingBlog(blog["id"].(string), []string{"id", "title"}, false)

	data["series_blogs"] = series_blogs
	data["bread_crumb"] = GetBreadCrumbForBlog(blog["title"].(string))
	data["blog"] = blog
	data["series"] = series
	data["tags"] = tags
	SetCurCateInCategory()
	utils.SetSingleGlobalDataForTemplate(data, "hot_blogs", hot_blogs)
}

// 栏目页
func GetCategoryData(data map[interface{}]interface{}, categoryRequest *request.IndexCategoryRequest) {
	m_blogs := new(models.Blogs)
	cate_name := categoryRequest.CateName
	page, _ := strconv.Atoi(categoryRequest.Page)
	pageRows := 20
	pageOffset := pageRows * (page - 1)

	// 获取当前栏目
	GetCurrentCate(categoryRequest.CateName)

	// 获取栏目对应的文章和标签
	blogs, err := m_blogs.GetBlogsByCateName([]string{cate_name}, pageRows, pageOffset)
	if err != nil {
		logs.Error("GetCategoryData - GetBlogsByCateName failed: " + err.Error())
	}

	// 获取文章标签
	blog_ids := utils.GetFieldSliceFromData(blogs, "id")
	if len(blog_ids) > 0 {
		data["blogs_tags"] = GetBlogsTagsMap(blog_ids)
	}

	// 更改seo
	Seo = SetSeo(CurrentCate["title"].(string), CurrentCate["keywords"].(string), CurrentCate["description"].(string))

	// 更改侧边文章
	hot_blogs := GetHotBlogs([]string{CurrentCate["id"].(string)}, 15)
	data["bread_crumb"] = GetBreadCrumbForCategory()
	data["blogs"] = blogs

	// 获取分页信息
	rows_num, err := m_blogs.GetRowsNumForBlogs(CurrentCate["id"].(string))
	if err != nil {
		logs.Error("GetCategoryData - GetRowsNumForBlogs failed: " + err.Error())
	}
	data["page_info"] = utils.GetPageInfo(page, pageRows, utils.TranInt64ToInt(rows_num))

	// 更新模板中的全局变量
	utils.SetSingleGlobalDataForTemplate(data, "hot_blogs", hot_blogs)
	utils.SetSingleGlobalDataForTemplate(data, "curCate", CurrentCate)
	utils.SetSingleGlobalDataForTemplate(data, "seo", Seo)
	SetCurCateInCategory()
}

// 标签文章页
func GetTagBlogsDetail(data map[interface{}]interface{}, tagBlogsRequest *request.IndexTagBlogsRequest){
	page, _ := strconv.Atoi(tagBlogsRequest.Page)
	pageRows := 15
	pageOffset := pageRows * (page - 1)
	m_tag := new(models.Tags)
	m_tags_blogs := new(models.TagsBlogs)
	m_blogs := new(models.Blogs)

	// 加载栏目信息
	tag, err := m_tag.GetTagsById(tagBlogsRequest.TagId)
	if err != nil {
		logs.Error("GetTagBlogsDetail - GetTagsById failed: " + err.Error())
	}

	// 加载栏目下所有文章信息
	tag_ids := []interface{}{tag["id"]}
	tag_blog_ids, err := m_tags_blogs.GetTBByTagIds(tag_ids)
	if err != nil {
		logs.Error("GetTagBlogsDetail - GetTBByTagIds failed: " + err.Error())
	}
	blog_ids := utils.TranSliceInterface2String(utils.GetFieldSliceFromData(tag_blog_ids, "blogs_id"))
	blog_nums,_ := m_blogs.GetRowsNumForTags(blog_ids)
	blog_nums_int := utils.TranInt64ToInt(blog_nums)
	blogs, err := m_blogs.GetBlogsByIds(blog_ids, "create_time " + tagBlogsRequest.OrderType, strconv.Itoa(pageRows), strconv.Itoa(pageOffset))

	// 更改seo
	Seo = SetSeo("", "", SetTagBlogsDescription(tagBlogsRequest.TagName))

	// 获取分页信息
	data["page_info"] = utils.GetPageInfo(page, pageRows, blog_nums_int)
	data["blogs"] = blogs
	data["tag"] = tag
	data["bread_crumb"] = GetBreadCrumbForTag(tagBlogsRequest.TagId, tagBlogsRequest.TagName)
	data["arts_num"] = blog_nums
	data["order_type"] = tagBlogsRequest.OrderType
	data["category_map"] = utils.GetModelDataMap(Category, "id")
	utils.SetSingleGlobalDataForTemplate(data, "seo", Seo)
}

// 搜索页
func GetSearchData(data map[interface{}]interface{}, searchRequest *request.IndexSearchRequest){
	search_type := searchRequest.SearchType
	kw := searchRequest.Kw
	m_blogs := new(models.Blogs)
	page, _ := strconv.Atoi(searchRequest.Page)
	pageRows := 20
	pageOffset := pageRows * (page - 1)

	// 获取搜索数据
	if strings.TrimSpace(kw) == "" {
		return
	}
	blogs, blog_nums, err := m_blogs.GetBlogsByKw(kw, search_type, pageRows, pageOffset, "id", "title", "tid", "thumb","description", "create_time")
	if err != nil{
		logs.Error("GetSearchData - GetBlogsByKw failed: " + err.Error())
	}

	// 获取搜索数据的tag
	blog_ids := utils.GetFieldSliceFromData(blogs, "id")
	blogs_tags_info := GetBlogsTagsMap(blog_ids)

	// 拼接分页连接
	params_map := make(map[string]string)
	params_map["search_type"] = search_type
	//params_map["page"] = searchRequest.Page
	params_map["kw"] = kw

	// 获取分页信息
	blog_nums_int := utils.TranInt64ToInt(blog_nums)

	data["page_info"] = utils.GetPageInfo(page, pageRows, blog_nums_int)
	//data["page_link"] = "/search?" + utils.BuildReqStr(params_map)
	data["blogs"] = blogs
	data["blogs_tags"] = blogs_tags_info
	utils.SetSingleGlobalDataForTemplate(data, "kw", strings.TrimSpace(kw))
}

/*********************************** 辅助方法 ************************************/

// 加载栏目信息
func GetCategory() (cates []orm.Params) {
	category := &models.Category{}
	cates, err := category.GetAllCategory() // 加载栏目
	if err != nil {
		logs.Error("Get all category fail:" + err.Error())
	}
	Category = cates
	return cates
}

// 加载SEO
func GetSeo() orm.Params {
	site := models.Site{}
	Seo = site.GetSeo()
	return Seo
}

// 加载当前栏目
func GetCurrentCate(curCate string) orm.Params {
	var err error
	if curCate != "" {
		category := &models.Category{}
		CurrentCate, err = category.GetCategoryByName(curCate)
		if err != nil {
			logs.Error("Get current category fail:" + err.Error())
		}
		return CurrentCate
	}

	return CurrentCate
}

// 加载热门文章
func GetHotBlogs(cate_ids []string, num int) (hot_blogs []orm.Params) {
	blog := &models.Blogs{}
	hot_blogs, err := blog.GetHotBlogs(cate_ids, num)
	if err != nil {
		logs.Error("GetHotBlogs failed:" + err.Error())
	}
	return hot_blogs
}

// 加载所有友情链接
func GetAllFriendLinks() (friend_links []orm.Params) {
	return friend_links
}

// 加载标签
func GetTags(limit int) (tags []orm.Params) {
	tag := &models.Tags{}
	tags, err := tag.GetTags(limit)
	if err != nil {
		logs.Error("GetHotBlogs failed:" + err.Error())
	}
	return tags
}

// 加载详情页breadcrumb  该方法只会在index控制器下调用
func GetBreadCrumbForBlog(blog_title string) (bread string) {
	bread = GetBreadCrumbForCategory()
	bread += fmt.Sprintf(" > <span>%s</span>", blog_title)
	return bread
}

// 加载栏目页breadcrumb 该方法只会在index控制器下调用
func GetBreadCrumbForCategory() (bread string) {
	bread = "<a href='/'>首页</a> > "
	bread += fmt.Sprintf("<a href='/cate-%s'>%s</a>", CurrentCate["en_name"], CurrentCate["zh_name"])
	return bread
}

// 加载标签页breadcrumb 该方法只会在index控制器下调用
func GetBreadCrumbForTag(tagId string, tagName string) (bread string) {
	bread = "<a href='/'>首页</a> > "
	bread += fmt.Sprintf("<a href='/tag%s-%s'>%s</a>", tagId, tagName, tagName)
	return bread
}

// 设置Category中的当前栏目
func SetCurCateInCategory() {
	if CurrentCate == nil {
		return
	}

	for i, c := range Category {
		if c["id"] == CurrentCate["id"] {
			c["is_current"] = true
			Category[i] = c
		}
	}
}

// 根据blog_ids获取标签map
func GetBlogsTagsMap(blog_ids []interface{}) map[string][]orm.Params {
	blog_tags_info := make(map[string][]orm.Params)

	// 查询标签和文章中间表数据，以 [blogs_id=>[tag_id1, tag_id2]]的形式返回
	blogs_tags := &models.TagsBlogs{}
	blogs_tags_data, err := blogs_tags.GetTBByBlogIds(blog_ids) //  得到 [["tags_id":1,"blogs_id":1], [], ...]形式的数据
	if err != nil {
		logs.Error("Get tags_blogs failed: " + err.Error())
		return blog_tags_info
	}
	//spew.Dump(blogs_tags_data)
	blogs_tags_map := utils.GetModelFieldDataMap(blogs_tags_data, "blogs_id", "tags_id") // 得到 ["123":[1,2,3], ...]	形式的数据， key是文章id，value是tag的id(无需查表)

	// 获取指定文章的tag的id（无需查表)
	tags_ids := []interface{}{}
	for _, bt := range blogs_tags_data {
		tags_ids = append(tags_ids, bt["tags_id"])
	}
	tags_ids = utils.UniqueSlice(tags_ids)

	// 获取每篇文章的tags并合并到blog_data中
	m_tags := &models.Tags{}
	tags, _ := m_tags.GetTagsByIds(tags_ids)
	tags_map := utils.TranDataSlice2Map(tags, "Id") // [tags_id1:orm.Params, tags_id2:orm.Params, ...]

	// blogs_tags_info是一个key为blog_id，值为tags详情的map
	blogs_tags_info := make(map[string][]orm.Params)
	for blogs_id, tags_ids := range blogs_tags_map {
		for _, tags_id := range tags_ids {
			blogs_tags_info[blogs_id] = append(blogs_tags_info[blogs_id], (tags_map[tags_id.(string)]))
		}
	}

	return blogs_tags_info
}

// 设置栏目页SEO
func SetSeo(title string, keywords string, description string) (seo map[string]interface{}) {
	seo = make(map[string]interface{})

	if title != "" {
		seo["title"] = title
	} else {
		seo["title"] = Seo["title"]
	}

	if keywords != "" {
		seo["keywords"] = keywords
	} else {
		seo["keywords"] = Seo["keywords"]
	}

	if description != "" {
		seo["description"] = description
	} else {
		seo["description"] = Seo["description"]
	}

	return seo
}

// 更新详情页的访问量
func UpdateViewTimeForBlog(blog_id string, req *http.Request) {
	// 客户端远程ip
	ip_port := req.RemoteAddr
	if ip_port == "" {
		return
	}

	ip_port_arr := strings.Split(ip_port, ":")
	if len(ip_port_arr) == 0{
		return
	}

	ip := ip_port_arr[0]
	prefix := "blog_view_"
	key := prefix + blog_id + "_" + ip
	redis_conn := cache.RedisPool.Get()

	has_view, err := redis.Int(redis_conn.Do("get", key))
	if err != nil{
		logs.Error("Get has view blog failed, key: %s, error: %s", key, err.Error())
	}

	if has_view == 0 {
		m_blog := new(models.Blogs)
		redis_conn.Do("setex", key, 3600 * 24, 1)
		err := m_blog.AddBlogView(blog_id)
		if err != nil {
			logs.Error("AddBlogView failed, error: %s", err.Error())
		}
	}
}

// 设置tag页description
func SetTagBlogsDescription(tagName string) (description string){
	return fmt.Sprintf("本页面是%s博客专题页,为用户展现有关%s方面知识的IT博客文章,与你分享和探讨更多%s相关的知识", tagName, tagName, tagName)
}

// 获取单个博客的标签
func GetTagsForSingleBlog(blog_id interface{}) ([]orm.Params){
	m_tags := new(models.Tags)
	m_tags_blogs := new(models.TagsBlogs)
	var blog_ids []interface{}
	blog_ids = append(blog_ids, blog_id)
	tags_blogs, _ := m_tags_blogs.GetTBByBlogIds(blog_ids)
	tag_ids := utils.GetFieldSliceFromData(tags_blogs, "tags_id")
	tags, err := m_tags.GetTagsByIds(tag_ids)
	if err != nil {
		logs.Error("GetBlogDetail - GetTagsByIds fail: " + err.Error())
	}

	return tags
}