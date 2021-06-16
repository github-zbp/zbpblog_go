package service

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"strconv"
	"strings"
	"zbpblog_go/errInfo"
	"zbpblog_go/models"
	"zbpblog_go/request"
	"zbpblog_go/response"
	"zbpblog_go/utils"
)

var all_category []orm.Params

// 校验登录态
func CheckLoginStatus(beegoController *beego.Controller) (user orm.Params, res bool){
	base_admin_request := &request.BaseAdminRequest{}
	base_admin_request.BuildReq(beegoController)
	name := base_admin_request.UserName
	sign := base_admin_request.Sign

	if name == "" || sign == ""{
		return user, false
	}

	// 只需校验用户名对应的sign是否正确即可
	m_admin_user := new(models.AdminUser)
	user, err := m_admin_user.GetUserByName(name)

	if err != nil {
		logs.Error("CheckLoginStatus - GetUserByName failed:" + err.Error())
		return user, false
	}

	if user == nil || user["sign"] != sign {
		return user, false
	}

	return user, true
}

// 校验登录态（针对api接口）
func CheckLoginStatusForApi(c *beego.Controller) (res bool){
	_, res = CheckLoginStatus(c)
	if !res {
		resp := new(response.CommonResponse)
		resp.SetError(errInfo.CODE_AUTH_ERR)
		c.Data["json"] = resp
		c.ServeJSON()
	}
	return res
}

// 获取文章数据
func GetAdminBlogs(data map[interface{}]interface{}, req *request.BlogListRequest) {
	// 获取文章
	var (
		blog_data []orm.Params
		rows_num int64
		err error
	)
	m_blogs := new(models.Blogs)
	page, _ := strconv.Atoi(req.Page)
	kw := req.Kw
	pageRows := 50
	pageOffset := pageRows * (page - 1)
	cateNames := []string{}

	if(strings.Trim(kw, " ") != ""){
		blog_data, rows_num, err = m_blogs.GetBlogsByKw(kw, "", pageRows, pageOffset)
	}else{
		blog_data, err = m_blogs.GetBlogsByCateName(cateNames, pageRows, pageOffset)
		rows_num, err = m_blogs.GetRowsNumForBlogs("")
	}

	if err != nil {
		logs.Error("GetAdminBlogs - GetBlogsByCateName failed: " + err.Error())
		return
	}
	data["blogs"] = blog_data

	// 获取分页信息
	if err != nil {
		logs.Error("GetAdminBlogs - GetRowsNumForBlogs failed: " + err.Error())
	}
	data["page_info"] = utils.GetPageInfo(page, pageRows, utils.TranInt64ToInt(rows_num))
	return
}

// 博客添加页数据
func GetAddBlogData(data map[interface{}]interface{}){
	m_category := new(models.Category)
	m_tags := new(models.Tags)

	// 加载所有栏目
	Category,_ = m_category.GetAllCategory("zh_name", "en_name", "id")
//spew.Dump(Category)
	//加载所有标签
	Tags, _ = m_tags.GetAllTags()

	data["categories"] = Category
	data["tags"] = Tags
}

// 博客修改页数据
func GetEditBlogData(data map[interface{}]interface{}, req *request.EditByIdRequest){
	m_blogs := new(models.Blogs)
	m_series := new(models.Series)
	m_tagsblogs := new(models.TagsBlogs)
	GetAddBlogData(data)		// 加载栏目和tag标签信息

	// 获取文章详情
	blog, err := m_blogs.GetBlogById(req.Id)
	if err != nil{
		logs.Error("GetEditBlogData - GetBlogById failed:" + err.Error())
		return
	}
	data["blog"] = blog

	// 获取文章相关联的栏目
	for _, category := range Category {
		if blog["tid"] == category["id"] {
			data["category"] = category
		}
	}

	// 加载文章标签
	var blog_ids = []interface{}{req.Id}
	tags_blogs, err := m_tagsblogs.GetTBByBlogIds(blog_ids)
	if err != nil{
		logs.Error("GetEditBlogData - GetTBByBlogIds failed:" + err.Error())
		return
	}
	tag_ids := utils.GetFieldSliceFromData(tags_blogs, "tags_id")
	data["tag_ids"] = tag_ids

	// 加载文章系列
	series, err := m_series.GetSeriesById(blog["series_id"].(string))
	if err != nil{
		logs.Error("GetEditBlogData - GetSeriesById failed:" + err.Error())
		return
	}
	data["series"] = series
}

// 编辑文章
func EditBlog(req *request.BlogDoAddEditRequest) (err error){
	m_blog := new(models.Blogs)
	req.SetEditData(m_blog)

	orm := models.Orm
	orm.Begin()		// 开启事务
	defer func() {
		if p := recover(); p != nil{		// 如果发生内部异常则回滚
			orm.Rollback()
			logs.Error(fmt.Sprintf("error: %v", p))
		}
	}()

	// 修改文章
	edit_fields := []string{"Title", "Tid", "Keywords", "Description", "Thumb", "Writer", "Content", "SeriesId", "Type", "UpdateTime"}
	_, err = models.Orm.Update(m_blog, edit_fields...)
	if err != nil {
		logs.Error("EditBlog - Update failed:" + err.Error())
		orm.Rollback()
		return err
	}

	// 绑定标签
	insert_rows, succ_rows, err := models.BindTagsBlog(m_blog.Id, req.TagsId)
	if err != nil{
		errmsg := fmt.Sprintf("EditBlog - BindTagsBlog failed, rows success: %d, rows to insert: %d, reason:%s", succ_rows, insert_rows, err.Error())
		logs.Error(errmsg)
		orm.Rollback()
		return err
	}

	orm.Commit()
	return err
}

// 栏目列表
func GetCategoryAdmin(data map[interface{}]interface{}){
	m_cate := new(models.Category)
	m_blog := new(models.Blogs)
	cates, err := m_cate.GetAllCategory2()
	if err != nil{
		logs.Error("GetCategoryAdmin - GetAllCategory2 failed:" + err.Error())
	}

	// 获取每个栏目下的文章数量
	blog_num_info, err := m_blog.GetBlogsNumGroupByTid()
	if err != nil{
		logs.Error("GetCategoryAdmin - GetBlogsNumGroupByTid failed:" + err.Error())
	}
	blog_num_info_map := utils.TranDataSlice2Map(blog_num_info, "tid")

	data["cates"] = cates
	data["blog_nums"] = blog_num_info_map
}

// 单个栏目
func GetSingleCateogryAdmin(data map[interface{}]interface{}, req *request.EditByIdRequest){
	m_cate := new(models.Category)
	cate, err := m_cate.GetCategoryById(req.Id)
	if err != nil{
		logs.Error("GetSingleCateogryAdmin - GetCategoryById failed:" + err.Error())
	}

	cates, err := m_cate.GetAllCategory2()
	if err != nil{
		logs.Error("GetSingleCateogryAdmin - GetAllCategory2 failed:" + err.Error())
	}

	data["cate"] = cate
	data["cates"] = cates
}

// 显示站点信息
func ShowSiteInfo(data map[interface{}]interface{}){
	m_site := new(models.Site)
	data["site"] = m_site.GetSeo()
}

// 获取标签页列表
func GetTagsListsAdmin(data map[interface{}]interface{}, req *request.PageRequest){
	m_tags := new(models.Tags)
	page := req.Page
	pageRows := req.P
	pageOffset := pageRows * (page - 1)

	//var err error
	tags, err := m_tags.GetTagsPage(pageRows, pageOffset)
	if err != nil{
		logs.Error("GetTagsListsAdmin - GetTagsPage failed:" + err.Error())
	}

	totalRows, _ := m_tags.GetTagsNum()

	data["page_info"] = utils.GetPageInfo(page, pageRows, utils.TranInt64ToInt(totalRows))
	data["tags"] = tags
}

// 获取标签详情
func GetAdminTagDetail(data map[interface{}]interface{}, req *request.EditByIdRequest){
	m_tags := new(models.Tags)
	tag, err :=  m_tags.GetTagsById(req.Id)
	if err != nil {
		logs.Error("GetAdminTagDetail - GetTagsById failed:" + err.Error())
	}
	data["tag"] = tag
}

// 添加标签
func AddTag(req *request.TagsDoAddEditRequest) {
	m_tags := new(models.Tags)
	req.SetAddEditData(m_tags)
	_, err := models.Orm.Insert(m_tags)
	if err != nil{
		logs.Error("AddTag - Insert failed:" + err.Error())
	}
}

// 编辑标签
func EditTag(req *request.TagsDoAddEditRequest) {
	m_tags := new(models.Tags)
	req.SetAddEditData(m_tags)
	_, err := models.Orm.Update(m_tags, req.EditField()...)
	if err != nil{
		logs.Error("EditTag - Update failed:" + err.Error())
	}
}

// 后台系列列表
func GetSeriesListsAdmin(data map[interface{}]interface{}, req *request.PageRequest){
	m_series := new(models.Series)
	page := req.Page
	pageRows := req.P
	pageOffset := pageRows * (page - 1)

	//var err error
	series, err := m_series.GetSeriesPage(pageRows, pageOffset)
	if err != nil{
		logs.Error("GetSeriesListsAdmin - GetSeriesPage failed:" + err.Error())
	}

	totalRows, _ := m_series.GetSeriesNum()

	data["page_info"] = utils.GetPageInfo(page, pageRows, utils.TranInt64ToInt(totalRows))
	data["series"] = series
}

// 后台系列详情
func GetAdminSeriesDetail(data map[interface{}]interface{}, req *request.EditByIdRequest){
	m_series := new(models.Series)
	series, err :=  m_series.GetSeriesById(req.Id)
	if err != nil {
		logs.Error("GetAdminSeriesDetail - GetSeriesById failed:" + err.Error())
	}
	data["series"] = series
}

// 添加系列
func AddSeries(req *request.SeriesDoAddEditRequest) {
	m_series := new(models.Series)
	req.SetAddEditData(m_series)
	_, err := models.Orm.Insert(m_series)
	if err != nil{
		logs.Error("AddSeries - Insert failed:" + err.Error())
	}
}

// 编辑系列
func EditSeries(req *request.SeriesDoAddEditRequest) {
	m_series := new(models.Series)
	req.SetAddEditData(m_series)
	_, err := models.Orm.Update(m_series, req.EditField()...)
	if err != nil{
		logs.Error("EditTag - Update failed:" + err.Error())
	}
}

// 查询所有管理员
func GetUserAdmin(data map[interface{}]interface{}) {
	m_user := new(models.AdminUser)
	users, err := m_user.GetAll()
	if err != nil {
		logs.Error("GetUserAdmin - GetAll:" + err.Error())
	}
	data["users"] = users
}