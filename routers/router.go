package routers

import (
	"kylin/controllers"
	"kylin/controllers/fyws"

	"github.com/astaxie/beego"
)

func init() {
	//首页
	beego.Router("/", &controllers.MainController{}, "*:Index")
	beego.Router("/public/index", &controllers.MainController{}, "*:Index")
	beego.Router("/public/login", &controllers.MainController{}, "*:Login")
	beego.Router("/public/logout", &controllers.MainController{}, "*:Logout")
	beego.Router("/public/changepwd", &controllers.MainController{}, "post:Changepwd")
	//用户
	beego.Router("/admin/user/index/?:id", &controllers.UserController{}, "*:Index")
	beego.Router("/admin/user/AddUser", &controllers.UserController{}, "*:AddUser")
	beego.Router("/admin/user/DelUser", &controllers.UserController{}, "*:DelUser")
	beego.Router("/admin/user/UpdateUser", &controllers.UserController{}, "*:UpdateUser")
	//组
	beego.Router("/admin/group/index/?:id", &controllers.GroupController{}, "*:Index")
	beego.Router("/admin/group/AddGroup", &controllers.GroupController{}, "post:AddGroup")
	beego.Router("/admin/group/UpdateGroup", &controllers.GroupController{}, "*:UpdateGroup")
	beego.Router("/admin/group/DelGroup", &controllers.GroupController{}, "*:DelGroup")
	//角色
	beego.Router("/admin/role/index/?:id", &controllers.RoleController{}, "*:Index")
	beego.Router("/admin/role/AddRole", &controllers.RoleController{}, "post:AddRole")
	beego.Router("/admin/role/DelRole", &controllers.RoleController{}, "post:DelRole")
	beego.Router("/admin/role/UpdateRole", &controllers.RoleController{}, "post:UpdateRole")
	beego.Router("/admin/role/AccessToNode", &controllers.RoleController{}, "*:AccessToNode")
	beego.Router("/admin/role/AddAccess", &controllers.RoleController{}, "*:AddAccess")
	beego.Router("/admin/role/RoleToUserList", &controllers.RoleController{}, "*:RoleToUserList")
	beego.Router("/admin/role/AddRoleToUser", &controllers.RoleController{}, "*:AddRoleToUser")

	//节点
	beego.Router("/admin/node/index/?:id", &controllers.NodeController{}, "*:Index")
	beego.Router("/admin/node/UpdateNode", &controllers.NodeController{}, "post:UpdateNode")
	beego.Router("/admin/node/DelNode", &controllers.NodeController{}, "post:DelNode")
	beego.Router("/admin/node/AddNode", &controllers.NodeController{}, "post:AddNode")
}

//风云无双
func init() {
	beego.Router("/fyws/host/index/?:id", &fyws.FywsController{}, "*:Index")
	beego.Router("/fyws/game/index/?:id", &fyws.FywsController{}, "*:GameIndex")
	beego.Router("/fyws/uniongame/index/?:id", &fyws.FywsController{}, "*:UnionIndex")
	beego.Router("/fyws/migregame/index/?:id", &fyws.FywsController{}, "*:MigreIndex")
}
