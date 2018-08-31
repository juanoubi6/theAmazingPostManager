package common

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"theAmazingPostManager/app/config"
	"github.com/gomodule/redigo/redis"
	"time"
)

var db *gorm.DB
var RedisPool *redis.Pool

func ConnectToDatabase() {
	var err error
	dbname := config.GetConfig().DB_NAME
	dbhost := config.GetConfig().DB_HOST
	dbport := config.GetConfig().DB_PORT
	dbuser := config.GetConfig().DB_USERNAME
	dbpass := config.GetConfig().DB_PASSWORD

	db, err = gorm.Open("mysql", dbuser+":"+dbpass+"@"+"tcp("+dbhost+":"+dbport+")"+"/"+dbname+"?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		panic(err)
	}

}

func GetDatabase() *gorm.DB {
	return db
}

func CreateRedisConnectionPool() {

	RedisPool = &redis.Pool{
		MaxIdle:   3,
		MaxActive: 10, // max number of connections
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", config.GetConfig().REDIS_ADDR)
			if err != nil {
				panic(err.Error())
			}
			return c, err
		},
	}

	newConn := RedisPool.Get()
	defer newConn.Close()

	_,err := newConn.Do("ping")
	if err != nil{
		panic("Couldn't check redis pool connection: " + err.Error())
	}

}