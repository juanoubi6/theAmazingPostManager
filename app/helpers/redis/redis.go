package redis

import (
	"theAmazingPostManager/app/common"
	"github.com/sirupsen/logrus"
)

func InsertIntoCappedList (data []byte,listName string, limit int){

	newConn := common.RedisPool.Get()

	//Push element
	_,err := newConn.Do("LPUSH",listName,data)
	if err != nil{
		logrus.WithFields(logrus.Fields{
			"operation": "inserting value in list",
		}).Info(err.Error())
	}

	//Trim list
	_,err = newConn.Do("LTRIM",listName,0,limit)
	if err != nil{
		logrus.WithFields(logrus.Fields{
			"operation": "triming list",
		}).Info(err.Error())
	}

}

func RetrieveFromCappedList(listName string,amount int)([]interface{},error){

	newConn := common.RedisPool.Get()

	//Get elements
	values,err := newConn.Do("LRANGE",listName,0,amount-1)
	if err != nil{
		return []interface{}{},err
	}

	return values.([]interface{}),nil

}