package routers

import (
	"myproject/controllers"
	"github.com/astaxie/beego"

)

func init() {

	beego.Router("/home", &controllers.MainController{})
	beego.Router("/cam", &controllers.ImageController{})
	// beego.Router("/camera", &controllers.ImageController{})
}
