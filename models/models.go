package models

import (
	"database/sql"
	"fmt"
	"kylin/lib"
	"log"
	"os"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

func init() {
	//初始化数据库
	initDB()
	DBConnect()
}
func initDB() {
	args := os.Args
	for _, arg := range args {
		if arg == "-syncdb" {
			SyncDb()
			os.Exit(0)
		}
	}
}
func DbConfig() (db_host, db_port, db_user, db_pass, db_name string) {
	db_host = beego.AppConfig.String("db_host")
	db_port = beego.AppConfig.String("db_port")
	db_user = beego.AppConfig.String("db_user")
	db_pass = beego.AppConfig.String("db_pass")
	db_name = beego.AppConfig.String("db_name")
	return
}
func DBConnect() {
	db_host, db_port, db_user, db_pass, db_name := DbConfig()
	db_url := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&loc=Local", db_user, db_pass, db_host, db_port, db_name)
	err := orm.RegisterDataBase("default", "mysql", db_url)
	if err != nil {
		log.Fatal("model.DBConnect RegisterDataBase is error!", err)
	}

}

var o orm.Ormer

//初试化所有的数据库信息
func SyncDb() {
	beego.Informational("Database init starting!")
	//创建数据库
	CreateDb()

	DBConnect()
	//数据库别名
	name := "default"
	//true -> drop table 后在建表   false -> 表存在
	force := false
	//打印执行sql的过程
	verbose := true

	err := orm.RunSyncdb(name, force, verbose)
	if err != nil {
		log.Fatal("models.SyncDb RunSyncdb is error!", err)
	}

	//初始化user表
	insertUser()
	//初始化role表
	insertRole()
	//初始化node表
	insertNodes()
	//初始化grup表
	insertGroup()

	beego.Informational("Database init is Complete.")
}

//创建项目的数据库
func CreateDb() {
	db_host, db_port, db_user, db_pass, db_name := DbConfig()
	db_url := fmt.Sprintf("%s:%s@tcp(%s:%s)/?charset=utf8", db_user, db_pass, db_host, db_port)
	db, err := sql.Open("mysql", db_url)
	defer db.Close()
	if err != nil {
		log.Fatal("model.CreateDb open database error!", err)
	}
	sqlString := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS `%s` DEFAULT CHARACTER SET UTF8", db_name)
	_, err = db.Exec(sqlString)
	if err != nil {
		log.Fatal("model.CreateDb run sql error!", err)

	}
	beego.Informational("Create database Success!", db_name)
}

//初试化 User 表
func insertUser() {
	beego.Informational("insert admin user ...")
	u := new(User)
	u.Username = beego.AppConfig.String("kylin_admin_user")
	u.Password = lib.Pwdhash("admin")
	u.Email = "wenbin@mokylin.com"
	u.Status = 2
	o = orm.NewOrm()
	o.Insert(u)
	beego.Informational("insert admin user end")
}

//初始化role表
func insertRole() {
	fmt.Println("insert role ...")
	r := new(Role)
	r.Name = "超级管理员组"
	r.Remark = "拥有所有项目的管理权限"
	r.Status = 2
	o.Insert(r)
	fmt.Println("insert role end")
}

//初始化group表
func insertGroup() {
	fmt.Println("insert group ...")
	g := new(Group)
	g.Name = "Project"
	g.Title = "admin"
	g.Status = 2
	o.Insert(g)
	g.Name = "fenyun"
	g.Title = "风云无双"
	g.Status = 2
	o.Insert(g)
	fmt.Println("insert group end")
}

//初始化node表
func insertNodes() {
	fmt.Println("insert node ...")
	g := new(Group)
	g.Id = 1
	nodes := []Node{
		{Name: "admin", Title: "后台管理", Remark: "", Level: 1, Pid: 0, Status: 2, Group: g},
		{Name: "node/index", Title: "节点管理", Remark: "", Level: 2, Pid: 1, Status: 2, Group: g},
		//{Name: "index", Title: "node list", Remark: "", Level: 3, Pid: 2, Status: 2, Group: g},
		//{Name: "AddAndEdit", Title: "add or edit", Remark: "", Level: 3, Pid: 2, Status: 2, Group: g},
		//{Name: "DelNode", Title: "del node", Remark: "", Level: 3, Pid: 2, Status: 2, Group: g},
		{Name: "user/index", Title: "用户管理", Remark: "", Level: 2, Pid: 1, Status: 2, Group: g},
		//{Name: "Index", Title: "user list", Remark: "", Level: 3, Pid: 6, Status: 2, Group: g},
		//{Name: "AddUser", Title: "add user", Remark: "", Level: 3, Pid: 6, Status: 2, Group: g},
		//{Name: "UpdateUser", Title: "update user", Remark: "", Level: 3, Pid: 6, Status: 2, Group: g},
		//{Name: "DelUser", Title: "del user", Remark: "", Level: 3, Pid: 6, Status: 2, Group: g},
		{Name: "role/index", Title: "角色管理", Remark: "", Level: 2, Pid: 1, Status: 2, Group: g},
		//{Name: "index", Title: "role list", Remark: "", Level: 3, Pid: 16, Status: 2, Group: g},
		//		{Name: "AddAndEdit", Title: "add or edit", Remark: "", Level: 3, Pid: 16, Status: 2, Group: g},
		//		{Name: "DelRole", Title: "del role", Remark: "", Level: 3, Pid: 16, Status: 2, Group: g},
		//		{Name: "Getlist", Title: "get roles", Remark: "", Level: 3, Pid: 16, Status: 2, Group: g},
		//		{Name: "AccessToNode", Title: "show access", Remark: "", Level: 3, Pid: 16, Status: 2, Group: g},
		//		{Name: "AddAccess", Title: "add accsee", Remark: "", Level: 3, Pid: 16, Status: 2, Group: g},
		//		{Name: "RoleToUserList", Title: "show role to userlist", Remark: "", Level: 3, Pid: 16, Status: 2, Group: g},
		//		{Name: "AddRoleToUser", Title: "add role to user", Remark: "", Level: 3, Pid: 16, Status: 2, Group: g},
		{Name: "group/index", Title: "项目管理", Remark: "", Level: 2, Pid: 1, Status: 2, Group: g},
	}
	for _, v := range nodes {
		n := new(Node)
		n.Name = v.Name
		n.Title = v.Title
		n.Remark = v.Remark
		n.Level = v.Level
		n.Pid = v.Pid
		n.Status = v.Status
		n.Group = v.Group
		o.Insert(n)
	}
	g.Id = 2
	nodes = []Node{
		{Name: "fyws", Title: "风云无双", Remark: "", Level: 10, Pid: 0, Status: 2, Group: g},
		{Name: "host/index", Title: "服务器信息", Remark: "", Level: 2, Pid: 10, Status: 2, Group: g},
		{Name: "game/index", Title: "游戏服管理", Remark: "", Level: 2, Pid: 10, Status: 2, Group: g},
		{Name: "uniongame/index", Title: "游戏服合服", Remark: "", Level: 2, Pid: 10, Status: 2, Group: g},
		{Name: "migregame/index", Title: "游戏服迁服", Remark: "", Level: 2, Pid: 10, Status: 2, Group: g},
	}
	for _, v := range nodes {
		n := new(Node)
		n.Name = v.Name
		n.Title = v.Title
		n.Remark = v.Remark
		n.Level = v.Level
		n.Pid = v.Pid
		n.Status = v.Status
		n.Group = v.Group
		o.Insert(n)
	}
	fmt.Println("insert node end")
}
