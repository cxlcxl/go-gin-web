package main

import (
	_ "gin-web/bootstrap"
	"gin-web/router"
	"log"
)

func main() {
	if err := router.Router(); err != nil {
		log.Fatal("web 启动失败：", err)
	}
}
