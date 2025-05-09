package migration

import (
	"tablelink/src/entity"
	"time"

	"gorm.io/gorm"
)

func Init(gormDB *gorm.DB) error {
	err := gormDB.AutoMigrate(&entity.User{})
	if err != nil {
		return err
	}

	roleRight := entity.RoleRight{
		RoleName:    "admin",
		RightCreate: true,
		RightRead:   true,
		RightUpdate: true,
		RightDelete: true,
		CreatedAt:   time.Now(),
	}
	err = gormDB.Save(&roleRight).Error
	if err != nil {
		return err
	}

	user := entity.User{
		Name:      "testing",
		Email:     "testing@mail.com",
		Password:  "$2a$14$7M/f0Z/GxcXpyiaaz1PgDO/9V/mf6vrqowI3bqm.rmXeAezyFXW1q",
		CreatedAt: time.Now(),
	}
	err = gormDB.Save(&user).Error
	if err != nil {
		return err
	}

	userRole := entity.UserRoles{
		UserId:      user.Id,
		RoleRightId: roleRight.Id,
	}
	err = gormDB.Save(&userRole).Error
	if err != nil {
		return err
	}
	return nil
}
