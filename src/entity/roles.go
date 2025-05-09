package entity

type UserRoles struct {
	Id          int64 `gorm:"column:id;primaryKey"`
	UserId      int64 `gorm:"column:user_id"`
	RoleRightId int64 `gorm:"column:role_right_id"`
}

func (UserRoles) TableName() string {
	return "roles"
}
