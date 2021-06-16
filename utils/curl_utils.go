package utils

import (
	"bytes"
	"encoding/json"
	"github.com/astaxie/beego/logs"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

type Curl struct{
	Url string
	Header map[string]string
	Params map[string]string
	Data interface{}
	ContentType string
}


func (curl *Curl) Get() (body []byte, err error){
	client := &http.Client{}
	queryStr := ""
	for key, value := range curl.Params{
		queryStr += key + "=" + value + "&"
	}
	queryStr = strings.TrimRight(queryStr, "&")
	if queryStr != "" {
		if strings.Contains(curl.Url, "?"){
			curl.Url = curl.Url + "&" +queryStr
		}else{
			curl.Url = curl.Url + "?" +queryStr
		}
	}

	req, err := http.NewRequest("GET", curl.Url, nil)
	if err != nil{
		return nil, err
	}

	curl.setHeader(req)

	resp, err := client.Do(req)		// 发起请求
	if err != nil{
		return nil, err
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

func (curl *Curl) Post() (body []byte, err error){
	client := &http.Client{}

	request_body, err := curl.setPostData()	// 将要传输的字节封装为io.Reader类型，因为下面的http.NewRequest方法要求传一个io.Reader类型
	if err != nil{
		return nil, err
	}

	req, err := http.NewRequest("POST", curl.Url, request_body)		// method必须大写，要写为POST不能协程Post
	if err != nil{
		return nil, err
	}

	curl.setHeader(req)

	resp, err := client.Do(req)		// 发起请求
	if err != nil{
		return nil, err
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

// 在 Curl.Get 方法的基础上，响应结果的 []byte类型 转为 string
func (curl *Curl) GetForString() (resp_str string, err error){
	resp_bytes, err := curl.Get()
	resp_str = string(resp_bytes)
	return resp_str, err
}

// 在 Curl.Get 方法的基础上，响应结果的 []byte类型 转为 string
func (curl *Curl) PostForString() (resp_str string, err error){
	resp_bytes, err := curl.Post()
	resp_str = string(resp_bytes)
	return resp_str, err
}

// 在 Curl.Post 方法的基础上，响应结果的 []byte类型的json 转为 map 格式
func (curl *Curl) PostForJson() (resp_json_map map[string]interface{}, err error){
	resp_bytes, resp_err := curl.Post()
	resp_json_map, handle_err := curl.handleRespForJson(resp_bytes)
	if resp_err != nil {
		err = resp_err
	}else if handle_err != nil{
		err = handle_err
	}
	return resp_json_map, err
}

// 在 Curl.Get 方法的基础上，响应结果的 []byte类型的json 转为 map 格式
func (curl *Curl) GetForJson() (resp_json_map map[string]interface{}, err error){
	resp_bytes, resp_err := curl.Get()
	resp_json_map, handle_err := curl.handleRespForJson(resp_bytes)
	if resp_err != nil {
		err = resp_err
	}else if handle_err != nil{
		err = handle_err
	}
	return resp_json_map, err
}

// 将响应结果的 []byte类型的json 转为 map 格式
func (curl *Curl) handleRespForJson(resp_bytes []byte) (resp_json_map map[string]interface{}, err error){
	if len(resp_bytes) > 0{
		err = json.Unmarshal(resp_bytes, &resp_json_map)
		if err != nil{
			logs.Error("PostForJson - Unmarshal failed:" + err.Error())
		}
	}
	return resp_json_map, err
}

func (curl *Curl) setHeader(req *http.Request){
	for header_key, header_val := range curl.Header{
		req.Header.Set(header_key, header_val)
	}
	if curl.ContentType != ""{
		req.Header.Set("Content-Type", curl.ContentType)
	}
}

// post请求根据不同的content-type处理数据
func (curl *Curl) setPostData() (post_data io.Reader, err error){
	var content_type string
	if curl.ContentType == ""{
		content_type = curl.Header["Content-Type"]
	}else{
		content_type = curl.ContentType
	}
	if content_type == "" || content_type == "application/x-www-form-urlencoded" || content_type == "application/json"{
		data_bytes, err := json.Marshal(curl.Data)			// 将map类型的curl.Data转为[]byte类型的json
		if err != nil{
			return nil, err
		}

		post_data = bytes.NewReader(data_bytes)
	}else if content_type == "text/plain"{
		post_data = bytes.NewReader([]byte(curl.Data.(string)))    // 将string类型的curl.Data转为[]byte
	}

	return post_data, err
}