package main

import nebula "github.com/vesoft-inc/nebula-go/v3"

const (
	SPACE    = "my_space"
	address  = "127.0.0.1"
	port     = 9669
	username = "root"
	password = "nebula"
	useHTTP2 = false
)

var SessionPool *nebula.SessionPool
