package controllers

import (
	"github.com/astaxie/beego"
	"kylin/models"
	"strconv"
)

type UserController struct {
	beego.Controller
}

func (self *UserController) Index() {
	uid := self.Ctx.Input.Param(":id")
	if uid == "" {
		users := models.GetUserList()
		self.Data["Users"] = users
		self.TplName = "admin/users.html"
	} else {
		userid, err := strconv.Atoi(uid)
		if err != nil {
			beego.Error("controllers.GroupController strconv.Atoi is error!", err.Error())
			self.Abort("500")

		}
		users := models.GetUserListById(userid)
		self.Data["Users"] = users
		self.TplName = "admin/userinfos.html"
	}
}

func (self *UserController) AddUser() {
	u := models.User{}
	if err := self.ParseForm(&u); err != nil {
		self.Ctx.Output.Body([]byte(err.Error()))
		return
	}
	if u.Password != u.Repassword {
		self.Ctx.Output.Body([]byte("两次输入的密码不一致！"))
		return
	}
	uid, err := models.AddUser(&u)
	if err == nil && uid > 0 {
		self.Ctx.Output.Body([]byte("1"))
		return
	} else {
		self.Ctx.Output.Body([]byte(err.Error()))
		return
	}
}
func (self *UserController) DelUser() {
	uid, err := self.GetInt64("Id")
	if err != nil {
		beego.Error("UserController.DelUser is error!", err.Error())
		return
	}
	_, err = models.DelUser(uid)
	if err != nil {
		beego.Error("UserController.DelUser is error!", err.Error())
	}
	users := models.GetUserList()
	self.Data["Users"] = users
	self.TplName = "admin/users.html"
}
func (self *UserController) UpdateUser() {
	u := models.User{}
	if err := self.ParseForm(&u); err != nil {
		self.Ctx.Output.Body([]byte(err.Error()))
		return
	}
	num, err := models.UpdateUser(&u)
	if err == nil && num > 0 {
		self.Ctx.Output.Body([]byte("1"))
		return
	} else {
		self.Ctx.Output.Body([]byte(err.Error()))
		return
	}
}
