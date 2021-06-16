package service

import "github.com/astaxie/beego/orm"

// 以下几个Service包下的全局变量供所有其他Service文件使用，为了提高其可见性，将他们提取到了单独的baseService中
var Seo orm.Params         // 当前SEO
var CurrentCate orm.Params // 当前栏目
var Category []orm.Params  // 所有栏目
var Tags []orm.Params		// 所有标签
