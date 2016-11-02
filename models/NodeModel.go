package models

import (
	"errors"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

//节点表
type Node struct {
	Id     int64
	Title  string  `orm:"size(100)" form:"Title"  valid:"Required"`
	Name   string  `orm:"size(100)" form:"Name"  valid:"Required"`
	Level  int64   `form:"Level"  valid:"Required"` //自定义level的等级，默认admin的等级是level
	Pid    int64   `form:"Pid"  valid:"Required"`
	Remark string  `orm:"null;size(200)" form:"Remark" valid:"MaxSize(200)"`
	Status int     `orm:"default(2)" form:"Status" valid:"Range(1,2)"` //1 表示停用 2 表示启用
	Group  *Group  `orm:"rel(fk)"`
	Role   []*Role `orm:"rel(m2m)"`
}

func (self *Node) TableName() string {
	return beego.AppConfig.String("kylin_node_table")
}

func init() {
	orm.RegisterModel(new(Node))
}

func GetNodeTree(pid int64) ([]orm.Params, error) {
	o := orm.NewOrm()
	node := new(Node)
	var nodes []orm.Params
	_, err := o.QueryTable(node).Filter("pid", pid).Values(&nodes, "Id", "Title", "Name", "Group__Id")
	//	_, err := o.QueryTable(node).Filter("pid", pid).Filter("Status", 2).Values(&nodes, "Id", "Title", "Name", "Group__Id")
	if err != nil {
		return nodes, err
	}
	return nodes, nil
}
func GetNodeTree1(groupId, pid int64) ([]orm.Params, error) {
	o := orm.NewOrm()
	node := new(Node)
	var nodes []orm.Params
	_, err := o.QueryTable(node).Exclude("pid", pid).Filter("Group__Id", groupId).Values(&nodes)
	if err != nil {
		return nodes, err
	}
	return nodes, nil
}
func GetNodeList() []orm.Params {
	var nodes []orm.Params
	o := orm.NewOrm()
	node := new(Node)
	_, err := o.QueryTable(node).Values(&nodes, "Id", "Title", "Name", "Status", "Level", "Pid", "Remark", "Group__Title", "Group__Id")
	//	_, err := o.QueryTable(node).Values(&nodes)
	if err != nil {
		beego.Error("NodeModels GetNodeList is error!", err.Error())
		return nil
	}
	return nodes
}
func GetNodeListById(nid int) []orm.Params {
	var nodes []orm.Params
	o := orm.NewOrm()
	node := new(Node)
	_, err := o.QueryTable(node).Filter("Id", nid).Values(&nodes, "Id", "Title", "Name", "Status", "Pid", "Remark", "Group__Title", "Group__Id")
	if err != nil {
		beego.Error("NodeModels GetUserListById is error!", err.Error())
		return nil
	}
	return nodes
}
func UpdateNode(n *Node) (int64, error) {
	o := orm.NewOrm()
	node := make(orm.Params)
	if len(n.Title) > 0 {
		node["Title"] = n.Title
	}
	if len(n.Name) > 0 {
		node["Name"] = n.Name
	}
	if len(n.Remark) > 0 {
		node["Remark"] = n.Remark
	}
	if n.Status != 0 {
		node["Status"] = n.Status
	}
	if len(node) == 0 {
		return 0, errors.New("NodeMode.UpdateNode update field is empty")
	}
	var table Node
	num, err := o.QueryTable(table).Filter("Id", n.Id).Update(node)
	return num, err
}

//删除节点
func DelNode(nid int64) (int64, error) {
	o := orm.NewOrm()
	status, err := o.Delete(&Node{Id: nid})
	return status, err
}
func GetNodelistByGroupid(Groupid int64) (nodes []orm.Params, count int64) {
	o := orm.NewOrm()
	node := new(Node)
	count, _ = o.QueryTable(node).Filter("Group", Groupid).RelatedSel().Values(&nodes)
	return nodes, count
}
func AddNode(n *Node) (int64, error) {
	o := orm.NewOrm()
	node := new(Node)
	node.Title = n.Title
	node.Name = n.Name
	node.Level = n.Level
	node.Pid = n.Pid
	node.Remark = n.Remark
	node.Status = n.Status
	node.Group = n.Group
	id, err := o.Insert(node)
	return id, err
}
