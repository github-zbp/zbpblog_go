package service

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"time"
	"zbpblog_go/errInfo"
	"zbpblog_go/models"
	"zbpblog_go/request"
	"zbpblog_go/utils"
)

// 执行登录操作
func ExecuteLogin(c *beego.Controller, loginRequest *request.LoginRequest) (success bool, err error){
	m_admin_user := new(models.AdminUser)
	user, err := m_admin_user.GetUserByName(loginRequest.Username)
	if err != nil {
		logs.Error("ExecuteLogin - GetUserByName failed:" + err.Error())
		return
	}

	if user == nil{
		return false, errInfo.CODE_USER_NOT_FIND
	}

	if user["password"] != utils.EncryptMd5(loginRequest.Password) {
		return false, errInfo.CODE_PASSWORD_WRONG
	}

	// 用户身份确认后，生成登录态sign
	sign := CreateUserSign(loginRequest.Username)
	if err = m_admin_user.UpdateSign(loginRequest.Username, sign); err != nil{
		logs.Error("ExecuteLogin - UpdateSign failed:" + err.Error())
		return false, err
	}

	// 生成cookie
	expire := 3600 * 24 * 7
	c.Ctx.SetCookie("user_name", loginRequest.Username, expire)
	c.Ctx.SetCookie("sign", sign, expire)

	return true, err
}

func CreateUserSign(username string) (sign string){
	t := time.Now().String()
	sign = utils.EncryptMd5(username + t)
	return sign
}