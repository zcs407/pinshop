package controllers

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/gomodule/redigo/redis"
	"pinshop/models"
	"strconv"
)

type CarController struct {
	beego.Controller
}

func respFunc(this *CarController, resp map[string]interface{}) {
	this.Data["json"] = resp
	this.ServeJSON()
}

//添加购物车数据
func (this *CarController) AddCart() {
	fmt.Println("已经调用添加购物车函数")
	skuid, err := this.GetInt("skuid")
	fmt.Println("skuid=========", skuid)
	if err != nil {
		panic(err)
	}
	count, err := this.GetInt("count")
	if err != nil {
		panic(err)
	}
	resp := make(map[string]interface{})
	defer respFunc(this, resp)

	//判断用户是否登录
	if this.GetSession("userName") == nil {
		resp["status"] = 403
		resp["msg"] = "用户没有登录"
	}
	//用redis的hash存储购物车数据
	conn, err := redis.Dial("tcp", "172.16.10.11:6379")
	if err != nil {
		resp["status"] = 401
		resp["msg"] = "redis connect failed"
		return
	}
	//获取username 写入数据库
	username := this.GetSession("userName").(string)
	//

	result, _ := redis.Int(conn.Do("hget", username+"_cart", skuid))
	_, err = redis.Int(conn.Do("hset", username+"_cart", skuid, count+result))
	resp["status"] = 200
	resp["msg"] = "OK"
}

func (this *CarController) ShowCart() {
	var totalPrice, totalNum int
	conn, err := redis.Dial("tcp", "172.16.10.11:6379")
	if err != nil {
		panic(err)
	}
	username := this.GetSession("userName").(string)
	result, err := redis.IntMap(conn.Do("hgetall", username+"_cart"))
	if err != nil {
		panic(err)
	}
	var goods []map[string]interface{}
	for skuid, count := range result {
		goodsMap := make(map[string]interface{})
		var goodSku models.GoodsSKU
		skuId, _ := strconv.Atoi(skuid)
		goodSku.Id = skuId
		o := orm.NewOrm()
		o.Read(&goodSku)
		goodsMap["goodsSku"] = goodSku
		goodsMap["count"] = count
		goodsMap["littleSum"] = goodSku.Price * count

		totalPrice += goodSku.Price * count
		totalNum += count
		goods = append(goods, goodsMap)
	}

	this.Data["goods"] = goods
	this.Data["totalPrice"] = totalPrice
	this.Data["totalNum"] = totalNum
	this.TplName = "cart.html"
}

func (this *CarController) DeleteCart() {
	userName := this.GetSession("userName")
	skuid, err := this.GetInt("skuid")
	resp := make(map[string]interface{})
	defer respFunc(this, resp)
	if userName == nil {
		resp["status"] = 400
		resp["msg"] = "用户登录失效"
		return
	}
	if err != nil {
		resp["status"] = 401
		resp["msg"] = "获取商品id出错"
		return
	}
	//connect redis server
	conn, err := redis.Dial("tcp", "172.16.10.11:6379")
	if err != nil {
		resp["status"] = 402
		resp["msg"] = "连接redis失败"
		return
	}
	_, err = conn.Do("hdel", userName.(string)+"_cart", skuid)
	if err != nil {
		resp["status"] = 403
		resp["msg"] = "删除redis商品失败"
		return
	}
	resp["status"] = 200
	resp["msg"] = "OK"
}
