package main

import (
	"fmt"
	"log"
	"time"

	nebula "github.com/vesoft-inc/nebula-go/v3"
)

func CheckResultSet(prefix string, res *nebula.ResultSet) {
	if !res.IsSucceed() {
		log.Fatal(fmt.Sprintf("%s, ErrorCode: %v, ErrorMsg: %s", prefix, res.GetErrorCode(), res.GetErrorMsg()))
	}
}

func NewSessionPool() {
	hostAddress := nebula.HostAddress{Host: address, Port: port}

	// Create configs for session pool
	config, err := nebula.NewSessionPoolConf(
		"root",
		"nebula",
		[]nebula.HostAddress{hostAddress},
		SPACE,
		nebula.WithHTTP2(useHTTP2),
	)
	if err != nil {
		log.Fatal(fmt.Sprintf("failed to create session pool config, %s", err.Error()))
	}

	// create session pool
	SessionPool, err = nebula.NewSessionPool(*config, nebula.DefaultLogger{})
	if err != nil {
		log.Fatal(fmt.Sprintf("failed to initialize session pool, %s", err.Error()))
	}
}

// Just a helper function to create a space for this example to run.
func prepareSpace() {
	hostAddress := nebula.HostAddress{Host: address, Port: port}
	hostList := []nebula.HostAddress{hostAddress}
	// Create configs for connection pool using default values
	testPoolConfig := nebula.GetDefaultConf()
	testPoolConfig.UseHTTP2 = useHTTP2

	// Initialize connection pool
	pool, err := nebula.NewConnectionPool(hostList, testPoolConfig, nebula.DefaultLogger{})
	if err != nil {
		log.Fatal(fmt.Sprintf("Fail to initialize the connection pool, host: %s, port: %d, %s", address, port, err.Error()))
	}
	// Close all connections in the pool
	defer pool.Close()

	// Create session
	session, err := pool.GetSession(username, password)
	if err != nil {
		log.Fatal(fmt.Sprintf("Fail to create a new session from connection pool, username: %s, password: %s, %s",
			username, password, err.Error()))
	}
	// Release session and return connection back to connection pool
	defer session.Release()

	checkResultSet := func(prefix string, res *nebula.ResultSet) {
		if !res.IsSucceed() {
			log.Fatal(fmt.Sprintf("%s, ErrorCode: %v, ErrorMsg: %s", prefix, res.GetErrorCode(), res.GetErrorMsg()))
		}
	}

	{
		// Prepare the query
		createSchema := fmt.Sprintf(`
			CREATE SPACE IF NOT EXISTS %s (vid_type=FIXED_STRING(20)); 
		`, SPACE)

		// Execute a query
		resultSet, err := session.Execute(createSchema)
		if err != nil {
			fmt.Print(err.Error())
			return
		}
		checkResultSet(createSchema, resultSet)
	}

	time.Sleep(5 * time.Second)

	log.Println("Space example_space was created")
}
