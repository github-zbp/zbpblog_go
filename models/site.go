package models

import (
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql" // 记得引入这个驱动
)

type Site struct{
	Id int				`json:"id"	orm:"column(id)"`
	Uid string			`json:"uid"	orm:"column(uid)"`
	Username string		`json:"username"	orm:"column(username)"`
	SiteName string		`json:"site_name"	orm:"column(site_name)"`
	Domain string		`json:"domain"	orm:"column(domain)"`
	Title string		`json:"title"	orm:"column(title)"`
	Keywords string		`json:"keywords"	orm:"column(keywords)"`
	Description string	`json:"description"	orm:"column(description)"`
	EndTime int			`json:"end_time"	orm:"column(end_time)"`
	NewsKeywords string	`json:"news_keywords"	orm:"column(news_keywords)"`
	BlogKeywords string	`json:"blog_keywords"	orm:"column(blog_keywords)"`
	Html int			`json:"html"	orm:"column(html)"`
}

func (m_site *Site) GetSeo() (site orm.Params){
	var sites []orm.Params
	sql := "select * from site"
	Orm.Raw(sql).Values(&sites)		// 这里请传入site而不是&site，因为site本身就是一个引用
	if len(sites) > 0{
		site = sites[0]
	}
	return site
}