package controllers

import (
	"fmt"
	"myproject/models"

	"github.com/astaxie/beego"
)

type UserController struct {
	beego.Controller
}

func (this *UserController) CheckUser() {
	// usr := this.GetString()
}
func (this *UserController) GetUsrCam() {
	usr := this.GetString("usr")
	res := struct{ CamID []string }{models.DefaultUserManager.GetCam(usr)}
	this.Data["json"] = res
	fmt.Println(res)
	this.ServeJSON()
}
