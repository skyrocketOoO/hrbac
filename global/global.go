package global

import nebula "github.com/vesoft-inc/nebula-go/v3"

const (
	SPACE    = "my_space"
	Address  = "127.0.0.1"
	Port     = 9669
	Username = "root"
	Password = "nebula"
	UseHTTP2 = false
)

var SessionPool *nebula.SessionPool
