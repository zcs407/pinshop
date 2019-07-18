package routers

import (
	"github.com/astaxie/beego"
	"pinshop/controllers"
)

func init() {
	beego.Router("/", &controllers.Index{}, "get:ShowIndex")
	//创建用户注册路由
	beego.Router("/register", &controllers.UserController{}, "get:ShowRegister;post:HandleCreateUser")
	beego.Router("/codeSend", &controllers.UserController{}, "post:CodeSend")
	//创建用户登录路由
	beego.Router("/login", &controllers.UserController{}, "get:ShowLogin;post:HandleLogin")
	beego.Router("/active", &controllers.UserController{}, "get:ShowActive;post:ActiveEmail")
	beego.Router("/activeUser", &controllers.UserController{}, "get:ActiveUser")
}
