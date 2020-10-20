package routers

import (
	"myproject/controllers"

	"github.com/astaxie/beego"
)

func init() {
	beego.Router("/report", &controllers.ReportController{}, "get:GetReport")
	beego.Router("/history/?:usr/?:camID/?:page", &controllers.ImageController{}, "get:GetPage")
	beego.Router("/monitor/data", &controllers.ImageController{}, "get:GetImageOnInit")
	beego.Router("/search", &controllers.ImageController{}, "get:GetImageInTime")
	beego.Router("/count/history", &controllers.ImageController{}, "get:GetCount")
	beego.Router("/count/monitor", &controllers.ImageController{}, "get:GetCountMonitor")
	beego.Router("/count/search", &controllers.ImageController{}, "get:GetCountQuery")
	beego.Router("/monitor/ws/?:usr/?:camid", &controllers.WebsocketController{}, "get:WsPage")
	beego.Router("/cam", &controllers.UserController{}, "get:GetUsrCam")

}
