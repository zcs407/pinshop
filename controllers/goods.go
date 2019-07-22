package controllers

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"pinshop/models"
)

type GoodsController struct {
	beego.Controller
}

func (this *GoodsController) ShowIndexSX() {
	//获取所有类型
	o := orm.NewOrm()
	var goodsTypes []models.GoodsType
	o.QueryTable("GoodsType").All(&goodsTypes)
	this.Data["goodsTypes"] = goodsTypes
	//获取所有轮播图
	var indexGoodsBanners []models.IndexGoodsBanner
	o.QueryTable("IndexGoodsBanner").RelatedSel("GoodsSKU").OrderBy("Index").All(&indexGoodsBanners)
	this.Data["indexGoodsBanners"] = indexGoodsBanners
	//获取所有促销图片oodsBanner struct { //首页分类商品展示表
	//	Id
	var indexPromotionBanners []models.IndexPromotionBanner
	o.QueryTable("IndexPromotionBanner").OrderBy("Index").All(&indexPromotionBanners)
	this.Data["indexPromotionBanners"] = indexPromotionBanners

	//获取所有首页图片
	var allGoods []map[string]interface{}
	for _, goodsType := range goodsTypes {
		//根据首页分类商品表获取到文字类的商品信息和图片类的商品信息

		qs := o.QueryTable("IndexTypeGoodsBanner").RelatedSel("GoodsType", "GoodsSKU")
		//获取文字类商品信息，DisplayType 0代表文字，1代表图片
		var textGoods []models.IndexTypeGoodsBanner
		qs.Filter("GoodsType__Id", goodsType.Id).Filter("DisplayType", 0).All(&textGoods)

		//获取图片类型
		var imgGoods []models.IndexTypeGoodsBanner
		qs.Filter("GoodsType__Id", goodsType.Id).Filter("DisplayType", 1).All(&imgGoods)

		tempCont := make(map[string]interface{})
		tempCont["goodsTypes"] = goodsType
		tempCont["textGoods"] = textGoods
		tempCont["imgGoods"] = imgGoods

		allGoods = append(allGoods, tempCont)
	}
	this.Data["allGoods"] = allGoods

	this.TplName = "index_sx.html"
}

func (this *GoodsController) ShowSxDetail() {
	id, err := this.GetInt("id")
	if err != nil {
		panic(err)
	}
	fmt.Println("生鲜点击后的id是：", id)

	this.TplName = "detail.html"
}
