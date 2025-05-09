package repository

import (
	"context"
	"tablelink/src/entity"

	"gorm.io/gorm"
)

func NewUserRoleRepository(gorm *gorm.DB) UserRoleRepository {
	return &userRoleRepository{
		gorm: gorm,
	}
}

func (repo *userRoleRepository) FindByUserId(ctx context.Context, userId int64) (*entity.UserRoles, error) {
	var data entity.UserRoles
	err := repo.gorm.
		WithContext(ctx).
		Where("user_id=?", userId).
		First(&data).Error
	if err != nil {
		return nil, err
	}

	return &data, nil
}

func (repo *userRoleRepository) FindByUserIds(ctx context.Context, userIds []int64) ([]*entity.UserRoles, error) {
	var data []*entity.UserRoles
	err := repo.gorm.
		WithContext(ctx).
		Where("user_id IN (?)", userIds).
		First(&data).Error
	if err != nil {
		return nil, err
	}

	return data, nil
}
