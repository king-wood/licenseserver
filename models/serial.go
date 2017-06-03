package models

import log "github.com/cihub/seelog"

const (
	SERIAL_ID           = "id"
	SERIAL_PHONE_NUMBER = "phone_number"
	SERIAL_PCID         = "pc_id"
	SERIAL_SERIAL       = "serial"
	SERIAL_STATUS       = "status"
	SERIAL_EXPIRE_DAY   = "expire_day"
	SERIAL_EXPORT_TIMES = "export_times"

	SERIAL_UNREGISTERED = 0
	SERIAL_REGISTERED   = 1
)

type Serial struct {
	ID          int64  `db:"id" json:"ID"`
	PhoneNumber string `db:"phone_number" json:"PhoneNumber"`
	Serial      string `db:"serial" json:"Serial"`
	Status      int    `db:"status" json:"Status"`
	ExpireDay   string `db:"expire_day" json:"ExpireDay"`
	ExportTimes int    `db:"export_times" json:"ExportTimes"`
	PCID        string `db:"pc_id" json:"PCID"`
}

func AddNewSerial(serial Serial) (int64, error) {
	var id int64
	var err error

	sql := "Insert Into tbl_serial (" +
		SERIAL_PHONE_NUMBER + ", " +
		SERIAL_SERIAL + ", " +
		SERIAL_STATUS + ", " +
		SERIAL_EXPIRE_DAY + ", " +
		SERIAL_PCID + ", " +
		SERIAL_EXPORT_TIMES +
		") Values (?, ?, ?, ?, ?, ?)"

	id, err = Sql_Insert(dbconn, SERIAL_ID, dbconn.Rebind(sql),
		serial.PhoneNumber,
		serial.Serial,
		serial.Status,
		serial.ExpireDay,
		serial.PCID,
		serial.ExportTimes)
	return id, err
}

func ExtendSerialExpireDate(serial string, phonenumber string, expiredate string, times int) error {
	sql := `Update tbl_serial Set ` +
		SERIAL_EXPIRE_DAY + ` = ?, ` +
		SERIAL_EXPORT_TIMES + " = ? " +
		`Where ` + SERIAL_SERIAL + `=? And ` +
		SERIAL_PHONE_NUMBER + " = ? "
	_, err := Sql_UpdDel(dbconn, dbconn.Rebind(sql), expiredate, times, serial, phonenumber)
	return err
}

func GetSerialBySerial(serial string) (*Serial, error) {
	var record Serial
	sql := "Select " + SERIAL_ID + ", " +
		SERIAL_PHONE_NUMBER + ", " +
		SERIAL_SERIAL + ", " +
		SERIAL_STATUS + ", " +
		SERIAL_EXPORT_TIMES + " , " +
		`strftime('%Y-%m-%d',` + SERIAL_EXPIRE_DAY + ") As " + SERIAL_EXPIRE_DAY +
		" From tbl_serial" +
		" Where " + SERIAL_SERIAL + " = ?"
	e := Sql_Get(dbconn, &record, dbconn.Rebind(sql), serial)
	return &record, e
}

func UpdateSerialStatus(serial string, status int) error {
	sql := "Update tbl_serial Set " + SERIAL_STATUS + " = ?" +
		" Where " + SERIAL_SERIAL + " = ?"
	_, err := Sql_UpdDel(dbconn, dbconn.Rebind(sql), serial, status)
	log.Info("sql command", sql)
	return err
}

func CheckSerialUse(serial string) (bool, error) {
	var exist bool
	sql := "Select Exists (Select 1 From tbl_serial Where " +
		SERIAL_SERIAL + " = ? And " +
		SERIAL_STATUS + " = ? )"
	e := Sql_Get(dbconn, &exist, dbconn.Rebind(sql), serial, SERIAL_UNREGISTERED)
	return exist, e
}

func RegisterSerial(serial string, phonenumber string, pcid string) error {
	sql := "Update tbl_serial Set " + SERIAL_STATUS + " = ?, " +
		SERIAL_PHONE_NUMBER + " = ?, " +
		SERIAL_PCID + " = ? " +
		" Where " + SERIAL_SERIAL + " = ? And " +
		SERIAL_STATUS + " = ?"
	_, err := Sql_UpdDel(dbconn, dbconn.Rebind(sql), SERIAL_REGISTERED, phonenumber, pcid, serial, SERIAL_UNREGISTERED)
	log.Info("sql command", sql)
	return err
}

func CheckIfRegistered(serial string, phonenumber string) (bool, error) {
	var exist bool
	sql := "Select Exists (Select 1 From tbl_serial Where " +
		SERIAL_SERIAL + " = ? And " +
		SERIAL_PHONE_NUMBER + " = ? And " +
		SERIAL_STATUS + " = ? )"
	e := Sql_Get(dbconn, &exist, dbconn.Rebind(sql), serial, phonenumber, SERIAL_REGISTERED)
	return exist, e
}

func GetSerialBySerialAndPhoneNumber(serial string, phonenumber string) (*Serial, error) {
	var record Serial
	sql := "Select " + SERIAL_ID + ", " +
		SERIAL_PHONE_NUMBER + ", " +
		SERIAL_SERIAL + ", " +
		SERIAL_STATUS + ", " +
		SERIAL_EXPORT_TIMES + " , " +
		SERIAL_EXPIRE_DAY +
		" From tbl_serial" +
		" Where " + SERIAL_SERIAL + " = ? And " +
		SERIAL_PHONE_NUMBER + " = ?"
	e := Sql_Get(dbconn, &record, dbconn.Rebind(sql), serial, phonenumber)
	return &record, e
}

func GetAvailableSerial() ([]Serial, error) {
	var records []Serial
	sql := "Select " + SERIAL_ID + ", " +
		SERIAL_PHONE_NUMBER + ", " +
		SERIAL_SERIAL + ", " +
		SERIAL_STATUS + ", " +
		SERIAL_EXPORT_TIMES + " , " +
		SERIAL_EXPIRE_DAY +
		" From tbl_serial" +
		" Where " + SERIAL_STATUS + " = ? "
	e := Sql_Select(dbconn, &records, dbconn.Rebind(sql), SERIAL_UNREGISTERED)
	return records, e
}

func IsSerialMatchPC(serial string, phonenumber string, pcid string) (bool, error) {
	var exist bool
	sql := "Select Exists (Select 1 From tbl_serial Where " +
		SERIAL_SERIAL + " = ? And " +
		SERIAL_PHONE_NUMBER + " = ? And " +
		SERIAL_PCID + " = ? And " +
		SERIAL_STATUS + " = ? )"
	e := Sql_Get(dbconn, &exist, dbconn.Rebind(sql), serial, phonenumber, pcid, SERIAL_REGISTERED)
	return exist, e
}

func GetSerialByPhoneNumber(phonenumber string) ([]Serial, error) {
	var serials []Serial
	sql := "SELECT * FROM tbl_serial WHERE " + SERIAL_PHONE_NUMBER + " = ?"
	e := Sql_Select(dbconn, &serials, dbconn.Rebind(sql), phonenumber)
	return serials, e
}
