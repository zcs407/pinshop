package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"pinshop/controllers"
)

func init() {
	beego.InsertFilter("/st/*", beego.BeforeExec, filters)

	beego.Router("/", &controllers.Index{}, "get:ShowIndex")
	//创建用户注册路由
	beego.Router("/register", &controllers.UserController{}, "get:ShowRegister;post:HandleCreateUser")
	beego.Router("/codeSend", &controllers.UserController{}, "post:CodeSend")
	//创建用户登录路由
	beego.Router("/login", &controllers.UserController{}, "get:ShowLogin;post:HandleLogin")
	beego.Router("/logout", &controllers.UserController{}, "get:ShowLogout")

	beego.Router("/active", &controllers.UserController{}, "get:ShowActive;post:ActiveEmail")
	beego.Router("/activeUser", &controllers.UserController{}, "get:ActiveUser")

	//用户信息中心
	beego.Router("/st/userCenterInfo", &controllers.UserController{}, "get:ShowUserCenterInfo")
	//用户地址管理
	beego.Router("/st/user_center_site", &controllers.UserController{}, "get:ShowUserAddress;post:HandleAddress")
	//生鲜页面展示
	beego.Router("/index_sx", &controllers.GoodsController{}, "get:ShowIndexSX")
	//生鲜信息页面展示
	beego.Router("/sxDetail", &controllers.GoodsController{}, "get:ShowSxDetail")
	//生鲜列表页面展示
	beego.Router("/sxlist", &controllers.GoodsController{}, "get:ShowSxList")
	//搜索商品
	beego.Router("/search", &controllers.GoodsController{}, "post:SearchGoods")
	//显示购物车
	beego.Router("/st/cart", &controllers.CarController{}, "get:ShowCart")
	//加入购物车
	beego.Router("/addCart", &controllers.CarController{}, "post:AddCart")
	//删除购物车商品
	beego.Router("/deleteCart", &controllers.CarController{}, "post:DeleteCart")
	//处理订单
	beego.Router("/st/commitOrder", &controllers.OrderController{}, "post:CommitOrder")
}

func filters(cxt *context.Context) {
	userName := cxt.Input.Session("userName")
	if userName == nil {
		cxt.Redirect(302, "/login")
		return
	}
}
