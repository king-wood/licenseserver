package routers

import (
	"licenseserver/controllers"

	"github.com/astaxie/beego"
)

func init() {
	beego.Router("/v1/regist", &controllers.LicenseController{}, "post:Regist")
	beego.Router("/v1/getlicense", &controllers.LicenseController{}, "get:GetLicense")
}
