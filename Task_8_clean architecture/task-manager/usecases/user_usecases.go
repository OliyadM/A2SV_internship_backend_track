package usecases

import (
	"context"
	"errors"
	"task_manager/domain"
	"github.com/google/uuid"
)

type UserUsecaseImpl struct {
	userRepo    domain.UserRepository
	passwordSvc domain.PasswordService
	jwtSvc      domain.JWTService
}

func NewUserUsecase(userRepo domain.UserRepository, passwordSvc domain.PasswordService, jwtSvc domain.JWTService) domain.UserUsecase {
	return &UserUsecaseImpl{
		userRepo:    userRepo,
		passwordSvc: passwordSvc,
		jwtSvc:      jwtSvc,
	}
}

func (u *UserUsecaseImpl) Register(ctx context.Context, user domain.User) error {
	if existing, _ := u.userRepo.FindUserByUsername(ctx, user.Username); existing != nil {
		return errors.New("username already exists")
	}

	
	user.ID = uuid.New().String()

	hashed, err := u.passwordSvc.HashPassword(user.Password)
	if err != nil {
		return err
	}
	user.Password = hashed

	isFirst, err := u.userRepo.IsFirstUser(ctx)
	if err != nil {
		return err
	}
	if isFirst {
		user.Role = "admin"
	} else {
		user.Role = "user"
	}

	return u.userRepo.CreateUser(ctx, user)
}

func (u *UserUsecaseImpl) Login(ctx context.Context, username, password string) (string, error) {
	user, err := u.userRepo.FindUserByUsername(ctx, username)
	if err != nil || user == nil {
		return "", errors.New("invalid credentials")
	}

	if err := u.passwordSvc.ComparePassword(user.Password, password); err != nil {
		return "", errors.New("invalid credentials")
	}

	token, err := u.jwtSvc.GenerateToken(user.ID, user.Username, user.Role)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (u *UserUsecaseImpl) PromoteUser(ctx context.Context, username string) error {
	user, err := u.userRepo.FindUserByUsername(ctx, username)
	if err != nil || user == nil {
		return errors.New("user not found")
	}
	return u.userRepo.PromoteUser(ctx, username)
}


func (u *UserUsecaseImpl) GetAllUsers(ctx context.Context) ([]*domain.User, error) {
	users, err := u.userRepo.GetAllUsers(ctx)
	if err != nil || len(users) == 0 {
		return nil, errors.New("no users found")
	}
	return users, nil
}


