package memory

import (
	"context"
	"github.com/zbitech/common/pkg/errs"
	"time"

	"github.com/zbitech/common/pkg/logger"
	"github.com/zbitech/common/pkg/model/entity"
	"github.com/zbitech/common/pkg/model/ztypes"
	"github.com/zbitech/common/pkg/rctx"
)

func (m *AdminMemoryRepository) GetExpiringInvitations(ctx context.Context, date time.Time) ([]entity.TeamMember, error) {

	ctx = rctx.BuildContext(ctx, rctx.Context(rctx.Component, "GetExpiringInvitations"), rctx.Context(rctx.StartTime, time.Now()))
	defer logger.LogComponentTime(ctx)

	items := m.members.GetItems()
	members := make([]entity.TeamMember, 0)

	for _, item := range items {
		member := item.(*entity.TeamMember)
		if member.ExpiresOn.After(time.Now()) {
			members = append(members, *member)
		}
	}

	return members, nil
}

func (m *AdminMemoryRepository) PurgeExpiredInvitations(ctx context.Context) (int64, error) {

	ctx = rctx.BuildContext(ctx, rctx.Context(rctx.Component, "PurgeExpiredInvitations"), rctx.Context(rctx.StartTime, time.Now()))
	defer logger.LogComponentTime(ctx)

	items := m.members.GetItems()
	var deleted int64 = 0

	for _, item := range items {
		member := item.(*entity.TeamMember)
		if member.IsExpired() {
			m.members.RemoveItem(member.Key)
			deleted++
		}
	}

	return deleted, nil
}

func (m *AdminMemoryRepository) CreateTeam(ctx context.Context, team entity.Team) error {
	ctx = rctx.BuildContext(ctx, rctx.Context(rctx.Component, "CreateTeam"), rctx.Context(rctx.StartTime, time.Now()))
	defer logger.LogComponentTime(ctx)

	_, err := m.teams.GetItem(team.TeamId)
	if err != nil {
		return nil
	}

	m.teams.StoreItem(team.TeamId, &team)
	return nil
}

func (m *AdminMemoryRepository) GetTeams(ctx context.Context) ([]entity.Team, error) {

	ctx = rctx.BuildContext(ctx, rctx.Context(rctx.Component, "getTeams"), rctx.Context(rctx.StartTime, time.Now()))
	defer logger.LogComponentTime(ctx)

	items := m.teams.GetItems()
	teams := make([]entity.Team, 0)

	for _, item := range items {
		team := item.(*entity.Team)
		teams = append(teams, *team)
	}

	return teams, nil
}

func (m *AdminMemoryRepository) GetTeam(ctx context.Context, teamId string) (*entity.Team, error) {

	ctx = rctx.BuildContext(ctx, rctx.Context(rctx.Component, "GetTeam"), rctx.Context(rctx.StartTime, time.Now()))
	defer logger.LogComponentTime(ctx)

	item, err := m.teams.GetItem(teamId)
	if err != nil {
		return nil, err
	}

	return item.(*entity.Team), nil
}

func (m *AdminMemoryRepository) GetTeamByOwner(ctx context.Context, userid string) (*entity.Team, error) {
	ctx = rctx.BuildContext(ctx, rctx.Context(rctx.Component, "GetTeamByOwner"), rctx.Context(rctx.StartTime, time.Now()))
	defer logger.LogComponentTime(ctx)

	items := m.teams.GetItems()
	//	teams := make([]entity.Team, 0)

	for _, item := range items {
		team := item.(*entity.Team)
		if team.Owner == userid {
			return team, nil
		}
		//		teams = append(teams, *team)
	}

	//	return teams, nil
	return nil, errs.ErrDBItemNotFound
}

func (m *AdminMemoryRepository) GetTeamMembers(ctx context.Context, teamId string) ([]entity.TeamMember, error) {

	ctx = rctx.BuildContext(ctx, rctx.Context(rctx.Component, "GetTeamMembers"), rctx.Context(rctx.StartTime, time.Now()))
	defer logger.LogComponentTime(ctx)

	items := m.members.GetItems()
	members := make([]entity.TeamMember, 0)

	for _, item := range items {
		member := item.(*entity.TeamMember)
		if member.TeamId == teamId {
			members = append(members, *member)
		}
	}

	return members, nil
}

func (m *AdminMemoryRepository) AddTeamMember(ctx context.Context, teamId string, member entity.TeamMember) error {

	ctx = rctx.BuildContext(ctx, rctx.Context(rctx.Component, "AddTeamMembers"), rctx.Context(rctx.StartTime, time.Now()))
	defer logger.LogComponentTime(ctx)

	m.members.StoreItem(member.Key, member)
	return nil
}

func (m *AdminMemoryRepository) RemoveTeamMembers(ctx context.Context, teamId string, keys []string) error {

	ctx = rctx.BuildContext(ctx, rctx.Context(rctx.Component, "RemoveTeamMembers"), rctx.Context(rctx.StartTime, time.Now()))
	defer logger.LogComponentTime(ctx)

	for _, key := range keys {
		m.members.RemoveItem(key)
	}

	return nil
}

func (m *AdminMemoryRepository) RemoveTeamMember(ctx context.Context, teamId string, key string) error {

	ctx = rctx.BuildContext(ctx, rctx.Context(rctx.Component, "RemoveTeamMember"), rctx.Context(rctx.StartTime, time.Now()))
	defer logger.LogComponentTime(ctx)

	m.members.RemoveItem(key)
	return nil

}

func (m *AdminMemoryRepository) UpdateTeamMemberEmail(ctx context.Context, teamId, key, email string) error {

	ctx = rctx.BuildContext(ctx, rctx.Context(rctx.Component, "UpdateTeamMemberEmail"), rctx.Context(rctx.StartTime, time.Now()))
	defer logger.LogComponentTime(ctx)

	item, err := m.members.GetItem(key)
	if err != nil {
		return err
	}

	member := item.(*entity.TeamMember)
	member.Email = email

	m.members.StoreItem(key, m)
	return nil
}

func (m *AdminMemoryRepository) UpdateTeamMemberRole(ctx context.Context, teamId, key string, role ztypes.Role) error {

	ctx = rctx.BuildContext(ctx, rctx.Context(rctx.Component, "UpdateTeamMemberRole"), rctx.Context(rctx.StartTime, time.Now()))
	defer logger.LogComponentTime(ctx)

	item, err := m.members.GetItem(key)
	if err != nil {
		return err
	}

	member := item.(*entity.TeamMember)
	member.Role = role

	m.members.StoreItem(key, m)
	return nil
}

func (m *AdminMemoryRepository) UpdateTeamMemberStatus(ctx context.Context, teamId, key string, status ztypes.InvitationStatus) error {

	ctx = rctx.BuildContext(ctx, rctx.Context(rctx.Component, "UpdateTeamMemberStatus"), rctx.Context(rctx.StartTime, time.Now()))
	defer logger.LogComponentTime(ctx)

	item, err := m.members.GetItem(key)
	if err != nil {
		return err
	}

	member := item.(*entity.TeamMember)
	member.Status = status

	m.members.StoreItem(key, m)
	return nil
}

func (m *AdminMemoryRepository) UpdateTeam(ctx context.Context, team entity.Team) error {

	ctx = rctx.BuildContext(ctx, rctx.Context(rctx.Component, "UpdateTeam"), rctx.Context(rctx.StartTime, time.Now()))
	defer logger.LogComponentTime(ctx)

	m.teams.StoreItem(team.TeamId, team)
	return nil
}

func (m *AdminMemoryRepository) DeleteTeam(ctx context.Context, teamId string) error {

	ctx = rctx.BuildContext(ctx, rctx.Context(rctx.Component, "DeleteTeam"), rctx.Context(rctx.StartTime, time.Now()))
	defer logger.LogComponentTime(ctx)

	m.teams.RemoveItem(teamId)
	return nil
}

func (m *AdminMemoryRepository) GetAllMemberships(ctx context.Context) ([]entity.TeamMember, error) {
	ctx = rctx.BuildContext(ctx, rctx.Context(rctx.Component, "GetAllTeamMemberships"), rctx.Context(rctx.StartTime, time.Now()))
	defer logger.LogComponentTime(ctx)

	items := m.members.GetItems()
	members := make([]entity.TeamMember, 0)

	for _, item := range items {
		member := item.(*entity.TeamMember)
		members = append(members, *member)
	}

	return members, nil
}

func (m *AdminMemoryRepository) GetTeamMemberships(ctx context.Context, email string) ([]entity.TeamMember, error) {
	ctx = rctx.BuildContext(ctx, rctx.Context(rctx.Component, "GetTeamMemberships"), rctx.Context(rctx.StartTime, time.Now()))
	defer logger.LogComponentTime(ctx)

	items := m.members.GetItems()
	members := make([]entity.TeamMember, 0)

	for _, item := range items {
		member := item.(*entity.TeamMember)
		if member.Email == email {
			members = append(members, *member)
		}
	}

	return members, nil
}

func (m *AdminMemoryRepository) GetTeamMembership(ctx context.Context, key string) (*entity.TeamMember, error) {
	ctx = rctx.BuildContext(ctx, rctx.Context(rctx.Component, "GetTeamMembership"), rctx.Context(rctx.StartTime, time.Now()))
	defer logger.LogComponentTime(ctx)

	item, err := m.members.GetItem(key)
	if err != nil {
		return nil, nil
	}

	return item.(*entity.TeamMember), nil
}

func (m *AdminMemoryRepository) GetTeamMembershipByEmail(ctx context.Context, teamId, email string) (*entity.TeamMember, error) {
	ctx = rctx.BuildContext(ctx, rctx.Context(rctx.Component, "GetTeamMembershipByEmail"), rctx.Context(rctx.StartTime, time.Now()))
	defer logger.LogComponentTime(ctx)

	items := m.members.GetItems()

	for _, item := range items {
		member := item.(*entity.TeamMember)
		if member.Email == email && member.TeamId == teamId {
			return member, nil
		}
	}

	return nil, errs.ErrDBItemNotFound
}
