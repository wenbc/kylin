package controllers

import (
	"fmt"
	"github.com/astaxie/beego"
	"kylin/models"
	"strings"
)

func init() {
	//验证权限
	AccessRegister()
}

type MainController struct {
	beego.Controller
}

type Tree struct {
	Id         int64      `json:"id"`
	Text       string     `json:"text"`
	IconCls    string     `json:"iconCls"`
	Checked    string     `json:"checked"`
	State      string     `json:"state"`
	Children   []Tree     `json:"children"`
	Attributes Attributes `json:"attributes"`
}
type Attributes struct {
	Url   string `json:"url"`
	Price int64  `json:"price"`
}

func (self *MainController) Index() {
	userInfo := self.GetSession(SessionName)
	if userInfo == nil {
		self.Ctx.Redirect(302, beego.AppConfig.String("kylin_gateway_url"))
	}
	loginUser := userInfo.(models.User).Username
	tree := self.GetTree()
	if loginUser == beego.AppConfig.String("kylin_admin_user") {
		//超级管理员登录的树
		self.Data["Tree"] = &tree
	} else {
		//普通用户登录的树
		accesslist, err := GetAccessList(userInfo.(models.User).Id)
		if err != nil {
			beego.Error("MainController Index GetAccessList is error!", err.Error(), loginUser)
		}
		utree := self.GetUserTree(accesslist, tree)
		self.Data["Tree"] = utree
	}
	self.Data["ProjectName"] = ProjectName
	self.Data["UserName"] = userInfo.(models.User).Username
	self.Data[SessionName] = userInfo
	self.TplName = "index.html"
}

//登录
func (self *MainController) Login() {
	method := self.Ctx.Request.Method
	if method == "POST" {
		username := self.GetString("username")
		password := self.GetString("password")
		user, err := CheckLogin(username, password)
		beego.Informational("user login:", user.Username)
		if err == nil {
			_, err := models.UpdateLastLogTime(user)
			if err != nil {
				beego.Error("MainController.UpdateLastLogTime is error!", err.Error())
			}
			self.SetSession(SessionName, user)
			//			accesslist, _ := GetAccessList(user.Id)
			//			self.SetSession("accesslist", accesslist)
			beego.Informational("login Success!", user.Username)
		} else {
			self.Data["IsErr"] = true
			self.Data["ErrorInfos"] = err.Error()
			beego.Warning(user.Username, " login Faild!", err.Error())
		}

	}
	userinfo := self.GetSession(SessionName)
	if userinfo != nil {
		self.Ctx.Redirect(302, "/public/index")
	}
	self.Data["ProjectName"] = ProjectName
	self.Layout = "layout_login.html"
	self.TplName = "head_login.html"
}

//退出
func (self *MainController) Logout() {
	self.DelSession("userinfo")
	self.Ctx.Redirect(302, "/public/login")
}

func (self *MainController) GetTree() []Tree {
	nodes, err := models.GetNodeTree(0)
	if err != nil {
		beego.Error("MainController.GetTree is error!", err.Error())
	}
	tree := make([]Tree, len(nodes))
	for k, v := range nodes {
		tree[k].Id = v["Id"].(int64)
		tree[k].Text = v["Title"].(string)
		children, err := models.GetNodeTree1(v["Group__Id"].(int64), 0)
		if err != nil {
			beego.Error("MainController.GetTree is error!", v["Id"].(int64), err.Error())
		}
		tree[k].Children = make([]Tree, len(children))
		for k1, v1 := range children {
			tree[k].Children[k1].Id = v1["Id"].(int64)
			tree[k].Children[k1].Text = v1["Title"].(string)
			tree[k].Children[k1].Attributes.Url = "/" + v["Name"].(string) + "/" + v1["Name"].(string)
		}
	}
	return tree
}
func (self *MainController) GetUserTree(accesslist map[string]bool, allTree []Tree) []Tree {
	tree := make([]Tree, 0)
	for _, v := range allTree {
		children := make([]Tree, 0)
		for _, v1 := range v.Children {
			//如果有游戏项目的权限，就有控制所有该项目下的所有操作权限
			urlList := strings.Split(v1.Attributes.Url, "/")
			//控制两层url权限
			purl := fmt.Sprintf("/%s/%s", urlList[1], urlList[2])
			_, ok := accesslist[purl]
			if ok {

				attr := Attributes{
					Url: v1.Attributes.Url,
				}
				ch := Tree{
					Id:         v1.Id,
					Text:       v1.Text,
					Attributes: attr,
				}
				children = append(children, ch)
			}
		}
		if len(children) == 0 {
			continue
		}

		ptree := Tree{
			Id:       v.Id,
			Text:     v.Text,
			Children: children,
		}
		tree = append(tree, ptree)
	}
	return tree
}

//修改密码
func (self *MainController) Changepwd() {
	userinfo := self.GetSession(SessionName)
	if userinfo == nil {
		self.Ctx.Redirect(302, beego.AppConfig.String("kylin_gateway_url"))
	}

	oldpassword := self.GetString("oldpassword")
	newpassword := self.GetString("newpassword")
	repeatpassword := self.GetString("repeatpassword")
	if newpassword != repeatpassword {
		self.Ctx.Output.Body([]byte("新密码输入的密码不一致！"))
		return
	}
	user, err := CheckLogin(userinfo.(models.User).Username, oldpassword)
	if err == nil {
		var u models.User
		u.Id = user.Id
		u.Password = newpassword
		id, err := models.UpdateUser(&u)
		if err == nil && id > 0 {
			self.Ctx.Output.Body([]byte("密码修改成功"))
			return
		} else {
			self.Ctx.Output.Body([]byte(err.Error()))
			return
		}
	}
	self.Ctx.Output.Body([]byte(err.Error()))
}
