package controllers

import(
	"github.com/astaxie/beego"
	"myproject/models"
	// "fmt"

)

type ImageController struct {
	beego.Controller
}

func (this *ImageController) Get(){
	res := struct{Images []*models.Image}{models.DefaultImageList.All()}
	this.Data["json"] = res
	this.TplName = "index.html"
	this.Render()

}
// func (this *ImageController) Get(){

// 	res := struct{Images []*models.Image}{models.DefaultImageList.All()}
// 	this.Data["json"] = res
// 	this.ServeJSON()
// }
