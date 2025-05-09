package repository

import (
	"context"
	"tablelink/src/entity"
	"tablelink/transaction"

	"gorm.io/gorm"
)

func NewUserRepository(gorm *gorm.DB) UserRepository {
	return &userRepository{
		gorm: gorm,
	}
}

func (repo *userRepository) DeleteById(ctx context.Context, userId int64) error {
	err := transaction.UnwrapContext(ctx, repo.gorm).Model(entity.User{}).Where("id=?", userId).Delete(nil).Error
	if err != nil {
		return err
	}

	return nil
}

func (repo *userRepository) Upsert(ctx context.Context, data *entity.User) error {
	err := repo.gorm.WithContext(ctx).Save(data).Error
	if err != nil {
		return err
	}

	return nil
}

func (repo *userRepository) FindById(ctx context.Context, userId int64) (*entity.User, error) {
	var data entity.User
	err := repo.gorm.WithContext(ctx).Where("id=?", userId).First(&data).Error
	if err != nil {
		return nil, err
	}

	return &data, nil
}

func (repo *userRepository) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
	var data entity.User
	err := repo.gorm.
		WithContext(ctx).
		Where("email=?", email).
		First(&data).Error
	if err != nil {
		return nil, err
	}

	return &data, nil
}

func (repo *userRepository) Find(ctx context.Context) ([]*entity.UserWithRole, error) {
	var data []*entity.UserWithRole
	err := repo.gorm.
		Select("users.id, users.name, users.email, role_rights.role_name, role_rights.id as role_id").
		Joins("JOIN roles ON roles.user_id=users.id").
		Joins("JOIN role_rights ON role_rights.id=roles.role_right_id").
		WithContext(ctx).
		Find(&data).Error
	if err != nil {
		return nil, err
	}

	return data, nil
}
