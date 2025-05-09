package repository

import (
	"context"
	"tablelink/src/entity"

	"gorm.io/gorm"
)

func NewRoleRightRepository(gorm *gorm.DB) RoleRightRepository {
	return &roleRightRepository{
		gorm: gorm,
	}
}

func (repo *roleRightRepository) FindById(ctx context.Context, roleId int64) (*entity.RoleRight, error) {
	var data entity.RoleRight
	err := repo.gorm.
		WithContext(ctx).
		Where("id=?", roleId).
		First(&data).Error
	if err != nil {
		return nil, err
	}

	return &data, nil
}

func (repo *roleRightRepository) FindByIds(ctx context.Context, roleIds []int64) ([]*entity.RoleRight, error) {
	var data []*entity.RoleRight
	err := repo.gorm.
		WithContext(ctx).
		Where("id IN (?)", roleIds).
		First(&data).Error
	if err != nil {
		return nil, err
	}

	return data, nil
}
