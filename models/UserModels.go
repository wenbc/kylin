package models

import (
	"errors"
	//	"log"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/validation"
	"kylin/lib"
	"time"
)

//用户表
type User struct {
	Id            int64
	Username      string    `orm:"unique;size(32)" form:"Username"  valid:"Required;MaxSize(20);MinSize(6)"`
	Password      string    `orm:"size(32)" form:"Password" valid:"Required;MaxSize(20);MinSize(6)"`
	Repassword    string    `orm:"-" form:"Repassword" valid:"Required"`
	Email         string    `orm:"size(32)" form:"Email" valid:"Email"`
	Status        int       `orm:"default(2)" form:"Status" valid:"Range(1,2)"` //1 表示停用 2 表示启用
	Lastlogintime time.Time `orm:"null;type(datetime)" form:"-"`
	Createtime    time.Time `orm:"type(datetime);auto_now_add" `
	Role          []*Role   `orm:"rel(m2m)"`
}

func (self *User) TableName() string {
	return beego.AppConfig.String("kylin_user_table")
}

func init() {
	orm.RegisterModel(new(User))
}

func GetUserByUsername(username string) (user User) {
	user = User{Username: username}
	o := orm.NewOrm()
	o.Read(&user, "Username")
	return user
}
func checkUser(u *User) (err error) {
	valid := validation.Validation{}
	b, err := valid.Valid(u)
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
func UpdateUser(u *User) (int64, error) {
	o := orm.NewOrm()
	user := make(orm.Params)
	if len(u.Username) > 0 {
		user["Username"] = u.Username
	}
	if len(u.Email) > 0 {
		user["Email"] = u.Email
	}
	if len(u.Password) > 0 {
		user["Password"] = lib.Strtomd5(u.Password)
	}
	if u.Status != 0 {
		user["Status"] = u.Status
	}
	if len(user) == 0 {
		return 0, errors.New("update field is empty")
	}
	var table User
	num, err := o.QueryTable(table).Filter("Id", u.Id).Update(user)
	return num, err
}
func UpdateLastLogTime(u User) (int64, error) {
	o := orm.NewOrm()
	user := make(orm.Params)
	user["Lastlogintime"] = time.Now()
	var table User
	num, err := o.QueryTable(table).Filter("Id", u.Id).Update(user)
	return num, err
}
func GetUserList() []orm.Params {
	var users []orm.Params
	o := orm.NewOrm()
	user := new(User)
	_, err := o.QueryTable(user).Values(&users)
	if err != nil {
		beego.Error("UserModels GetUserList is error!", err.Error())
		return nil
	}
	return users
}
func GetUserListById(uid int) []orm.Params {
	var users []orm.Params
	o := orm.NewOrm()
	user := new(User)
	_, err := o.QueryTable(user).Filter("Id", uid).Values(&users)
	if err != nil {
		beego.Error("UserModels GetUserListById is error!", err.Error())
		return nil
	}
	return users
}

//添加用户
func AddUser(u *User) (int64, error) {
	if err := checkUser(u); err != nil {
		return 0, err
	}
	o := orm.NewOrm()
	user := new(User)
	user.Username = u.Username
	user.Password = lib.Strtomd5(u.Password)
	user.Email = u.Email
	user.Status = u.Status
	id, err := o.Insert(user)
	return id, err
}

//删除用户
func DelUser(uid int64) (int64, error) {
	o := orm.NewOrm()
	status, err := o.Delete(&User{Id: uid})
	return status, err
}
