package controllers

import (
	"fmt"
	"github.com/astaxie/beego"
)

type UserController struct {
	beego.Controller
}

func (this *UserController) ShowLogin() {
	this.TplName = "login.html"
}

func (this *UserController) Login() {
	fmt.Println("")
}
func (this *UserController) ShowRegister() {
	this.TplName = "register.html"
}

func (this *UserController) HandleCreateUser() {
	fmt.Println("")
}
