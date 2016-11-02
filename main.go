package main

import (
	"kylin/controllers"
	_ "kylin/routers"

	"github.com/astaxie/beego"
)

func main() {
	beego.ErrorController(&controllers.ErrorController{})
	beego.Run()
}
