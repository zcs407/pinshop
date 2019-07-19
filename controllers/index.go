package controllers

import (
	"github.com/astaxie/beego"
)

type Index struct {
	beego.Controller
}

func (this *Index) ShowIndex() {
	//session查询方式
	userName := this.GetSession("userName")
	if userName != nil {
		this.Data["userName"] = userName.(string)
	}

	//数据库查询方式
	//id, _ := this.GetInt("id")
	//o := orm.NewOrm()
	//var user models.User
	//user.Id = id
	//o.Read(&user, "id")
	//this.Data["userName"] = user.Name
	//fmt.Println("===================", id, user.Name)

	this.TplName = "index.html"
}
