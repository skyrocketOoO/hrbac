package global

import (
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

const (
	ApiVersion = "/v1"
)

var (
	// env | flag
	Database    string
	AutoMigrate bool = false
	Env         string

	// instance
	DB        *gorm.DB
	Validator *validator.Validate // use a single instance of Validate, it caches struct info
)
