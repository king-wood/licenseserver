DROP TABLE IF EXISTS tbl_license;
CREATE TABLE IF NOT EXISTS tbl_license (
        id integer NOT NULL PRIMARY KEY AUTOINCREMENT, 
        phone_number varchar(255) NOT NULL,
        guid varchar(255) NOT NULL,
        company_name varchar(255) NOT NULL DEFAULT '',
        expire_day date NOT NULL,
        export_times integer NOT NULL
);

DROP TABLE IF EXISTS tbl_serial;
CREATE TABLE IF NOT EXISTS tbl_serial (
        id integer NOT NULL PRIMARY KEY AUTOINCREMENT, 
        phone_number varchar(255),
        serial varchar(255) NOT NULL,
        status integer NOT NULL DEFAULT 0,
        expire_day varchar(255) NOT NULL,
        export_times integer NOT NULL DEFAULT 0,
        pc_id varchar(255) 
);

DROP TABLE IF EXISTS tbl_administrator;
CREATE TABLE IF NOT EXISTS tbl_administrator (
        id integer NOT NULL PRIMARY KEY AUTOINCREMENT,
        user_name varchar(255) NOT NULL,
        password varchar(255) NOT NULL
);
