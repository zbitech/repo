package basic

import (
	"context"
	"github.com/golang-jwt/jwt"
	"github.com/zbitech/common/interfaces"
	"github.com/zbitech/common/pkg/errs"
	"github.com/zbitech/common/pkg/model/entity"
	"github.com/zbitech/common/pkg/model/object"
	"github.com/zbitech/common/pkg/model/ztypes"
	"github.com/zbitech/common/pkg/vars"
	"time"
)

type BasicIAMService struct {
	jwtServer interfaces.JwtServerIF
}

func NewBasicIAMService(jwtServer interfaces.JwtServerIF) interfaces.IAMServiceIF {
	return &BasicIAMService{jwtServer: jwtServer}
}

func (b *BasicIAMService) RegisterUser(ctx context.Context, user *entity.User, pass *entity.UserPassword) error {

	adminRepo := vars.RepositoryFactory.GetAdminRepository()
	return adminRepo.RegisterUser(ctx, user, pass)
}

func (b *BasicIAMService) DeactivateUser(ctx context.Context, userid string) error {

	user, err := b.GetUser(ctx, userid)
	if err != nil {
		return err
	}
	user.Active = false
	user.LastUpdate = time.Now()

	return b.UpdateUser(ctx, user)
}

func (b *BasicIAMService) ReactivateUser(ctx context.Context, userid string) error {

	user, err := b.GetUser(ctx, userid)
	if err != nil {
		return err
	}
	user.Active = true
	user.LastUpdate = time.Now()

	return b.UpdateUser(ctx, user)
}

func (b *BasicIAMService) GetUsers(ctx context.Context) []entity.User {

	adminRepo := vars.RepositoryFactory.GetAdminRepository()
	return adminRepo.GetUsers(ctx)
}

func (b *BasicIAMService) GetUser(ctx context.Context, userId string) (*entity.User, error) {

	adminRepo := vars.RepositoryFactory.GetAdminRepository()
	return adminRepo.GetUser(ctx, userId)
}

func (b *BasicIAMService) UpdateUser(ctx context.Context, user *entity.User) error {

	adminRepo := vars.RepositoryFactory.GetAdminRepository()
	return adminRepo.UpdateUser(ctx, user)
}

func (b *BasicIAMService) ChangePassword(ctx context.Context, userid string, pass *entity.UserPassword) error {

	adminRepo := vars.RepositoryFactory.GetAdminRepository()
	return adminRepo.UpdatePassword(ctx, userid, pass)
}

func (b *BasicIAMService) AuthenticateUser(ctx context.Context, userId, password string) (*string, error) {

	adminRepo := vars.RepositoryFactory.GetAdminRepository()
	return adminRepo.AuthenticateUser(ctx, userId, password)
}

func (b *BasicIAMService) ValidateAuthToken(ctx context.Context, tokenString string) (jwt.Claims, *entity.User, error) {

	token, err := jwt.ParseWithClaims(tokenString, b.jwtServer.GetPayload(), func(token *jwt.Token) (interface{}, error) {
		err := b.jwtServer.ValidateToken(token)
		if err != nil {
			return nil, err
		}

		return b.jwtServer.GetKey()
	})

	if err != nil {
		return nil, nil, err
	}

	if !token.Valid {
		return nil, nil, errs.ErrInvalidToken
	}

	claims := token.Claims.(*object.ZBIBasicClaims)
	user, err := b.GetUser(ctx, claims.Subject)
	if err != nil {
		return nil, nil, errs.ErrUnregisteredUser
	}

	return claims, user, nil
}

func (b *BasicIAMService) GetAPIKeys(ctx context.Context, userId string) ([]string, error) {

	adminRepo := vars.RepositoryFactory.GetAdminRepository()
	return adminRepo.GetAPIKeys(ctx, userId)
}

func (b *BasicIAMService) GetAPIKey(ctx context.Context, apiKey string) (*entity.APIKey, error) {

	adminRepo := vars.RepositoryFactory.GetAdminRepository()
	return adminRepo.GetAPIKey(ctx, apiKey)
}

func (b *BasicIAMService) CreateAPIKey(ctx context.Context, userid string) (*entity.APIKey, error) {

	adminRepo := vars.RepositoryFactory.GetAdminRepository()
	return adminRepo.CreateAPIKey(ctx, userid)
}

func (b *BasicIAMService) DeleteAPIKey(ctx context.Context, apiKey string) error {

	adminRepo := vars.RepositoryFactory.GetAdminRepository()
	return adminRepo.DeleteAPIKey(ctx, apiKey)
}

func (b *BasicIAMService) StoreUserPolicy(ctx context.Context, p entity.UserPolicy) error {

	adminRepo := vars.RepositoryFactory.GetAdminRepository()
	return adminRepo.StoreUserPolicy(ctx, p)
}

func (b *BasicIAMService) GetUserPolicy(ctx context.Context, userId string) (*entity.UserPolicy, error) {

	adminRepo := vars.RepositoryFactory.GetAdminRepository()
	return adminRepo.GetUserPolicy(ctx, userId)
}

func (b *BasicIAMService) StoreAPIKeyPolicy(ctx context.Context, p entity.APIKeyPolicy) error {

	adminRepo := vars.RepositoryFactory.GetAdminRepository()
	return adminRepo.StoreAPIKeyPolicy(ctx, p)
}

func (b *BasicIAMService) GetAPIKeyPolicy(ctx context.Context, key string) (*entity.APIKeyPolicy, error) {

	adminRepo := vars.RepositoryFactory.GetAdminRepository()
	return adminRepo.GetAPIKeyPolicy(ctx, key)
}

func (b *BasicIAMService) StoreInstancePolicy(ctx context.Context, p entity.InstancePolicy) error {

	adminRepo := vars.RepositoryFactory.GetAdminRepository()
	return adminRepo.StoreInstancePolicy(ctx, p)
}

func (b *BasicIAMService) StoreInstancePolicies(ctx context.Context, p []entity.InstancePolicy) error {

	adminRepo := vars.RepositoryFactory.GetAdminRepository()
	return adminRepo.StoreInstancePolicies(ctx, p)
}

func (b *BasicIAMService) GetInstancePolicy(ctx context.Context, project, instance string) (*entity.InstancePolicy, error) {

	adminRepo := vars.RepositoryFactory.GetAdminRepository()
	return adminRepo.GetInstancePolicy(ctx, project, instance)
}

func (b *BasicIAMService) GetInstanceMethodPolicy(ctx context.Context, project, instance, methodName string) (*entity.MethodPolicy, error) {

	adminRepo := vars.RepositoryFactory.GetAdminRepository()
	return adminRepo.GetInstanceMethodPolicy(ctx, project, instance, methodName)
}

func (b *BasicIAMService) GetInstanceMethodPolicies(ctx context.Context, project, instance, methodCategory string) ([]entity.MethodPolicy, error) {

	adminRepo := vars.RepositoryFactory.GetAdminRepository()
	return adminRepo.GetInstanceMethodPolicies(ctx, project, instance, methodCategory)
}

func (b *BasicIAMService) GetExpiringInvitations(ctx context.Context, date time.Time) ([]entity.TeamMember, error) {

	adminRepo := vars.RepositoryFactory.GetAdminRepository()
	return adminRepo.GetExpiringInvitations(ctx, date)
}

func (b *BasicIAMService) PurgeExpiredInvitations(ctx context.Context) (int64, error) {

	adminRepo := vars.RepositoryFactory.GetAdminRepository()
	return adminRepo.PurgeExpiredInvitations(ctx)
}

func (b *BasicIAMService) CreateTeam(ctx context.Context, team entity.Team) error {

	adminRepo := vars.RepositoryFactory.GetAdminRepository()
	return adminRepo.CreateTeam(ctx, team)
}

func (b *BasicIAMService) GetTeams(ctx context.Context) ([]entity.Team, error) {

	adminRepo := vars.RepositoryFactory.GetAdminRepository()
	return adminRepo.GetTeams(ctx)
}

func (b *BasicIAMService) GetTeam(ctx context.Context, teamId string) (*entity.Team, error) {

	adminRepo := vars.RepositoryFactory.GetAdminRepository()
	return adminRepo.GetTeam(ctx, teamId)
}

func (b *BasicIAMService) UpdateTeam(ctx context.Context, team entity.Team) error {

	adminRepo := vars.RepositoryFactory.GetAdminRepository()
	return adminRepo.UpdateTeam(ctx, team)
}

func (b *BasicIAMService) DeleteTeam(ctx context.Context, teamId string) error {

	adminRepo := vars.RepositoryFactory.GetAdminRepository()
	return adminRepo.DeleteTeam(ctx, teamId)
}

func (b *BasicIAMService) GetTeamByOwner(ctx context.Context, owner string) (*entity.Team, error) {

	adminRepo := vars.RepositoryFactory.GetAdminRepository()
	return adminRepo.GetTeamByOwner(ctx, owner)
}

func (b *BasicIAMService) GetTeamMembers(ctx context.Context, teamId string) ([]entity.TeamMember, error) {

	adminRepo := vars.RepositoryFactory.GetAdminRepository()
	return adminRepo.GetTeamMembers(ctx, teamId)
}

func (b *BasicIAMService) AddTeamMember(ctx context.Context, teamId string, member entity.TeamMember) error {

	adminRepo := vars.RepositoryFactory.GetAdminRepository()
	return adminRepo.AddTeamMember(ctx, teamId, member)
}

func (b *BasicIAMService) RemoveTeamMembers(ctx context.Context, teamId string, key []string) error {

	adminRepo := vars.RepositoryFactory.GetAdminRepository()
	return adminRepo.RemoveTeamMembers(ctx, teamId, key)
}

func (b *BasicIAMService) RemoveTeamMember(ctx context.Context, teamId string, key string) error {

	adminRepo := vars.RepositoryFactory.GetAdminRepository()
	return adminRepo.RemoveTeamMember(ctx, teamId, key)
}

func (b *BasicIAMService) UpdateTeamMemberEmail(ctx context.Context, teamId, key, email string) error {

	adminRepo := vars.RepositoryFactory.GetAdminRepository()
	return adminRepo.UpdateTeamMemberEmail(ctx, teamId, key, email)
}

func (b *BasicIAMService) UpdateTeamMemberRole(ctx context.Context, teamId, key string, role ztypes.Role) error {

	adminRepo := vars.RepositoryFactory.GetAdminRepository()
	return adminRepo.UpdateTeamMemberRole(ctx, teamId, key, role)
}

func (b *BasicIAMService) UpdateTeamMemberStatus(ctx context.Context, teamId, key string, status ztypes.InvitationStatus) error {

	adminRepo := vars.RepositoryFactory.GetAdminRepository()
	return adminRepo.UpdateTeamMemberStatus(ctx, teamId, key, status)
}

func (b *BasicIAMService) GetAllMemberships(ctx context.Context) ([]entity.TeamMember, error) {

	adminRepo := vars.RepositoryFactory.GetAdminRepository()
	return adminRepo.GetAllMemberships(ctx)
}

func (b *BasicIAMService) GetTeamMemberships(ctx context.Context, email string) ([]entity.TeamMember, error) {

	adminRepo := vars.RepositoryFactory.GetAdminRepository()
	return adminRepo.GetTeamMemberships(ctx, email)
}

func (b *BasicIAMService) GetTeamMembership(ctx context.Context, key string) (*entity.TeamMember, error) {

	adminRepo := vars.RepositoryFactory.GetAdminRepository()
	return adminRepo.GetTeamMembership(ctx, key)
}

func (b *BasicIAMService) GetTeamMembershipByEmail(ctx context.Context, team, email string) (*entity.TeamMember, error) {

	adminRepo := vars.RepositoryFactory.GetAdminRepository()
	return adminRepo.GetTeamMembershipByEmail(ctx, team, email)
}
