package repository

import (
	"context"
	"tablelink/src/entity"

	"gorm.io/gorm"
)

type UserRepository interface {
	FindByEmail(context.Context, string) (*entity.User, error)
	Find(ctx context.Context) ([]*entity.UserWithRole, error)
}

type userRepository struct {
	gorm *gorm.DB
}

type RoleRightRepository interface {
	FindById(context.Context, int64) (*entity.RoleRight, error)
	FindByIds(context.Context, []int64) ([]*entity.RoleRight, error)
}

type roleRightRepository struct {
	gorm *gorm.DB
}

type UserRoleRepository interface {
	FindByUserId(context.Context, int64) (*entity.UserRoles, error)
	FindByUserIds(context.Context, []int64) ([]*entity.UserRoles, error)
}
type userRoleRepository struct {
	gorm *gorm.DB
}
