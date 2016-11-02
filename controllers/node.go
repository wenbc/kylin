package controllers

import (
	"kylin/models"
	"strconv"

	"github.com/astaxie/beego"
)

type NodeController struct {
	beego.Controller
}

func (self *NodeController) Index() {
	nodeid := self.Ctx.Input.Param(":id")
	if nodeid == "" {
		nodes := models.GetNodeList()
		grouplist := models.GetGrouplist()
		self.Data["Nodes"] = nodes
		self.Data["Groups"] = grouplist
		self.TplName = "admin/nodes.html"
	} else {
		nid, err := strconv.Atoi(nodeid)
		if err != nil {
			beego.Error("controllers.NodeController strconv.Atoi is error!", err.Error())
			self.Abort("500")

		}
		nodes := models.GetNodeListById(nid)
		self.Data["Nodes"] = nodes
		self.TplName = "admin/nodeinfos.html"
	}
}
func (self *NodeController) UpdateNode() {
	n := models.Node{}
	if err := self.ParseForm(&n); err != nil {
		self.Ctx.Output.Body([]byte(err.Error()))
		return
	}
	num, err := models.UpdateNode(&n)
	if err == nil && num > 0 {
		self.Ctx.Output.Body([]byte("1"))
		return
	} else {
		self.Ctx.Output.Body([]byte(err.Error()))
		return
	}
}
func (self *NodeController) DelNode() {
	nid, err := self.GetInt64("Id")
	if err != nil {
		beego.Error("NodeController.DelNode is error!", err.Error())
		return
	}
	_, err = models.DelNode(nid)
	if err != nil {
		beego.Error("NodeController.DelNode is error!", err.Error())
	}
	nodes := models.GetNodeList()
	self.Data["Nodes"] = nodes
	self.TplName = "admin/nodes.html"
}
func (self *NodeController) AddNode() {
	n := models.Node{}
	if err := self.ParseForm(&n); err != nil {
		beego.Error("NodeController AddNode self.ParseForm is error!", err.Error())
		self.Ctx.Output.Body([]byte(err.Error()))
		return
	}
	groupid, err := self.GetInt64("groupid")
	if err != nil {
		beego.Error("NodeController AddNode GetInt64(groupid) is error!", err.Error())
		self.Ctx.Output.Body([]byte(err.Error()))
		return
	}
	group := new(models.Group)
	group.Id = groupid
	n.Group = group

	id, err := models.AddNode(&n)
	if err == nil && id > 0 {
		self.Ctx.Output.Body([]byte("1"))
		return
	} else {
		self.Ctx.Output.Body([]byte(err.Error()))
		return
	}

}
