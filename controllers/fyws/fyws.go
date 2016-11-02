package fyws

import (
	"github.com/astaxie/beego"
)

type FywsController struct {
	beego.Controller
}

func (self *FywsController) Index() {
	self.TplName = "fyws/host.html"
}

func (self *FywsController) GameIndex() {
	self.TplName = "fyws/game.html"
}
func (self *FywsController) UnionIndex() {
	self.TplName = "fyws/union.html"
}
func (self *FywsController) MigreIndex() {
	self.TplName = "fyws/migre.html"
}
