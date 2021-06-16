package main

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	_ "github.com/go-sql-driver/mysql" // 记得引入这个驱动
	"os"
	"strconv"
	"zbpblog_go/cache"
	_ "zbpblog_go/routers"
	"zbpblog_go/template_utils"
	"zbpblog_go/utils"
)

func main() {
	beego.Run()
}

func init() {

	//初始化数据库
	//orm.RunSyncdb("default", false, true)		// 异步执行db命令行命令。只有当我们需要自动创建模型是才需要调用

	// 根据参数选择配置文件
	choiceConf()

	// 初始化日志配置
	initLogs()

	// 初始化Redis连接池
	cache.InitRedis()

	// 设置静态文件目录
	beego.SetStaticPath("/statics", "statics")
	beego.SetStaticPath("/bootstrap", "bootstrap")
	beego.SetStaticPath("/uploads", "uploads")
	registerTemplateFunc()
}

// 注册模板函数
func registerTemplateFunc() {
	beego.AddFuncMap("TranInterface2String", template_utils.TranInterface2String)
	beego.AddFuncMap("Add", template_utils.Add)
	beego.AddFuncMap("Sub", template_utils.Sub)
	beego.AddFuncMap("Multiple", template_utils.Multiple)
	beego.AddFuncMap("Divide", template_utils.Divide)
	beego.AddFuncMap("InSlice", utils.InSlice)
}

// 选择配置文件(用于开启多个项目进程监听不同端口)
func choiceConf(){
	var port string
	port = "8080"		// 默认使用8080端口
	if arg_len := len(os.Args); arg_len > 1{
		if _, err := strconv.Atoi(os.Args[arg_len - 1]);err == nil{
			port = os.Args[arg_len - 1]
		}
	}
	confDir := "./conf/"
	confPath :=  confDir + "app-" + port + ".conf"
	beego.LoadAppConfig("ini", confPath)	// 选择配置文件的路径
}

// 初始化日志配置
func initLogs(){
	log_path := beego.AppConfig.String("log_path")
	log_config := fmt.Sprintf(`{"filename":"%s"}`, log_path)
	logs.SetLogger("file", log_config)			// 指定日志输出到文件，文件路径为log_path, 此时日志会同时输出到控制台和文件
	//logs.BeeLogger.DelLogger("console")		// 不让日志输出到控制台
	logs.Async(1e3)		// 设置日志为异步输出, chan缓冲区为1000
	//logs.SetLogFuncCall(true)
	logs.SetLevel(logs.LevelInformational)
}