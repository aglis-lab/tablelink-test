package repository

import (
	"context"
	"tablelink/src/entity"
	"tablelink/transaction"

	"gorm.io/gorm"
)

func NewUserRoleRepository(gorm *gorm.DB) UserRoleRepository {
	return &userRoleRepository{
		gorm: gorm,
	}
}

func (repo *userRoleRepository) Create(ctx context.Context, data *entity.UserRoles) error {
	err := repo.gorm.WithContext(ctx).Save(data).Error
	if err != nil {
		return err
	}

	return nil
}

func (repo *userRoleRepository) DeleteByUserId(ctx context.Context, userId int64) error {
	err := transaction.UnwrapContext(ctx, repo.gorm).Model(entity.UserRoles{}).Where("user_id=?", userId).Delete(nil).Error
	if err != nil {
		return err
	}

	return nil
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
