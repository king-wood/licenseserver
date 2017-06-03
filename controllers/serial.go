package controllers

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"licenseserver/controllers/internalerrors"
	"licenseserver/models"
	"net/http"
	"time"

	"io/ioutil"

	"github.com/astaxie/beego/validation"
	log "github.com/cihub/seelog"
)

func NewSerial(w http.ResponseWriter, r *http.Request) {
	h := md5.New()
	io.WriteString(h, fmt.Sprintf("%s", time.Now()))
	serial := fmt.Sprintf("%x", h.Sum(nil))

	record := models.Serial{
		Serial:      serial,
		Status:      models.SERIAL_UNREGISTERED,
		ExportTimes: 0,
		ExpireDay:   time.Now().Format("2006-01-02"),
	}
	_, err := models.AddNewSerial(record)
	if err != nil {
		log.Warn("text: add new serial error: ", err)
		handleError(w, internalerrors.NewLogicError(internalerrors.InternalError, err.Error()))
		return
	}
	Response(w, http.StatusOK, map[string]interface{}{
		"serial": serial,
	})
}

type RegisterRequest struct {
	PhoneNumber string `json:"PhoneNumber" valid:"Required"`
	Serial      string `json:"Serial" valid:"Required"`
	PCID        string `json:"PCID" valid:"Required"`
}

func RegisterSerial(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body)
	var request RegisterRequest
	if err := json.Unmarshal(body, &request); err != nil {
		handleError(w, internalerrors.NewLogicError(internalerrors.RequestError, err.Error()))
		return
	}
	check := validation.Validation{}
	if b, _ := check.Valid(&request); !b {
		handleError(w, internalerrors.NewLogicError(internalerrors.RequestError, "query peremeter error"))
		return
	}

	if canuse, err := models.CheckSerialUse(request.Serial); err != nil {
		handleError(w, internalerrors.NewLogicError(internalerrors.InternalError, "query db error"))
		return
	} else if !canuse {
		handleError(w, internalerrors.NewLogicError(internalerrors.RequestError, "serial has already registered or not yet created"))
		return
	}

	if registered, err := models.CheckIfRegistered(request.Serial, request.Serial); err != nil {
		handleError(w, internalerrors.NewLogicError(internalerrors.InternalError, "query db error"))
		return
	} else if registered {
		handleError(w, internalerrors.NewLogicError(internalerrors.RequestError, "phone number and serial have already registered"))
		return
	}

	if err := models.RegisterSerial(request.Serial, request.PhoneNumber, request.PCID); err != nil {
		handleError(w, internalerrors.NewLogicError(internalerrors.InternalError, "update db error"))
		return
	} else {
		Response(w, http.StatusOK, map[string]string{
			"status": "succeed",
		})
	}
}

func Sync(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body)
	var request RegisterRequest
	if err := json.Unmarshal(body, &request); err != nil {
		handleError(w, internalerrors.NewLogicError(internalerrors.RequestError, err.Error()))
		return
	}
	check := validation.Validation{}
	if b, _ := check.Valid(&request); !b {
		handleError(w, internalerrors.NewLogicError(internalerrors.RequestError, "query peremeter error"))
		return
	}

	if registered, err := models.CheckIfRegistered(request.Serial, request.PhoneNumber); err != nil {
		handleError(w, internalerrors.NewLogicError(internalerrors.InternalError, "query db error"))
		return
	} else if !registered {
		handleError(w, internalerrors.NewLogicError(internalerrors.RequestError, "phone number and serial not register"))
		return
	}

	if exist, err := models.IsSerialMatchPC(request.Serial, request.PhoneNumber, request.PCID); err != nil {
		handleError(w, internalerrors.NewLogicError(internalerrors.InternalError, "query db error"))
		return
	} else if !exist {
		handleError(w, internalerrors.NewLogicError(internalerrors.RequestError, "Current Serial and PhoneNumber have used by other PC"))
		return
	}

	record, err := models.GetSerialBySerialAndPhoneNumber(request.Serial, request.PhoneNumber)
	if err != nil {
		handleError(w, internalerrors.NewLogicError(internalerrors.InternalError, "query db error"))
		return
	}
	Response(w, http.StatusOK, map[string]interface{}{
		"PhoneNumber": record.PhoneNumber,
		"Serial":      record.Serial,
		"ExpireDay":   record.ExpireDay,
		"ExportTimes": record.ExportTimes,
	})
}

type ExtendRequest struct {
	PhoneNumber string `json:"PhoneNumber" valid:"Required"`
	Serial      string `json:"Serial" valid:"Required"`
	ExpireDate  string `json:"ExpireDate" valid:"Required;Match(/\d{4}-\d{1,2}-\d{1,2}/)"`
	ExportTimes int    `json:"ExportTimes" valid:"Required"`
}

func Extend(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body)
	var request ExtendRequest
	if err := json.Unmarshal(body, &request); err != nil {
		handleError(w, internalerrors.NewLogicError(internalerrors.RequestError, err.Error()))
		return
	}
	check := validation.Validation{}
	if b, _ := check.Valid(&request); !b {
		handleError(w, internalerrors.NewLogicError(internalerrors.RequestError, "query peremeter error"))
		return
	}

	if err := models.ExtendSerialExpireDate(request.Serial,
		request.PhoneNumber,
		request.ExpireDate,
		request.ExportTimes); err != nil {
		handleError(w, internalerrors.NewLogicError(internalerrors.InternalError, "update db error"))
		return
	}

	Response(w, http.StatusOK, map[string]string{
		"status": "succeed",
	})
}

func GetAvailableSerial(w http.ResponseWriter, r *http.Request) {
	records, e := models.GetAvailableSerial()
	if e != nil {
		handleError(w, internalerrors.NewLogicError(internalerrors.InternalError, "query db error"))
		return
	}

	var result []string
	for _, v := range records {
		result = append(result, v.Serial)
	}
	Response(w, http.StatusOK, map[string]interface{}{
		"list": result,
	})
}

func GetSerialByPhone(w http.ResponseWriter, r *http.Request) {
	phone := r.FormValue("phonenumber")
	if phone == "" {
		handleError(w, internalerrors.NewLogicError(internalerrors.RequestError, "query peremeter error"))
	}

	result, e := models.GetSerialByPhoneNumber(phone)
	if e != nil {
		handleError(w, internalerrors.NewLogicError(internalerrors.InternalError, "query db error"))
		return
	}
	Response(w, http.StatusOK, map[string]interface{}{
		"list": result,
	})
}
