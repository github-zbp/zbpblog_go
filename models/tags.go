package models

import (
	"fmt"
	"github.com/astaxie/beego/orm"
)

type Tags struct {
	Id    int    `orm:"column(id)"`
	Name  string `orm:"column(name)"`
	Type  int    `orm:"column(type)"`
	Color string `orm:"column(color)"`
	Order int    `orm:"column(order)"`
	Thumb string `orm:"column(thumb)"`
	Top   int    `orm:"column(top)"`
}

func (m_tags *Tags) GetTagsById(id string) (tag orm.Params, err error) {
	var tags []orm.Params
	sql := fmt.Sprintf("select * from tags where id = %s", id)
	_, err = Orm.Raw(sql).Values(&tags)
	if len(tags) > 0 {
		tag = tags[0]
	}
	return tag, err
}

func (m_tags *Tags) GetTagsByIds(ids []interface{}, fields ...string) (tags []orm.Params, err error) {
	qs := Orm.QueryTable("tags").Filter("id__in", ids)
	if len(fields) > 0{
		_, err = qs.Values(&tags, fields...)
	}else{
		_, err = qs.Values(&tags)
	}
	return tags, err
}

func (m_tags *Tags) GetTags(limit int) (tags []orm.Params, err error) {
	_, err = Orm.QueryTable("tags").OrderBy("-top").Limit(limit).Values(&tags)
	return tags, err
}

func (m_tags *Tags) GetAllTags() (tags []orm.Params, err error){
	sql := "select * from tags order by top desc"
	_, err = Orm.Raw(sql).Values(&tags)
	return tags, err
}

func (m_tags *Tags) GetTagsPage(limit int, offset int) (tags []orm.Params, err error) {
	limit_str := fmt.Sprintf("limit %d, %d", offset, limit)
	sql := fmt.Sprintf("select * from tags order by id desc %s", limit_str)
	_, err = Orm.Raw(sql).Values(&tags)
	return tags, err
}

func (m_tags *Tags) GetTagsNum() (num int64, err error){
	num, err = Orm.QueryTable("tags").Count()
	return num, err
}