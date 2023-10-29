package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/FarStep131/go-jwt/domain/repository"
	"github.com/FarStep131/go-jwt/ent"
	"github.com/FarStep131/go-jwt/myerror"
	"github.com/FarStep131/go-jwt/util"
)

type UseCase interface {
	Singup(c context.Context, username, email, password string) (*ent.User, error)
	Login(c context.Context, email, password string) (string, *ent.User, error)
}

type useCase struct {
	repository repository.UserRepository
	timeout    time.Duration
}

func NewUseCase(userRepo repository.UserRepository) UseCase {
	return &useCase{
		repository: userRepo,
		timeout:    time.Duration(2) * time.Second,
	}
}

// Singup implements UseCase.
func (uc *useCase) Singup(c context.Context, username string, email string, password string) (*ent.User, error) {
	ctx, cancel := context.WithTimeout(c, uc.timeout)
	defer cancel()

	exsiteUser, err := uc.repository.GetUserByEmail(ctx, email)

	if err != nil {
		return nil, &myerror.BadRequestError{Err: err}
	}

	if exsiteUser.ID != 0 {
		return nil, &myerror.BadRequestError{Err: errors.New("user already exists")}
	}

	hashedPassword, err := util.HashPassword(password)
	if err != nil {
		return nil, &myerror.InternalSeverError{Err: err}
	}

	u := &ent.User{
		UserName: username,
		Email:    email,
		Password: hashedPassword,
	}

	user, err := uc.repository.CreateUser(ctx, u)
	if err != nil {
		return nil, &myerror.InternalSeverError{Err: err}
	}

	return user, err
}

// Login implements UseCase.
func (uc *useCase) Login(c context.Context, email string, password string) (string, *ent.User, error) {
	ctx, cancel := context.WithTimeout(c, uc.timeout)
	defer cancel()

	user, err := uc.repository.GetUserByEmail(ctx, email)
	if err != nil {
		return "", nil, &myerror.InternalSeverError{Err: err}
	}
	if user.ID == 0 {
		return "", nil, &myerror.BadRequestError{Err: errors.New("user is not exist")}
	}

	err = util.CheckPasswrod(user.Password, password)
	if err != nil {
		return "", nil, &myerror.BadRequestError{Err: errors.New("password is incorrect")}
	}

	signedString, err := util.GenerateSignedString(int64(user.ID), user.UserName)
	if err != nil {
		return "", nil, &myerror.InternalSeverError{Err: err}
	}

	return signedString, user, nil
}
