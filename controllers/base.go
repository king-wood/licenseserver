package controllers

import (
	"licenseserver/controllers/internalerrors"
	"net/http"

	"github.com/astaxie/beego"
)

type BaseController struct {
	beego.Controller
}

func (this *BaseController) handleError(err *internalerrors.LogicError) {
	if err.Type == internalerrors.RequestError {
		this.Data["json"] = map[string]string{
			"error": err.Description,
		}
		this.Ctx.Output.SetStatus(http.StatusBadRequest)
	} else {
		this.Data["json"] = map[string]string{
			"error": "server internal error",
		}
		this.Ctx.Output.SetStatus(http.StatusInternalServerError)
	}
	this.ServeJson()
}

func (this *BaseController) Response(code int, content interface{}) {
	this.Data["json"] = content
	this.Ctx.Output.SetStatus(code)
	this.ServeJson()
}

func (this *BaseController) ServeJson() {
	this.ServeJSON()
}
