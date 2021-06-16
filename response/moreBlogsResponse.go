package response

import "github.com/astaxie/beego/orm"

/********   更多文章api协议    ***********/
type MoreBlogsResponse struct {
	CommonResponse		// 包含 Errno 和 Errmsg
	//Data *DataMoreBlogsResponse		`json:"data"`
}

type DataMoreBlogsResponse struct {
	CommonPageResponse
	Data []*DataMoreBlogsPageResponse		`json:"data"`
}

type DataMoreBlogsPageResponse struct{
	Id string			`json:"id"`
	CreateTime string 	`json:"create_time"`
	Description string 	`json:"description"`
	Title string		`json:"title"`
	Thumb string		`json:"thumb"`
	Cate orm.Params		`json:"cate"`
}

func (more_blog_resp *DataMoreBlogsResponse) SetData(blogs []orm.Params, cate map[string]orm.Params){
	var data []*DataMoreBlogsPageResponse
	for _, blog := range blogs {
		data_more_blogs_page := DataMoreBlogsPageResponse{
			Id: blog["id"].(string),
			CreateTime: blog["create_time"].(string),
			Description: blog["description"].(string),
			Title: blog["title"].(string),
			Thumb: blog["thumb"].(string),
			Cate: cate[blog["tid"].(string)],
		}
		data = append(data, &data_more_blogs_page)
	}
	more_blog_resp.Data = data
}

func (more_blog_resp *DataMoreBlogsResponse) SetPage(page_info map[string]interface{}){
	more_blog_resp.CurrentPage = page_info["curPage"].(int)
	more_blog_resp.LastPage = page_info["totalPage"].(int)
	more_blog_resp.PerPage = page_info["pageRows"].(int)
	more_blog_resp.Total = page_info["totalRow"].(int)
}

func (more_blog_resp *MoreBlogsResponse) Set(blogs []orm.Params, cate map[string]orm.Params, page_info map[string]interface{}){
	data_more_blogs_resp := new(DataMoreBlogsResponse)
	data_more_blogs_resp.SetPage(page_info)
	data_more_blogs_resp.SetData(blogs, cate)

	more_blog_resp.SetData(data_more_blogs_resp)
}