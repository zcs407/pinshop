package controllers

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/gomodule/redigo/redis"
	"pinshop/models"
	"strconv"
)

type OrderController struct {
	beego.Controller
}

func (this *OrderController) CommitOrder() {
	skuids := this.GetStrings("skuidid")
	if len(skuids) == 0 {
		fmt.Println("skuids:======", len(skuids))
		this.Redirect("/st/cart", 302)
		return
	}
	username := this.GetSession("userName")
	if username == nil {
		panic("用户没登录")
	}
	var addrs []models.Address
	o := orm.NewOrm()
	o.QueryTable("Address").RelatedSel("User").Filter("User__Name", username.(string)).All(&addrs)
	//连接redis
	conn, err := redis.Dial("tcp", "172.16.10.11:6379")
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	//定义大容器，接收商品的信息，商品的数量，商品的小计
	//定义商品数量和小计的变量
	var allgoods []map[string]interface{}
	var totalCount, totalPrice, truePrice int
	//循环获取相关数据
	for _, skuid := range skuids {
		tempMap := make(map[string]interface{})
		var goodsku models.GoodsSKU
		id, _ := strconv.Atoi(skuid)
		goodsku.Id = id
		o.Read(&goodsku)
		//获取redis中商品数量
		count, err := redis.Int(conn.Do("hget", username.(string)+"_cart", id))
		if err != nil {
			panic(err)
		}

		totalCount += count
		totalPrice += count * goodsku.Price

		tempMap["count"] = count
		tempMap["goodsku"] = goodsku
		tempMap["subTotal"] = goodsku.Price * count
		allgoods = append(allgoods, tempMap)

	}
	truePrice = totalPrice + 10
	this.Data["kuaidi"] = 10
	this.Data["totalCount"] = totalCount
	this.Data["totalPrice"] = totalPrice
	this.Data["truePrice"] = truePrice
	this.Data["allgoods"] = allgoods
	this.Data["addrs"] = addrs
	this.TplName = "place_order.html"
}
