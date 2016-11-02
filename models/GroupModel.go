package models

import (
	"errors"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/validation"
)

//分组表
type Group struct {
	Id     int64
	Name   string  `orm:"unique;size(100)" form:"Name"  valid:"Required"`
	Title  string  `orm:"size(100)" form:"Title"  valid:"Required"`
	Status int     `orm:"default(2)" form:"Status" valid:"Range(1,2)"`
	Nodes  []*Node `orm:"reverse(many)"`
}

func (g *Group) TableName() string {
	return beego.AppConfig.String("kylin_group_table")
}

func init() {
	orm.RegisterModel(new(Group))
}
func GetGroupListById(gid int) []orm.Params {
	var groups []orm.Params
	o := orm.NewOrm()
	group := new(Group)
	_, err := o.QueryTable(group).Filter("Id", gid).Values(&groups)
	if err != nil {
		beego.Error("models.GetGroupListById is error!", err.Error())
		return nil
	}
	return groups
}
func GetGrouplist() []orm.Params {
	var groups []orm.Params
	o := orm.NewOrm()
	group := new(Group)
	_, err := o.QueryTable(group).Values(&groups)
	if err != nil {
		beego.Error("models.GetGrouplist is error!", err.Error())
		return nil
	}
	return groups
}
func checkGroup(g *Group) (err error) {
	valid := validation.Validation{}
	b, err := valid.Valid(g)
	if err != nil {
		return err
	}
	if !b {
		for _, err := range valid.Errors {
			beego.Notice(err.Key, err.Message)
			return errors.New(err.Key + ":" + err.Message)
		}
	}
	return nil
}
func AddGroupModel(g *Group) (int64, error) {
	if err := checkGroup(g); err != nil {
		return 0, err
	}
	o := orm.NewOrm()
	group := new(Group)
	group.Name = g.Name
	group.Title = g.Title
	group.Status = g.Status
	id, err := o.Insert(group)
	return id, err
}
func UpdateGroup(g *Group) (int64, error) {
	if err := checkGroup(g); err != nil {
		return 0, err
	}
	o := orm.NewOrm()
	group := make(orm.Params)
	if len(g.Name) > 0 {
		group["Name"] = g.Name
	}
	if len(g.Title) > 0 {
		group["Title"] = g.Title
	}
	if g.Status != 0 {
		group["Status"] = g.Status
	}
	if len(group) == 0 {
		return 0, errors.New("update field is empty")
	}
	var table Group
	num, err := o.QueryTable(table).Filter("Id", g.Id).Update(group)
	return num, err
}
func DelGroupById(Id int64) (int64, error) {
	o := orm.NewOrm()
	status, err := o.Delete(&Group{Id: Id})
	return status, err
}

//这个组是启用状态
func IsEnableGroup(Id int64) bool {
	o := orm.NewOrm()
	group := Group{Id: Id}
	err := o.Read(&group)
	if err != nil {
		beego.Error("GroupModel IsEnableGroup is error!", err.Error())
		return false
	}
	if group.Status == 2 {
		return true
	} else {
		return false
	}
}
