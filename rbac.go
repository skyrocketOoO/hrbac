package main

import (
	"hrbac/global"

	nebula "github.com/vesoft-inc/nebula-go/v3"
)

const BasicSchema = `
	CREATE TAG IF NOT EXISTS user();
	CREATE TAG IF NOT EXISTS role();
	CREATE TAG IF NOT EXISTS device();

	CREATE EDGE IF NOT EXISTS has_permission(type STRING NOT NULL);
	CREATE EDGE IF NOT EXISTS belongs_to();
	CREATE EDGE IF NOT EXISTS leader_of();

	// INSERT VERTEX user() VALUES "user_1":(), "user_2":();
	// INSERT VERTEX role() VALUES "role_1":(), "role_2":();
	// INSERT VERTEX device() VALUES "device_1":(), "device_2":(), "device_3":();

	// INSERT EDGE has_permission(type) VALUES "role_2"->"device_3":("write");
	// INSERT EDGE belongs_to() VALUES "user_1"->"role_1":();
	// INSERT EDGE belongs_to() VALUES "user_2"->"role_2":();
	// INSERT EDGE has_permission(type) VALUES "role_1"->"device_1":("write");
	// INSERT EDGE has_permission(type) VALUES "role_1"->"device_2":("write");
	// INSERT EDGE leader_of() VALUES "role_1"->"role_2":();
`

func Exec(schema string) (res *nebula.ResultSet, err error) {
	res, err = global.SessionPool.Execute(schema)
	CheckResultSet(schema, res)
	return
}
