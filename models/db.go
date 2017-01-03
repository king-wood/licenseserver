package models

import (
	"errors"

	"fmt"

	log "github.com/cihub/seelog"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

const (
	DB_TYPE_SQLITE   = "sqlite3"
	DB_TYPE_POSTGRES = "postgres"
)

var dbconn *sqlx.DB
var dbType string
var dbUser *sqlx.DB
var dbUserType string

func InitDB(databaseType string, conStr string) {
	log.Info("text:connecting database ", databaseType)
	dbType = databaseType
	var err error
	dbconn, err = sqlx.Connect(dbType, conStr)
	if err != nil {
		log.Error("text:panic! failed to initialize db")
		panic(err)
	}
	log.Info("text:connected database ", databaseType)
}

func InitUserDB(databaseType string, conStr string) {
	log.Info("text:connecting account database ", databaseType)
	dbUserType = databaseType
	var err error
	dbUser, err = sqlx.Connect(dbUserType, conStr)
	if err != nil {
		log.Error("text:panic! failed to initialize account db")
		panic(err)
	}
	log.Info("text:connected account database ", databaseType)
}

func sql_InsertPg(db *sqlx.DB, sql string, param ...interface{}) (int64, error) {
	var id int64
	err := db.QueryRow(sql, param...).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func sql_InsertSqlite(db *sqlx.DB, sql string, param ...interface{}) (int64, error) {
	stmt, err := db.Prepare(sql)
	if err != nil {
		return 0, err
	}

	res, err := stmt.Exec(param...)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}

func Sql_AddMinute(off int) (sqlstatement string) {
	if off == 0 {
		sqlstatement = "CURRENT_TIMESTAMP"
	} else {
		switch dbType {
		case DB_TYPE_POSTGRES:
			sqlstatement = fmt.Sprintf(` CURRENT_TIMESTAMP + '%d min' `, off)
		case DB_TYPE_SQLITE:
			sqlstatement = fmt.Sprintf(` datetime('now','%d minute') `, off)
		}
	}
	return
}

func Sql_Insert(db *sqlx.DB, pk string, sql string, param ...interface{}) (int64, error) {
	var id int64
	var err error
	switch dbType {
	case DB_TYPE_POSTGRES:
		sql = sql + " RETURNING " + pk
		id, err = sql_InsertPg(db, sql, param...)
	case DB_TYPE_SQLITE:
		id, err = sql_InsertSqlite(db, sql, param...)
	default:
		err = errors.New("unknown db type")
	}

	return id, err
}

func Sql_UpdDel(db *sqlx.DB, sql string, param ...interface{}) (int64, error) {
	stmt, err := db.Prepare(sql)
	if err != nil {
		return 0, err
	}

	res, err := stmt.Exec(param...)
	if err != nil {
		return 0, err
	}

	rowCnt, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}

	return rowCnt, nil
}

func Sql_Get(db *sqlx.DB, dest interface{}, sql string, param ...interface{}) error {
	err := db.Get(dest, sql, param...)
	return err
}

func Sql_Select(db *sqlx.DB, dest interface{}, sql string, param ...interface{}) error {
	err := db.Select(dest, sql, param...)
	return err
}

func Sql_GetColumn(db *sqlx.DB, dest interface{}, sql string, param ...interface{}) error {
	err := db.QueryRow(sql, param...).Scan(dest)
	return err
}
