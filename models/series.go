package models

import (
	"fmt"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql" // 记得引入这个驱动
	"zbpblog_go/utils"
)

//+-------------+------------------+------+-----+---------+----------------+
//| Field       | Type             | Null | Key | Default | Extra          |
//+-------------+------------------+------+-----+---------+----------------+
//| id          | int(10) unsigned | NO   | PRI | NULL    | auto_increment |
//| name        | varchar(90)      | NO   | UNI |         |                |
//| thumb       | varchar(255)     | NO   |     |         |                |
//| create_time | int(10)          | NO   |     | NULL    |                |
//| update_time | int(10)          | NO   |     | NULL    |                |
//+-------------+------------------+------+-----+---------+----------------+

type Series struct {
	Id         int    `orm:"column(id)"`
	Name       string `orm:"column(Name)"`
	Thumb      string `orm:"column(thumb)"`
	ModelTimer
	//Blogs []*Blogs `orm:"reverse(many)"`
}

func (m_series *Series) GetSeriesById(id string) (series orm.Params, err error) {
	var series_slice []orm.Params
	sql := fmt.Sprintf("select * from series where id = %s", id)
	_, err = Orm.Raw(sql).Values(&series_slice)
	if len(series_slice) > 0 {
		series = series_slice[0]
	}
	return series, err
}

func (m_series *Series) GetSeriesByName(name string, fields ...string) (series []orm.Params, err error) {
	field_str := utils.BuildFieldStr(fields...)
	like_str := fmt.Sprintf("%%%s%%", name)
	sql := fmt.Sprintf("select %s from series where name like ? order by create_time desc", field_str)
	_, err = Orm.Raw(sql, like_str).Values(&series)
	return series, err
}

func (m_series *Series) GetSeriesPage(limit int, offset int) (series []orm.Params, err error) {
	limit_str := fmt.Sprintf("limit %d, %d", offset, limit)
	sql := fmt.Sprintf("select * from series order by id desc %s", limit_str)
	//spew.Dump(sql)
	_, err = Orm.Raw(sql).Values(&series)
	utils.TranTimeFields(&series) // 格式化时间
	return series, err
}

func (m_series *Series) GetSeriesNum() (num int64, err error){
	num, err = Orm.QueryTable("series").Count()
	return num, err
}