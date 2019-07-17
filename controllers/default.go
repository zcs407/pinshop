package controllers

import (
	"github.com/astaxie/beego"
)

type MainController struct {
	beego.Controller
}

func (c *MainController) ShowIndex() {
	c.TplName = "index.html"
}
