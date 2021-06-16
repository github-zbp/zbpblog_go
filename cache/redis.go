package cache

import (
	"github.com/astaxie/beego"
	"github.com/garyburd/redigo/redis"
	"strconv"
	"time"
)

var RedisPool *redis.Pool

// 创建redis连接池
func InitRedis() *redis.Pool {
	host := beego.AppConfig.String("redis_host")
	password := beego.AppConfig.String("redis_password")
	max_idle,_ := strconv.Atoi(beego.AppConfig.String("redis_max_idle"))
	max_active,_ := strconv.Atoi(beego.AppConfig.String("redis_max_active"))
	RedisPool = &redis.Pool{
		MaxIdle: max_idle,
		IdleTimeout: 240 * time.Second,
		MaxActive: max_active,
		Dial: func() (conn redis.Conn, err error){
			conn, err = redis.Dial("tcp", host)
			if err != nil{
				return conn, err
			}
			if password != ""{
				if _, err := conn.Do("AUTH", password); err != nil{
					conn.Close()
					return conn, err
				}
			}
			return conn, err
		},
		TestOnBorrow: func(conn redis.Conn, t time.Time) (err error){
			_, err = conn.Do("PING")
			return err
		},
	}
	return RedisPool
}

