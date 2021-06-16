package models

import (
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"strconv"
	"strings"
	"zbpblog_go/utils"
)

type TagsBlogs struct {
	TagsId  int `orm:"column(tags_id);pk"`
	BlogsId int `orm:"column(blogs_id)"`
}

func (tb *TagsBlogs) GetTBByTagIds(tag_ids []interface{}) (tb_data []orm.Params, err error) {
	tag_ids_str := strings.Join(utils.TranSliceInterface2String(tag_ids), ",")
	sql := fmt.Sprintf("select * from tags_blogs where tags_id in (%s)", tag_ids_str)
	_, err = Orm.Raw(sql).Values(&tb_data)
	return tb_data, err
}

func (tb *TagsBlogs) GetTBByBlogIds(blog_ids []interface{}) (tb_data []orm.Params, err error) {
	blog_ids_str := strings.Join(utils.TranSliceInterface2String(blog_ids), ",")
	sql := fmt.Sprintf("select * from tags_blogs where blogs_id in (%s)", blog_ids_str)
	_, err = Orm.Raw(sql).Values(&tb_data)
	return tb_data, err
}

// 为一篇博客绑定多个标签
func BindTagsBlog(blog_id int, tag_ids []string) (insert_rows int, succ_rows int64 ,err error){
	slice_tags_blogs := []TagsBlogs{}
	for _, tag_id := range tag_ids{
		tag_id_int, _ := strconv.Atoi(tag_id)
		m_tags_blogs := TagsBlogs{
			TagsId: tag_id_int,
			BlogsId: blog_id,
		}
		slice_tags_blogs = append(slice_tags_blogs, m_tags_blogs)
	}

	// 先删除后新增
	err = DelTagsBlog(blog_id)
	if err != nil{
		logs.Error("BindTagsBlog - DelTagsBlog failed:" + err.Error())
		return
	}

	insert_rows = len(slice_tags_blogs)
	if(insert_rows > 0){
		succ_rows, err =Orm.InsertMulti(insert_rows, slice_tags_blogs)
	}
	return insert_rows, succ_rows, err
}

// 为一篇博客删除多个标签
func DelTagsBlog(blog_id int, tag_ids ...string) (err error){
	var sql_delete string
	if len(tag_ids) == 0{
		sql_delete = fmt.Sprintf("delete from tags_blogs where blogs_id = %d", blog_id)
	}else{
		tag_ids_str := strings.Join(tag_ids, ",")
		sql_delete = fmt.Sprintf("delete from tags_blogs where blogs_id = %d and tags_id in (%s)", blog_id, tag_ids_str)
	}
	_, err = Orm.Raw(sql_delete).Exec()
	return err
}