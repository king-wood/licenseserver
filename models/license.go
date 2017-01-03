package models

import (
	"strconv"
	"time"
)

const (
	LICENSE_ID          = "id"
	LICENSE_PHONENUMBER = "phone_number"
	LICENSE_GUID        = "guid"
	LICENSE_COMPANYNAME = "company_name"
	LICENSE_EXPIREDAY   = "expire_day"
	LICENSE_EXPORTTIMES = "export_times"
)

type License struct {
	ID          int64  `db:"id" json:"id"`
	PhoneNumber string `db:"phone_number" json:"phonenumber"`
	GUID        string `db:"guid" json:"guid"`
	CompanyName string `db:"company_name" json:"companyname"`
	ExpireDay   string `db:"expire_day" json:"expireday"`
	ExportTimes int64  `db:"export_times" json:"exporttimes"`
}

func CheckLicensePhone(phone string) (bool, error) {
	var exists bool
	sql := "SELECT EXISTS (SELECT 1 FROM tbl_license WHERE " + LICENSE_PHONENUMBER + " = ?)"
	err := Sql_GetColumn(dbconn, &exists, dbconn.Rebind(sql), phone)
	return exists, err
}

func AddNewLicense(license License) (int64, error) {
	var id int64
	var err error

	sql := "Insert Into tbl_license (" +
		LICENSE_PHONENUMBER + ", " +
		LICENSE_GUID + ", " +
		LICENSE_COMPANYNAME + ", " +
		LICENSE_EXPIREDAY + ", " +
		LICENSE_EXPORTTIMES +
		") Values (?, ?, ?, ?, ?)"

	id, err = Sql_Insert(dbconn, LICENSE_ID, dbconn.Rebind(sql),
		license.PhoneNumber,
		license.GUID,
		license.CompanyName,
		license.ExpireDay,
		license.ExportTimes)
	return id, err
}

func ExtendExpireDay(phonenumber string, days int, times int) error {
	sql := `Update tbl_license Set ` + LICENSE_EXPIREDAY + ` = ` +
		`date(` + LICENSE_EXPIREDAY + `, 'start of day', ` + strconv.Itoa(days) + ` days)` + `, ` +
		`LICENSE_EXPORTTIMES = LICENSE_EXPORTTIMES+` + strconv.Itoa(times) + ` ` +
		`Where ` + LICENSE_PHONENUMBER + `=?`
	_, err := Sql_UpdDel(dbconn, dbconn.Rebind(sql), phonenumber)
	return err
}

func ExtendExpireYear(phonenumber string, years int, times int) error {
	sql := `Update tbl_license Set ` + LICENSE_EXPIREDAY + ` = ` +
		`date(` + LICENSE_EXPIREDAY + `, 'start of day', ` + strconv.Itoa(years) + ` year)` + `, ` +
		`LICENSE_EXPORTTIMES = LICENSE_EXPORTTIMES+` + strconv.Itoa(times) + ` ` +
		`Where ` + LICENSE_PHONENUMBER + `=?`
	_, err := Sql_UpdDel(dbconn, dbconn.Rebind(sql), phonenumber)
	return err
}

func ExtendExpireDate(phonenumber string, expiredate time.Time, times int) error {
	sql := `Update tbl_license Set ` + LICENSE_EXPIREDAY + ` = ` +
		expiredate.Format("2006-01-02") + `, ` +
		`LICENSE_EXPORTTIMES = LICENSE_EXPORTTIMES+` + strconv.Itoa(times) + ` ` +
		`Where ` + LICENSE_PHONENUMBER + `=?`
	_, err := Sql_UpdDel(dbconn, dbconn.Rebind(sql), phonenumber)
	return err
}

func GetLicenseByPhoneNumber(phonenumber string) (*License, error) {
	var record License
	sql := "Select " + LICENSE_ID + ", " +
		LICENSE_PHONENUMBER + ", " +
		LICENSE_GUID + ", " +
		LICENSE_COMPANYNAME + ", " +
		LICENSE_EXPIREDAY + " From tbl_license" +
		" Where " + LICENSE_PHONENUMBER + " = ?"
	e := Sql_Get(dbconn, &record, dbconn.Rebind(sql), phonenumber)
	return &record, e
}
