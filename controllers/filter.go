package controllers

import (
	"errors"
	"fmt"
	"kylin/lib"
	"kylin/models"
	"strings"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
)

func AccessRegister() {
	var Check = func(ctx *context.Context) {
		authGateway := beego.AppConfig.String("kylin_gateway_url")
		urlParams := strings.Split(strings.ToLower(ctx.Request.RequestURI), "/")
		if CheckUrl(urlParams) {
			//需要认证权限
			sessionInfo := ctx.Input.Session(SessionName)
			if sessionInfo == nil {
				//未登陆
				ctx.Redirect(302, authGateway)
			}
			//admin用户不需要认证权限
			adminUser := beego.AppConfig.String("kylin_admin_user")
			if sessionInfo.(models.User).Username == adminUser {
				return
			}
			accessList, _ := GetAccessList(sessionInfo.(models.User).Id)
			ret := AccessDecision(urlParams, accessList)
			if !ret {
				ctx.Output.Body([]byte("<h1 class=\"alert alert-danger\">权限不足</h1>"))
				//				ctx.Output.JSON(&map[string]interface{}{"status": false, "info": "权限不足"}, true, false)
			}
		}

	}
	beego.InsertFilter("/*", beego.BeforeRouter, Check)
}
func CheckUrl(urlParams []string) bool {
	if len(urlParams) < 3 {
		return false
	}
	for _, urlPrefix := range strings.Split(beego.AppConfig.String("not_auth_package"), ",") {
		if urlPrefix == urlParams[1] {
			return false
		}
	}
	return true
}

type AccessNode struct {
	Id        int64
	Level     int64
	Name      string
	Childrens []*AccessNode
}

func GetAccessList(uid int64) (map[string]bool, error) {
	list, err := models.AccessList(uid)
	if err != nil {
		return nil, err
	}
	alist := make([]*AccessNode, 0)
	//游戏项目的权限
	for _, l := range list {
		//		if l["Pid"].(int64) == 0 && l["Level"].(int64) == 1 {
		isEnableGroup := models.IsEnableGroup(l["Group"].(int64))
		if l["Pid"].(int64) == 0 && isEnableGroup {
			anode := new(AccessNode)
			anode.Id = l["Id"].(int64)
			anode.Level = l["Level"].(int64)
			anode.Name = l["Name"].(string)
			alist = append(alist, anode)
		}
	}
	//游戏项目的内容权限
	for _, l := range list {
		for _, an := range alist {
			if an.Level == l["Pid"].(int64) {
				anode := new(AccessNode)
				anode.Id = l["Id"].(int64)
				anode.Name = l["Name"].(string)
				an.Childrens = append(an.Childrens, anode)
			}
		}
	}
	//游戏项目内容权限
	//	for _, l := range list {
	//		if l["Level"].(int64) == 3 {
	//			for _, an := range alist {
	//				for _, an1 := range an.Childrens {
	//					if an1.Id == l["Pid"].(int64) {
	//						anode := new(AccessNode)
	//						anode.Id = l["Id"].(int64)
	//						anode.Name = l["Name"].(string)
	//						an1.Childrens = append(an1.Childrens, anode)
	//					}
	//				}

	//			}
	//		}
	//	}

	accesslist := make(map[string]bool)
	for _, v := range alist {
		for _, v1 := range v.Childrens {
			//			for _, v2 := range v1.Childrens {
			vname := strings.Split(v.Name, "/")
			v1name := strings.Split(v1.Name, "/")
			//				v2name := strings.Split(v2.Name, "/")
			//控制两层url权限
			str := fmt.Sprintf("/%s/%s", strings.ToLower(vname[0]), strings.ToLower(v1name[0]))
			accesslist[str] = true
		}
		//		}
	}
	return accesslist, nil
}
func AccessDecision(urlParams []string, accesslist map[string]bool) bool {
	if CheckUrl(urlParams) {
		s := fmt.Sprintf("/%s/%s", urlParams[1], urlParams[2])
		if len(accesslist) < 1 {
			return false
		}
		_, ok := accesslist[s]
		if ok != false {
			return true
		}
	} else {
		return true
	}
	return false
}

//check login
func CheckLogin(username string, password string) (user models.User, err error) {
	user = models.GetUserByUsername(username)
	if user.Id == 0 {
		return user, errors.New("用户不存在")
	}
	if user.Password != lib.Pwdhash(password) {
		return user, errors.New("密码错误")
	}
	if user.Status != 2 {
		//用户启用状态
		return user, errors.New("用户未启用")
	}
	return user, nil
}
