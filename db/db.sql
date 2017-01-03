DROP TABLE IF EXISTS tbl_license;
CREATE TABLE IF NOT EXISTS tbl_license(
        id integer NOT NULL PRIMARY KEY AUTOINCREMENT, 
        phone_number varchar(255) NOT NULL,
        guid varchar(255) NOT NULL,
        company_name varchar(255) NOT NULL DEFAULT '',
        expire_day date NOT NULL,
        export_times integer NOT NULL
);
