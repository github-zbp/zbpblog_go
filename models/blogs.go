package models

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql" // 记得引入这个驱动
	"strconv"
	"strings"
	"time"
	"zbpblog_go/utils"
)

//+-------------+------------------+------+-----+---------+----------------+
//| Field       | Type             | Null | Key | Default | Extra          |
//+-------------+------------------+------+-----+---------+----------------+
//| id          | int(11)          | NO   | PRI | NULL    | auto_increment |
//| tid         | int(11)          | NO   | MUL | NULL    |                |
//| title       | varchar(255)     | NO   |     | NULL    |                |
//| keywords    | varchar(255)     | YES  |     | NULL    |                |
//| description | text             | YES  |     | NULL    |                |
//| thumb       | varchar(255)     | YES  |     | NULL    |                |
//| writer      | varchar(255)     | YES  |     | NULL    |                |
//| content     | mediumtext       | YES  |     | NULL    |                |
//| status      | int(2)           | NO   |     | 1       |                |
//| series_id   | int(10) unsigned | NO   | MUL | 0       |                |
//| click       | int(11)          | NO   |     | 0       |                |
//| create_time | int(10)          | NO   |     | NULL    |                |
//| update_time | int(11)          | YES  |     | 0       |                |
//| type        | tinyint(4)       | YES  |     | 1       |                |
//| real_click  | int(10) unsigned | YES  |     | 0       |                |
//| is_push     | tinyint(4)       | YES  |     | 0       |                |
//+-------------+------------------+------+-----+---------+----------------+

const CODE_STATUS_YES = 1 // 状态：正常
const CODE_STATUS_NO = 0  // 状态：禁用

const CODE_TYPE_BLOG = 1     // 文章类型：博客
const CODE_TYPE_VIDEO = 2    // 文章类型：视频
const CODE_TYPE_TEMPLATE = 3 // 文章类型：模板

type Blogs struct {
	Id          int    `orm:"column(id)"`
	Tid         int    `orm:"column(tid)"`
	Title       string `orm:"column(title)"`
	Keywords    string `orm:"column(keywords)"`
	Description string `orm:"column(description)"`
	Thumb       string `orm:"column(thumb)"`
	Writer      string `orm:"column(writer)"`
	Content     string `orm:"column(content)"`
	Status      int    `orm:"column(status)"`
	SeriesId    int    `orm:"column(series_id)"`
	Type        int    `orm:"column(type)"`
	RealClick   int    `orm:"column(real_click)"`
	IsPush      int    `orm:"column(is_push)"`
	Click       int    `orm:"column(click)"`
	ModelTimer
	//Series *Series `orm:"rel(fk)"`
	//Category *Category `orm:"rel(fk)"`
}

/* cateName 表示获取哪些栏目的文章，空切片表示部分栏目
 * limit 获取文字的条数限制
 */
func (m_blog *Blogs) GetBlogsByCateName(cate_names []string, limit int, offset int) (blogs []orm.Params, err error) {
	// 根据栏目名称获取栏目id
	var categories []*Category
	var tids []string
	var limit_str string
	if len(cate_names) > 0 {
		Orm.QueryTable("category").Filter("en_name__in", cate_names).All(&categories, "Id") //	All和One方法除第一参以外后面几参表示只查 指定的字段
	}

	for _, category := range categories {
		tids = append(tids, strconv.Itoa(category.Id))
	}

	// 使用原生查询进行关联查询
	var cond string
	if(beego.AppConfig.String("runmode") == "dev"){
		cond = fmt.Sprintf(" where b.type = %d ", CODE_TYPE_BLOG)
	}else{
		cond = fmt.Sprintf(" where b.status = %d and b.type = %d ", CODE_STATUS_YES, CODE_TYPE_BLOG)
	}

	if len(tids) > 0 {
		cond += " and tid in (%s) "
		tids_str := strings.Join(tids, ",")
		cond = fmt.Sprintf(cond, tids_str)
	}
	order_str := " order by create_time desc "
	if offset == 0 {
		limit_str = fmt.Sprintf(" limit %d ", limit)
	} else {
		limit_str = fmt.Sprintf(" limit %d,%d ", offset, limit)
	}
	sql := "select b.id, b.title, b.keywords, b.description, b.thumb, b.writer, b.status, b.click, b.create_time, b.real_click, c.zh_name, c.en_name from blogs b inner join category c on b.tid = c.id " + cond + order_str + limit_str

	_, err = Orm.Raw(sql).Values(&blogs)
	utils.TranTimeFields(&blogs) // 格式化时间
	return blogs, err
}

/**
 * 根据多个id获取文章详情
 * orderby 根据某字段排序，如 "create_time asc"
 */
func (m_blog *Blogs) GetBlogsByIds(ids []string, orderby string, limit string, offset string) (blogs []orm.Params,err error){
	fields := []string{"title", "create_time", "id", "type", "description", "tid", "thumb"}
	fields_str := strings.Join(fields, ",")
	where_str := strings.Join(ids, ",")
	limit_str := fmt.Sprintf("limit %s, %s", offset, limit)
	sql := fmt.Sprintf("select %s from blogs where id in (%s) order by %s %s", fields_str, where_str, orderby, limit_str)

	_, err = Orm.Raw(sql).Values(&blogs)
	utils.TranTimeFields(&blogs) // 格式化时间
	return blogs, err
}

/**
 * 获取热门文章
 * cateIds 表示获取哪些栏目的文章，空切片表示部分栏目
 * limit 获取文字的条数限制
 */
func (m_blog *Blogs) GetHotBlogs(cate_ids []string, limit int) (blogs []orm.Params, err error) {
	where_str := ""
	if len(cate_ids) > 0 {
		cate_ids_str := strings.Join(cate_ids, ", ")
		where_str = fmt.Sprintf("where tid in (%s)", cate_ids_str)
	}
	order_str := " order by real_click desc, create_time desc "
	limit_str := fmt.Sprintf(" limit %d ", limit)
	fields := []string{"title", "create_time", "id", "type"}
	fields_str := strings.Join(fields, ", ")
	sql := fmt.Sprintf("select %s from blogs %s %s %s", fields_str, where_str, order_str, limit_str)

	_, err = Orm.Raw(sql).Values(&blogs)
	utils.TranTimeFields(&blogs) // 格式化时间
	return blogs, err
}

/**
 * 获取单个文章
 * id 文章id
 */
func (m_blog *Blogs) GetBlogById(id string) (blog_detail orm.Params, err error) {
	var blogs []orm.Params
	sql := fmt.Sprintf("select * from blogs where id = %s", id)
	_, err = Orm.Raw(sql).Values(&blogs)
	if len(blogs) > 0 {
		blog_detail = blogs[0]
		utils.TranTimeFieldsForOne(&blog_detail) // 格式化时间
	}
	return blog_detail, err
}

/**
 * 根据系列id获取文章
 * series_id 系列id
 */
func (m_blog *Blogs) GetBlogsBySeriesId(series_id string) (blogs []orm.Params, err error) {
	if series_id == "0" || series_id == "" {
		return
	}
	fields := []string{"id", "title"}
	field_str := fmt.Sprintf("%s", strings.Join(fields, ","))
	sql := fmt.Sprintf("select %s from blogs where series_id = %s", field_str, series_id)
	_, err = Orm.Raw(sql).Values(&blogs)
	return blogs, err
}

/**
 * 获取上一篇或下一篇文章
 * id 当前的文章id
 * fields 要查询的字段
 * is_next 是否是下一篇 true为下一篇，false为上一篇
 */
func (m_blog *Blogs) GetSublingBlog(id string, fields []string, is_next bool) (blog *orm.Params, err error) {
	var cond_str string
	var blogs []orm.Params
	field_str := fmt.Sprintf("%s", strings.Join(fields, ","))
	if is_next {
		cond_str = fmt.Sprintf(" where id > %s order by id asc ", id)
	} else {
		cond_str = fmt.Sprintf(" where id < %s order by id desc ", id)
	}
	sql := fmt.Sprintf("select %s from blogs %s limit 1", field_str, cond_str)
	_, err = Orm.Raw(sql).Values(&blogs)
	if len(blogs) > 0 {
		blog = &blogs[0]
	}
	return blog, err
}

/**
 * 获取栏目页的文章总条数
 */
func (m_blogs *Blogs) GetRowsNumForBlogs(tid string) (num int64, err error) {
	if tid == ""{
		num, err = Orm.QueryTable("blogs").Count()
	}else{
		num, err = Orm.QueryTable("blogs").Filter("tid", tid).Count()
	}
	return num, err
}

/**
 * 根据文章id获取标签页的文章总条数
 */
func (m_blogs *Blogs) GetRowsNumForTags(ids []string) (num int64, err error) {
	num, err = Orm.QueryTable("blogs").Filter("id__in", ids).Count()
	return num, err
}

/**
 * 更改文章阅读数
 */
func (m_blogs *Blogs) AddBlogView(id string) (err error){
	sql := fmt.Sprintf("update blogs set real_click = real_click + 1 where id = %s", id)
	_, err = Orm.Raw(sql).Exec()
	return err
}

/**
 * 获取文章分页信息
 */
func (m_blogs *Blogs) GetBlogsPage(fields []string, offset int, limit int) (blogs []orm.Params, num int64, err error){
	fields_str := strings.Join(fields, ",")
	limit_str := fmt.Sprintf("limit %s, %s", strconv.Itoa(offset), strconv.Itoa(limit))
	sql := fmt.Sprintf("select %s from blogs where type = 1 and status = 1 order by create_time desc %s", fields_str, limit_str)
	num, err = Orm.Raw(sql).Values(&blogs)
	utils.TranTimeFields(&blogs) // 格式化时间
	return blogs, num, err
}

/*
 * 根据关键词获取文章
 */
func (m_blogs *Blogs) GetBlogsByKw(kw string, search_type string, limit int, offset int, fields ...string) (blogs []orm.Params, total_num int64, err error){
	var fields_str string
	var where_str string
	var limit_str string

	if len(fields) > 0{
		fields_str = strings.Join(fields, ",")
	}else{
		fields_str = "*"
	}

	where_str = " title like ? "
	if search_type != ""{
		where_str += fmt.Sprintf(" and type = %s", search_type)
	}

	limit_str = fmt.Sprintf(" limit %d, %d", offset, limit)
	sql := fmt.Sprintf("select %s from blogs where %s order by create_time desc %s", fields_str, where_str, limit_str)

	_, err = Orm.Raw(sql, "%"+kw+"%").Values(&blogs)
	if err != nil {
		return
	}
	utils.TranTimeFields(&blogs)

	total_num, err = Orm.QueryTable("blogs").Filter("title__icontains",kw).Count()
	return
}

// 截取正文内容获取描述
func (m_blogs *Blogs) SetDescription(){
	if m_blogs.Description != ""{
		return
	}
	content_without_html := utils.StripHtml(m_blogs.Content)
	content_rune := []rune(content_without_html)
	if(len(content_rune) > 100){
		m_blogs.Description = string(content_rune[:100])
	}else{
		m_blogs.Description = content_without_html
	}
}

// 删除文章
func (m_blogs *Blogs) Delete(id string) (err error){
	sql := "delete from blogs where id = ?"
	_, err = Orm.Raw(sql, id).Exec()
	return err
}

// 批量删除文章
func (m_blogs *Blogs) DeleteMulti(ids []string) (err error){
	ids_str := strings.Join(ids, ",")
	sql := fmt.Sprintf("delete from blogs where id in (%s)", ids_str)
	_, err = Orm.Raw(sql).Exec()
	return err
}

// 获取所有分类下文章的数量
func (m_blogs *Blogs) GetBlogsNumGroupByTid() (blog_info []orm.Params, err error){
	sql := "select count(id) num, tid from blogs group by tid"
	_, err = Orm.Raw(sql).Values(&blog_info)
	return blog_info, err
}

// 获取所有博客图片（分页）
func (m_blogs *Blogs)  GetBlogsThumbs(limit int, offset int) (thumbs []orm.Params, err error){
	limit_str := fmt.Sprintf(" limit %d, %d ", offset, limit)
	sql := "select thumb from blogs group by thumb order by create_time desc " + limit_str
	_, err = Orm.Raw(sql).Values(&thumbs)
	return thumbs, err
}

// 推送百度的文章
func (m_blogs *Blogs) GetBlogsIdPostBaidu() (blog_ids []interface{}, err error){
	expire := time.Now().Unix() - 86400
	var blogs []orm.Params
	sql := "select id from blogs where create_time > ? and status > 0"
	_, err = Orm.Raw(sql, expire).Values(&blogs)
	blog_ids = utils.GetFieldSliceFromData(blogs, "id")
	return blog_ids, err
}