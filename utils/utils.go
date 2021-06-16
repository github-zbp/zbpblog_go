package utils

import (
	"crypto/md5"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"io"
	"math"
	"mime/multipart"
	"os"
	"path"
	"regexp"
	"strconv"
	"strings"
	"time"
	"zbpblog_go/errInfo"
	"zbpblog_go/response"

	//"github.com/davecgh/go-spew/spew"
)

// formatString是格式化的格式 Y 年， m 月， d 日，H 小时，i 分钟， s 秒
func TimeToStr(timeStamp interface{}, formatStr string) (timeFormat string) {
	var timeStampInt int
	value, ok := timeStamp.(string)
	if ok { // 如果是字符串则转为数字
		timeStampInt, _ = strconv.Atoi(value)
	} else {
		timeStampInt, _ = timeStamp.(int)
	}
	timeObj := time.Unix(int64(timeStampInt), 0)
	formatStr = strings.Replace(formatStr, "Y", "2006", -1)
	formatStr = strings.Replace(formatStr, "m", "01", -1)
	formatStr = strings.Replace(formatStr, "d", "02", -1)
	formatStr = strings.Replace(formatStr, "H", "15", -1)
	formatStr = strings.Replace(formatStr, "i", "04", -1)
	formatStr = strings.Replace(formatStr, "s", "05", -1)
	return timeObj.Format(formatStr)
}

// 将数据中的指定时间字段从时间戳转为time.Time类型
func TranTimeFieldsForOne(data *orm.Params, fields ...interface{}) {
	if len(fields) == 0 {
		fields = append(fields, "create_time", "update_time")
	}

	for _, field := range fields {
		field_str, ok1 := field.(string) // 确保field是字符串类型

		if ok1 {
			(*data)[field_str] = TimeToStr((*data)[field_str], "Y-m-d H:i:s")
		}
	}
}

// 将数据中的指定时间字段从时间戳转为time.Time类型
func TranTimeFields(data *[]orm.Params, fields ...interface{}) {
	for _, value := range *data {
		TranTimeFieldsForOne(&value, fields...)
	}
}

// 从数据列表中获取某一个字段slice
func GetFieldSliceFromData(data []orm.Params, field_name string) (field_slice []interface{}) {
	for _, value := range data {
		if field_value, ok := value[field_name]; ok {
			field_slice = append(field_slice, field_value)
		}
	}

	return field_slice
}

/**
 * 把数据库的数据从 []orm.Params 变为 map[id]orm.Params
 */
func GetModelDataMap(data []orm.Params, field_key string) (map[string]orm.Params){
	data_map := make(map[string]orm.Params)
	for _, row := range data {
		var key string
		switch row[field_key].(type){
		case string:
			key = row[field_key].(string)
		case int:
			key = strconv.Itoa(row[field_key].(int))
		default:
			return nil
		}
		data_map[key] = row
	}
	return data_map
}

/**
把从数据库里取出的数据用一个字段作为key来返回map
data 要处理的数据集
field_name	作为key的字段
field_value	作为value的字段
*/
func GetModelFieldDataMap(data []orm.Params, field_key string, field_value string) (data_map map[string][]interface{}) {
	data_map = make(map[string][]interface{}) // 这里必须要make一下，不然后面对data_map的元素赋值时会报错说不能对一个nil的map的元素赋值
	//var data_slice []interface{}			// 某个key下的数据
	for _, value := range data {
		if key, ok := value[field_key]; ok {
			key_str := key.(string)
			data_map[key_str] = append(data_map[key_str], value[field_value])
		}
	}

	return data_map
}

// 将 interface转为string
func TranInterface2String(value interface{}) (value_string string){
	switch v := value.(type) {
	case string:
		return v
	case int:
		return strconv.Itoa(v)
	default:
		panic("tranSliceInterface2String: params type not supported")
	}
}

// 将interface转为int
func TranInterface2Int(value interface{}) (value_int int){
	switch v := value.(type) {
	case string:
		value_int,_ = strconv.Atoi(v)
	case int:
		value_int = v
	default:
		panic("tranSliceInterface2String: params type not supported")
	}
	return value_int
}

// 将 []interface{}转为[]string
func TranSliceInterface2String(interface_slice []interface{}) (string_slice []string) {
	for _, itf := range interface_slice {
		string_slice = append(string_slice, TranInterface2String(itf))
	}
	return string_slice
}

// 将 []string 转为 []interface{}
func TranSliceString2Interface(string_slice []string) (interface_slice []interface{}){
	for _, val := range string_slice{
		interface_slice = append(interface_slice, val)
	}

	return interface_slice
}

// 将map[string]interface{}转为map[string]string
func TranMapInterface2String(inteface_map map[string]interface{}) (string_map map[string]string) {
	string_map = make(map[string]string)
	for key, value := range inteface_map {
		switch value.(type) {
		case string:
			string_map[key] = value.(string)
		case int:
			string_map[key] = strconv.Itoa(value.(int))
		case int64:
			string_map[key] = strconv.FormatInt(value.(int64), 10)
		}
	}
	return string_map
}

// map复合类型添加子元素
func AddMapInner(m map[interface{}]interface{}, outerKey string, innerKey string, value interface{}) {
	innerData := m[outerKey].(map[string]interface{})
	innerData[innerKey] = value
	m[outerKey] = innerData
}

// 将[]orm.Params转为 map[string]orm.Params形式
func TranDataSlice2Map(data []orm.Params, field_name string) (data_map map[string]orm.Params) {
	data_map = make(map[string]orm.Params)
	for _, ormParam := range data {
		if fn, ok := ormParam[field_name]; ok {
			if fn_int, is_int := fn.(int); is_int {
				fn = strconv.Itoa(fn_int)
			}
			if fn_int, is_int64 := fn.(int64); is_int64 {
				fn = strconv.Itoa(int(fn_int))
			}
			data_map[fn.(string)] = ormParam
		}
	}

	return data_map
}

// 切片元素去重
func UniqueSlice(s []interface{}) (unique_slice []interface{}) {
	tmp_map := make(map[interface{}]interface{})
	for _, value := range s {
		length := len(tmp_map)
		tmp_map[value] = struct{}{}
		if len(tmp_map) > length {
			unique_slice = append(unique_slice, value)
		}
	}
	return unique_slice
}

// 设置单个全局变量
func SetSingleGlobalDataForTemplate(data map[interface{}]interface{}, key string, value interface{}) {
	_, ok := data["global"]
	var global map[string]interface{}
	if !ok || data["global"] == nil {
		global = make(map[string]interface{})
	} else {
		global = data["global"].(map[string]interface{})
	}

	global[key] = value
	data["global"] = global
}

// 获取分页信息，返回当前页数，当前条数，总页数，总条数，上一页页数，下一页页数
func GetPageInfo(curPage int, pageRows int, totalRows int) map[string]interface{} {
	pageInfo := make(map[string]interface{})
	pageInfo["curPage"] = curPage
	totalPage := int(math.Ceil(float64(totalRows) / float64(pageRows)))

	pageInfo["totalPage"] = totalPage
	pageInfo["totalRow"] = totalRows
	pageInfo["pageRows"] = pageRows
	if curPage != totalPage {
		pageInfo["curRow"] = pageRows
	} else {
		pageInfo["curRow"] = totalRows % pageRows
	}
	if curPage-1 > 0 {
		pageInfo["prevPage"] = curPage - 1
	}
	if curPage+1 <= totalPage {
		pageInfo["nextPage"] = curPage + 1
	}
	return pageInfo
}

// 将int64转为int
func TranInt64ToInt(num_int64 int64) (num_int int) {
	num_str := strconv.FormatInt(num_int64, 10)
	num_int, _ = strconv.Atoi(num_str)
	return num_int
}


// 将map替换为get请求参数字符串
func BuildReqStr(params_map map[string]string) (params_str string){
	for key, val := range params_map{
		params_str += "&" + key + "=" + val
	}

	return strings.Trim(params_str, "&")
}

// 设置页面错误信息（info或error模板的变量信息）
func SetErrInfoForTemplate(c *beego.Controller, info string, others ...interface{}) {
	c.Data["info"] = info

	if len(others) >= 1{
		c.Data["redirect"] = others[0]
		c.TplName = "public/error_back.html"
	}

	if len(others) >= 2{
		c.TplName = others[1].(string)
	}
}

// 生成md5
func EncryptMd5(str string) (encrypt_str string) {
	data := []byte(str)
	has := md5.Sum(data)
	encrypt_str = fmt.Sprintf("%x", has) //将[]byte转成16进制
	return encrypt_str
}

// model拼接字段字符串
func BuildFieldStr(fields ...string) (field_str string){
	if len(fields) > 0{
		field_str = strings.Join(fields, ",")
	}else{
		field_str = "*"
	}
	return field_str
}

/**
 * 生成上传文件的相对项目根目录的绝对路径
 * @param fn 要上传的文件名，含后缀
 * @param destination 文件要上传到的目标目录路径（相对路径）, 如 "./statics/img"
 *
 * @return fp 图片上传后相对项目根目录的绝对路径
 * @return fp_real 图片上传后的物理绝对路径
 */
func CreateUploadPath(fn string, destination ...string) (fp string, fp_real string){
	var dir_path string
	if len(destination) == 0{
		dir_path = "./uploads/images/thumb/"
	}else{
		dir_path = destination[0]
	}

	root_path := beego.AppConfig.String("rootpath")
	dir_path_real := strings.TrimRight(root_path, "/") + "/" + strings.Trim(strings.TrimLeft(dir_path, "."), "/") + "/"
	if is_exist, _ := PathExists(dir_path_real); !is_exist {
		os.MkdirAll(dir_path_real, 0777)
	}
	fp = strings.TrimLeft(dir_path, ".") + fn
	fp_real = dir_path_real + fn
	return fp, fp_real
}

// 获取文件md5, 可传入文件名或者os.File类型
func GetFileMd5(file interface{}) (md5sum string, err error){
	var f io.ReadCloser
	var md5_byte []byte
	switch fp := file.(type) {
		case string:
			if is_exist, err := PathExists(fp); !is_exist{
				return "", err
			}
			f, err = os.Open(fp)
			defer f.Close()
		case *os.File:
			f = fp
		case multipart.File :
			f = fp
	}

	md5h := md5.New()
	_, err = io.Copy(md5h, f)
	if err != nil {
		return "", err
	}
	md5_byte = md5h.Sum(md5_byte)		// 生成md5

	return fmt.Sprintf("%x", md5_byte), err
}

// 判断一个文件或目录是否存在
func PathExists(path string) (is_exist bool, err error){
	_, err = os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

/**
 * 文件拷贝
 * @param file 文件对象
 * @param fp 目标文件绝对路径
 */
func CopyFile(file multipart.File, fp string) (err error){

	f_target, err := os.OpenFile(fp, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0766)
	if err != nil{
		return err
	}
	defer f_target.Close()


	file.Seek(0, 0)		// 必须先将要上传的文件指针调回开头
	for{
		var n int
		buf := make([]byte, 4096 * 4)		// 16kb的运输频率
		n, err = file.Read(buf)
		if n > 0{
			_, err = f_target.Write(buf[0:n])

		}else{
			if err == io.EOF{
				err = nil
			}
			break
		}
	}

	return err
}

/**
 * 通用上传文件的方法
 * @param file 要上传的文件
 * @param ext  文件的后缀,格式
 * @param destination 文件要上传到的目标目录路径（相对路径）, 如 "./statics/img"
 *
 * @return fp 图片上传后相对项目根目录的绝对路径
 * @notice 该方法不会关闭file资源，需要在调用UploadFile后手动关闭file资源
 */
func UploadFile(file multipart.File, ext string, destination string) (fp string, err error){
	// 获取文件md5
	file_md5, err := GetFileMd5(file)
	if err != nil {
		return
	}

	if file_md5 != ""{
		var fp_real string
		file_name := file_md5 + "." + ext
		fp, fp_real = CreateUploadPath(file_name, destination)		// 生成上传文件的绝对路径
		if file_exists, _ := PathExists(fp_real); !file_exists{
			err = CopyFile(file, fp_real)							// 上传文件
		}
	}

	return fp, err
}

/**
 * 上传指定名称的file控件到目标路径
 * @param file_input_name 上传文件的input标签的name属性值
 * @param destination 文件要上传到的目标目录路径（相对路径）, 如 "./statics/img"
 *
 * @return fp 图片上传后相对项目根目录的绝对路径
 */
func UploadFileByName(c *beego.Controller, file_input_name string, destination string) (fp string, err error){
	file, file_header, err := c.GetFile(file_input_name)
	defer file.Close()

	// 上传
	ext := strings.TrimLeft(path.Ext(file_header.Filename), ".")
	fp, err = UploadFile(file, ext, destination)

	return fp, err
}

// 取出html标签
func StripHtml(content string) string {
	reg, _ := regexp.Compile("</?[^>]+>")
	content_new := reg.ReplaceAllString(content, "")

	reg2, _ := regexp.Compile("&\\w+;")
	content_new = reg2.ReplaceAllString(content_new, "")
	return content_new
}

// 判断切片是否存在某元素(仅限字符串和整型两种类型)
func InSlice(ele interface{}, s []interface{})  bool{
	for _, s_ele:=range s{
		if s_ele == ele{
			return true
		}
	}
	return false
}

// 设置发生报错时的接口信息
func SetErrorApiResp(resp *response.CommonResponse, c *beego.Controller){
	if p := recover(); p != nil{
		logs.Error(fmt.Sprintf("error: %v", p))
		resp.SetError(errInfo.CODE_PROG_ERR)
	}
	c.Data["json"] = resp
	c.ServeJSON()
}