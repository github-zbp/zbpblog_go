package service

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"os"
	"strconv"
	"strings"
	"time"
	"zbpblog_go/errInfo"
	"zbpblog_go/models"
	"zbpblog_go/request"
	"zbpblog_go/response"
	"zbpblog_go/utils"
)

// 加载更多文章
func GetMoreBlogsData(data map[interface{}]interface{}, moreBlogRequest *request.IndexMoreBlogsRequest){
	current_page,_ := strconv.Atoi(moreBlogRequest.Page)
	per_page_rows := 15		// 每页多少条
	page_offset := per_page_rows * (current_page - 1)
	m_blogs := new(models.Blogs)
	m_category := new(models.Category)
	resp := new(response.MoreBlogsResponse)

	// 获取文章总数
	total_num_int64, _ := m_blogs.GetRowsNumForBlogs("")
	total_num := utils.TranInt64ToInt(total_num_int64)
	page_info := utils.GetPageInfo(current_page, per_page_rows, total_num)

	// 获取文章信息
	blogs, _, err := m_blogs.GetBlogsPage([]string{"id", "create_time", "thumb", "description", "title", "tid"}, page_offset, per_page_rows)
	if err != nil{
		logs.Error("GetMoreBlogsData - GetBlogsPage failed: " + err.Error())
		resp.SetError(errInfo.CODE_MORE_BLOG_FAILED)
	}

	category, _ := m_category.GetAllCategory("en_name", "zh_name", "id")
	category_map := utils.GetModelDataMap(category, "id")
	resp.Set(blogs, category_map, page_info)
	data["json"] = resp
}

// 加载标签图片
func LoadTagsPic(data map[interface{}]interface{}, req *request.BlogGetTagsPicRequest){
	resp := new(response.GetTagsPicResponse)
	m_tags := new(models.Tags)

	tag_ids := utils.TranSliceString2Interface(req.TagIds)
	tags, err := m_tags.GetTagsByIds(tag_ids, "thumb")
	if err != nil {
		logs.Error("LoadTagsPic - GetTagsByIds failed:" + err.Error())
		resp.SetError(errInfo.CODE_LOAD_TAG_PIC_FAIL)
	}
	tag_thumbs := utils.GetFieldSliceFromData(tags, "Thumb")
	resp.SetData(tag_thumbs)
	data["json"] = resp
}

// 搜索系列
func SearchSeriesByName(data map[interface{}]interface{}, req *request.SearchSeriesRequest){
	resp := new(response.CommonResponse)
	m_series := new(models.Series)

	series, err := m_series.GetSeriesByName(req.Search)
	if err != nil {
		logs.Error("SearchSeriesByName - GetSeriesByName failed:" + err.Error())
		resp.SetError(errInfo.CODE_SEARCH_SERIES_FAIL)
	}
	resp.SetData(series)
	data["json"] = resp
}

// 后台添加博客文章
func AddBlog(data map[interface{}]interface{}, req *request.BlogDoAddEditRequest){
	resp := new(response.CommonResponse)
	m_blog := new(models.Blogs)
	req.SetAddData(m_blog)

	orm := models.Orm
	orm.Begin()		// 开启事务
	defer func() {
		if p := recover(); p != nil{		// 如果发生内部异常则回滚
			orm.Rollback()
			logs.Error(fmt.Sprintf("error: %v", p))
		}
		data["json"] = resp
	}()

	// 写入文章表
	id, err := models.Orm.Insert(m_blog)
	if err != nil {
		logs.Error("AddBlog - Insert failed:" + err.Error())
		orm.Rollback()
		resp.SetError(errInfo.CODE_INSERT_BLOG_FAILED)
		return
	}

	// 绑定标签
	blog_id := utils.TranInt64ToInt(id)
	insert_rows, succ_rows, err := models.BindTagsBlog(blog_id, req.TagsId)
	if err != nil{
		errmsg := fmt.Sprintf("AddBlog - BindTagsBlog failed, rows success: %d, rows to insert: %d, reason:%s", succ_rows, insert_rows, err.Error())
		logs.Error(errmsg)
		resp.SetError(errInfo.CODE_BIND_BLOG_TAGS_FAILED)
		orm.Rollback()
		return
	}

	orm.Commit()
}


// 删除操作
func Del(req *request.DelRequest) (err error){
	if !IsCorrectModel(req.Model){
		return errInfo.CODE_PASS_WRONG_MODEL
	}
	if req.Id != ""{
		err = models.Delete(req.Model, req.Id)
	}else{
		err = models.DeleteMulti(req.Model, req.Ids)
	}

	return err
}

// 判断前端传递的模型名称是否正确
func IsCorrectModel(model string) bool{
	return utils.InSlice(model, utils.TranSliceString2Interface(models.Models))
}

// 添加栏目
func AddCategory(req *request.CategoryDoAddEditRequest) (errCode errInfo.ErrCode){
	m_cate := new(models.Category)
	req.SetAddData(m_cate)
	orm := models.Orm
	_, err := orm.Insert(m_cate)
	if err != nil {
		logs.Error("AddCategory - Insert failed:" + err.Error())
		return errInfo.CODE_ADD_CATEGORY_FAILED
	}
	return errInfo.CODE_OK
}

// 修改栏目
func EditCategory(req *request.CategoryDoAddEditRequest) (errCode errInfo.ErrCode){
	m_cate := new(models.Category)
	req.SetEditData(m_cate)
	orm := models.Orm
	_, err := orm.Update(m_cate)
	if err != nil {
		logs.Error("EditCategory - Update failed:" + err.Error())
		return errInfo.CODE_EDIT_CATEGORY_FAILED
	}
	return errInfo.CODE_OK
}

// 更改状态操作
func TranStatus(req *request.TranStatusRequest) (err error){
	if !IsCorrectModel(req.Model){
		return errInfo.CODE_PASS_WRONG_MODEL
	}

	return models.TranStatus(req.Model, req.Id, req.Status)
}

// 排序
func Sort(req *request.SortRequest) (err error){
	if !IsCorrectModel(req.Model){
		return errInfo.CODE_PASS_WRONG_MODEL
	}

	return models.Sort(req.Model, req.Sort)
}

// 修改站点信息
func EditSite(req *request.SiteEditRequest) (err errInfo.ErrCode){
	m_site := new(models.Site)
	req.SetEditData(m_site)
	edit_fields := req.EditField()
	_, update_err := models.Orm.Update(m_site, edit_fields...)
	if update_err != nil{
		logs.Error("EditSite - Update failed:" + update_err.Error())
		return errInfo.CODE_EDIT_SITE_FAILED
	}

	return
}

// 后端获取所有博客缩略图
func GetBlogThumbs(req *request.PageRequest) (interface{}, errInfo.ErrCode){
	m_blog := new(models.Blogs)
	page_offset := req.P * (req.Page - 1)

	thumbs, model_err := m_blog.GetBlogsThumbs(req.P, page_offset)
	if model_err != nil {
		logs.Error("GetBlogThumbs - GetBlogsThumbs:" + model_err.Error())
		return thumbs, errInfo.CODE_GET_BLOG_THUMBS_FAILED
	}

	return thumbs, errInfo.CODE_OK
}

// 添加用户
func AddUser(req *request.UserDoAddEditRequest) (err errInfo.ErrCode){
	m_user := new(models.AdminUser)
	req.SetAddEditData(m_user)
	_, insert_err := models.Orm.Insert(m_user)
	if insert_err != nil {
		logs.Error("AddUser - Insert faled:" + insert_err.Error())
		return errInfo.CODE_ADD_USER_FAILED
	}

	return errInfo.CODE_OK
}

func PostBaidu() (resp_string string, errCode errInfo.ErrCode){
	m_blog := new(models.Blogs)

	// 获取创建时间在一天之前的文章
	blog_ids, err := m_blog.GetBlogsIdPostBaidu()
	if err != nil{
		logs.Error("PostBaidu - GetBlogsIdPostBaidu failed: " + err.Error())
		return "", errInfo.CODE_GET_BLOGS_POST_BD_FAILED
	}

	if len(blog_ids) == 0{
		return "", errInfo.CODE_GET_BLOGS_POST_BD_EMPTY
	}

	var post_urls []string
	for _, blog_id := range blog_ids{
		post_url := fmt.Sprintf("http://%s/blog-%s.html", beego.AppConfig.String("domain"), blog_id)
		post_urls = append(post_urls, post_url)
	}

	post_api := "http://data.zz.baidu.com/urls?site=http://www.zbpblog.com&token=6StPClSC0W0ItoZG"
	curl := utils.Curl{
		Url : post_api,
		Data:strings.Join(post_urls, "\n"),
		ContentType: "text/plain",
	}
	resp_string, err = curl.PostForString()
	if err != nil {
		logs.Error("PostBaidu - curl Post failed: " + err.Error())
		return "", errInfo.CODE_GET_BLOGS_POST_BD_FAILED
	}

	if resp_string != ""{
		fs, err := os.OpenFile(beego.AppConfig.String("bd_push_blogs_path"), os.O_WRONLY | os.O_APPEND | os.O_CREATE, 0777)
		if err == nil{
			defer fs.Close()
		}
		curTime := time.Now().Format(beego.AppConfig.String("format_time_str"))
		fs.WriteString(fmt.Sprintf("%s : %s\n", curTime, resp_string))
	}
	return resp_string, errInfo.CODE_OK
}