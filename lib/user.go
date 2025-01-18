package lib

import (
	"fmt"

	"hrbac/global"
)

func GetUsers()    {}
func DeleteUsers() {}
func AddUser(id string) error {
	sql := fmt.Sprintf(`INSERT VERTEX IF NOT EXISTS user() VALUES "%s":();`, id)

	resp, err := global.SessionPool.Execute(sql)
	if err != nil {
		return fmt.Errorf("query execution failed: %v", err)
	}

	if !resp.IsSucceed() {
		return fmt.Errorf("query failed: %s", resp.GetErrorMsg())
	}
	return nil
}
func GetUserPermissions() {}
