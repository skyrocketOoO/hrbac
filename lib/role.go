package lib

import (
	"fmt"

	"hrbac/global"
)

func GetRoles() {}
func AddRole(id string) error {
	sql := fmt.Sprintf(`INSERT VERTEX IF NOT EXISTS role() VALUES "%s":();`, id)

	resp, err := global.SessionPool.Execute(sql)
	if err != nil {
		return fmt.Errorf("query execution failed: %v", err)
	}

	if !resp.IsSucceed() {
		return fmt.Errorf("query failed: %s", resp.GetErrorMsg())
	}
	return nil
}
func DeleteRoles() {}

func GetRolePermissions() {}

func AddLeader(leaderID string, roleID string) error {
	sql := fmt.Sprintf(`INSERT EDGE IF NOT EXISTS leader_of() VALUES "%s"->"%s":();`, leaderID, roleID)

	resp, err := global.SessionPool.Execute(sql)
	if err != nil {
		return fmt.Errorf("query execution failed: %v", err)
	}

	if !resp.IsSucceed() {
		return fmt.Errorf("query failed: %s", resp.GetErrorMsg())
	}
	return nil
}
func RemoveLeader() {}

func AssignRole(userID, roleID string) error {
	sql := fmt.Sprintf(`INSERT EDGE IF NOT EXISTS belongs_to() VALUES "%s"->"%s":();`, userID, roleID)

	resp, err := global.SessionPool.Execute(sql)
	if err != nil {
		return fmt.Errorf("query execution failed: %v", err)
	}

	if !resp.IsSucceed() {
		return fmt.Errorf("query failed: %s", resp.GetErrorMsg())
	}
	return nil
}
func UnassignRole() {}

func GetRoleMembers() {}
