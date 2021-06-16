package models

import (
	"fmt"
	"github.com/astaxie/beego/orm"
	"strings"
)

type Category struct {
	Id             int    `orm:column(id)`
	ZhName         string `orm:"column(zh_name)"`
	EnName         string `orm:"column(en_name)"`
	Title          string `orm:"column(title)"`
	Keywords       string `orm:"column(keywords)"`
	Description    string `orm:"column(description)"`
	Order          int    `orm:"column(order)"`
	Status         int    `orm:"column(status)"`
	Pid            int    `orm:"column(pid)"`
	ListTemplate   string `orm:"column(list_template)"`
	DetailTemplate string `orm:"column(detail_template)"`
	UseTable       string `orm:"column(use_table)"`
}

func (category *Category) GetAllCategory(fields ...string) (topCates []orm.Params, err error) {
	// 查询所有顶级栏目
	var fields_str string
	if len(fields) > 0{
		fields_str = strings.Join(fields, ",")
	}else{
		fields_str = "*"
	}
	sql := fmt.Sprintf("select %s from category where pid > 0 and status > 0", fields_str)
	_, err = Orm.Raw(sql).Values(&topCates)
	return topCates, err
}

// 获取所有的父子栏目和状态为0的栏目
func (category *Category) GetAllCategory2(fields ...string) (topCates []orm.Params, err error) {
	// 查询所有顶级栏目
	var fields_str string
	if len(fields) > 0{
		fields_str = strings.Join(fields, ",")
	}else{
		fields_str = "*"
	}
	sql := fmt.Sprintf("select %s from category order by `order`", fields_str)
	_, err = Orm.Raw(sql).Values(&topCates)
	return topCates, err
}

func (category *Category) GetCategoryByName(name string) (cate orm.Params, err error) {
	var cates []orm.Params
	sql := fmt.Sprintf("select * from category where en_name = '%s'", name)
	_, err = Orm.Raw(sql).Values(&cates)
	if len(cates) > 0 {
		cate = cates[0]
	}
	return cate, err
}

func (category *Category) GetCategoryById(id string) (cate orm.Params, err error) {
	var cates []orm.Params
	sql := fmt.Sprintf("select * from category where id = %s", id)
	_, err = Orm.Raw(sql).Values(&cates)
	if len(cates) > 0 {
		cate = cates[0]
	}
	return cate, err
}