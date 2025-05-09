package repository

import (
	"context"
	"tablelink/src/entity"

	"gorm.io/gorm"
)

type UserRepository interface {
	FindByEmail(context.Context, string) (*entity.User, error)
	Find(context.Context) ([]*entity.UserWithRole, error)
	FindById(context.Context, int64) (*entity.User, error)
	Upsert(context.Context, *entity.User) error
	DeleteById(context.Context, int64) error
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
	Create(context.Context, *entity.UserRoles) error
	DeleteByUserId(context.Context, int64) error
}
type userRoleRepository struct {
	gorm *gorm.DB
}
