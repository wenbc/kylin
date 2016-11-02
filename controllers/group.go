package controllers

import (
	"github.com/astaxie/beego"
	"kylin/models"
	"strconv"
)

type GroupController struct {
	beego.Controller
}

func (self *GroupController) Index() {
	groupId := self.Ctx.Input.Param(":id")
	if groupId == "" {
		//获取所有的分组
		groups := models.GetGrouplist()
		self.Data["Groups"] = groups
		self.TplName = "admin/groups.html"
	} else {
		gid, err := strconv.Atoi(groupId)
		if err != nil {
			beego.Error("controllers.GroupController strconv.Atoi is error!", err.Error())
			self.Abort("500")

		}
		groups := models.GetGroupListById(gid)
		self.Data["Groups"] = groups
		self.TplName = "admin/groupinfos.html"
	}
}

//添加项目控制器
func (self *GroupController) AddGroup() {
	g := models.Group{}
	if err := self.ParseForm(&g); err != nil {
		self.Ctx.Output.Body([]byte(err.Error()))
		return
	}
	id, err := models.AddGroupModel(&g)
	if err == nil && id > 0 {
		self.Ctx.Output.Body([]byte("1"))
		return
	} else {

		self.Ctx.Output.Body([]byte(err.Error()))
		return
	}
}

//更新项目  控制器
func (self *GroupController) UpdateGroup() {
	g := models.Group{}
	if err := self.ParseForm(&g); err != nil {
		self.Abort("500")
		return
	}
	id, err := models.UpdateGroup(&g)
	if err == nil && id > 0 {
		self.Ctx.Output.Body([]byte("1"))
		return
	} else {
		self.Ctx.Output.Body([]byte(err.Error()))
		return
	}
}

//删除项目  控制器
func (self *GroupController) DelGroup() {
	Id, _ := self.GetInt64("Id")
	_, err := models.DelGroupById(Id)
	if err != nil {
		beego.Error("GroupController DelGroup is error!", err.Error())
	}
	groups := models.GetGrouplist()
	self.Data["Groups"] = groups
	self.TplName = "admin/groups.html"
}
