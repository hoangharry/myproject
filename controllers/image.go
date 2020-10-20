package controllers

import (
	"fmt"
	"myproject/models"
	"strconv"
	"time"

	"github.com/astaxie/beego"
)

type ImageController struct {
	beego.Controller
}

// @router /history
func (this *ImageController) GetPage() {

	page, _ := strconv.Atoi(this.Ctx.Input.Param(":page"))
	usr := this.Ctx.Input.Param(":usr")
	camID := this.Ctx.Input.Param(":camID")
	res := struct{ Images []*models.Image }{models.DefaultImageList.PageImg(usr, camID, page)}
	this.Data["json"] = res
	this.ServeJSON()
}

// @router /count
func (this *ImageController) GetCount() {
	usr := this.GetString("usr")
	camID := this.GetString("camID")
	res := struct{ Number int64 }{models.DefaultImageList.CountImg(usr, camID)}
	this.Data["json"] = res
	this.ServeJSON()
}

func (this *ImageController) GetCountMonitor() {
	usr := this.GetString("usr")
	camID := this.GetString("camID")
	fromTime, err := time.Parse(time.RFC3339, this.GetString("fTime"))
	if err != nil {
		fmt.Println(err)
	}
	res := struct{ Number int64 }{models.DefaultImageList.CountImgMonitor(usr, camID, fromTime)}
	this.Data["json"] = res
	this.ServeJSON()
}

func (this *ImageController) GetCountQuery() {
	usr := this.GetString("usr")
	camID := this.GetString("camID")
	fromTime, err := time.Parse(time.RFC3339, this.GetString("fTime"))
	if err != nil {
		fmt.Println(err)
	}
	toTime, err := time.Parse(time.RFC3339, this.GetString("tTime"))
	if err != nil {
		fmt.Println(err)
	}
	res := struct{ Number int64 }{models.DefaultImageList.CountImgFromQuery(usr, camID, fromTime, toTime)}
	this.Data["json"] = res
	this.ServeJSON()
}

func (this *ImageController) GetImageOnInit() {
	usr := this.GetString("usr")
	camID := this.GetString("camID")
	fmt.Println(this.GetString("fTime"))
	fromTime, err := time.Parse(time.RFC3339, this.GetString("fTime"))
	if err != nil {
		fmt.Println(err)
	}
	res := struct{ Images []*models.Image }{models.DefaultImageList.ImgOnInit(usr, camID, fromTime)}
	this.Data["json"] = res
	this.ServeJSON()
}

func (this *ImageController) GetImageInTime() {
	usr := this.GetString("usr")
	camID := this.GetString("camID")
	page, _ := strconv.Atoi(this.GetString("page"))
	fromTime, err := time.Parse(time.RFC3339, this.GetString("fTime"))
	if err != nil {
		fmt.Println(err)
	}
	toTime, err := time.Parse(time.RFC3339, this.GetString("tTime"))
	if err != nil {
		fmt.Println(err)
	}
	res := struct{ Images []*models.Image }{models.DefaultImageList.ImgInTime(usr, camID, fromTime, toTime, page)}
	this.Data["json"] = res
	this.ServeJSON()

}
