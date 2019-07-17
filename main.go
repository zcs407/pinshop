package main

import (
	"github.com/astaxie/beego"
	_ "pinshop/models"
	_ "pinshop/routers"
)

func main() {
	beego.Run()
}
