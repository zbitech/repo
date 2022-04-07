package auth

import (
	"context"
	"time"

	"github.com/zbitech/common/interfaces"
	"github.com/zbitech/common/pkg/errs"
	"github.com/zbitech/common/pkg/logger"
	"github.com/zbitech/common/pkg/model/config"
	"github.com/zbitech/common/pkg/model/entity"
	"github.com/zbitech/common/pkg/model/ztypes"
	"github.com/zbitech/common/pkg/rctx"
	"github.com/zbitech/common/pkg/vars"
)

type AccessAuthorizer struct {
	iamService interfaces.IAMServiceIF
}

func NewAccessAuthorizer(iamService interfaces.IAMServiceIF) interfaces.AccessAuthorizerIF {
	return &AccessAuthorizer{iamService: iamService}
}

func (a *AccessAuthorizer) getTeam(ctx context.Context, teamName, userId string) (*entity.Team, *entity.TeamMember, error) {

	ctx = rctx.BuildContext(ctx, rctx.Context(rctx.Component, "getTeam"), rctx.Context(rctx.StartTime, time.Now()))
	defer logger.LogComponentTime(ctx)

	team, err := a.iamService.GetTeam(ctx, teamName)
	if err != nil {
		return nil, nil, err
	}

	user, err := a.iamService.GetUser(ctx, userId)
	if err != nil {
		return nil, nil, err
	}

	ut := user.GetTeam(team.TeamId)
	if ut == nil {
		return nil, nil, errs.ErrDBItemNotFound
	}

	mbr, err := a.iamService.GetTeamMembership(ctx, ut.Key)
	if err != nil {
		return nil, nil, err
	}

	return team, mbr, nil
}

func (a *AccessAuthorizer) getOwnerInfo(ctx context.Context, ownerId string) (*entity.User, *entity.ResourceSummary, *config.SubscriptionPolicy, error) {

	ctx = rctx.BuildContext(ctx, rctx.Context(rctx.Component, "getOwnerInfo"), rctx.Context(rctx.StartTime, time.Now()))
	defer logger.LogComponentTime(ctx)

	projRepo := vars.RepositoryFactory.GetProjectRepository()
	owner, err := a.iamService.GetUser(ctx, ownerId)
	if err != nil {
		logger.Errorf(ctx, "Failed to get owner information - %s", err)
		return nil, nil, nil, err
	}

	ownerSummary, err := projRepo.GetUserSummary(ctx, ownerId)
	if err != nil {
		logger.Errorf(ctx, "Failed to get owner summary - %s", err)
		return nil, nil, nil, err
	}

	sPolicy := vars.AppConfig.Policy.GetSubscriptionPolicy(owner.Level)
	return owner, ownerSummary, sPolicy, nil
}

func (a *AccessAuthorizer) ValidateProjectAction(ctx context.Context, project string, action ztypes.ZBIAction) error {

	ctx = rctx.BuildContext(ctx, rctx.Context(rctx.Component, "ValidateProjectAction"), rctx.Context(rctx.StartTime, time.Now()))
	defer logger.LogComponentTime(ctx)

	currUser := rctx.GetCurrentUser(ctx)

	var ownerId string
	var proj *entity.Project
	var err error
	var mbr *entity.TeamMember

	if action == ztypes.ACTION_CREATE {
		ownerId = currUser.UserId
	} else {
		projRepo := vars.RepositoryFactory.GetProjectRepository()
		proj, err = projRepo.GetProject(ctx, project)
		if err != nil {
			return err
		}
		ownerId = proj.Owner
	}

	_, ownerSummary, sPolicy, err := a.getOwnerInfo(ctx, ownerId)
	if err != nil {
		logger.Errorf(ctx, "Failed to get owner information - %s", err)
		return errs.ErrProjectCreateNotAllowed
	}

	if proj != nil {
		_, mbr, err = a.getTeam(ctx, proj.TeamId, currUser.UserId)
		if err != nil {
			logger.Errorf(ctx, "Unable to get team information - %s", err)
		}
	}

	switch action {
	case ztypes.ACTION_CREATE:
		if currUser.IsOwner() || (mbr != nil && mbr.IsJoined() && mbr.IsAdmin()) {
			logger.Infof(ctx, "Evaluating if owner or team admin can create project. Summary (%v), Policy (%v)", ownerSummary, sPolicy)
			if ownerSummary.TotalProjects < sPolicy.MaxProjects {
				return nil
			} else {
				return errs.ErrMaxProjectsCreated
			}
		}

		return errs.ErrProjectCreateNotAllowed

	case ztypes.ACTION_UPDATE:
		if proj.Owner == currUser.UserId ||
			(mbr != nil && mbr.IsJoined() && mbr.IsAdmin()) {
			return nil
		}

		return errs.ErrProjectUpdateNotAllowed

	case ztypes.ACTION_DELETE:
		if currUser.IsAdmin() || proj.Owner == currUser.UserId || (mbr != nil && mbr.IsJoined() && mbr.IsAdmin()) {
			return nil
		}

		return errs.ErrProjectDeleteNotAllowed

	case ztypes.ACTION_ACCESS:
		if currUser.IsAdmin() || proj.Owner == currUser.UserId || (mbr != nil && mbr.IsJoined()) {
			return nil
		}

		return errs.ErrProjectAccessNotAllowed

	}

	return errs.ErrProjectAccessError
}

func (a *AccessAuthorizer) ValidateInstanceAction(ctx context.Context, project, instance string, action ztypes.ZBIAction) error {

	ctx = rctx.BuildContext(ctx, rctx.Context(rctx.Component, "ValidateInstanceAction"), rctx.Context(rctx.StartTime, time.Now()))
	defer logger.LogComponentTime(ctx)

	currUser := rctx.GetCurrentUser(ctx)

	var ownerId string
	var proj *entity.Project
	var err error
	var mbr *entity.TeamMember

	if action == ztypes.ACTION_CREATE {
		ownerId = currUser.UserId
	} else {
		projRepo := vars.RepositoryFactory.GetProjectRepository()
		proj, err = projRepo.GetProject(ctx, project)
		if err != nil {
			return err
		}
		ownerId = proj.Owner
	}

	_, ownerSummary, sPolicy, err := a.getOwnerInfo(ctx, ownerId)
	if err != nil {
		logger.Errorf(ctx, "Failed to get owner information - %s", err)
		return errs.ErrProjectCreateNotAllowed
	}

	if proj != nil {
		_, mbr, err = a.getTeam(ctx, proj.TeamId, currUser.UserId)
		if err != nil {
			logger.Errorf(ctx, "Unable to get team information - %s", err)
		}
	}

	switch action {
	case ztypes.ACTION_CREATE:
		if currUser.IsOwner() || (mbr != nil && mbr.IsJoined() && mbr.IsAdmin()) {
			if ownerSummary.TotalInstances < sPolicy.MaxInstances {
				return nil
			} else {
				return errs.ErrMaxInstancesCreated
			}
		}

		return errs.ErrInstanceCreateNotAllowed

	case ztypes.ACTION_UPDATE:
		if proj.Owner == currUser.UserId ||
			(mbr != nil && mbr.IsJoined() && mbr.IsAdmin()) {
			return nil
		}

		return errs.ErrInstanceUpdateNotAllowed

	case ztypes.ACTION_DELETE:
		if currUser.IsAdmin() || proj.Owner == currUser.UserId || (mbr != nil && mbr.IsJoined() && mbr.IsAdmin()) {
			return nil
		}

		return errs.ErrInstanceDeleteNotAllowed

	case ztypes.ACTION_ACCESS:
		if currUser.IsAdmin() || proj.Owner == currUser.UserId || (mbr != nil && mbr.IsJoined()) {
			return nil
		}

		return errs.ErrInstanceAccessNotAllowed

	}

	return errs.ErrInstanceAccessError
}

func (a *AccessAuthorizer) ValidateTeamAction(ctx context.Context, teamId string, action ztypes.ZBIAction) error {

	ctx = rctx.BuildContext(ctx, rctx.Context(rctx.Component, "ValidateTeamAction"), rctx.Context(rctx.StartTime, time.Now()))
	defer logger.LogComponentTime(ctx)

	currUser := rctx.GetCurrentUser(ctx)

	var ownerId string
	var err error
	var mbr *entity.TeamMember

	if action == ztypes.ACTION_CREATE {
		ownerId = currUser.UserId
	}

	_, ownerSummary, sPolicy, err := a.getOwnerInfo(ctx, ownerId)
	if err != nil {
		logger.Errorf(ctx, "Failed to get owner information - %s", err)
		return errs.ErrTeamAccessError
	}

	team, err := a.iamService.GetTeam(ctx, teamId)
	if err != nil {
		logger.Errorf(ctx, "Failed to get team details - %s", err)
		return errs.ErrTeamAccessError
	}

	ut := currUser.User.GetTeam(team.TeamId)
	if ut != nil {
		mbr, err = a.iamService.GetTeamMembership(ctx, ut.Key)
		if err != nil {
			logger.Errorf(ctx, "Failed to get team membership info - %s", err)
		}
	}

	switch action {
	case ztypes.ACTION_CREATE:
		if currUser.IsOwner() || (mbr != nil && mbr.IsJoined() && mbr.IsAdmin()) {
			if ownerSummary.TotalTeams < sPolicy.MaxTeams {
				return nil
			} else {
				return errs.ErrMaxTeamsCreated
			}
		}

		return errs.ErrTeamCreateNotAllowed

	case ztypes.ACTION_UPDATE:
		if team.Owner == currUser.UserId ||
			(mbr != nil && mbr.IsJoined() && mbr.IsAdmin()) {
			return nil
		}

		return errs.ErrTeamUpdateNotAllowed

	case ztypes.ACTION_DELETE:
		if currUser.IsAdmin() || team.Owner == currUser.UserId || (mbr != nil && mbr.IsJoined() && mbr.IsAdmin()) {
			return nil
		}

		return errs.ErrProjectDeleteNotAllowed

	case ztypes.ACTION_ACCESS:
		if currUser.IsAdmin() || team.Owner == currUser.UserId || (mbr != nil && mbr.IsJoined()) {
			return nil
		}

		return errs.ErrTeamAccessNotAllowed

	}

	return errs.ErrTeamAccessError
}

func (a *AccessAuthorizer) ValidateUserInstanceMethodAccess(ctx context.Context, project, instance, method string) (ztypes.SubscriptionLevel, error) {

	ctx = rctx.BuildContext(ctx, rctx.Context(rctx.Component, "ValidateUserInstanceMethodAccess"), rctx.Context(rctx.StartTime, time.Now()))
	defer logger.LogComponentTime(ctx)

	projRepo := vars.RepositoryFactory.GetProjectRepository()

	proj, err := projRepo.GetProject(ctx, project)
	if err != nil {
		return ztypes.NO_SUB_LEVEL, err
	}

	instancePolicy, err := a.iamService.GetInstanceMethodPolicy(ctx, project, instance, method)
	if err != nil {
		return ztypes.NO_SUB_LEVEL, err
	}

	if !instancePolicy.Allow {
		return ztypes.NO_SUB_LEVEL, errs.ErrInstanceAccessNotAllowed
	}

	currUser := rctx.GetCurrentUser(ctx)

	team, mbr, err := a.getTeam(ctx, proj.TeamId, currUser.UserId)
	if err != nil || !team.IsOwner(proj.Owner) || !mbr.IsJoined() {
		return ztypes.NO_SUB_LEVEL, err
	}

	userPolicy, err := a.iamService.GetUserPolicy(ctx, proj.Owner)
	if err != nil {
		return ztypes.NO_SUB_LEVEL, err
	}

	userInstanceAccess := userPolicy.GetInstanceAccess(project, instance)
	methodPolicy := userInstanceAccess.GetMethodByName(method)
	if !userInstanceAccess.Allow || (methodPolicy != nil && !methodPolicy.Allow) {
		return ztypes.NO_SUB_LEVEL, err
	}

	// Return project owner's subscription level
	owner, err := a.iamService.GetUser(ctx, proj.Owner)
	if err != nil {
		return ztypes.NO_SUB_LEVEL, err
	}

	return owner.Level, nil
}

func (a *AccessAuthorizer) ValidateAPIKeyInstanceMethodAccess(ctx context.Context, project, instance, method, apikey string) (ztypes.SubscriptionLevel, error) {

	ctx = rctx.BuildContext(ctx, rctx.Context(rctx.Component, "ValidateAPIKeyInstanceMethodAccess"), rctx.Context(rctx.StartTime, time.Now()))
	defer logger.LogComponentTime(ctx)

	projRepo := vars.RepositoryFactory.GetProjectRepository()

	proj, err := projRepo.GetProject(ctx, project)
	if err != nil {
		return ztypes.NO_SUB_LEVEL, err
	}

	instancePolicy, err := a.iamService.GetInstanceMethodPolicy(ctx, project, instance, method)
	if err != nil {
		return ztypes.NO_SUB_LEVEL, err
	}

	if !instancePolicy.Allow {
		return ztypes.NO_SUB_LEVEL, errs.ErrInstanceAccessNotAllowed
	}

	apiKey, err := a.iamService.GetAPIKey(ctx, apikey)
	if err != nil {
		return ztypes.NO_SUB_LEVEL, err
	}

	team, mbr, err := a.getTeam(ctx, proj.TeamId, apiKey.UserId)
	if err != nil || !team.IsOwner(proj.Owner) || !mbr.IsJoined() {
		return ztypes.NO_SUB_LEVEL, err
	}

	keyPolicy, err := a.iamService.GetAPIKeyPolicy(ctx, apikey)
	if err != nil {
		return ztypes.NO_SUB_LEVEL, nil
	}

	keyInstance := keyPolicy.GetInstanceAccess(project, instance)
	methodPolicy := keyInstance.GetMethodByName(method)
	if !keyInstance.Allow || methodPolicy == nil || !methodPolicy.Allow {
		return ztypes.NO_SUB_LEVEL, err
	}

	// Return project owner's subscription level
	owner, err := a.iamService.GetUser(ctx, proj.Owner)
	if err != nil {
		return ztypes.NO_SUB_LEVEL, err
	}

	return owner.Level, nil
}
