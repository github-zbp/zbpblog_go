package models

import (
	"github.com/astaxie/beego/orm"
	"zbpblog_go/utils"
)

type AdminUser struct {
	Id int				`orm:column(id)`
	Username string		`orm:column(user_name)`
	Password string		`orm:column(password)`
	LastLogtime	int		`orm:column(last_logtime)`
	LastIp string		`orm:column(last_ip)`
	Level int			`orm:column(level)`
	Status int			`orm:column(status)`
	ModelTimer
	Sign string			`orm:column(sign)`
}

func (au *AdminUser) GetAll() (users []orm.Params, err error){
	sql := "select * from admin_user"
	_, err = Orm.Raw(sql).Values(&users)
	utils.TranTimeFields(&users, "last_logtime")
	return users, err
}

func (au *AdminUser) GetUserByName(name string) (user orm.Params, err error){
	var users []orm.Params
	sql := "select * from admin_user where username = ? and status = 1 limit 1"
	_, err = Orm.Raw(sql, name).Values(&users)

	if len(users) > 0{
		user = users[0]
	}
	return user, err
}

func (au *AdminUser) UpdateSign(name string, sign string) (err error){
	sql := "update admin_user set sign = ? where username = ? and status"
	_, err = Orm.Raw(sql, sign, name).Exec()
	return err
}