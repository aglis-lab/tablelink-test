package service

import (
	"context"
	"errors"
	"tablelink/src/entity"
	"tablelink/src/grpc"
	"tablelink/src/utils"
)

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

func (service *usersService) guard(ctx context.Context, accessToken string) (*entity.UserToken, *entity.UserRoles, *entity.RoleRight, error) {
	// Check Access Token
	user, err := utils.ParseAccessToken(accessToken)
	if err != nil {
		return nil, nil, nil, err
	}

	// Check Have Role
	userRole, roleRight, err := service.authenticationService.ValidRoleByUserId(ctx, user.Id, entity.RoleRight{RightRead: true})
	if err != nil {
		return nil, nil, nil, err
	}

	return user, userRole, roleRight, nil
}

func (service *usersService) GetUser(ctx context.Context, req *grpc.GetUserRequest) (*grpc.GetUserResponse, error) {
	// Check if user have access
	user, userRole, roleRight, err := service.guard(ctx, req.AccessToken)
	if err != nil {
		return nil, err
	}

	// Get User
	return &grpc.GetUserResponse{
		Status:  true,
		Message: "Successfully",
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
	_, _, _, err := service.guard(ctx, req.AccessToken)
	if err != nil {
		return nil, err
	}

	// Get all users
	users, err := service.userRepository.Find(ctx)
	if err != nil {
		return nil, err
	}

	resp := grpc.FetchUserResponse{}
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

// func (service *usersService) GetUser(ctx context.Context, _ *emptypb.Empty) (*grpc.GetUserResponse, error) {
// 	// fmt.Printf("--- BidirectionalStreamingEcho ---\n")

// 	// fmt.Printf("--- UnaryEcho ---\n")

// 	// md, ok := metadata.FromIncomingContext(ctx)
// 	// if !ok {
// 	// 	return nil, status.Errorf(codes.Internal, "UnaryEcho: missing incoming metadata in rpc context")
// 	// }

// 	// // Read and print metadata added by the interceptor.
// 	// if v, ok := md["key1"]; ok {
// 	// 	fmt.Printf("key1 from metadata: \n")
// 	// 	for i, e := range v {
// 	// 		fmt.Printf(" %d. %s\n", i, e)
// 	// 	}
// 	// }

// 	return &grpc.GetUserResponse{
// 		Status:  true,
// 		Message: "success",
// 		Data:    &grpc.UserResponse{},
// 	}, nil
// }
