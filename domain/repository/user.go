package repository

import (
	"context"

	"github.com/FarStep131/go-jwt/ent"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *ent.User) (*ent.User, error)
	GetUserByEmail(ctx context.Context, email string) (*ent.User, error)
}
