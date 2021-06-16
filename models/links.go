package models

import (
	_ "github.com/go-sql-driver/mysql"		// 记得引入这个驱动
)

//+-------+--------------+------+-----+---------+----------------+
//| Field | Type         | Null | Key | Default | Extra          |
//+-------+--------------+------+-----+---------+----------------+
//| id    | int(11)      | NO   | PRI | NULL    | auto_increment |
//| name  | varchar(255) | NO   |     | NULL    |                |
//| url   | varchar(255) | NO   |     | NULL    |                |
//| title | varchar(255) | YES  |     | NULL    |                |
//+-------+--------------+------+-----+---------+----------------+

type Links struct {
	Id int `orm:"column(id)"`
	Name string `orm:"column(Name)"`
	Title string `orm:"column(title)"`
	Url string `orm:"column(url)"`
}