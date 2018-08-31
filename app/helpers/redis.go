package helpers

import (
	"theAmazingPostManager/app/common"
)

func InsertIntoCappedList (data []byte,listName string, limit int)error{

	newConn := common.RedisPool.Get()

	//Push element
	_,err := newConn.Do("LPUSH",listName,string(data))
	if err != nil{
		return err
	}

	//Trim list
	_,err = newConn.Do("LTRIM",listName,0,limit)
	if err != nil{
		return err
	}

	return nil

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