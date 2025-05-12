package handlers

import (
    "context"
    "log"
    pb "user-service/pkg/api"
    "user-service/pkg/jwt"
    "user-service/pkg/service"

    "google.golang.org/grpc/codes"
    "google.golang.org/grpc/status"
)

type UserHandler struct {
    pb.UnimplementedUserServiceServer
    Service *service.UserService
}

func NewUserHandler(svc *service.UserService) *UserHandler {
    return &UserHandler{Service: svc}
}

func (h *UserHandler) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
    userID, err := h.Service.Register(req.Username, string(req.Password), req.Email)
    if err != nil {
        log.Printf("Ошибка при регистрации: %v", err)
        return nil, err
    }
    return &pb.RegisterResponse{
        UserId: userID,
    }, nil
}

func (h *UserHandler) VerifyOTP(ctx context.Context, req *pb.VerifyOTPRequest) (*pb.VerifyOTPResponse, error) {
    if !h.Service.VerifyOTP(req.UserId, req.Otp) {
        log.Printf("Ошибка при верификации OTP для пользователя %s", req.UserId)
        return &pb.VerifyOTPResponse{
            Success: false,
            Message: "Invalid OTP",
        }, nil
    }
    return &pb.VerifyOTPResponse{
        Success: true,
        Message: "OTP verified successfully",
    }, nil
}

func (h *UserHandler) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
    user, err := h.Service.Login(req.Email, string(req.Password))
    if err != nil {
        log.Printf("Ошибка при логине: %v", err)
        return nil, status.Errorf(codes.Unauthenticated, "login failed: %v", err)
    }
    token, err := jwt.GenerateToken(user.ID)
    if err != nil {
        log.Printf("Ошибка при генерации токена: %v", err)
        return nil, status.Errorf(codes.Internal, "token generation failed: %v", err)
    }
    return &pb.LoginResponse{
        UserId: user.ID,
        Token:  token,
    }, nil
}

func (h *UserHandler) ForgotPassword(ctx context.Context, req *pb.ForgotPasswordRequest) (*pb.ForgotPasswordResponse, error) {
    err := h.Service.ForgotPassword(req.Email)
    if err != nil {
        log.Printf("Ошибка при отправке OTP: %v", err)
        return nil, status.Errorf(codes.Internal, "failed to send OTP: %v", err)
    }
    return &pb.ForgotPasswordResponse{
        Message: "OTP sent to your email",
    }, nil
}

func (h *UserHandler) ResetPassword(ctx context.Context, req *pb.ResetPasswordRequest) (*pb.ResetPasswordResponse, error) {
    err := h.Service.ResetPassword(req.Email, req.Otp, req.NewPassword)
    if err != nil {
        log.Printf("Ошибка при сбросе пароля: %v", err)
        return nil, status.Errorf(codes.Internal, "reset password failed: %v", err)
    }
    return &pb.ResetPasswordResponse{
        Message: "Password successfully reset",
    }, nil
}

func (h *UserHandler) ListUsers(ctx context.Context, req *pb.ListUsersRequest) (*pb.ListUsersResponse, error) {
    users, err := h.Service.ListUsers()
    if err != nil {
        log.Printf("Ошибка при получении пользователей: %v", err)
        return nil, status.Errorf(codes.Internal, "ошибка при получении пользователей: %v", err)
    }

    var pbUsers []*pb.User
    for _, u := range users {
        pbUsers = append(pbUsers, &pb.User{
            UserId:   u.ID,
            Username: u.Username,
            Email:    u.Email,
        })
    }

    return &pb.ListUsersResponse{Users: pbUsers}, nil
}

func (h *UserHandler) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.UserResponse, error) {
    user, err := h.Service.GetUser(req.UserId)
    if err != nil {
        log.Printf("Ошибка при получении пользователя: %v", err)
        return &pb.UserResponse{
            Success: false,
            Message: "User not found",
        }, status.Errorf(codes.NotFound, "user not found: %v", err)
    }
    return &pb.UserResponse{
        User: &pb.User{
            UserId:   user.ID,
            Username: user.Username,
            Email:    user.Email,
        },
        Success: true,
        Message: "User found",
    }, nil
}

func (h *UserHandler) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UserResponse, error) {
    err := h.Service.UpdateUser(req.UserId, req.Username, req.Email)
    if err != nil {
        log.Printf("Ошибка при обновлении пользователя: %v", err)
        return &pb.UserResponse{
            Success: false,
            Message: "Update failed",
        }, status.Errorf(codes.Internal, "update failed: %v", err)
    }
    user, _ := h.Service.GetUser(req.UserId)
    return &pb.UserResponse{
        User: &pb.User{
            UserId:   user.ID,
            Username: user.Username,
            Email:    user.Email,
        },
        Success: true,
        Message: "User updated",
    }, nil
}

func (h *UserHandler) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*pb.DeleteUserResponse, error) {
    err := h.Service.DeleteUser(req.UserId)
    if err != nil {
        log.Printf("Ошибка при удалении пользователя: %v", err)
        return &pb.DeleteUserResponse{
            Success: false,
            Message: "Delete failed",
        }, status.Errorf(codes.Internal, "delete failed: %v", err)
    }
    return &pb.DeleteUserResponse{
        Success: true,
        Message: "User deleted",
    }, nil
}