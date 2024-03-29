package models

import (
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

func init() {
	//获取DB的配置文件
	dbInfo := DbConfigInfo("conf/db.conf")
	//注册数据库
	orm.RegisterDataBase("default", "mysql", dbInfo["dbuser"]+":"+
		dbInfo["dbpwd"]+"@("+dbInfo["dbhost"]+":"+dbInfo["dbport"]+")/"+dbInfo["dbname"]+"?loc=Local")

	//注册表
	orm.RegisterModel(new(User), new(Address), new(TpshopCategory), new(Goods), new(GoodsType),
		new(GoodsSKU), new(GoodsImage), new(IndexGoodsBanner), new(IndexTypeGoodsBanner),
		new(IndexPromotionBanner), new(OrderInfo), new(OrderGoods))
	//创建表
	orm.RunSyncdb("default", false, true)
}
