package template_utils

import (
	"zbpblog_go/utils"
)

func TranInterface2String(data interface{}) string {
	if data == nil {
		return ""
	}
	return data.(string)
}

// 加减乘除运算，用于template模板
func Add(n1 interface{}, n2 interface{}) (res interface{}){
	return utils.TranInterface2Int(n1) + utils.TranInterface2Int(n2)
}

func Sub(n1 interface{}, n2 interface{}) (res interface{}){
	return utils.TranInterface2Int(n1) - utils.TranInterface2Int(n2)
}

func Multiple(n1 interface{}, n2 interface{}) (res interface{}){
	return utils.TranInterface2Int(n1) * utils.TranInterface2Int(n2)
}

func Divide(n1 interface{}, n2 interface{}) (res interface{}){
	return utils.TranInterface2Int(n1) / utils.TranInterface2Int(n2)
}