package repositoryimpl

import (
	"context"

	"go-jwt/ent"
	"go-jwt/ent/user"
	"go-jwt/pkg/domain/repository"
	"go-jwt/pkg/infrastructure/database"
)

type userRepo struct {
	dbClient *database.DBClient
}

func NewUserRepo(dbClient *database.DBClient) repository.UserRepository {
	return &userRepo{
		dbClient: dbClient,
	}
}

// CreateUser implements repository.UserRepository.
func (r *userRepo) CreateUser(ctx context.Context, user *ent.User) (*ent.User, error) {
	newUser, err := r.dbClient.Client.User.Create().
		SetUserName(user.UserName).
		SetEmail(user.Email).
		SetPassword(user.Password).
		Save(context.Background())
	if err != nil {
		return nil, err
	}
	return newUser, nil
}

// GetUserByEmail implements repository.UserRepository.
func (r *userRepo) GetUserByEmail(ctx context.Context, email string) (*ent.User, error) {
	user, err := r.dbClient.Client.User.Query().
		Where(user.Email(email)).
		Only(context.Background())
	if err != nil {
		return nil, err
	}
	return user, nil
}
