package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"kylin/models"
	"strconv"
	"strings"
)

type RoleController struct {
	beego.Controller
}

func (self *RoleController) Index() {
	roleid := self.Ctx.Input.Param(":id")
	if roleid == "" {
		roles := models.GetRoleList()
		self.Data["Roles"] = roles
		self.TplName = "admin/roles.html"
	} else {
		rid, err := strconv.Atoi(roleid)
		if err != nil {
			beego.Error("controllers.RoleController strconv.Atoi is error!", err.Error())
			self.Abort("500")
		}
		roles := models.GetRoleListById(rid)
		self.Data["Roles"] = roles
		self.TplName = "admin/roleinfos.html"
	}
}
func (self *RoleController) AddRole() {
	r := models.Role{}
	if err := self.ParseForm(&r); err != nil {
		self.Ctx.Output.Body([]byte(err.Error()))
		return
	}
	rid, err := models.AddRole(&r)
	if err == nil && rid > 0 {
		self.Ctx.Output.Body([]byte("1"))
		return
	} else {
		self.Ctx.Output.Body([]byte(err.Error()))
		return
	}
}
func (self *RoleController) DelRole() {
	rid, err := self.GetInt64("Id")
	if err != nil {
		beego.Error("RoleController.DelRole is error!", err.Error())
		return
	}
	_, err = models.DelRole(rid)
	if err != nil {
		beego.Error("RoleController models.DelRole is error!", err.Error())
	}
	roles := models.GetRoleList()
	self.Data["Roles"] = roles
	self.TplName = "admin/roles.html"
}
func (self *RoleController) UpdateRole() {
	r := models.Role{}
	if err := self.ParseForm(&r); err != nil {
		self.Ctx.Output.Body([]byte(err.Error()))
		return
	}
	num, err := models.UpdateRole(&r)
	if err == nil && num > 0 {
		self.Ctx.Output.Body([]byte("1"))
		return
	} else {
		self.Ctx.Output.Body([]byte(err.Error()))
		return
	}
}

func (self *RoleController) AccessToNode() {
	if self.Ctx.Input.Method() == "POST" {
		groupid, _ := self.GetInt64("group_id")
		roleId, _ := self.GetInt64("Id")
		nodes, _ := models.GetNodelistByGroupid(groupid)
		list, _ := models.GetNodelistByRoleId(roleId)
		for i := 0; i < len(nodes); i++ {
			for l := 0; l < len(list); l++ {
				if nodes[i]["Id"] == list[l]["Id"] {
					nodes[i]["Checked"] = 1
				}
			}
		}
		if len(nodes) < 1 {
			nodes = []orm.Params{}
		}
		self.Data["Nodes"] = nodes
		self.TplName = "admin/accessnode.html"
	} else {
		roleId, err := self.GetInt64("Id")
		if err != nil {
			beego.Error("RoleController AccessToNode.GetInt64 is error!", err.Error())
			self.Abort("500")
		}
		role, err := models.GetRowRole(roleId)
		if err != nil {
			beego.Error("RoleController AccessToNode.GetRowRole is error!", err.Error())
			self.Abort("500")
		}
		self.Data["RoleId"] = roleId
		grouplist := models.GetGrouplist()
		self.Data["RoleName"] = role.Name
		self.Data["GroupList"] = grouplist
		self.TplName = "admin/accesstonode.html"
	}
}
func (self *RoleController) AddAccess() {
	groupId, _ := self.GetInt64("group_id")
	roleId, _ := self.GetInt64("rid")
	ids := self.GetString("nids")
	err := models.DelGroupNode(roleId, groupId)
	if err != nil {
		beego.Error("RoleController Access DelGroupNode is error!", err.Error())
		self.Ctx.Output.Body([]byte(err.Error()))
		return
	}
	nodeids := strings.Split(ids, ",")
	for _, v := range nodeids {
		if v == "" {
			continue
		}
		id, err := strconv.Atoi(v)
		if err != nil {
			beego.Error("RoleController Access strconv.Atoi is error!", err.Error())
		}
		_, err = models.AddRoleNode(roleId, int64(id))
		if err != nil {
			beego.Error("RoleController Access  AddRoleNode is error!", err.Error())
			self.Ctx.Output.Body([]byte(err.Error()))
		}
	}
	self.Ctx.Output.Body([]byte("1"))
}
func (self *RoleController) RoleToUserList() {
	roleId, err := self.GetInt64("Id")
	if err != nil {
		beego.Error("RoleController RoleToUserList.GetInt64 is error!", err.Error())
		self.Abort("500")
	}
	role, err := models.GetRowRole(roleId)
	if err != nil {
		beego.Error("RoleController RoleToUserList.GetRowRole is error!", err.Error())
		self.Abort("500")
	}
	users := models.GetUserList()
	list, _ := models.GetUserByRoleId(roleId)
	for u := 0; u < len(users); u++ {
		for l := 0; l < len(list); l++ {
			if users[u]["Id"] == list[l]["Id"] {
				users[u]["Checked"] = 1
			}
		}
	}
	if len(users) < 1 {
		users = []orm.Params{}
	}
	self.Data["RoleName"] = role.Name
	self.Data["Users"] = users
	self.Data["RoleId"] = roleId
	self.TplName = "admin/roletouserlist.html"
}
func (self *RoleController) AddRoleToUser() {
	roleid, _ := self.GetInt64("Id")
	ids := self.GetString("ids")
	userids := strings.Split(ids, ",")
	err := models.DelUserRole(roleid)
	if err != nil {
		beego.Error("RoleController AddRoleToUser  DelUserRole is error!", err.Error())
		self.Abort("500")
		return
	}
	if len(ids) > 0 {
		for _, v := range userids {
			id, err := strconv.Atoi(v)
			if err != nil {
				beego.Error("RoleController AddRoleToUser strconv.Atoi is error!", v, err.Error())
				self.Abort("500")
				return
			}
			_, err = models.AddRoleUser(roleid, int64(id))
			if err != nil {
				beego.Error("RoleController AddRoleToUser AddRoleUser is error!", v, err.Error())
				self.Abort("500")
				return
			}
		}
	}
	self.Ctx.Output.Body([]byte("1"))
}
