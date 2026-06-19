package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/findardi/Wadi/server/internal/access/dto"
	accessdb "github.com/findardi/Wadi/server/internal/access/repository/sqlc"
	"github.com/findardi/Wadi/server/internal/platform/permission"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
)

const (
	MemberStatusInvited   = "invited"
	MemberStatusActive    = "active"
	MemberStatusSuspended = "suspended"
)

var (
	ErrRoleNameTaken    = errors.New("role name already taken")
	ErrRoleNotFound     = errors.New("role not found")
	ErrRoleInUse        = errors.New("role is still assigned to members")
	ErrMemberAlreadyAdd = errors.New("user already a member of this workspace")
	ErrMemberNotFound   = errors.New("member not found")
)

type AccessService struct {
	repo AccessRepository
	mail MailService
}

func NewAccessService(repo AccessRepository, mail MailService) *AccessService {
	return &AccessService{
		repo: repo,
		mail: mail,
	}
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

func isUniqueViolation(err error, constraint string) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code == "23505" && pgErr.ConstraintName == constraint
	}
	return false
}

func isForeignKeyViolation(err error) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code == "23503"
	}
	return false
}

func (s *AccessService) InsertRole(ctx context.Context, req dto.CreateWorkspaceRoleRequest) (dto.WorkspaceRoleResponse, error) {
	var uid pgtype.UUID
	if err := uid.Scan(req.WorkspaceID); err != nil {
		return dto.WorkspaceRoleResponse{}, fmt.Errorf("parse owner id: %w", err)
	}

	for _, r := range req.Permission {
		if ok := permission.IsValid(r); !ok {
			return dto.WorkspaceRoleResponse{}, fmt.Errorf("invalid permission: %s", r)
		}
	}

	role, err := s.repo.InsertRole(ctx, accessdb.InsertRoleParams{
		WorkspaceID: uid,
		Name:        req.Name,
		Permissions: req.Permission,
		IsSystem:    req.IsSystem,
	})

	if isUniqueViolation(err, "workspace_roles_name_key") {
		return dto.WorkspaceRoleResponse{}, ErrRoleNameTaken
	}
	if err != nil {
		return dto.WorkspaceRoleResponse{}, fmt.Errorf("insert role: %w", err)
	}

	return dto.WorkspaceRoleResponse{
		ID:          uuidString(role.ID),
		WorkspaceID: uuidString(role.WorkspaceID),
		Name:        role.Name,
		Permissions: role.Permissions,
		IsSystem:    role.IsSystem,
		CreatedAt:   role.CreatedAt.Time,
		UpdatedAt:   role.UpdatedAt.Time,
	}, nil
}

func (s *AccessService) ProvisionWorkspace(ctx context.Context, tx pgx.Tx, workspaceID, ownerID pgtype.UUID) error {
	q := accessdb.New(tx)

	var ownerRoleID pgtype.UUID
	for _, r := range permission.DefaultSystemRoles() {
		role, err := q.InsertRole(ctx, accessdb.InsertRoleParams{
			WorkspaceID: workspaceID,
			Name:        r.Name,
			Permissions: r.Permissions,
			IsSystem:    true,
		})
		if err != nil {
			return fmt.Errorf("seed role %s: %w", r.Name, err)
		}
		if r.Name == permission.RoleOwner {
			ownerRoleID = role.ID
		}
	}

	if _, err := q.AddMember(ctx, accessdb.AddMemberParams{
		WorkspaceID: workspaceID,
		UserID:      ownerID,
		RoleID:      ownerRoleID,
		Status:      MemberStatusActive,
	}); err != nil {
		return fmt.Errorf("add owner member: %w", err)
	}

	return nil
}

func (s *AccessService) AddMember(ctx context.Context, req dto.CreateWorkspaceMemberRequest) (dto.WorkspaceMemberResponse, error) {
	var wsID, userID, roleID pgtype.UUID
	if err := wsID.Scan(req.WorkspaceId); err != nil {
		return dto.WorkspaceMemberResponse{}, fmt.Errorf("parse workspace id: %w", err)
	}
	if err := userID.Scan(req.UserId); err != nil {
		return dto.WorkspaceMemberResponse{}, fmt.Errorf("parse user id: %w", err)
	}
	if err := roleID.Scan(req.RoleId); err != nil {
		return dto.WorkspaceMemberResponse{}, fmt.Errorf("parse role id: %w", err)
	}

	status := req.Status
	if status == "" {
		status = MemberStatusActive
	}

	member, err := s.repo.AddMember(ctx, accessdb.AddMemberParams{
		WorkspaceID: wsID,
		UserID:      userID,
		RoleID:      roleID,
		Status:      status,
	})
	if isUniqueViolation(err, "workspace_members_user_key") {
		return dto.WorkspaceMemberResponse{}, ErrMemberAlreadyAdd
	}
	if err != nil {
		return dto.WorkspaceMemberResponse{}, fmt.Errorf("add member: %w", err)
	}

	return dto.WorkspaceMemberResponse{
		ID:          uuidString(member.ID),
		WorkspaceID: uuidString(member.WorkspaceID),
		UserID:      uuidString(member.UserID),
		RoleID:      uuidString(member.RoleID),
		Status:      member.Status,
		CreatedAt:   member.CreatedAt.Time,
		UpdatedAt:   member.UpdatedAt.Time,
	}, nil
}

func (s *AccessService) GetRoles(ctx context.Context, workspaceID string) ([]dto.WorkspaceRoleResponse, error) {
	var roles []dto.WorkspaceRoleResponse

	var wsID pgtype.UUID
	if err := wsID.Scan(workspaceID); err != nil {
		return roles, fmt.Errorf("parse workspace id: %w", err)
	}

	wsRoles, err := s.repo.GetRoles(ctx, wsID)
	if err != nil {
		return roles, fmt.Errorf("get roles: %w", err)
	}

	for _, r := range wsRoles {
		role := dto.WorkspaceRoleResponse{
			ID:          uuidString(r.ID),
			WorkspaceID: uuidString(r.WorkspaceID),
			Name:        r.Name,
			Permissions: r.Permissions,
			IsSystem:    r.IsSystem,
			CreatedAt:   r.CreatedAt.Time,
			UpdatedAt:   r.UpdatedAt.Time,
		}

		roles = append(roles, role)
	}

	return roles, nil
}

func (s *AccessService) GetRole(ctx context.Context, roleId string) (dto.WorkspaceRoleResponse, error) {
	var roleID pgtype.UUID
	if err := roleID.Scan(roleId); err != nil {
		return dto.WorkspaceRoleResponse{}, fmt.Errorf("parse role id: %w", err)
	}

	role, err := s.repo.GetRole(ctx, roleID)
	if errors.Is(err, pgx.ErrNoRows) {
		return dto.WorkspaceRoleResponse{}, ErrRoleNotFound
	}
	if err != nil {
		return dto.WorkspaceRoleResponse{}, fmt.Errorf("get role: %w", err)
	}

	return dto.WorkspaceRoleResponse{
		ID:          uuidString(role.ID),
		WorkspaceID: uuidString(role.WorkspaceID),
		Name:        role.Name,
		Permissions: role.Permissions,
		IsSystem:    role.IsSystem,
		CreatedAt:   role.CreatedAt.Time,
		UpdatedAt:   role.UpdatedAt.Time,
	}, nil
}

func (s *AccessService) UpdateRole(ctx context.Context, req dto.UpdateWorkspaceRoleRequest) (dto.WorkspaceRoleResponse, error) {
	var roleID pgtype.UUID
	if err := roleID.Scan(req.RoleID); err != nil {
		return dto.WorkspaceRoleResponse{}, fmt.Errorf("parse role id: %w", err)
	}

	role, err := s.repo.EditRole(ctx, accessdb.EditRoleParams{
		ID:          roleID,
		Name:        req.Name,
		Permissions: req.Permission,
	})

	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return dto.WorkspaceRoleResponse{}, ErrRoleNotFound
	case isUniqueViolation(err, "workspace_roles_name_key"):
		return dto.WorkspaceRoleResponse{}, ErrRoleNameTaken
	case err != nil:
		return dto.WorkspaceRoleResponse{}, fmt.Errorf("update role: %w", err)
	}

	return dto.WorkspaceRoleResponse{
		ID:          uuidString(role.ID),
		WorkspaceID: uuidString(role.WorkspaceID),
		Name:        role.Name,
		Permissions: role.Permissions,
		IsSystem:    role.IsSystem,
		CreatedAt:   role.CreatedAt.Time,
		UpdatedAt:   role.UpdatedAt.Time,
	}, nil
}

func (s *AccessService) DeleteRole(ctx context.Context, roleId string) error {
	var roleID pgtype.UUID
	if err := roleID.Scan(roleId); err != nil {
		return fmt.Errorf("parse role id: %w", err)
	}

	err := s.repo.DeleteRole(ctx, roleID)
	if isForeignKeyViolation(err) {
		return ErrRoleInUse
	}
	if err != nil {
		return fmt.Errorf("delete role: %w", err)
	}

	return nil
}

func (s *AccessService) GetMembers(ctx context.Context, workspaceID string) ([]dto.GetMemberResponse, error) {
	var members []dto.GetMemberResponse
	var wsID pgtype.UUID
	if err := wsID.Scan(workspaceID); err != nil {
		return members, fmt.Errorf("parse workspace id: %w", err)
	}

	wsMembers, err := s.repo.GetMembers(ctx, wsID)
	if err != nil {
		return members, fmt.Errorf("get members: %w", err)
	}

	for _, w := range wsMembers {
		member := dto.GetMemberResponse{
			ID:          uuidString(w.ID),
			WorkspaceID: uuidString(w.WorkspaceID),
			UserID:      uuidString(w.UserID),
			RoleID:      uuidString(w.RoleID),
			Status:      w.Status,
			CreatedAt:   w.CreatedAt.Time,
			UpdatedAt:   w.UpdatedAt.Time,
			RoleName:    deref(w.RoleName),
			Username:    deref(w.Username),
			Email:       deref(w.Email),
			GroupNames:  w.GroupNames,
		}

		members = append(members, member)
	}

	return members, nil
}

func (s *AccessService) GetMember(ctx context.Context, memberID string) (dto.GetMemberResponse, error) {
	var mID pgtype.UUID
	if err := mID.Scan(memberID); err != nil {
		return dto.GetMemberResponse{}, fmt.Errorf("parse member id: %w", err)
	}

	w, err := s.repo.GetMember(ctx, mID)
	if errors.Is(err, pgx.ErrNoRows) {
		return dto.GetMemberResponse{}, ErrMemberNotFound
	}
	if err != nil {
		return dto.GetMemberResponse{}, fmt.Errorf("get member: %w", err)
	}

	return dto.GetMemberResponse{
		ID:          uuidString(w.ID),
		WorkspaceID: uuidString(w.WorkspaceID),
		UserID:      uuidString(w.UserID),
		RoleID:      uuidString(w.RoleID),
		Status:      w.Status,
		CreatedAt:   w.CreatedAt.Time,
		UpdatedAt:   w.UpdatedAt.Time,
		RoleName:    deref(w.RoleName),
		Username:    deref(w.Username),
		Email:       deref(w.Email),
		GroupNames:  w.GroupNames,
	}, nil
}

func (s *AccessService) UpdateMemberRole(ctx context.Context, req dto.UpdateMemberRoleRequest) (dto.GetMemberResponse, error) {
	var mID, rID pgtype.UUID
	if err := mID.Scan(req.MemberID); err != nil {
		return dto.GetMemberResponse{}, fmt.Errorf("parse member id: %w", err)
	}
	if err := rID.Scan(req.RoleId); err != nil {
		return dto.GetMemberResponse{}, fmt.Errorf("parse role id: %w", err)
	}

	_, err := s.repo.UpdateRole(ctx, accessdb.UpdateRoleParams{
		ID:     mID,
		RoleID: rID,
	})
	if errors.Is(err, pgx.ErrNoRows) {
		return dto.GetMemberResponse{}, ErrMemberNotFound
	}
	if err != nil {
		return dto.GetMemberResponse{}, fmt.Errorf("update member role: %w", err)
	}

	return s.GetMember(ctx, req.MemberID)
}

func (s *AccessService) DeleteMember(ctx context.Context, memberID string) error {
	var mID pgtype.UUID
	if err := mID.Scan(memberID); err != nil {
		return fmt.Errorf("parse member id: %w", err)
	}

	if err := s.repo.DeleteMember(ctx, mID); err != nil {
		return fmt.Errorf("delete member: %w", err)
	}

	return nil
}
