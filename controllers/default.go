package controllers

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"pinshop/models"
)

type MainController struct {
	beego.Controller
}

func (this *MainController) ShowIndex() {
	id, _ := this.GetInt("id")
	fmt.Println("登录后获得的id是：", id)
	o := orm.NewOrm()
	var user models.User
	user.Id = id
	o.Read(&user)
	this.Data["userName"] = user.Name

	this.TplName = "index.html"
}
