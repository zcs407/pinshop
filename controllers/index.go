package controllers

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"pinshop/models"
)

type Index struct {
	beego.Controller
}

func (this *Index) ShowIndex() {
	id, _ := this.GetInt("id")
	o := orm.NewOrm()
	var user models.User
	user.Id = id
	o.Read(&user, "id")
	this.Data["userName"] = user.Name
	fmt.Println("===================", id, user.Name)
	this.TplName = "index.html"
}
