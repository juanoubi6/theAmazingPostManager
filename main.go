package main

import (
	"theAmazingPostManager/app/common"
	"theAmazingPostManager/app/router"
)

func main() {
	common.ConnectToDatabase()
	common.CreateRedisConnectionPool()
	router.CreateRouter()
	router.RunRouter()
}
