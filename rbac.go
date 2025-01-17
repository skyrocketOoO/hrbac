package main

import (
	"fmt"

	nebula "github.com/vesoft-inc/nebula-go/v3"
)

const BasicSchema = `
	CREATE TAG IF NOT EXISTS user();
	CREATE TAG IF NOT EXISTS role();
	CREATE TAG IF NOT EXISTS device();

	CREATE EDGE IF NOT EXISTS has_permission(type STRING NOT NULL);
	CREATE EDGE IF NOT EXISTS belongs_to();
	CREATE EDGE IF NOT EXISTS leader_of();

	INSERT VERTEX user() VALUES "user_1":(), "user_2":();
	INSERT VERTEX role() VALUES "role_1":(), "role_2":();
	INSERT VERTEX device() VALUES "device_1":(), "device_2":(), "device_3":();

	INSERT EDGE has_permission(type) VALUES "role_2"->"device_3":("write");
	INSERT EDGE belongs_to() VALUES "user_1"->"role_1":();
	INSERT EDGE belongs_to() VALUES "user_2"->"role_2":();
	INSERT EDGE has_permission(type) VALUES "role_1"->"device_1":("write");
	INSERT EDGE has_permission(type) VALUES "role_1"->"device_2":("write");
	INSERT EDGE leader_of() VALUES "role_1"->"role_2":();

	MATCH (v)-[belongs_to]->(:role)-[leader_of*0..]->(:role)-[has_permission]->(d:device) WHERE id(v) == 'user_1' RETURN d;

`

func CheckPermission(userID, permissionType, deviceID string) (ok bool, err error) {
	schema := fmt.Sprintf(
		`MATCH (v)-[belongs_to]->(:role)-[leader_of*0..]->(:role)`+
			`-[has_permission{type:'%s'}]->(d:device) WHERE id(v) == '%s' AND id(d) == '%s' RETURN d;`,
		permissionType, userID, deviceID,
	)

	resp, err := SessionPool.Execute(schema)
	if err != nil {
		return false, fmt.Errorf("query execution failed: %v", err)
	}

	if !resp.IsSucceed() {
		return false, fmt.Errorf("query failed: %s", resp.GetErrorMsg())
	}

	if resp.GetRowSize() > 0 {
		return true, nil
	}

	return false, nil
}

func Exec(schema string) (res *nebula.ResultSet, err error) {
	res, err = SessionPool.Execute(schema)
	CheckResultSet(schema, res)
	return
}
