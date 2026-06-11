package service

import (
	"context"

	authdb "github.com/findardi/Wadi/server/internal/auth/repository/sqlc"
)

type UserRepository interface {
	GetUserByEmail(ctx context.Context, email string) (authdb.User, error)
	GetUserByUsername(ctx context.Context, username *string) (authdb.User, error)
	CreateUser(ctx context.Context, arg authdb.CreateUserParams) (authdb.User, error)
	CreateUserToken(ctx context.Context, arg authdb.CreateUserTokenParams) (authdb.UserToken, error)
	ExecTx(ctx context.Context, fn func(q *authdb.Queries) error) error
}

type OTPService interface {
	Generate() string
	Hash(code string) string
	Compare(hash, code string) bool
}
