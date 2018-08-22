package main

import (
	"theAmazingPostManager/app/common"
	"theAmazingPostManager/app/router"
)

func main() {
	common.ConnectToDatabase()
	router.CreateRouter()
	router.RunRouter()
}
