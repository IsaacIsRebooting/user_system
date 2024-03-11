package main

import (
	"my_user_system/conf"
	"my_user_system/router"
)

func Init() {
	conf.InitConfig()
}

func main() {
	Init()
	router.InitRouterAndServe()
}
