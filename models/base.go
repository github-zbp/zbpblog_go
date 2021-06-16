package models

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql" // 记得引入这个驱动
	"strings"
	"time"
	"zbpblog_go/utils"
)

var Orm orm.Ormer
var Qb orm.QueryBuilder
var Models []string

func init() {
	InitDb()
	Orm = orm.NewOrm()
	Qb, _ = orm.NewQueryBuilder("mysql")
}

// 初始化数据库配置
func InitDb() {
	// 数据库配置
	dbhost := beego.AppConfig.String("dbhost")
	dbuser := beego.AppConfig.String("dbuser")
	dbpassword := beego.AppConfig.String("dbpassword")
	dbname := beego.AppConfig.String("dbname")
	dbport := beego.AppConfig.String("dbport")
	dbcharset := beego.AppConfig.String("dbcharset")

	// root:123456@tcp(127.0.0.1:3306)/zbpblog?charset=utf8
	orm.RegisterModel(getDbTables()...) // 注册模型表
	dsn := dbuser + ":" + dbpassword + "@tcp(" + dbhost + ":" + dbport + ")/" + dbname + "?charset=" + dbcharset
	orm.RegisterDataBase("default", "mysql", dsn) // 注册数据库（可以注册多个数据库，但是必须注册一个别名为 default 的数据库，作为默认使用）
	orm.RegisterDriver("default", orm.DRMySQL)    // 注册数据库驱动。第一参是自定义驱动名称，第二参是驱动类型这里使用mysql驱动

	runmode := beego.AppConfig.String("runmode")
	orm.Debug = runmode == "dev"
}

func getDbTables() (tables []interface{}) {
	tables = append(tables,
		new(Site),
		new(Category),
		new(Blogs),
		new(Series),
		new(Tags),
		new(TagsBlogs),
		new(AdminUser),
	)
	Models = []string{"site", "category", "blogs", "series", "tags", "links", "admin_user"}
	return tables
}


// 数据表共有字段
type ModelTimer struct{
	CreateTime  int    `orm:"column(create_time)"`
	UpdateTime  int    `orm:"column(update_time)"`
}

func (mt *ModelTimer) SetTimeWhenAdd(){
	timeStamp := utils.TranInt64ToInt(time.Now().Unix())
	mt.UpdateTime = timeStamp
	mt.CreateTime = timeStamp
}

func (mt *ModelTimer) SetTimeWhenEdit(){
	mt.UpdateTime = utils.TranInt64ToInt(time.Now().Unix())
}

// 模型删除数据通用方法
func Delete(model string, id string) (err error){
	sql := fmt.Sprintf("delete from %s where id = ?", model)
	_, err = Orm.Raw(sql, id).Exec()
	return err
}

// 模型批量删除数据通用方法
func DeleteMulti(model string, ids []string) (err error){
	ids_str := strings.Join(ids, ",")
	sql := fmt.Sprintf("delete from %s where id in (%s)", model, ids_str)
	_, err = Orm.Raw(sql).Exec()
	return err
}

// 模型修改数据状态通用方法
func TranStatus(model string, id string, status int) (err error){
	sql := fmt.Sprintf("update %s set status=? where id=%s", model, id)
	_, err = Orm.Raw(sql, status).Exec()
	return err
}

// 模型排序通用方法
func Sort(model string, sort []map[string]string) (err error){
	update_str := ""
	id_str := ""
	for _, item := range sort{
		update_str += fmt.Sprintf(" when %s then %s ", item["id"], item["order"])
		id_str += fmt.Sprintf("%s,", item["id"])
	}
	//update_str = strings.TrimRight(update_str, ",")
	id_str = strings.TrimRight(id_str, ",")
	sql := fmt.Sprintf("update %s set `order` = case `id` %s end where id in (%s)", model, update_str, id_str)
	_, err = Orm.Raw(sql).Exec()
	return err
}