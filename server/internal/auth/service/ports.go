package service

import (
	"context"

	authdb "github.com/findardi/Wadi/server/internal/auth/repository/sqlc"
	"github.com/findardi/Wadi/server/internal/platform/token"
	"github.com/jackc/pgx/v5/pgtype"
)

type UserRepository interface {
	GetUserByEmail(ctx context.Context, email string) (authdb.User, error)
	GetUserByUsername(ctx context.Context, username *string) (authdb.User, error)
	GetValidUserToken(ctx context.Context, arg authdb.GetValidUserTokenParams) (authdb.UserToken, error)

	CreateUser(ctx context.Context, arg authdb.CreateUserParams) (authdb.User, error)
	CreateUserToken(ctx context.Context, arg authdb.CreateUserTokenParams) (authdb.UserToken, error)

	UpdateStatus(ctx context.Context, arg authdb.UpdateStatusParams) (authdb.User, error)
	UpdateUser(ctx context.Context, arg authdb.UpdateUserParams) (authdb.User, error)
	UpdateUserToken(ctx context.Context, arg authdb.UpdateUserTokenParams) (authdb.UserToken, error)

	DeleteUserToken(ctx context.Context, arg authdb.DeleteUserTokenParams) error
	DeleteExpiredUserTokens(ctx context.Context, userID pgtype.UUID) error

	ExecTx(ctx context.Context, fn func(q *authdb.Queries) error) error
}

type OTPService interface {
	Generate() string
	Hash(code string) string
	Compare(hash, code string) bool
}

type JWTService interface {
	CreateToken(claims token.JwtClaims, tokenType token.TokenType) (string, error)
	VerifyToken(tokenString string) (*token.JwtClaims, error)
}
