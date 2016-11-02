package models

import (
	"errors"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/validation"
)

//角色表
type Role struct {
	Id     int64
	Name   string  `orm:"size(100)" form:"Name"  valid:"Required"`
	Remark string  `orm:"null;size(200)" form:"Remark" valid:"MaxSize(200)"`
	Status int     `orm:"default(2)" form:"Status" valid:"Range(1,2)"` //1 表示停用 2 表示启用
	Node   []*Node `orm:"reverse(many)"`
	User   []*User `orm:"reverse(many)"`
}

func (self *Role) TableName() string {
	return beego.AppConfig.String("kylin_role_table")
}

func init() {
	orm.RegisterModel(new(Role))
}

//获取用户能访问项目node
func AccessList(uid int64) (list []orm.Params, err error) {
	var roles []orm.Params
	o := orm.NewOrm()
	role := new(Role)
	_, err = o.QueryTable(role).Filter("User__User__Id", uid).Filter("Status", 2).Values(&roles)
	if err != nil {
		return nil, err
	}

	var nodes []orm.Params
	node := new(Node)
	for _, r := range roles {
		_, err := o.QueryTable(node).Filter("Role__Role__Id", r["Id"]).Filter("Status", 2).Values(&nodes)
		if err != nil {
			return nil, err
		}
		for _, n := range nodes {
			list = append(list, n)
		}
	}
	return list, nil
}
func GetRoleList() []orm.Params {
	var roles []orm.Params
	o := orm.NewOrm()
	role := new(Role)
	_, err := o.QueryTable(role).Values(&roles)
	if err != nil {
		beego.Error("RoleModels GetRoleList is error!", err.Error())
		return nil
	}
	return roles
}
func GetRoleListById(rid int) []orm.Params {
	var roles []orm.Params
	o := orm.NewOrm()
	role := new(Role)
	_, err := o.QueryTable(role).Filter("Id", rid).Values(&roles)
	if err != nil {
		beego.Error("RoleModels GetRoleListById is error!", err.Error())
		return nil
	}
	return roles
}
func GetRowRole(rid int64) (Role, error) {
	o := orm.NewOrm()
	role := Role{Id: rid}
	err := o.Read(&role)
	if err != nil {
		return Role{}, err
	}
	return role, err
}
func checkRole(r *Role) (err error) {
	valid := validation.Validation{}
	b, err := valid.Valid(r)
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
func AddRole(r *Role) (int64, error) {
	if err := checkRole(r); err != nil {
		return 0, err
	}
	o := orm.NewOrm()
	role := new(Role)
	role.Name = r.Name
	role.Remark = r.Remark
	role.Status = r.Status
	id, err := o.Insert(role)
	return id, err
}
func DelRole(rid int64) (int64, error) {
	o := orm.NewOrm()
	status, err := o.Delete(&Role{Id: rid})
	return status, err
}
func UpdateRole(r *Role) (int64, error) {
	o := orm.NewOrm()
	role := make(orm.Params)
	if len(r.Name) > 0 {
		role["Name"] = r.Name
	}
	if len(r.Remark) > 0 {
		role["Remark"] = r.Remark
	}
	if r.Status != 0 {
		role["Status"] = r.Status
	}
	if len(role) == 0 {
		return 0, errors.New("RoleModel.UpdateRole update field is empty")
	}
	var table Role
	num, err := o.QueryTable(table).Filter("Id", r.Id).Update(role)
	return num, err
}
func GetNodelistByRoleId(Id int64) (nodes []orm.Params, count int64) {
	o := orm.NewOrm()
	node := new(Node)
	count, _ = o.QueryTable(node).Filter("Role__Role__Id", Id).Values(&nodes)
	return nodes, count
}
func DelGroupNode(roleid int64, groupid int64) error {
	var nodes []*Node
	var node Node
	role := Role{Id: roleid}
	o := orm.NewOrm()
	num, err := o.QueryTable(node).Filter("Group", groupid).RelatedSel().All(&nodes)
	if err != nil {
		return err
	}
	if num < 1 {
		return nil
	}
	for _, n := range nodes {
		m2m := o.QueryM2M(n, "Role")
		_, err1 := m2m.Remove(&role)
		if err1 != nil {
			return err1
		}
	}
	return nil
}
func AddRoleNode(roleid int64, nodeid int64) (int64, error) {
	o := orm.NewOrm()
	role := Role{Id: roleid}
	node := Node{Id: nodeid}
	m2m := o.QueryM2M(&node, "Role")
	num, err := m2m.Add(&role)
	return num, err
}
func GetUserByRoleId(roleid int64) (users []orm.Params, count int64) {
	o := orm.NewOrm()
	user := new(User)
	count, _ = o.QueryTable(user).Filter("Role__Role__Id", roleid).Values(&users)
	return users, count
}
func DelUserRole(roleid int64) error {
	o := orm.NewOrm()
	_, err := o.QueryTable("user_roles").Filter("role_id", roleid).Delete()
	return err
}
func AddRoleUser(roleid int64, userid int64) (int64, error) {
	o := orm.NewOrm()
	role := Role{Id: roleid}
	user := User{Id: userid}
	m2m := o.QueryM2M(&user, "Role")
	num, err := m2m.Add(&role)
	return num, err
}
