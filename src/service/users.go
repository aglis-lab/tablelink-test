package service

import (
	"context"
	"errors"
	"tablelink/src/entity"
	"tablelink/src/grpc"
	"tablelink/src/utils"
	"tablelink/transaction"
)

const successMessage = "Successfully"

func (service *usersService) Login(ctx context.Context, req *grpc.LoginRequest) (*grpc.LoginResponse, error) {
	// Find User by email
	data, err := service.userRepository.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}

	if !utils.VerifyPassword(req.Password, data.Password) {
		return nil, errors.New("password is wrong")
	}

	accessToken, err := utils.SignAccessToken(entity.UserToken{
		Id:        data.Id,
		Name:      data.Name,
		Email:     data.Email,
		CreatedAt: data.CreatedAt,
	})
	if err != nil {
		return nil, err
	}

	return &grpc.LoginResponse{
		Status:  true,
		Message: "success to login",
		Data: &grpc.Token{
			AccessToken: accessToken,
		},
	}, nil
}

func (service *usersService) guard(ctx context.Context, accessToken string, targetRole entity.RoleRight) (*entity.UserToken, *entity.UserRoles, *entity.RoleRight, error) {
	// Check Access Token
	user, err := utils.ParseAccessToken(accessToken)
	if err != nil {
		return nil, nil, nil, err
	}

	// Check Have Role
	userRole, roleRight, err := service.authenticationService.ValidRoleByUserId(ctx, user.Id, targetRole)
	if err != nil {
		return nil, nil, nil, err
	}

	return user, userRole, roleRight, nil
}

func (service *usersService) GetUser(ctx context.Context, req *grpc.GetUserRequest) (*grpc.GetUserResponse, error) {
	// Check if user have access
	user, userRole, roleRight, err := service.guard(ctx, req.AccessToken, entity.RoleRight{RightRead: true})
	if err != nil {
		return nil, err
	}

	// Get User
	return &grpc.GetUserResponse{
		Status:  true,
		Message: successMessage,
		Data: &grpc.UserResponse{
			User: &grpc.User{
				RoleId:   int32(userRole.Id),
				Name:     user.Name,
				Email:    user.Email,
				RoleName: roleRight.RoleName,
			},
		},
	}, nil
}

func (service *usersService) FetchUser(ctx context.Context, req *grpc.FetchUserRequest) (*grpc.FetchUserResponse, error) {
	// Check if user have access
	_, _, _, err := service.guard(ctx, req.AccessToken, entity.RoleRight{RightRead: true})
	if err != nil {
		return nil, err
	}

	// Get all users
	users, err := service.userRepository.Find(ctx)
	if err != nil {
		return nil, err
	}

	resp := grpc.FetchUserResponse{
		Status:  true,
		Message: successMessage,
	}
	for _, v := range users {
		resp.Data = append(resp.Data, &grpc.User{
			RoleId:   int32(v.RoleId),
			RoleName: v.RoleName,
			Name:     v.Name,
			Email:    v.Email,
		})
	}

	return &resp, nil
}

func (service *usersService) CreateUser(ctx context.Context, req *grpc.CreateUserRequest) (*grpc.CreateUserResponse, error) {
	// Check if user have access
	_, _, _, err := service.guard(ctx, req.AccessToken, entity.RoleRight{RightCreate: true})
	if err != nil {
		return nil, err
	}

	// Create User
	hashPass, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	user := entity.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: hashPass,
	}
	err = service.userRepository.Upsert(ctx, &user)
	if err != nil {
		return nil, err
	}

	// Create User Role
	userRole := entity.UserRoles{
		UserId:      user.Id,
		RoleRightId: req.RoleId,
	}
	err = service.userRoleRepository.Create(ctx, &userRole)
	if err != nil {
		return nil, err
	}

	return &grpc.CreateUserResponse{
		Status:  true,
		Message: successMessage,
	}, nil
}

func (service *usersService) UpdateUser(ctx context.Context, req *grpc.UpdateUserRequest) (*grpc.UpdateUserResponse, error) {
	// Check if user have access
	_, _, _, err := service.guard(ctx, req.AccessToken, entity.RoleRight{RightUpdate: true})
	if err != nil {
		return nil, err
	}

	// Get User
	user, err := service.userRepository.FindById(ctx, req.UserId)
	if err != nil {
		return nil, err
	}

	// Update User
	user.Name = req.Name
	err = service.userRepository.Upsert(ctx, user)
	if err != nil {
		return nil, err
	}

	return &grpc.UpdateUserResponse{
		Status:  true,
		Message: successMessage,
	}, nil
}

func (service *usersService) DeleteUser(ctx context.Context, req *grpc.DeleteUserRequest) (*grpc.DeleteUserResponse, error) {
	// Check if user have access
	_, _, _, err := service.guard(ctx, req.AccessToken, entity.RoleRight{RightDelete: true})
	if err != nil {
		return nil, err
	}

	err = service.gormTransaction.WithGormTransaction(ctx, func(tx transaction.GormContext) error {
		err = service.userRepository.DeleteById(ctx, req.UserId)
		if err != nil {
			return err
		}

		err = service.userRoleRepository.DeleteByUserId(ctx, req.UserId)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return &grpc.DeleteUserResponse{
		Status:  true,
		Message: successMessage,
	}, nil
}
