package response

import (
	err "zbpblog_go/errInfo"
)

/* response包用于定义输出api接口的协议
 * 规则如下：
 * 1.所有输出协议都应包含Data，Errmsg和Errno这3个成员表示输出数据、错误信息和错误码，命名规范为：{控制器方法名}Response
 * 2.公共的response对象命名规范为 Common{自定义名称}Response，驼峰命名
 * 3.所有代码Data成员属性的response对象命名规范为 Data{控制器方法名}Response，驼峰命名
 * 4.response文件命名规范是 {控制器名}{方法名}Response
 */

// 错误信息设置接口
type ResponseInterface interface{
	SetError(errno int)
	SetErrmsg(errmsg string)
	SetData(data interface{})
}

// 公共接口协议 所有api的response对象都必须包含该公共response
type CommonResponse struct{
	Errno err.ErrCode	`json:"errno"`	// 错误码
	Errmsg string 	`json:"errmsg"`	// 错误信息
	Data interface{}		`json:"data"`
}

// 分页协议
type CommonPageResponse struct {
	CurrentPage int	`json:"current_page"`	// 当前页
	LastPage int		`json:"last_page"`	// 总页数
	PerPage int		`json:"per_page"`	// 每页多少条
	Total int 		`json:"total"`		// 总条数
}

func (cr *CommonResponse) SetError(errno err.ErrCode){
	cr.Errno = errno
	cr.Errmsg = errno.Error()
}

func (cr *CommonResponse) SetErrmsg(errmsg string){
	cr.Errmsg = errmsg
}

func (cr *CommonResponse) SetData(data interface{}){
	cr.Data = data
}