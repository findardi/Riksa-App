package service

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/findardi/Wadi/server/internal/auth/dto"
	authdb "github.com/findardi/Wadi/server/internal/auth/repository/sqlc"
	"github.com/findardi/Wadi/server/internal/platform/token"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrEmailUnique          = errors.New("email already registered")
	ErrUsernameUnique       = errors.New("username already registered")
	ErrInvalidCredentials   = errors.New("invalid email/username or password")
	ErrInvalidCodeOTP       = errors.New("invalid or expired OTP code")
	ErrEmailAlreadyVerified = errors.New("email already verified")
)

type AuthService struct {
	repo UserRepository
	otp  OTPService
	jwt  JWTService
}

func NewAuthService(repo UserRepository, otp OTPService, jwt JWTService) *AuthService {
	return &AuthService{
		repo: repo,
		otp:  otp,
		jwt:  jwt,
	}
}

func (s *AuthService) RegisterUser(ctx context.Context, req dto.RegisterRequest) (dto.RegisterResponse, error) {
	email := strings.ToLower(strings.TrimSpace(req.Email))
	username := strings.TrimSpace(req.Username)

	if _, err := s.repo.GetUserByEmail(ctx, email); err == nil {
		return dto.RegisterResponse{}, ErrEmailUnique
	} else if !errors.Is(err, pgx.ErrNoRows) {
		return dto.RegisterResponse{}, fmt.Errorf("check email: %w", err)
	}

	if _, err := s.repo.GetUserByUsername(ctx, &username); err == nil {
		return dto.RegisterResponse{}, ErrUsernameUnique
	} else if !errors.Is(err, pgx.ErrNoRows) {
		return dto.RegisterResponse{}, fmt.Errorf("check username: %w", err)
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), 12)
	if err != nil {
		return dto.RegisterResponse{}, fmt.Errorf("hash password: %w", err)
	}

	hashPassword := string(hash)
	code := s.otp.Generate()

	// Transaction process
	var user authdb.User
	err = s.repo.ExecTx(ctx, func(q *authdb.Queries) error {
		u, err := q.CreateUser(ctx, authdb.CreateUserParams{
			Email:        email,
			Username:     &username,
			PasswordHash: &hashPassword,
			Status:       "pending",
		})
		if err != nil {
			return fmt.Errorf("create user: %w", err)
		}

		if _, err := q.CreateUserToken(ctx, authdb.CreateUserTokenParams{
			UserID:   u.ID,
			Type:     "email_verification",
			CodeHash: s.otp.Hash(code),
			ExpiresAt: pgtype.Timestamptz{
				Time:  time.Now().Add(5 * time.Minute),
				Valid: true,
			},
		}); err != nil {
			return fmt.Errorf("create otp: %w", err)
		}

		user = u
		return nil
	})
	if err != nil {
		return dto.RegisterResponse{}, fmt.Errorf("register tx: %w", err)
	}

	// TODO: kirim OTP ke email SETELAH commit (di luar transaksi).
	// if err := s.mailer.SendVerificationOTP(ctx, email, code); err != nil { ... }
	// TODO: delete log when email already implementasi
	fmt.Println("Code:", code)

	return dto.RegisterResponse{
		ID:       uuidString(user.ID),
		Username: deref(user.Username),
	}, nil
}

func (s *AuthService) LoginUser(ctx context.Context, req dto.LoginRequest) (dto.LoginResponse, error) {
	var user authdb.User
	var err error

	switch {
	case req.Email != "":
		email := strings.ToLower(strings.TrimSpace(req.Email))
		user, err = s.repo.GetUserByEmail(ctx, email)
	case req.Username != "":
		username := strings.TrimSpace(req.Username)
		user, err = s.repo.GetUserByUsername(ctx, &username)
	default:
		return dto.LoginResponse{}, ErrInvalidCredentials
	}

	if errors.Is(err, pgx.ErrNoRows) {
		return dto.LoginResponse{}, ErrInvalidCredentials
	}
	if err != nil {
		return dto.LoginResponse{}, fmt.Errorf("get user: %w", err)
	}

	// User SSO without password
	if user.PasswordHash == nil {
		return dto.LoginResponse{}, ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(*user.PasswordHash), []byte(req.Password)); err != nil {
		return dto.LoginResponse{}, ErrInvalidCredentials
	}

	claims := token.JwtClaims{
		ID:       uuidString(user.ID),
		Username: deref(user.Username),
		Email:    user.Email,
		Status:   user.Status,
	}

	accessToken, err := s.jwt.CreateToken(claims, token.TokenType("token_login"))
	if err != nil {
		return dto.LoginResponse{}, fmt.Errorf("create access token: %w", err)
	}

	// delete unused token
	if err := s.repo.DeleteExpiredUserTokens(ctx, user.ID); err != nil {
		return dto.LoginResponse{}, fmt.Errorf("cleanup expired tokens: %w", err)
	}

	refreshToken := s.otp.Generate()

	if _, err = s.repo.CreateUserToken(ctx, authdb.CreateUserTokenParams{
		UserID:   user.ID,
		Type:     "refresh",
		CodeHash: s.otp.Hash(refreshToken),
		ExpiresAt: pgtype.Timestamptz{
			Time:  time.Now().Add(24 * time.Hour),
			Valid: true,
		},
	}); err != nil {
		return dto.LoginResponse{}, fmt.Errorf("create refresh token: %w", err)
	}

	return dto.LoginResponse{
		Token:        accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *AuthService) ResendOTP(ctx context.Context, email string) error {
	email = strings.ToLower(strings.TrimSpace(email))

	user, err := s.repo.GetUserByEmail(ctx, email)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("get user: %w", err)
	}

	if user.Status != "pending" {
		return nil
	}

	code := s.otp.Generate()
	if _, err := s.repo.UpdateUserToken(ctx, authdb.UpdateUserTokenParams{
		UserID:   user.ID,
		Type:     "email_verification",
		CodeHash: s.otp.Hash(code),
		ExpiresAt: pgtype.Timestamptz{
			Time:  time.Now().Add(5 * time.Minute),
			Valid: true,
		},
	}); err != nil {
		return fmt.Errorf("update otp: %w", err)
	}

	// TODO: kirim OTP ke email. TODO: hapus log saat mailer siap.
	fmt.Println("Code:", code)

	return nil
}

func (s *AuthService) VerifyEmail(ctx context.Context, req dto.VerifyOtpRequest) error {
	email := strings.ToLower(strings.TrimSpace(req.Email))

	user, err := s.repo.GetUserByEmail(ctx, email)
	if errors.Is(err, pgx.ErrNoRows) {
		return ErrInvalidCodeOTP
	}
	if err != nil {
		return fmt.Errorf("get user: %w", err)
	}

	if user.Status != "pending" {
		return ErrEmailAlreadyVerified
	}

	otp, err := s.repo.GetValidUserToken(ctx, authdb.GetValidUserTokenParams{
		UserID: user.ID,
		Type:   "email_verification",
	})
	if errors.Is(err, pgx.ErrNoRows) {
		return ErrInvalidCodeOTP
	}
	if err != nil {
		return fmt.Errorf("get otp: %w", err)
	}

	if !s.otp.Compare(otp.CodeHash, req.Code) {
		return ErrInvalidCodeOTP
	}

	err = s.repo.ExecTx(ctx, func(q *authdb.Queries) error {
		if _, err := q.UpdateStatus(ctx, authdb.UpdateStatusParams{
			ID:     user.ID,
			Status: "active",
		}); err != nil {
			return fmt.Errorf("update status: %w", err)
		}

		if err := q.DeleteUserToken(ctx, authdb.DeleteUserTokenParams{
			CodeHash: otp.CodeHash,
			UserID:   user.ID,
			Type:     "email_verification",
		}); err != nil {
			return fmt.Errorf("delete token: %w", err)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("verify tx: %w", err)
	}

	return nil
}

func (s *AuthService) LogoutUser(ctx context.Context, req dto.LogoutRequest) error {
	var uid pgtype.UUID
	if err := uid.Scan(req.UserID); err != nil {
		return fmt.Errorf("parse user id: %w", err)
	}

	if err := s.repo.DeleteUserToken(ctx, authdb.DeleteUserTokenParams{
		CodeHash: s.otp.Hash(req.RefreshToken),
		UserID:   uid,
		Type:     "refresh",
	}); err != nil {
		return fmt.Errorf("delete refresh token: %w", err)
	}

	return nil
}

func uuidString(u pgtype.UUID) string {
	v, err := u.Value()
	if err != nil || v == nil {
		return ""
	}
	s, _ := v.(string)
	return s
}

func deref(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
