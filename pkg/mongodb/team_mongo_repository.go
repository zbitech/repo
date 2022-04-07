package mongodb

import (
	"context"
	"github.com/zbitech/common/pkg/utils"
	"time"

	"github.com/zbitech/common/pkg/errs"
	"github.com/zbitech/common/pkg/logger"
	"github.com/zbitech/common/pkg/model/entity"
	"github.com/zbitech/common/pkg/model/ztypes"
	"github.com/zbitech/common/pkg/rctx"
	"github.com/zbitech/common/pkg/vars"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (m *AdminMongoRepository) GetExpiringInvitations(ctx context.Context, date time.Time) ([]entity.TeamMember, error) {

	ctx = rctx.BuildContext(ctx, rctx.Context(rctx.Component, "GetExpiringInvitations"), rctx.Context(rctx.StartTime, time.Now()))
	defer logger.LogComponentTime(ctx)

	collection := m.conn.GetCollection(vars.MONGODB_NAME, MONGODB_COLL_TEAM_MEMBERS)
	filter := bson.M{"expireson": date}

	result, err := collection.Find(ctx, filter)

	if err != nil {
		return nil, handleMongoError(ctx, err)
	}

	var invites []entity.TeamMember
	err = result.Decode(&invites)
	if err != nil {
		return nil, handleMongoError(ctx, err)
	}

	return invites, nil
}

func (m *AdminMongoRepository) PurgeExpiredInvitations(ctx context.Context) (int64, error) {

	ctx = rctx.BuildContext(ctx, rctx.Context(rctx.Component, "PurgeExpiredInvitations"), rctx.Context(rctx.StartTime, time.Now()))
	defer logger.LogComponentTime(ctx)

	collection := m.conn.GetCollection(vars.MONGODB_NAME, MONGODB_COLL_TEAM_MEMBERS)
	filter := bson.M{"status": ztypes.EXPIRED_INVITATION}

	result, err := collection.DeleteMany(ctx, filter)
	if err != nil {
		return 0, handleMongoError(ctx, err)
	}

	return result.DeletedCount, nil
}

func (m *AdminMongoRepository) CreateTeam(ctx context.Context, team entity.Team) error {
	ctx = rctx.BuildContext(ctx, rctx.Context(rctx.Component, "CreateTeam"), rctx.Context(rctx.StartTime, time.Now()))
	defer logger.LogComponentTime(ctx)

	collection := m.conn.GetCollection(vars.MONGODB_NAME, MONGODB_COLL_TEAMS)
	filter := bson.M{"_id": team.TeamId}
	result := collection.FindOne(ctx, filter)

	err := result.Err()
	if err == mongo.ErrNoDocuments {
		_, err = collection.InsertOne(ctx, team)
	}

	if err != nil {
		logger.Errorf(ctx, "Unable to insert team %s - %s", utils.MarshalObject(team), err)
		return errs.ErrDBItemInsertFailed
	}

	return nil
}

func (m *AdminMongoRepository) GetTeams(ctx context.Context) ([]entity.Team, error) {

	ctx = rctx.BuildContext(ctx, rctx.Context(rctx.Component, "getTeams"), rctx.Context(rctx.StartTime, time.Now()))
	defer logger.LogComponentTime(ctx)

	collection := m.conn.GetCollection(vars.MONGODB_NAME, MONGODB_COLL_TEAMS)
	filter := bson.D{}
	result, err := collection.Find(ctx, filter)

	if err != nil {
		return nil, handleMongoError(ctx, err)
	}

	teams := make([]entity.Team, 0)
	err = result.All(ctx, &teams)

	if err != nil {
		logger.Errorf(ctx, "Unable to marshal result - %s", err)
		return nil, errs.ErrMarshalFailed
	}

	return teams, nil
}

func (m *AdminMongoRepository) GetTeam(ctx context.Context, teamId string) (*entity.Team, error) {

	ctx = rctx.BuildContext(ctx, rctx.Context(rctx.Component, "GetTeam"), rctx.Context(rctx.StartTime, time.Now()))
	defer logger.LogComponentTime(ctx)

	collection := m.conn.GetCollection(vars.MONGODB_NAME, MONGODB_COLL_TEAMS)
	filter := bson.M{"_id": teamId}
	result := collection.FindOne(ctx, filter)

	err := result.Err()
	if err != nil {
		return nil, handleMongoError(ctx, err)
	}

	team := entity.Team{}
	err = result.Decode(&team)

	if err != nil {
		logger.Errorf(ctx, "Unable to marshal result - %s", err)
		return nil, errs.ErrMarshalFailed
	}

	return &team, nil
}

func (m *AdminMongoRepository) GetTeamByOwner(ctx context.Context, userid string) (*entity.Team, error) {
	ctx = rctx.BuildContext(ctx, rctx.Context(rctx.Component, "GetTeamByOwner"), rctx.Context(rctx.StartTime, time.Now()))
	defer logger.LogComponentTime(ctx)

	collection := m.conn.GetCollection(vars.MONGODB_NAME, MONGODB_COLL_TEAMS)
	filter := bson.M{"owner": userid}
	result := collection.FindOne(ctx, filter)

	err := result.Err()
	if err != nil {
		return nil, handleMongoError(ctx, err)
	}

	team := entity.Team{}
	err = result.Decode(&team)

	if err != nil {
		logger.Errorf(ctx, "Unable to marshal result - %s", err)
		return nil, errs.ErrMarshalFailed
	}

	return &team, nil
}

func (m *AdminMongoRepository) GetTeamMembers(ctx context.Context, teamId string) ([]entity.TeamMember, error) {

	ctx = rctx.BuildContext(ctx, rctx.Context(rctx.Component, "getTeamMembers"),
		rctx.Context(rctx.StartTime, time.Now()))
	defer logger.LogComponentTime(ctx)

	coll := m.conn.GetCollection(vars.MONGODB_NAME, MONGODB_COLL_TEAM_MEMBERS)
	filter := bson.M{"_id": teamId}

	result, err := coll.Find(ctx, filter)
	if err != nil {
		return nil, handleMongoError(ctx, err)
	}

	members := make([]entity.TeamMember, 0)
	err = result.All(ctx, &members)
	if err != nil {
		return nil, errs.ErrMarshalFailed
	}

	return members, nil
}

func (m *AdminMongoRepository) AddTeamMember(ctx context.Context, teamId string, member entity.TeamMember) error {

	ctx = rctx.BuildContext(ctx, rctx.Context(rctx.Component, "setTeamMembers"), rctx.Context(rctx.StartTime, time.Now()))
	defer logger.LogComponentTime(ctx)

	coll := m.conn.GetCollection(vars.MONGODB_NAME, MONGODB_COLL_TEAM_MEMBERS)
	filter := bson.M{"teamid": teamId, "email": member.Email}

	result := coll.FindOne(ctx, filter)
	err := result.Err()

	if err != nil {
		if err == mongo.ErrNoDocuments {
			i_result, err := coll.InsertOne(ctx, member)
			if err != nil {
				return handleMongoError(ctx, err)
			}

			logger.Infof(ctx, "Inserted new member: %s, key: %s, team: %s, id: %s", member.Email, member.Key, member.TeamId, i_result.InsertedID)
			return nil
		}

		return handleMongoError(ctx, err)
	}

	return errs.ErrDBKeyAlreadyExists
}

func (m *AdminMongoRepository) RemoveTeamMembers(ctx context.Context, teamId string, key []string) error {

	ctx = rctx.BuildContext(ctx, rctx.Context(rctx.Component, "removeTeamMembers"), rctx.Context(rctx.StartTime, time.Now()))
	defer logger.LogComponentTime(ctx)

	coll := m.conn.GetCollection(vars.MONGODB_NAME, MONGODB_COLL_TEAM_MEMBERS)
	filter := bson.M{"key": bson.M{"$in": key}}

	result, err := coll.DeleteMany(ctx, filter)
	if err != nil {
		return handleMongoError(ctx, err)
	}

	logger.Infof(ctx, "Deleted %d members from team %s", result.DeletedCount, teamId)

	if result.DeletedCount != int64(len(key)) {
		return errs.ErrDBItemDeleteFailed
	}

	return nil

}

func (m *AdminMongoRepository) RemoveTeamMember(ctx context.Context, teamId string, key string) error {

	ctx = rctx.BuildContext(ctx, rctx.Context(rctx.Component, "removeTeamMember"), rctx.Context(rctx.StartTime, time.Now()))
	defer logger.LogComponentTime(ctx)

	coll := m.conn.GetCollection(vars.MONGODB_NAME, MONGODB_COLL_TEAM_MEMBERS)
	filter := bson.M{"key": key}

	result, err := coll.DeleteOne(ctx, filter)
	if err != nil {
		return handleMongoError(ctx, err)
	}

	logger.Infof(ctx, "Deleted %d members from team %s", result.DeletedCount, teamId)

	if result.DeletedCount != 1 {
		return errs.ErrDBItemDeleteFailed
	}

	return nil

}

func (m *AdminMongoRepository) UpdateTeamMemberEmail(ctx context.Context, teamId, key, email string) error {

	ctx = rctx.BuildContext(ctx, rctx.Context(rctx.Component, "updateTeamMemberEmail"), rctx.Context(rctx.StartTime, time.Now()))
	defer logger.LogComponentTime(ctx)

	coll := m.conn.GetCollection(vars.MONGODB_NAME, MONGODB_COLL_TEAM_MEMBERS)
	filter := bson.M{"teamid": teamId, "key": key}
	update := bson.M{"$set": bson.M{"email": email, "lastupdate": time.Now()}}

	result := coll.FindOneAndUpdate(ctx, filter, update)
	err := result.Err()

	if err != nil {
		return handleMongoError(ctx, err)
	}

	return nil
}

func (m *AdminMongoRepository) UpdateTeamMemberRole(ctx context.Context, teamId, key string, role ztypes.Role) error {

	ctx = rctx.BuildContext(ctx, rctx.Context(rctx.Component, "updateTeamMemberRole"), rctx.Context(rctx.StartTime, time.Now()))
	defer logger.LogComponentTime(ctx)

	coll := m.conn.GetCollection(vars.MONGODB_NAME, MONGODB_COLL_TEAM_MEMBERS)
	filter := bson.M{"teamid": teamId, "key": key}
	update := bson.M{"$set": bson.M{"role": role, "lastupdate": time.Now()}}

	result := coll.FindOneAndUpdate(ctx, filter, update)
	err := result.Err()

	if err != nil {
		return handleMongoError(ctx, err)
	}

	return nil
}

func (m *AdminMongoRepository) UpdateTeamMemberStatus(ctx context.Context, teamId, key string, status ztypes.InvitationStatus) error {

	ctx = rctx.BuildContext(ctx, rctx.Context(rctx.Component, "updateTeamMemberStatus"), rctx.Context(rctx.StartTime, time.Now()))
	defer logger.LogComponentTime(ctx)

	coll := m.conn.GetCollection(vars.MONGODB_NAME, MONGODB_COLL_TEAM_MEMBERS)
	filter := bson.M{"teamid": teamId, "key": key}
	update := bson.M{"$set": bson.M{"status": status, "lastupdate": time.Now()}}

	result := coll.FindOneAndUpdate(ctx, filter, update)
	err := result.Err()

	if err != nil {
		return handleMongoError(ctx, err)
	}

	return nil
}

func (m *AdminMongoRepository) UpdateTeam(ctx context.Context, team entity.Team) error {

	ctx = rctx.BuildContext(ctx, rctx.Context(rctx.Component, "updateTeam"), rctx.Context(rctx.StartTime, time.Now()))
	defer logger.LogComponentTime(ctx)

	collection := m.conn.GetCollection(vars.MONGODB_NAME, MONGODB_COLL_TEAMS)
	filter := bson.M{"_id": team.TeamId}

	result := collection.FindOneAndUpdate(ctx, filter, team)
	err := result.Err()

	if err != nil {
		return handleMongoError(ctx, err)
	}

	return nil
}

func (m *AdminMongoRepository) DeleteTeam(ctx context.Context, teamId string) error {

	ctx = rctx.BuildContext(ctx, rctx.Context(rctx.Component, "deleteTeam"), rctx.Context(rctx.StartTime, time.Now()))
	defer logger.LogComponentTime(ctx)

	collection := m.conn.GetCollection(vars.MONGODB_NAME, MONGODB_COLL_TEAMS)
	filter := bson.M{"_id": teamId}

	result := collection.FindOneAndDelete(ctx, filter)
	err := result.Err()

	if err != nil {
		return handleMongoError(ctx, err)
	}

	return nil
}

func (m *AdminMongoRepository) GetAllMemberships(ctx context.Context) ([]entity.TeamMember, error) {
	ctx = rctx.BuildContext(ctx, rctx.Context(rctx.Component, "getTeamMemberships"), rctx.Context(rctx.StartTime, time.Now()))
	defer logger.LogComponentTime(ctx)

	coll := m.conn.GetCollection(vars.MONGODB_NAME, MONGODB_COLL_TEAM_MEMBERS)
	filter := bson.D{}

	result, err := coll.Find(ctx, filter)
	if err != nil {
		return nil, handleMongoError(ctx, err)
	}

	members := make([]entity.TeamMember, 0)
	if err = result.All(ctx, &members); err != nil {
		return nil, errs.ErrMarshalFailed
	}

	return members, nil
}

func (m *AdminMongoRepository) GetTeamMemberships(ctx context.Context, email string) ([]entity.TeamMember, error) {
	ctx = rctx.BuildContext(ctx, rctx.Context(rctx.Component, "getTeamMemberships"), rctx.Context(rctx.StartTime, time.Now()))
	defer logger.LogComponentTime(ctx)

	coll := m.conn.GetCollection(vars.MONGODB_NAME, MONGODB_COLL_TEAM_MEMBERS)
	filter := bson.M{"email": email}

	result, err := coll.Find(ctx, filter)
	if err != nil {
		return nil, handleMongoError(ctx, err)
	}

	members := make([]entity.TeamMember, 0)
	if err = result.All(ctx, &members); err != nil {
		return nil, errs.ErrMarshalFailed
	}

	return members, nil
}

func (m *AdminMongoRepository) GetTeamMembership(ctx context.Context, key string) (*entity.TeamMember, error) {
	ctx = rctx.BuildContext(ctx, rctx.Context(rctx.Component, "getTeamMembership"),
		rctx.Context(rctx.StartTime, time.Now()))
	defer logger.LogComponentTime(ctx)

	coll := m.conn.GetCollection(vars.MONGODB_NAME, MONGODB_COLL_TEAM_MEMBERS)
	filter := bson.M{"key": key}

	result := coll.FindOne(ctx, filter)
	err := result.Err()

	if err != nil {
		return nil, handleMongoError(ctx, err)
	}

	member := entity.TeamMember{}
	if err = result.Decode(&member); err != nil {
		return nil, errs.ErrMarshalFailed
	}

	return &member, nil
}

func (m *AdminMongoRepository) GetTeamMembershipByEmail(ctx context.Context, teamId, email string) (*entity.TeamMember, error) {
	ctx = rctx.BuildContext(ctx, rctx.Context(rctx.Component, "getTeamMembershipByEmail"),
		rctx.Context(rctx.StartTime, time.Now()))
	defer logger.LogComponentTime(ctx)

	coll := m.conn.GetCollection(vars.MONGODB_NAME, MONGODB_COLL_TEAM_MEMBERS)
	filter := bson.M{"teamid": teamId, "email": email}

	result := coll.FindOne(ctx, filter)
	err := result.Err()

	if err != nil {
		return nil, handleMongoError(ctx, err)
	}

	member := entity.TeamMember{}
	if err = result.Decode(&member); err != nil {
		return nil, errs.ErrMarshalFailed
	}

	return &member, nil
}
