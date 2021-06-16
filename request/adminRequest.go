package request

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/utils/captcha"
	"github.com/davecgh/go-spew/spew"
	"strconv"
	"zbpblog_go/errInfo"
	"zbpblog_go/models"
	"zbpblog_go/utils"
)

type BaseAdminRequest struct{
	UserName string		// 用户名
	Sign     string		// 用户登录态
}

func (bq *BaseAdminRequest) BuildReq(c *beego.Controller){
	bq.UserName = c.Ctx.GetCookie("user_name")
	bq.Sign = c.Ctx.GetCookie("sign")
}

/*** 执行登录Request  ***/
type LoginRequest struct {
	Username string		`form:"username"`
	Password string		`form:"password"`
	Cap	string			`form:"cap"`		// 验证码
}

func (loginRequest *LoginRequest) BuildReq(c *beego.Controller){
	if err := c.ParseForm(loginRequest); err != nil {
		logs.Error("LoginRequest - ParseForm failed:" + err.Error())
	}
}

func (loginRequest *LoginRequest) Valid(c *beego.Controller, cpt *captcha.Captcha) (is_valid bool, err error){
	//is_valid = true
	if loginRequest.Username == ""{
		return false, errInfo.CODE_USERNAME_EMPTY
	}

	if loginRequest.Password == ""{
		return false, errInfo.CODE_PASSWORD_EMPTY
	}

	if loginRequest.Cap == ""{
		return false, errInfo.CODE_CAPCHA_EMPTY
	}

	if !cpt.VerifyReq(c.Ctx.Request){
		return false, errInfo.CODE_CAPCHA_WRONG
	}

	return true, err
}


/*** 文章页Request ***/
type BlogListRequest struct {
	Kw string
	Page string
}

func (blogListRequest *BlogListRequest) BuildReq(c *beego.Controller){
	params := c.Ctx.Request.URL.Query()
	if params["kw"] != nil && len(params["kw"]) > 0 {
		blogListRequest.Kw = params["kw"][0]
	}else{
		blogListRequest.Kw = ""
	}

	if params["page"] != nil && len(params["page"]) > 0 {
		blogListRequest.Page = params["page"][0]
	}else{
		blogListRequest.Page = "1"
	}
}

/*** 获取标签图片接口Request  ***/
type BlogGetTagsPicRequest struct {
	TagIds []string	`form:"tag_ids"`
}

func (blogGetTagsPicRequest *BlogGetTagsPicRequest) BuildReq(c *beego.Controller){
	blogGetTagsPicRequest.TagIds = c.GetStrings("tag_ids[]")
}

/*** 搜索系列接口Request ***/
type SearchSeriesRequest struct{
	Search string `form:"search"`
}

func (req *SearchSeriesRequest) BuildReq(c *beego.Controller){
	req.Search = c.GetString("search")
}

/*** 添加栏目Request  ***/
type CategoryDoAddEditRequest struct {
	Id string 			`form:"id"`
	ZhName string		`form:"zh_name"`
	EnName string 		`form:"en_name"`
	Pid string 			`form:"pid"`
	Title string		`form:"title"`
	Keywords string		`form:"keywords"`
	Description string 	`form:"description"`
	ListTemplate string 		`form:"list_template"`
	DetailTemplate string		`form:"detail_template"`
	UseTable string		`form:"use_table"`
}

func (req *CategoryDoAddEditRequest) SetData(m_cate *models.Category){
	m_cate.ZhName = req.ZhName
	m_cate.EnName = req.EnName
	m_cate.Pid, _ = strconv.Atoi(req.Pid)
	m_cate.Title = req.Title
	m_cate.Keywords = req.Keywords
	m_cate.Description = req.Description
	m_cate.ListTemplate = req.ListTemplate
	m_cate.DetailTemplate = req.DetailTemplate
	m_cate.UseTable = req.UseTable
}

func (req *CategoryDoAddEditRequest) SetAddData(m_cate *models.Category){
	req.SetData(m_cate)

	if beego.AppConfig.String("runmode") == "dev" {
		m_cate.Status = 0
	}
}

func (req *CategoryDoAddEditRequest) SetEditData(m_cate *models.Category){
	m_cate.Id, _ = strconv.Atoi(req.Id)
	req.SetData(m_cate)
}

func (req *CategoryDoAddEditRequest) BuildReq(c *beego.Controller){
	c.ParseForm(req)
}

/*** 添加文章Request  ***/
type BlogDoAddEditRequest struct {
	Id string			`form:"id"`
	Title string 		`form:"title"`
	Tid string 			`form:"tid"`
	Type string 		`form:"type"`
	Keywords string		`form:"keywords"`
	Description string 	`form:"description"`
	Writer string 		`form:"writer"`
	TagsId []string		`form:"tags_id[]"`
	Thumb string		`form:"tag_img"`
	Content string		`form:"content"`
	SeriesId string		`form:"series_id"`
}

func (req *BlogDoAddEditRequest) SetData(m_blog *models.Blogs){
	m_blog.Title = req.Title
	m_blog.Tid, _ = strconv.Atoi(req.Tid)
	m_blog.Type,_ = strconv.Atoi(req.Type)
	m_blog.Keywords = req.Keywords
	m_blog.Description = req.Description
	m_blog.Writer = req.Writer
	m_blog.Thumb = req.Thumb
	m_blog.Content = req.Content
	m_blog.SeriesId,_ = strconv.Atoi(req.SeriesId)
	m_blog.SetDescription()
}

func (req *BlogDoAddEditRequest) SetAddData(m_blog *models.Blogs){
	req.SetData(m_blog)
	m_blog.SetTimeWhenAdd()

	if beego.AppConfig.String("runmode") == "dev" {
		m_blog.Status = 0
	}
}

func (req *BlogDoAddEditRequest) SetEditData(m_blog *models.Blogs){
	m_blog.Id, _ = strconv.Atoi(req.Id)
	req.SetData(m_blog)
	m_blog.SetTimeWhenEdit()
}

func (req *BlogDoAddEditRequest) BuildReq(c *beego.Controller){
	c.ParseForm(req)

	if req.Thumb == ""{		// 如果没有使用标签图片或者系列封面图，则接收上传图片
		destination := beego.AppConfig.String("thumbpath")
		fp, err := utils.UploadFileByName(c, "thumb", destination)
		if err != nil {
			logs.Error("BlogDoAddEditRequest - UploadFile failed:" + err.Error())
			return
		}
		req.Thumb = fp
	}
}

/*** 修改文章/栏目/标签等Request ***/
type EditByIdRequest struct{
	Id string
}

func (req *EditByIdRequest) BuildReq(c *beego.Controller){
	params := c.Ctx.Request.URL.Query()
	if params["id"] != nil {
		req.Id = params["id"][0]
	}else{
		req.Id = ""
	}
}

/*** 删除文章页Request ***/
type DelRequest struct{
	Model string	`form:"model"`
	Id string		`form:"id"`
	Ids []string	`form:"ids[]"`
}

func (req *DelRequest) BuildReq(c *beego.Controller){
	c.ParseForm(req)

	params := c.Ctx.Request.URL.Query()
	if params["id"] != nil && req.Id == ""{
		req.Id = params["id"][0]
	}
	spew.Dump(req)
}

/*** 变更状态Request ***/
type TranStatusRequest struct{
	Id string	`form:"id"`
	Status int  `form:"status"`
	Model string `form:"model"`
}

func (req *TranStatusRequest) BuildReq(c *beego.Controller){
	c.ParseForm(req)
}

/*** 排序Request ***/
type SortRequest struct{
	Sort []map[string]string  `form:"sort[][]"`
	Model string `form:"model"`
}

func (req *SortRequest) BuildReq(c *beego.Controller){
	c.ParseForm(req)
	c.Ctx.Input.Bind(&req.Sort, "sort")		// sort字段是一个复合类型，需要用Ctx.Input.Bind获取并赋值
}

/*** 站点设置Request ***/
type SiteEditRequest struct{
	Id	string		`form:"id"`
	Domain string `form:"domain"`
	SiteName string `form:"site_name"`
	Title string `form:"title"`
	Keywords string `form:"keywords"`
	BlogKeywords string `form:"blog_keywords"`
	NewsKeywords string `form:"news_keywords"`
	Description string `form:"description"`
	TjCode string `form:"tj_code"`
	Iframe string `form:"iframe"`
}

func (req *SiteEditRequest) BuildReq(c *beego.Controller){
	c.ParseForm(req)
}

func (req *SiteEditRequest) SetEditData(m_site *models.Site){
	m_site.Id, _ = strconv.Atoi(req.Id)
	m_site.Domain = req.Domain
	m_site.SiteName = req.SiteName
	m_site.Title = req.Title
	m_site.Keywords = req.Keywords
	m_site.BlogKeywords = req.BlogKeywords
	m_site.NewsKeywords = req.NewsKeywords
	m_site.Description = req.Description
}

func (req *SiteEditRequest) EditField() []string{
	return []string{"domain", "site_name", "title", "keywords", "blog_keywords", "news_keywords", "description"}
}

/*** 标签页Request ***/
type PageRequest struct{
	Page int		`form:"page"`
	P	int		`form:"p"`
}

func (req *PageRequest) BuildReq(c *beego.Controller){
	params := c.Ctx.Request.URL.Query()

	if params["page"] != nil && len(params["page"]) > 0 {
		req.Page, _ = strconv.Atoi(params["page"][0])
	}else{
		req.Page = 1
	}

	if params["p"] != nil && len(params["p"]) > 0 {
		req.P, _ = strconv.Atoi(params["p"][0])
	}else{
		req.P = 20
	}
}

/*** 添加系列Request ***/
type SeriesDoAddEditRequest struct{
	Id string 		`form:"id"`
	Name string 		`form:"name"`
	Thumb string 		`form:"tag_img"`		// 优先使用博客图片，次之是上传图片
}

func (req *SeriesDoAddEditRequest) BuildReq(c *beego.Controller){
	c.ParseForm(req)
	spew.Dump(req)
	if req.Thumb == ""{
		destination := beego.AppConfig.String("thumbpath")
		fp, err := utils.UploadFileByName(c, "thumb", destination)
		if err != nil {
			logs.Error("TagsDoAddEditRequest - UploadFileByName failed:" + err.Error())
			return
		}
		req.Thumb = fp
	}
}

func (req *SeriesDoAddEditRequest) SetAddEditData(m_series *models.Series){
	m_series.Id, _ = strconv.Atoi(req.Id)
	m_series.Name = req.Name
	m_series.Thumb = req.Thumb

	if m_series.Id == 0{
		m_series.SetTimeWhenAdd()
	}else{
		m_series.SetTimeWhenEdit()
	}
}

func (req *SeriesDoAddEditRequest) EditField() []string{
	return []string{"name", "thumb", "update_time"}
}

/*** 添加标签Request ***/
type TagsDoAddEditRequest struct{
	Id string 		`form:"id"`
	Name string 		`form:"name"`
	Type int 		`form:"type"`
	Top int 		`form:"top"`
	Thumb string 		`form:"thumb_hidden"`
}

func (req *TagsDoAddEditRequest) BuildReq(c *beego.Controller){
	c.ParseForm(req)
spew.Dump(req)
	if req.Thumb == ""{
		destination := beego.AppConfig.String("thumbpath")
		fp, err := utils.UploadFileByName(c, "thumb", destination)
		if err != nil {
			logs.Error("TagsDoAddEditRequest - UploadFileByName failed:" + err.Error())
			return
		}
		req.Thumb = fp
	}
}

func (req *TagsDoAddEditRequest) SetAddEditData(m_tags *models.Tags){
	m_tags.Id, _ = strconv.Atoi(req.Id)
	m_tags.Name = req.Name
	m_tags.Type = req.Type
	m_tags.Top = req.Top
	m_tags.Thumb = req.Thumb
}

func (req *TagsDoAddEditRequest) EditField() []string{
	return []string{"name", "type", "top", "thumb"}
}


/*** 添加管理员Request ***/
type UserDoAddEditRequest struct{
	Id string 		`form:"id"`
	Username string 		`form:"username"`
	Password string 		`form:"password"`
}

func (req *UserDoAddEditRequest) BuildReq(c *beego.Controller){
	c.ParseForm(req)
}

func (req *UserDoAddEditRequest) SetAddEditData(m_user *models.AdminUser){
	m_user.Id, _ = strconv.Atoi(req.Id)
	m_user.Username = req.Username
	m_user.Password = utils.EncryptMd5(req.Password)
	if m_user.Id == 0 {
		m_user.SetTimeWhenAdd()
	}else{
		m_user.SetTimeWhenEdit()
	}
}

func (req *UserDoAddEditRequest) EditField() []string{
	return []string{"name", "type", "top", "thumb"}
}