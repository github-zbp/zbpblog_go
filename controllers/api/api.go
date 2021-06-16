package api

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"zbpblog_go/errInfo"
	"zbpblog_go/request"
	"zbpblog_go/response"
	"zbpblog_go/service"
	"zbpblog_go/utils"
)

type ApiController struct {
	beego.Controller
}
var beegoController *beego.Controller

func (c *ApiController) Prepare() {
	beegoController = &c.Controller
}

// 加载更多文章（api）
func (c *ApiController) MoreBlogs(){
	// 初始话前端参数
	moreBlogRequest := &request.IndexMoreBlogsRequest{}
	moreBlogRequest.BuildReq(beegoController)

	// 加载标签文章页数据
	service.GetMoreBlogsData(c.Data, moreBlogRequest)
	c.ServeJSON()		// 以json格式输出数据
}

// 后台加载标签图片
func (c *ApiController) GetTagsPic(){
	req := new(request.BlogGetTagsPicRequest)
	req.BuildReq(beegoController)

	// 加载标签图片
	service.LoadTagsPic(c.Data, req)
	c.ServeJSON()
}

// 搜索系列接口
func (c *ApiController) SearchSeries(){
	req := new(request.SearchSeriesRequest)
	req.BuildReq(beegoController)

	service.SearchSeriesByName(c.Data, req)
	c.ServeJSON()
}

// 添加文章接口
func (c *ApiController) BlogDoAdd(){
	if check_login_res := service.CheckLoginStatusForApi(beegoController); !check_login_res {
		return
	}

	req := new(request.BlogDoAddEditRequest)
	req.BuildReq(beegoController)
	service.AddBlog(c.Data, req)
	c.ServeJSON()
}

// 通过id删除的通用接口
func (c *ApiController) ApiDel(){
	if check_login_res := service.CheckLoginStatusForApi(beegoController); !check_login_res {
		return
	}

	req := new(request.DelRequest)
	req.BuildReq(beegoController)
	resp := new(response.CommonResponse)
	if (req.Id == "" && len(req.Ids) == 0) || req.Model == ""{
		resp.SetError(errInfo.CODE_PARAM_ERR)
		c.Data["json"] = resp
		c.ServeJSON()
		return
	}

	defer utils.SetErrorApiResp(resp, beegoController)

	err := service.Del(req)
	if err != nil{
		logs.Error("ApiDel - Delete failed:" + err.Error())
		resp.SetError(errInfo.CODE_DELETE_FAILED)
		return
	}
}

// 添加栏目接口
func (c *ApiController) CategoryDoAdd(){
	if check_login_res := service.CheckLoginStatusForApi(beegoController); !check_login_res {
		return
	}

	resp := new(response.CommonResponse)
	defer utils.SetErrorApiResp(resp, beegoController)

	req := new(request.CategoryDoAddEditRequest)
	req.BuildReq(beegoController)
	errCode := service.AddCategory(req)
	resp.SetError(errCode)
}

// 修改栏目接口
func (c *ApiController) CategoryDoEdit(){
	if check_login_res := service.CheckLoginStatusForApi(beegoController); !check_login_res {
		return
	}

	resp := new(response.CommonResponse)
	defer utils.SetErrorApiResp(resp, beegoController)

	req := new(request.CategoryDoAddEditRequest)
	req.BuildReq(beegoController)
	errCode := service.EditCategory(req)
	resp.SetError(errCode)
}

// 修改状态（通用）
func (c *ApiController) ApiTranStatus(){
	if check_login_res := service.CheckLoginStatusForApi(beegoController); !check_login_res {
		return
	}

	req := new(request.TranStatusRequest)
	req.BuildReq(beegoController)
	resp := new(response.CommonResponse)
	if req.Id == ""  || req.Model == ""{
		resp.SetError(errInfo.CODE_PARAM_ERR)
		c.Data["json"] = resp
		c.ServeJSON()
		return
	}

	defer utils.SetErrorApiResp(resp, beegoController)

	err := service.TranStatus(req)
	if err != nil{
		logs.Error("ApiTranStatus - TranStatus failed:" + err.Error())
		resp.SetError(errInfo.CODE_TRAN_STATUS_FAILED)
		return
	}
}

// 排序（通用）
func (c *ApiController) ApiSort(){
	if check_login_res := service.CheckLoginStatusForApi(beegoController); !check_login_res {
		return
	}

	req := new(request.SortRequest)
	req.BuildReq(beegoController)
	resp := new(response.CommonResponse)
	if len(req.Sort) == 0  || req.Model == ""{
		resp.SetError(errInfo.CODE_PARAM_ERR)
		c.Data["json"] = resp
		c.ServeJSON()
		return
	}

	defer utils.SetErrorApiResp(resp, beegoController)

	err := service.Sort(req)
	if err != nil{
		logs.Error("ApiSort - Sort failed:" + err.Error())
		resp.SetError(errInfo.CODE_SORT_FAILED)
		return
	}
}

// 修改站点信息
func (c *ApiController) SiteDoEdit(){
	if check_login_res := service.CheckLoginStatusForApi(beegoController); !check_login_res {
		return
	}

	req := new(request.SiteEditRequest)
	req.BuildReq(beegoController)
	resp := new(response.CommonResponse)

	defer utils.SetErrorApiResp(resp, beegoController)

	err := service.EditSite(req)
	if err != errInfo.CODE_OK{
		resp.SetError(err)
	}
}

// 获取博客图片
func (c *ApiController) GetBlogThumbs(){
	if check_login_res := service.CheckLoginStatusForApi(beegoController); !check_login_res {
		return
	}
	req := new(request.PageRequest)
	req.BuildReq(beegoController)
	resp := new(response.CommonResponse)
	defer utils.SetErrorApiResp(resp, beegoController)

	data, err := service.GetBlogThumbs(req)
	resp.SetError(err)
	resp.SetData(data)
}

// 添加用户
func (c *ApiController) UserDoAdd(){
	if check_login_res := service.CheckLoginStatusForApi(beegoController); !check_login_res {
		return
	}

	resp := new(response.CommonResponse)
	defer utils.SetErrorApiResp(resp, beegoController)

	req := new(request.UserDoAddEditRequest)
	req.BuildReq(beegoController)
	errCode := service.AddUser(req)
	resp.SetError(errCode)
}

// 百度推送接口
func (c *ApiController) PostBaidu(){
	//if check_login_res := service.CheckLoginStatusForApi(beegoController); !check_login_res {
	//	return
	//}
	resp := new(response.CommonResponse)
	defer utils.SetErrorApiResp(resp, beegoController)
	resp_string, errCode := service.PostBaidu()
	resp.SetData(resp_string)
	resp.SetError(errCode)
}

// 测试接收请求接口
//func (c *ApiController) TestReceivePost(){
//	if check_login_res := service.CheckLoginStatusForApi(beegoController); !check_login_res {
//		return
//	}
//	resp := new(response.CommonResponse)
//	defer utils.SetErrorApiResp(resp, beegoController)
//}

// 测试发送请求接口
//func (c *ApiController) TestSendPost(){
//	resp := new(response.CommonResponse)
//	defer utils.SetErrorApiResp(resp, beegoController)
//	header := map[string]string{
//		"Cookie" : "Hm_lvt_080836300300be57b7f34f4b3e97d911=1619660610,1620610638,1620697137,1620869895; user_name=fx112; sign=d97030a03d8e2d07339ab152bb50187b; Hm_lvt_4898023a3a963db158d0eeee27c56c80=1620351594,1621154263,1621216108,1621302517; Hm_lpvt_4898023a3a963db158d0eeee27c56c80=1621302517",
//	}
//	url := "http://81.71.136.86:8080/zbpadmin/api/test_receive_post"
//	curl := utils.Curl{
//		Url:url,
//		Header:header,
//	}
//	body, err := curl.PostForJson()
//	if err != nil{
//		logs.Error("TestSendPost failed:" + err.Error())
//	}
//
//	resp.SetData(body)
//}