package controllers

import (
	"encoding/json"
	"licenseserver/controllers/internalerrors"
	"licenseserver/models"
	"net/http"

	"github.com/astaxie/beego/validation"
	log "github.com/cihub/seelog"
)

type LicenseController struct {
	BaseController
}

type RegisterLicenseRequest struct {
	PhoneNumber string `valid:"Required"`
	GUID        string `valid:"Required"`
	CompanyName string
	ExpireAt    string `valid:"Required;Match(/\d{4}-\d{1,2}-\d{1,2}/)"`
}

func (this *LicenseController) Regist() {
	var request RegisterLicenseRequest
	log.Info((string)(this.Ctx.Input.RequestBody))
	if err := json.Unmarshal(this.Ctx.Input.RequestBody, &request); err != nil {
		this.handleError(internalerrors.NewLogicError(internalerrors.RequestError, err.Error()))
		return
	}
	check := validation.Validation{}
	if b, _ := check.Valid(&request); !b {
		this.handleError(internalerrors.NewLogicError(internalerrors.RequestError, "request error1"))
		return
	}
	if err := addLicense(&request); err != nil {
		this.handleError(err)
		return
	} else {
		this.Response(http.StatusOK, "Succeed")
	}
}

func (this *LicenseController) GetLicense() {
	phonenumber := this.GetString("PhoneNumber")
	if phonenumber == "" {
		this.handleError(internalerrors.NewLogicError(internalerrors.RequestError, "request error"))
		return
	}
	license, err := models.GetLicenseByPhoneNumber(phonenumber)
	if err != nil {
		this.handleError(internalerrors.NewLogicError(internalerrors.RequestError, err.Error()))
		return
	}
	this.Response(http.StatusOK, license)
}

func addLicense(request *RegisterLicenseRequest) *internalerrors.LogicError {
	if exist, err := models.CheckLicensePhone(request.PhoneNumber); err != nil {
		return internalerrors.NewLogicError(internalerrors.InternalError, err.Error())
	} else if exist {
		return internalerrors.NewLogicError(internalerrors.RequestError, "Phone Number Existed")
	}
	license := models.License{
		PhoneNumber: request.PhoneNumber,
		GUID:        request.GUID,
		ExpireDay:   request.ExpireAt,
		CompanyName: request.CompanyName,
		ExportTimes: 0,
	}
	if _, err := models.AddNewLicense(license); err != nil {
		return internalerrors.NewLogicError(internalerrors.InternalError, err.Error())
	} else {
		return nil
	}
}
