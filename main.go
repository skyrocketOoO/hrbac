package main

import (
	"fmt"
	"log"
	"time"

	"hrbac/global"
	"hrbac/lib"
)

func main() {
	prepareSpace()
	NewSessionPool()
	defer global.SessionPool.Close()

	if _, err := Exec(BasicSchema); err != nil {
		log.Fatal(err)
	}
	time.Sleep(5 * time.Second)

	if err := lib.AddUser("user_1"); err != nil {
		log.Fatal(err)
	}

	if err := lib.AddUser("user_2"); err != nil {
		log.Fatal(err)
	}

	if err := lib.AddRole("role_1"); err != nil {
		log.Fatal(err)
	}

	if err := lib.AddRole("role_2"); err != nil {
		log.Fatal(err)
	}

	if err := lib.AddDevice("device_1"); err != nil {
		log.Fatal(err)
	}

	if err := lib.AddDevice("device_2"); err != nil {
		log.Fatal(err)
	}

	if err := lib.AddDevice("device_3"); err != nil {
		log.Fatal(err)
	}

	if err := lib.AssignPermission("device_3", "write", "role_2"); err != nil {
		log.Fatal(err)
	}

	if err := lib.AssignPermission("device_2", "write", "role_1"); err != nil {
		log.Fatal(err)
	}

	if err := lib.AssignPermission("device_1", "write", "role_1"); err != nil {
		log.Fatal(err)
	}

	if err := lib.AssignRole("user_1", "role_1"); err != nil {
		log.Fatal(err)
	}

	if err := lib.AssignRole("user_2", "role_2"); err != nil {
		log.Fatal(err)
	}

	if err := lib.AddLeader("role_1", "role_2"); err != nil {
		log.Fatal(err)
	}

	ok, err := lib.CheckPermission("user_1", "write", "device_1")
	if err != nil {
		log.Fatal(err)
	}
	if ok {
		fmt.Println("user_1 has write permission on device_1")
	} else {
		fmt.Println("user_1 does not have write permission on device_1")
	}

	ok, err = lib.CheckPermission("user_1", "write", "device_3")
	if err != nil {
		log.Fatal(err)
	}
	if ok {
		fmt.Println("user_1 has write permission on device_3")
	} else {
		fmt.Println("user_1 does not have write permission on device_3")
	}

	// Drop space
	{
		query := fmt.Sprintf(`DROP SPACE IF EXISTS %s`, global.SPACE)
		// Send query
		resultSet, err := global.SessionPool.Execute(query)
		if err != nil {
			fmt.Print(err.Error())
			return
		}
		CheckResultSet(query, resultSet)
	}
	fmt.Print("\n")
	log.Println("Nebula Go Client Session Pool Example Finished")
}
