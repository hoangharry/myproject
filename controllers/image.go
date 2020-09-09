package controllers

import(
	"github.com/astaxie/beego"
	// "go.mongodb.org/mongo-driver/bson"
	// "go.mongodb.org/mongo-driver/mongo"
	// "go.mongodb.org/mongo-driver/mongo/options"
	"myproject/models"
)

type ImageController struct {
	beego.Controller
}

func (this *ImageController) ListImg(){
	res := struct{Images []*models.Image }{models.DefaultImageList.All()}
	this.Data["json"] = res
	this.ServeJSON()
}
func (this *ImageController) Get() {
	this.TplName = "cam.html"
	this.Render()
}