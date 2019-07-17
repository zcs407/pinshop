package routers

import (
	"github.com/astaxie/beego"
	"pinshop/controllers"
)

func init() {
	beego.Router("/", &controllers.MainController{}, "get:ShowIndex")
	//创建用户注册路由
	beego.Router("/register", &controllers.UserController{}, "get:ShowRegister;post:HandleCreateUser")
	//创建用户登录路由
	beego.Router("/login", &controllers.UserController{}, "get:ShowLogin;post:Login")
}
