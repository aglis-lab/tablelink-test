package service

import (
	"context"
	"errors"
	"tablelink/src/entity"
	"tablelink/src/repository"
)

type AuthenticationService interface {
	ValidRoleByUserId(context.Context, int64, entity.RoleRight) (*entity.UserRoles, *entity.RoleRight, error)
}

type authenticationService struct {
	userRepository      repository.UserRepository
	userRoleRepository  repository.UserRoleRepository
	roleRightRepository repository.RoleRightRepository
}

func NewAuthenticationService(
	userRepository repository.UserRepository,
	userRoleRepository repository.UserRoleRepository,
	roleRightRepository repository.RoleRightRepository,
) AuthenticationService {
	return &authenticationService{
		userRepository:      userRepository,
		userRoleRepository:  userRoleRepository,
		roleRightRepository: roleRightRepository,
	}
}

func (service *authenticationService) ValidRoleByUserId(ctx context.Context, userId int64, role entity.RoleRight) (*entity.UserRoles, *entity.RoleRight, error) {
	// Check if have right to read
	userRole, err := service.userRoleRepository.FindByUserId(ctx, userId)
	if err != nil {
		return nil, nil, err
	}

	// Get role right
	roleRight, err := service.roleRightRepository.FindById(ctx, userRole.RoleRightId)
	if err != nil {
		return nil, nil, err
	}

	if role.RightCreate && !roleRight.RightCreate {
		return nil, nil, errors.New("doesn't have role to create")
	}

	if role.RightDelete && !roleRight.RightDelete {
		return nil, nil, errors.New("doesn't have role to delete")
	}

	if role.RightRead && !roleRight.RightRead {
		return nil, nil, errors.New("doesn't have role to read")
	}

	if role.RightUpdate && !roleRight.RightUpdate {
		return nil, nil, errors.New("doesn't have role to update")
	}

	return userRole, roleRight, nil
}
