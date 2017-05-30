package controllers

import (
	"encoding/json"
	"io/ioutil"
	"licenseserver/controllers/internalerrors"
	"licenseserver/models"
	"net/http"

	"github.com/astaxie/beego/validation"
	log "github.com/cihub/seelog"
)

type RegisterLicenseRequest struct {
	PhoneNumber string `valid:"Required"`
	GUID        string `valid:"Required"`
	CompanyName string
	ExpireAt    string `valid:"Required;Match(/\d{4}-\d{1,2}-\d{1,2}/)"`
}

func Regist(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body)
	var request RegisterLicenseRequest
	if err := json.Unmarshal(body, &request); err != nil {
		handleError(w, internalerrors.NewLogicError(internalerrors.RequestError, err.Error()))
		return
	}
	check := validation.Validation{}
	if b, _ := check.Valid(&request); !b {
		handleError(w, internalerrors.NewLogicError(internalerrors.RequestError, "request error1"))
		return
	}
	if err := addLicense(&request); err != nil {
		handleError(w, err)
		return
	} else {
		Response(w, http.StatusOK, "Succeed")
	}
}

func GetLicense(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	phonenumber := r.FormValue("PhoneNumber")
	guid := r.FormValue("GUID")
	if phonenumber == "" && guid == "" {
		handleError(w, internalerrors.NewLogicError(internalerrors.RequestError, "request error"))
		return
	}
	var license *models.License
	var err error
	if phonenumber != "" {
		license, err = models.GetLicenseByPhoneNumber(phonenumber)
	} else {
		license, err = models.GetLicenseByGUID(guid)
	}
	if err != nil {
		handleError(w, internalerrors.NewLogicError(internalerrors.RequestError, err.Error()))
		return
	}
	Response(w, http.StatusOK, license)
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

func CheckPhoneNumberExist(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	phonenumber := r.FormValue("PhoneNumber")
	if phonenumber == "" {
		handleError(w, internalerrors.NewLogicError(internalerrors.RequestError, "request error"))
		return
	}
	b, err := models.CheckLicensePhone(phonenumber)
	if err != nil {
		handleError(w, internalerrors.NewLogicError(internalerrors.InternalError, err.Error()))
		return
	}
	Response(w, http.StatusOK, map[string]interface{}{
		"Exist": b,
	})
}

type UpdateUserInfoRequest struct {
	PhoneNumber string `valid:"Required"`
	GUID        string `valid:"Required"`
	CompanyName string
}

func UpdateUserInfo(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body)
	var request UpdateUserInfoRequest
	if err := json.Unmarshal(body, &request); err != nil {
		handleError(w, internalerrors.NewLogicError(internalerrors.RequestError, err.Error()))
		return
	}
	check := validation.Validation{}
	if b, _ := check.Valid(&request); !b {
		handleError(w, internalerrors.NewLogicError(internalerrors.RequestError, "request error"))
		return
	}

	license := models.License{
		PhoneNumber: request.PhoneNumber,
		GUID:        request.GUID,
		CompanyName: request.CompanyName,
	}
	if err := models.UpdateLicenseInfo(&license); err != nil {
		log.Warn("error:", err)
		handleError(w, internalerrors.NewLogicError(internalerrors.RequestError, err.Error()))
		return
	}
	Response(w, http.StatusOK, "Succeed")
}
