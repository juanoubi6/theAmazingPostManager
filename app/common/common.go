package common

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"theAmazingPostManager/app/config"
	"github.com/gomodule/redigo/redis"
	"time"
	"errors"
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

	// If redis server is not up, there won't be a problem because the application doesn't depend on redis. A fake
	// connection was mocked, so when the pool cannot connect to redis, the fake connection is returned.
	RedisPool = &redis.Pool{
		MaxIdle:   3,
		MaxActive: 10, // max number of connections
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", config.GetConfig().REDIS_ADDR)
			if err != nil {
				return FakeRedisConn{},nil
			}
			return c, err
		},
	}

	// This code would be used when we are sure there is a redis server working and we want to make
	// sure that the connection was established. Also, we should remove the FakeRedisConn{] from the
	// pool Dial function and replace that line with a panic

	/*newConn := RedisPool.Get()
	defer newConn.Close()

	_,err := newConn.Do("ping")
	if err != nil{
		panic("Couldn't check redis pool connection: " + err.Error())
	}
	*/
}

type FakeRedisConn struct{}

func (f FakeRedisConn) Close() error{
	return errors.New("Fake error")
}
func (f FakeRedisConn) Err() error{
	return errors.New("Fake error")
}
func (f FakeRedisConn) Flush() error{
	return errors.New("Fake error")
}
func (f FakeRedisConn) Do(commandName string,args ...interface{}) (reply interface{},err error){
	var inter interface{}
	return inter,errors.New("Fake error")
}
func (f FakeRedisConn) Send(commandName string,args ...interface{}) (error){
	return errors.New("Fake error")
}
func (f FakeRedisConn) Receive() (reply interface{},err error){
	var inter interface{}
	return inter,errors.New("Fake error")
}