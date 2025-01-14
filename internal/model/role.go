package model

type Role struct {
	Base

	Name        string
	Permissions []*Permission `gorm:"many2many:role_permissions"`
}
