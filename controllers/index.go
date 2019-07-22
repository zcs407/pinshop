package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"pinshop/models"
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
	//1、首先将整个三级联动的菜单当做一个大的容器[]map[string]interface{[]map[string]interface{}}，那我们从后往前推
	//1.1、interface{}存的是二级菜单的[]map[string]interface{}
	//1.2、string指的是一级菜单的名字
	//1.3、map指的是将一级菜单的名字和二级菜单的map切片组合成一个map对象
	//1.4、切片指的是将多个一级菜单的对象放到一个切片中，方便前端的遍历
	var allMenusrq []map[string]interface{}
	//2.1、定义一个一级菜单切片,并查询到所有的一级菜单
	var firstMenusSlice []models.TpshopCategory
	o := orm.NewOrm()
	o.QueryTable("TpshopCategory").Filter("Pid", 0).All(&firstMenusSlice)
	//2.2根据一级菜单切片遍历查询到所有二级菜单的[]map[string]interface{}
	for _, firstMenu := range firstMenusSlice {
		var firstrq = make(map[string]interface{})
		var secondMenusSlice []models.TpshopCategory
		o.QueryTable("TpshopCategory").Filter("Pid", firstMenu.Id).All((&secondMenusSlice))
		//给一级菜单map赋值，一级菜单map内容包括一级菜单的对象和二级菜单的map切片
		firstrq["firstMenu"] = firstMenu
		firstrq["secondMenusSlice"] = secondMenusSlice

		//将一级容器放到大容器中
		allMenusrq = append(allMenusrq, firstrq)
	}
	//查询三级菜单并赋值给
	for _, first := range allMenusrq {
		var sectrq []map[string]interface{}

		for _, secondMap := range first["secondMenusSlice"].([]models.TpshopCategory) {
			var sanjiSlice []models.TpshopCategory
			var sanjirq = make(map[string]interface{})
			o.QueryTable("TpshopCategory").Filter("Pid", secondMap.Id).All(&sanjiSlice)
			sanjirq["secondMenusSlice"] = secondMap
			sanjirq["sajiMenus"] = sanjiSlice
			//将三级容器放到二级容器中
			sectrq = append(sectrq, sanjirq)
		}
		first["sajiMenus"] = sectrq
	}
	this.Data["allMenusrq"] = allMenusrq
	this.TplName = "index.html"
}
