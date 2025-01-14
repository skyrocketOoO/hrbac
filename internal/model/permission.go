package model

type Permission struct {
	Base

	Name string
	Roles []*Role `gorm:"many2many:role_permissions"`
}
