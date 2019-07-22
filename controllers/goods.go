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

func splitPage(setPage, count, pageSize, pageIndex int) (pages []int, prePage, nextPage int) {
	setPage = setPage
	if setPage <= 1 {
		setPage = 1
	}
	size := pageSize
	pageCount := count / size
	prePage = pageIndex - 1
	if prePage <= 1 {
		prePage = 1
	}
	nextPage = pageIndex + 1
	if nextPage >= pageCount {
		nextPage = pageCount
	}
	if pageCount <= setPage {
		for i := 1; i <= pageCount; i++ {
			pages = append(pages, i)
		}
	} else if pageIndex <= 3 {
		for i := 1; i <= setPage; i++ {
			pages = append(pages, i)
		}

	} else if pageIndex >= pageCount-2 {
		for i := pageCount - setPage + 1; i <= pageCount; i++ {
			pages = append(pages, i)
		}
	} else {
		for i := pageIndex - 2; i <= pageIndex+2; i++ {
			pages = append(pages, i)
		}
	}

	return pages, prePage, nextPage
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
	o := orm.NewOrm()
	//获取所有商品类型
	var goodsTypes []models.GoodsType
	o.QueryTable("GoodsType").All(&goodsTypes)

	//商品信息展示SKU
	var goods models.GoodsSKU

	o.QueryTable("GoodsSKU").RelatedSel("Goods", "GoodsType").Filter("Id", id).One(&goods)
	//新品展示，根据商品类型获取此类所有商品
	var newGoods []models.GoodsSKU
	qs := o.QueryTable("GoodsSKU").RelatedSel("GoodsType").Filter("GoodsType__GoodsSKU", goods)
	qs.OrderBy("-Time").Limit(2, 0).All(&newGoods)
	//一类商品的信息SPU，外键时商品的类型
	//o.QueryTable("GoodsSKU").RelatedSel("Goods").Filter("")

	this.Data["goodsTypes"] = goodsTypes
	this.Data["goods"] = goods
	this.Data["newGoods"] = newGoods
	fmt.Println(goods.Name)
	this.TplName = "detail.html"
}

func (this *GoodsController) ShowSxList() {
	//分页显示
	//设置页面显示产品的数量
	pageSize := 1
	goodsTypeId, err := this.GetInt("id")
	if err != nil {
		panic(err)
	}
	sort := this.GetString("sort")
	pageIndex, err := this.GetInt("pageIndex")
	if err != nil {
		pageIndex = 1
	}

	o := orm.NewOrm()
	//按商品类型查询所有商品，sku表
	var goodsSkus []models.GoodsSKU
	qs := o.QueryTable("GoodsSKU").RelatedSel("GoodsType").Filter("GoodsType__Id", goodsTypeId)
	qs.All(&goodsSkus)

	//调用分页函数
	count, _ := qs.Count()
	pages, prePage, nextPage := splitPage(6, int(count), pageSize, pageIndex)
	//查出相关类型的最新的两个商品
	var newGoodsSkus []models.GoodsSKU
	o.QueryTable("GoodsSKU").RelatedSel("GoodsType").Filter("GoodsType__Id", goodsTypeId).OrderBy("-Time").Limit(2, 0).All(&newGoodsSkus)

	if sort == "" {
		o.QueryTable("GoodsSKU").RelatedSel("GoodsType").Filter("GoodsType__Id", goodsTypeId).Limit(pageSize, pageIndex-1).All(&goodsSkus)
	} else if sort == "price" {
		o.QueryTable("GoodsSKU").RelatedSel("GoodsType").Filter("GoodsType__Id", goodsTypeId).OrderBy("Price").Limit(pageSize, pageIndex-1).All(&goodsSkus)
	} else {
		o.QueryTable("GoodsSKU").RelatedSel("GoodsType").Filter("GoodsType__Id", goodsTypeId).OrderBy("Sales").Limit(pageSize, pageIndex-1).All(&goodsSkus)
	}

	this.Data["sort"] = sort
	this.Data["typeId"] = goodsTypeId
	this.Data["newGoodsSkus"] = newGoodsSkus
	this.Data["goodsSkus"] = goodsSkus
	this.Data["pages"] = pages
	this.Data["pageIndex"] = pageIndex
	this.Data["prePage"] = prePage
	this.Data["nextPage"] = nextPage
	this.TplName = "list.html"
}
