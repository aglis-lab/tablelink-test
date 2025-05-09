package entity

import "time"

type RoleRight struct {
	Id          int64     `gorm:"column:id;primaryKey"`
	RoleName    string    `gorm:"column:role_name"`
	RightCreate bool      `gorm:"column:right_create"`
	RightRead   bool      `gorm:"column:right_read"`
	RightUpdate bool      `gorm:"column:right_update"`
	RightDelete bool      `gorm:"column:right_delete"`
	CreatedAt   time.Time `gorm:"column:created_at"`
}
