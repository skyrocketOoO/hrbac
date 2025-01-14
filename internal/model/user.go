package model

type User struct {
	Base

	Name     string `gorm:"unique"`
	Password string `gorm:"type:varchar(255)"`

	RoleID string `gorm:"type:char(36)"`
	Role   Role   `gorm:"foreignKey:RoleID"`
}
