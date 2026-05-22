package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Tencent/WeKnora/internal/config"
	"github.com/Tencent/WeKnora/internal/types"
)

type registerUserRepo struct {
	byEmail    map[string]*types.User
	byUsername map[string]*types.User
	created    *types.User
	updated    *types.User
}

func newRegisterUserRepo() *registerUserRepo {
	return &registerUserRepo{
		byEmail:    map[string]*types.User{},
		byUsername: map[string]*types.User{},
	}
}

func (r *registerUserRepo) CreateUser(_ context.Context, user *types.User) error {
	cp := *user
	r.created = &cp
	r.byEmail[user.Email] = &cp
	r.byUsername[user.Username] = &cp
	return nil
}

func (r *registerUserRepo) GetUserByID(_ context.Context, id string) (*types.User, error) {
	if r.created != nil && r.created.ID == id {
		cp := *r.created
		return &cp, nil
	}
	return nil, errors.New("not found")
}

func (r *registerUserRepo) GetUsersByIDs(context.Context, []string) (map[string]*types.User, error) {
	return map[string]*types.User{}, nil
}

func (r *registerUserRepo) GetUserByEmail(_ context.Context, email string) (*types.User, error) {
	if u := r.byEmail[email]; u != nil {
		cp := *u
		return &cp, nil
	}
	return nil, errors.New("not found")
}

func (r *registerUserRepo) GetUserByUsername(_ context.Context, username string) (*types.User, error) {
	if u := r.byUsername[username]; u != nil {
		cp := *u
		return &cp, nil
	}
	return nil, errors.New("not found")
}

func (r *registerUserRepo) GetUserByTenantID(context.Context, uint64) (*types.User, error) {
	return nil, errors.New("not found")
}

func (r *registerUserRepo) UpdateUser(_ context.Context, user *types.User) error {
	cp := *user
	r.updated = &cp
	r.byEmail[user.Email] = &cp
	r.byUsername[user.Username] = &cp
	return nil
}

func (r *registerUserRepo) DeleteUser(context.Context, string) error { return nil }
func (r *registerUserRepo) ListUsers(context.Context, int, int) ([]*types.User, error) {
	return nil, nil
}
func (r *registerUserRepo) SearchUsers(context.Context, string, int) ([]*types.User, error) {
	return nil, nil
}

type registerTokenRepo struct{}

func (registerTokenRepo) CreateToken(context.Context, *types.AuthToken) error { return nil }
func (registerTokenRepo) GetTokenByValue(context.Context, string) (*types.AuthToken, error) {
	return nil, errors.New("not found")
}
func (registerTokenRepo) GetTokensByUserID(context.Context, string) ([]*types.AuthToken, error) {
	return nil, nil
}
func (registerTokenRepo) UpdateToken(context.Context, *types.AuthToken) error { return nil }
func (registerTokenRepo) DeleteToken(context.Context, string) error           { return nil }
func (registerTokenRepo) DeleteExpiredTokens(context.Context) error           { return nil }
func (registerTokenRepo) RevokeTokensByUserID(context.Context, string) error  { return nil }

type registerTenantService struct {
	nextID uint64
}

func (s *registerTenantService) CreateTenant(_ context.Context, tenant *types.Tenant) (*types.Tenant, error) {
	if s.nextID == 0 {
		s.nextID = 1
	}
	cp := *tenant
	cp.ID = s.nextID
	cp.CreatedAt = time.Now()
	cp.UpdatedAt = cp.CreatedAt
	s.nextID++
	return &cp, nil
}

func (s *registerTenantService) GetTenantByID(_ context.Context, id uint64) (*types.Tenant, error) {
	return &types.Tenant{ID: id, Name: "Shared"}, nil
}

func (s *registerTenantService) GetTenantsByIDs(context.Context, []uint64) (map[uint64]*types.Tenant, error) {
	return map[uint64]*types.Tenant{}, nil
}
func (s *registerTenantService) ListTenants(context.Context) ([]*types.Tenant, error) {
	return nil, nil
}
func (s *registerTenantService) UpdateTenant(context.Context, *types.Tenant) (*types.Tenant, error) {
	return nil, nil
}
func (s *registerTenantService) DeleteTenant(context.Context, uint64) error { return nil }
func (s *registerTenantService) UpdateAPIKey(context.Context, uint64) (string, error) {
	return "", nil
}
func (s *registerTenantService) ExtractTenantIDFromAPIKey(string) (uint64, error) { return 0, nil }
func (s *registerTenantService) ListAllTenants(context.Context) ([]*types.Tenant, error) {
	return nil, nil
}
func (s *registerTenantService) SearchTenants(context.Context, string, uint64, int, int) ([]*types.Tenant, int64, error) {
	return nil, 0, nil
}
func (s *registerTenantService) GetTenantByIDForUser(context.Context, uint64, string) (*types.Tenant, error) {
	return nil, nil
}
func (s *registerTenantService) GetWeKnoraCloudCredentials(context.Context) *types.WeKnoraCloudCredentials {
	return nil
}

type registerMemberService struct {
	added []types.TenantMember
}

func (s *registerMemberService) AddMember(_ context.Context, userID string, tenantID uint64, role types.TenantRole, invitedBy *string) (*types.TenantMember, error) {
	member := types.TenantMember{UserID: userID, TenantID: tenantID, Role: role, Status: types.TenantMemberStatusActive}
	s.added = append(s.added, member)
	return &member, nil
}

func (s *registerMemberService) EnsureOwner(ctx context.Context, userID string, tenantID uint64) (*types.TenantMember, error) {
	return s.AddMember(ctx, userID, tenantID, types.TenantRoleOwner, nil)
}

func (s *registerMemberService) GetMembership(context.Context, string, uint64) (*types.TenantMember, error) {
	return nil, nil
}
func (s *registerMemberService) ListByUser(context.Context, string) ([]*types.TenantMember, error) {
	return nil, nil
}
func (s *registerMemberService) ListByTenant(context.Context, uint64) ([]*types.TenantMember, error) {
	return nil, nil
}
func (s *registerMemberService) ListMembersPage(context.Context, uint64, string, int, int) ([]*types.TenantMember, int64, error) {
	return nil, 0, nil
}
func (s *registerMemberService) HasAnyMembers(context.Context, uint64) (bool, error) {
	return false, nil
}
func (s *registerMemberService) UpdateRole(context.Context, string, uint64, types.TenantRole) error {
	return nil
}
func (s *registerMemberService) RemoveMember(context.Context, string, uint64) error { return nil }

func TestRegisterAutoJoinAddsViewerMembershipAndActivePreference(t *testing.T) {
	userRepo := newRegisterUserRepo()
	memberSvc := &registerMemberService{}
	svc := NewUserService(
		&config.Config{Auth: &config.AuthConfig{
			AutoJoinTenantID:       99,
			AutoJoinRole:           string(types.TenantRoleViewer),
			AutoJoinAsActiveTenant: true,
		}},
		userRepo,
		registerTokenRepo{},
		&registerTenantService{nextID: 1},
		memberSvc,
	)

	user, err := svc.Register(context.Background(), &types.RegisterRequest{
		Username: "alice",
		Email:    "alice@example.com",
		Password: "password123",
	})
	if err != nil {
		t.Fatalf("Register returned error: %v", err)
	}
	if user == nil {
		t.Fatalf("Register returned nil user")
	}
	if len(memberSvc.added) != 2 {
		t.Fatalf("members added = %d, want 2 (owner + auto-join viewer)", len(memberSvc.added))
	}
	auto := memberSvc.added[1]
	if auto.TenantID != 99 || auto.Role != types.TenantRoleViewer || auto.UserID != user.ID {
		t.Fatalf("auto-join member = %+v, want user=%s tenant=99 role=viewer", auto, user.ID)
	}
	if userRepo.updated == nil || userRepo.updated.Preferences.LastActiveTenantID == nil ||
		*userRepo.updated.Preferences.LastActiveTenantID != 99 {
		t.Fatalf("LastActiveTenantID was not persisted to 99")
	}
}
