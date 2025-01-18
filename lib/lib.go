package lib

import (
	"fmt"

	"hrbac/global"
)

func CheckPermission(userID, permissionType, deviceID string) (ok bool, err error) {
	sql := fmt.Sprintf(
		`MATCH (v)-[belongs_to]->(:role)-[leader_of*0..]->(:role)`+
			`-[has_permission{type:'%s'}]->(d:device) WHERE id(v) == '%s' AND id(d) == '%s' RETURN d LIMIT 1;`,
		permissionType, userID, deviceID,
	)

	resp, err := global.SessionPool.Execute(sql)
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
