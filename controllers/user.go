package controllers

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/utils"
	"github.com/gomodule/redigo/redis"
	"math/rand"
	"pinshop/models"
	"regexp"
	"strconv"
	"time"
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
	phone := this.GetString("phone")
	vscode := this.GetString("code")
	pwd := this.GetString("password")
	repwd := this.GetString("repassword")
	if phone == "" || vscode == "" || pwd == "" || repwd == "" {
		fmt.Println("信息不完整")
		this.Data["errmsg"] = "信息不完整,请重新填写"
		this.Redirect("/register", 302)
		return
	}
	if pwd != repwd {
		this.Data["errmsg"] = "两次密码不一致！"
		this.Redirect("/register", 302)
		return
	}
	conn, err := redis.Dial("tcp", "172.16.10.11:6379")
	if err != nil {
		this.Data["errmsg"] = "redis connect err"
		this.TplName = "register.html"
		return
	}
	result, err := redis.String(conn.Do("get", phone+"_code"))
	if err != nil {
		this.Data["errmsg"] = "redis don't has code"
		this.TplName = "register.html"
		return
	}
	fmt.Println("redis code is", result, "用户输入的是", vscode)
	if result != vscode {
		this.Data["errmsg"] = "验证码不一致"
		this.TplName = "register.html"
		return
	}

	//将新用户插入数据库
	o := orm.NewOrm()
	var user models.User
	user.Name = phone
	user.PassWord = pwd
	id, err := o.Insert(&user)

	if err != nil {
		this.Data["errmsg"] = "用户已存在"
		this.TplName = "register.html"
		return
	}
	this.Redirect("/active?id="+strconv.Itoa(int(id)), 302)

}
func (this *UserController) CodeSend() {
	rand.Seed(time.Now().UnixNano())
	vsCode := fmt.Sprintf("%05d", rand.Int31n(100000))
	fmt.Println(vsCode)
	phone := this.GetString("phone")
	if phone == "" {
		fmt.Println("电话号码为空")
		return
	}
	//电话格式校验
	reg, _ := regexp.Compile(`^1[3-9][0-9]{9}$`)
	result := reg.FindString(phone)
	if result == "" {
		fmt.Println("电话号码不匹配")
		return
	}
	//merr := MsgSend(phone, vsCode)
	//if merr != nil {
	//	fmt.Println("发送验证码失败")
	//	return
	//}

	conn, merr := redis.Dial("tcp", "172.16.10.11:6379")
	if merr != nil {
		resp := make(map[string]interface{})
		resp["statusCode"] = 401
		resp["msg"] = "redis链接失败"
		this.Data["json"] = resp
		this.ServeJSON()
	}
	defer conn.Close()
	conn.Do("setex", phone+"_code", 60*50, vsCode)
	this.ServeJSON()
	return
}
func (this *UserController) ShowActive() {
	id, err := this.GetInt("id")
	if err != nil {
		this.Data["errmsg"] = "获取id不正确"
		return
	}
	this.Data["id"] = id
	this.TplName = "register-email.html"
}
func (this *UserController) ActicveEmail() {
	id, err := this.GetInt("id")
	if err != nil {
		this.Data["errmsg"] = "获取id不正确"
		return
	}
	email := this.GetString("email")
	if err != nil || email == "" {
		this.Redirect("/active?id="+strconv.Itoa(id), 302)
		return
	}
	//校验邮箱格式
	reg, _ := regexp.Compile(`^[a-zA-Z0-9_-]+@[a-zA-Z0-9_-]+(\.[a-zA-Z0-9_-]+)+$`)
	result := reg.FindString(email)
	if result == "" {
		fmt.Println("邮箱格式不正确")
		this.Redirect("/active?id="+strconv.Itoa(id), 302)
		return
	}
	//发送邮件激活
	config := `{"username":"563364657@qq.com","password":"jabebzakmlqobaii","host":"smtp.qq.com","port":587}`
	sendEmail := utils.NewEMail(config)
	sendEmail.From = "563364657@qq.com"
	sendEmail.To = []string{email}
	sendEmail.Subject = "品优购用户激活"
	sendEmail.HTML = `<a href="http://127.0.0.1:8080/activeUser?email=` + email + `&id=` + strconv.Itoa(id) + `">点击激活用户</a>`

	err = sendEmail.Send()
	if err != nil {
		fmt.Println(err)
		return
	}
	this.Redirect("/", 302)
}
