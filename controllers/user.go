package controllers

import (
	"encoding/base64"
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
func (this *UserController) ActiveEmail() {
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
func (this *UserController) ActiveUser() {
	//获取数据
	id, _ := this.GetInt("id")
	email := this.GetString("email")
	//校验数据

	//处理数据
	//获取orm对象
	o := orm.NewOrm()
	//获取user对象
	var user models.User
	//赋值
	user.Id = id
	o.Read(&user)
	user.Email = email
	user.Active = true
	//更新数据
	o.Update(&user)
	this.Data["errmsg"] = "激活成功"
	this.Redirect("/login?id="+strconv.Itoa(id), 302)
}
func (this *UserController) ShowLogin() {
	//id, _ := this.GetInt("id")
	userName := this.Ctx.GetCookie("userName")

	if userName != "" {
		dec, _ := base64.StdEncoding.DecodeString(userName)
		this.Data["userName"] = string(dec)
		fmt.Println("获取到的cookie用户名", string(dec))
		this.Data["checked"] = "checked"
	} else {
		this.Data["userName"] = ""
		this.Data["checked"] = ""
	}
	this.TplName = "login.html"
}
func (this *UserController) HandleLogin() {
	id, _ := this.GetInt("id")
	userName := this.GetString("userName")
	password := this.GetString("password")
	if password == "" || userName == "" {
		this.Data["errmsg"] = "用户名和密码不能为空"
		this.Redirect("/st/login?id="+strconv.Itoa(id), 302)
		return
	}

	o := orm.NewOrm()
	var user models.User
	user.Name = userName
	err := o.Read(&user, "Name")
	if err != nil {
		this.Data["errmsg"] = "用户不存在"
		this.Redirect("/st/login?id="+strconv.Itoa(id), 302)
		return
	}
	if password != user.PassWord {
		this.Data["errmsg"] = "用户或密码不对"
		this.Redirect("/st/login?id="+strconv.Itoa(id), 302)
		return
	}
	remember := this.GetString("remember")
	if remember == "on" {
		enc := base64.StdEncoding.EncodeToString([]byte(userName))
		this.Ctx.SetCookie("userName", enc, 60*60)
		fmt.Println("执行了设置cookie")
	} else {
		this.Ctx.SetCookie("userName", userName, -1)
	}
	this.SetSession("userName", userName)

	conn, err := redis.Dial("tcp", "172.16.10.11:6379")
	if err != nil {
		fmt.Println("连接redis，错误：", err)
		return
	}
	_, err = conn.Do("set", userName, userName)
	if err != nil {
		fmt.Println("redis 写入session错误:", err)
	}
	defer conn.Close()
	id = user.Id

	this.Redirect("/st/?id="+strconv.Itoa(id), 302)

}
func (this *UserController) ShowLogout() {
	this.DelSession("userName")
	this.Redirect("/login", 302)
}
func (this *UserController) ShowUserCenterInfo() {
	userName := this.GetSession("userName")
	o := orm.NewOrm()
	var address models.Address
	qs := o.QueryTable("Address").RelatedSel("User").Filter("User__Name", userName.(string))
	qs.Filter("Isdefault", true).One(&address)

	//给页面赋值
	this.Data["userName"] = userName.(string)
	this.Data["phone"] = address.Phone
	this.Data["addr"] = address.Addr
	this.TplName = "user_center_info.html"
}
func (this *UserController) ShowUserAddress() {
	userName := this.GetSession("userName")
	o := orm.NewOrm()
	var uAddress models.Address

	qs := o.QueryTable("Address").RelatedSel("User").Filter("User__Name", userName.(string))
	qs.Filter("Isdefault", true).One(&uAddress)
	this.Data["userName"] = uAddress.Receiver
	this.Data["address"] = uAddress.Addr
	phnum := uAddress.Phone
	if len(phnum) > 10 {
		per := phnum[0:3]
		has := phnum[7:]
		this.Data["phone"] = per + "****" + has
	}

	fmt.Println(uAddress)
	this.TplName = "user_center_site.html"
}
func (this *UserController) HandleAddress() {
	//获取用户填写数据
	userName := this.GetSession("userName")
	Receiver := this.GetString("Receiver")
	address := this.GetString("address")
	zipCode := this.GetString("zipCode")
	phone := this.GetString("phone")
	if Receiver == "" || address == "" || zipCode == "" || phone == "" {
		this.Data["errmsg"] = "填写信息不完整"
		this.Redirect("/st/user_center_site", 302)
		return
	}
	o := orm.NewOrm()

	var user models.User
	user.Name = userName.(string)
	o.Read(&user, "Name")

	var uAddress models.Address
	qs := o.QueryTable("Address").RelatedSel("User").Filter("User__Name", userName.(string))
	qs.Filter("Isdefault", true).One(&uAddress)

	//先将之前所有地址更新为非默认地址
	uAddress.Isdefault = false
	o.Update(&uAddress)

	//获取新的收货人信息并更新为默认地址
	var newAddress models.Address
	newAddress.Receiver = Receiver
	newAddress.Addr = address
	newAddress.Phone = phone
	newAddress.Zipcode = zipCode
	newAddress.Isdefault = true
	newAddress.User = &user
	n, err := o.Insert(&newAddress)
	if err != nil {
		fmt.Println("插入地址数据失败", err)
		this.Redirect("/st/user_center_site", 302)
		return
	}
	fmt.Println("C插入数据为", n)
	this.Redirect("/st/user_center_site", 302)
}
