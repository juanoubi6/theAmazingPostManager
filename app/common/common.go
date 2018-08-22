package common

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/nats-io/go-nats"
	"github.com/streadway/amqp"
	"theAmazingPostManager/app/config"
)

var db *gorm.DB

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



