package lib

import (
	"fmt"

	"hrbac/global"
)

func GetDevices() {}
func AddDevice(id string) error {
	sql := fmt.Sprintf(`INSERT VERTEX IF NOT EXISTS device() VALUES "%s":();`, id)

	resp, err := global.SessionPool.Execute(sql)
	if err != nil {
		return fmt.Errorf("query execution failed: %v", err)
	}

	if !resp.IsSucceed() {
		return fmt.Errorf("query failed: %s", resp.GetErrorMsg())
	}
	return nil
}
func DeleteDevices() {}

func AssignPermission(deviceID, permissionType, roleID string) error {
	sql := fmt.Sprintf(
		`INSERT EDGE IF NOT EXISTS has_permission(type) VALUES "%s"->"%s":("%s");`,
		roleID, deviceID, permissionType,
	)

	resp, err := global.SessionPool.Execute(sql)
	if err != nil {
		return fmt.Errorf("query execution failed: %v", err)
	}

	if !resp.IsSucceed() {
		return fmt.Errorf("query failed: %s", resp.GetErrorMsg())
	}
	return nil
}
func UnassignPermission() {}
