package main

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Role struct {
	ID   string `gorm:"type:char(36);primaryKey"` // UUID as primary key
	Name string `gorm:"size:255;not null"`
}

type Permission struct {
	ID   string `gorm:"type:char(36);primaryKey"` // UUID as primary key
	Name string `gorm:"size:255;not null"`
}

type Object struct {
	ID   string `gorm:"type:char(36);primaryKey"` // UUID as primary key
	Name string `gorm:"size:255;not null"`
}

type RolePermission struct {
	RoleID       string     `gorm:"type:char(36);primaryKey"` // Composite primary key
	PermissionID string     `gorm:"type:char(36);primaryKey"` // Composite primary key
	ObjectID     string     `gorm:"type:char(36);primaryKey"` // Composite primary key
	Role         Role       `gorm:"foreignKey:RoleID"`
	Permission   Permission `gorm:"foreignKey:PermissionID"`
	Object       Object     `gorm:"foreignKey:ObjectID"`
}

func (r *Role) BeforeCreate(tx *gorm.DB) (err error) {
	r.ID = uuid.NewString() // Generate UUID
	return
}

func (o *Object) BeforeCreate(tx *gorm.DB) (err error) {
	o.ID = uuid.NewString() // Generate UUID
	return
}

func (p *Permission) BeforeCreate(tx *gorm.DB) (err error) {
	p.ID = uuid.NewString() // Generate UUID
	return
}
