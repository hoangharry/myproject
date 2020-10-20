package controllers

import (
	"fmt"
	"myproject/models"
	"time"

	"github.com/astaxie/beego"
)

type ReportController struct {
	beego.Controller
}

// @router /report
func (this *ReportController) GetReport() {

	fromDate, err := time.Parse(time.RFC3339, this.GetString("fDate"))
	if err != nil {
		fmt.Println(err)
	}
	toDate, err := time.Parse(time.RFC3339, this.GetString("tDate"))
	if err != nil {
		fmt.Println(err)
	}
	usr := this.GetString("usr")
	camID := this.GetString("cam_id")
	if camID == "" {
		res := struct{ Reports []*models.Report }{models.DefaultReportManager.FindReportOnUsr(fromDate, toDate, usr)}
		fmt.Println(res)
		this.Data["json"] = res
		this.ServeJSON()
	} else {
		res := struct{ Reports []*models.Report }{models.DefaultReportManager.FindReportOnCam(fromDate, toDate, usr, camID)}
		fmt.Println(res)
		this.Data["json"] = res
		this.ServeJSON()
	}

}
