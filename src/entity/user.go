package entity

import (
	"time"
)

type User struct {
	Id        int64     `gorm:"column:id;primaryKey"`
	Name      string    `gorm:"column:name"`
	Email     string    `gorm:"column:email"`
	Password  string    `gorm:"column:password"`
	CreatedAt time.Time `gorm:"column:created_at"`
}

type UserWithRole struct {
	Id        int64     `gorm:"column:id;primaryKey"`
	RoleId    int64     `gorm:"column:role_id"`
	RoleName  string    `gorm:"column:role_name"`
	Name      string    `gorm:"column:name"`
	Email     string    `gorm:"column:email"`
	CreatedAt time.Time `gorm:"column:created_at"`
}

func (UserWithRole) TableName() string {
	return "users"
}
