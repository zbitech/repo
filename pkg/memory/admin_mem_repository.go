package memory

import (
	"context"
	"fmt"
	"github.com/zbitech/common/pkg/errs"
	"github.com/zbitech/common/pkg/id"
	"github.com/zbitech/common/pkg/model/config"
	"github.com/zbitech/common/pkg/model/entity"
	"github.com/zbitech/common/pkg/rctx"
	"github.com/zbitech/common/pkg/utils"
	"github.com/zbitech/common/pkg/vars"
	"time"

	"github.com/zbitech/common/pkg/logger"
	mem "github.com/zbitech/common/pkg/memory"
	"github.com/zbitech/repo/internal/helper"
)

type AdminMemoryRepository struct {
	users            *mem.MemoryStore
	passwords        map[string]string
	apikeys          *mem.MemoryStore
	userPolicies     *mem.MemoryStore
	instancePolicies *mem.MemoryStore
	apikeyPolicies   *mem.MemoryStore
	teams            *mem.MemoryStore
	members          *mem.MemoryStore
}

func NewAdminMemoryRepository(ctx context.Context, path string) (*AdminMemoryRepository, error) {

	var store *AdminMemoryRepository = &AdminMemoryRepository{
		users:            mem.NewMemoryStore(),
		passwords:        make(map[string]string, 0),
		apikeys:          mem.NewMemoryStore(),
		userPolicies:     mem.NewMemoryStore(),
		instancePolicies: mem.NewMemoryStore(),
		apikeyPolicies:   mem.NewMemoryStore(),
		teams:            mem.NewMemoryStore(),
		members:          mem.NewMemoryStore(),
	}

	var adminConfig config.AdminConfig
	configPath := fmt.Sprintf("%s/admin.yaml", vars.ASSET_PATH_DIRECTORY)
	err := utils.ReadConfig(configPath, nil, &adminConfig)
	if err != nil {
		logger.Errorf(ctx, "Unable to load default admin data: %s", err)
		return nil, errs.ErrDBError
	} else {

		for _, user := range adminConfig.Users {
			if user.Memberships == nil {
				user.Memberships = make([]entity.UserTeam, 0)
			}
			user.Created = time.Now()
			user.LastUpdate = time.Now()
			user.Active = true

			userid := user.UserId
			password := adminConfig.Passwords[userid]
			userPolicy := entity.NewUserPolicy(user.UserId)

			store.users.StoreItem(user.UserId, user)
			store.users.StoreItem(user.UserId, userPolicy)
			store.passwords[user.UserId] = password
		}

		for _, apikey := range adminConfig.Keys {
			apikey.Created = time.Now()
			apikey.Expires = apikey.Created.Add(time.Hour * 8760)
			keyPolicy := entity.NewAPIKeyPolicy(apikey.Key, true)

			store.apikeys.StoreItem(apikey.Key, apikey)
			store.apikeyPolicies.StoreItem(apikey.Key, keyPolicy)
		}

		for _, team := range adminConfig.Teams {
			team.Created = time.Now()
			team.LastUpdate = time.Now()
			store.teams.StoreItem(team.TeamId, team)
		}

		for _, mbr := range adminConfig.Members {
			mbr.CreatedOn = time.Now()
			mbr.LastUpdate = time.Now()
			store.members.StoreItem(mbr.Key, mbr)
		}

	}

	return store, nil
}

func (m *AdminMemoryRepository) RegisterUser(ctx context.Context, user *entity.User, pass *entity.UserPassword) error {
	ctx = rctx.BuildContext(ctx, rctx.Context(rctx.Component, "RegisterUser"), rctx.Context(rctx.StartTime, time.Now()))
	defer logger.LogComponentTime(ctx)

	_, err := m.users.GetItem(user.UserId)
	if err != nil {
		m.users.StoreItem(user.UserId, user)
		m.passwords[user.UserId] = pass.Password

		return nil
	}

	return errs.ErrUserAlreadyExists
}

func (m *AdminMemoryRepository) GetUsers(ctx context.Context) []entity.User {

	ctx = rctx.BuildContext(ctx, rctx.Context(rctx.Component, "GetUsers"), rctx.Context(rctx.StartTime, time.Now()))
	defer logger.LogComponentTime(ctx)

	items := m.users.GetItems()
	users := make([]entity.User, len(items))
	index := 0
	for _, item := range items {
		user := item.(*entity.User)
		users[index] = *user
		index++
	}

	return users
}

func (m *AdminMemoryRepository) GetUser(ctx context.Context, userId string) (*entity.User, error) {

	ctx = rctx.BuildContext(ctx, rctx.Context(rctx.Component, "GetUser"), rctx.Context(rctx.StartTime, time.Now()))
	defer logger.LogComponentTime(ctx)

	item, err := m.users.GetItem(userId)
	if err != nil {
		return nil, err
	}

	user := item.(*entity.User)
	return user, nil
}

func (m *AdminMemoryRepository) UpdateUser(ctx context.Context, user *entity.User) error {
	ctx = rctx.BuildContext(ctx, rctx.Context(rctx.Component, "UpdateUser"), rctx.Context(rctx.StartTime, time.Now()))
	defer logger.LogComponentTime(ctx)

	_, err := m.users.GetItem(user.UserId)
	if err != nil {
		return err
	}

	m.users.StoreItem(user.UserId, user)
	return nil
}

func (m *AdminMemoryRepository) UpdatePassword(ctx context.Context, userid string, password *entity.UserPassword) error {

	ctx = rctx.BuildContext(ctx, rctx.Context(rctx.Component, "UpdatePassword"), rctx.Context(rctx.StartTime, time.Now()))
	defer logger.LogComponentTime(ctx)

	_, err := m.users.GetItem(userid)
	if err != nil {
		return nil
	}

	m.passwords[userid] = password.Password
	return nil
}

func (m *AdminMemoryRepository) AuthenticateUser(ctx context.Context, userId, password string) (*string, error) {

	ctx = rctx.BuildContext(ctx, rctx.Context(rctx.Component, "AuthenticateUser"), rctx.Context(rctx.StartTime, time.Now()))
	defer logger.LogComponentTime(ctx)

	userPass, ok := m.passwords[userId]
	if !ok {
		return nil, errs.ErrAuthFailed
	}

	if id.ValidatePassword(userPass, []byte(password)) {
		user, err := m.GetUser(ctx, userId)
		if err != nil {
			return nil, err
		}

		return helper.GenerateJwtToken(*user)
	}

	return nil, errs.ErrAuthFailed
}

func (m *AdminMemoryRepository) GetAPIKeys(ctx context.Context, userId string) ([]string, error) {

	ctx = rctx.BuildContext(ctx, rctx.Context(rctx.Component, "GetAPIKey"), rctx.Context(rctx.StartTime, time.Now()))
	defer logger.LogComponentTime(ctx)

	keys := make([]string, 0)
	items := m.apikeys.GetItems()
	for _, item := range items {
		apiKey := item.(*entity.APIKey)
		if apiKey.UserId == userId {
			keys = append(keys, apiKey.Key)
		}
	}

	return keys, nil
}

func (m *AdminMemoryRepository) GetAPIKey(ctx context.Context, apiKey string) (*entity.APIKey, error) {
	ctx = rctx.BuildContext(ctx, rctx.Context(rctx.Component, "GetAPIKey"), rctx.Context(rctx.StartTime, time.Now()))
	defer logger.LogComponentTime(ctx)

	item, err := m.apikeys.GetItem(apiKey)
	if err != nil {
		return nil, err
	}

	apikey := item.(*entity.APIKey)
	return apikey, nil
}

func (m *AdminMemoryRepository) CreateAPIKey(ctx context.Context, user_id string) (*entity.APIKey, error) {

	ctx = rctx.BuildContext(ctx, rctx.Context(rctx.Component, "CreateAPIKey"), rctx.Context(rctx.StartTime, time.Now()))
	defer logger.LogComponentTime(ctx)

	apikey := entity.NewAPIKey(user_id, vars.HOURS_IN_YEAR)
	m.apikeys.StoreItem(apikey.Key, &apikey)

	return &apikey, nil
}

func (m *AdminMemoryRepository) DeleteAPIKey(ctx context.Context, apiKey string) error {
	ctx = rctx.BuildContext(ctx, rctx.Context(rctx.Component, "DeleteAPIKey"), rctx.Context(rctx.StartTime, time.Now()))
	defer logger.LogComponentTime(ctx)

	_, err := m.apikeys.GetItem(apiKey)
	if err != nil {
		return err
	}

	m.apikeys.RemoveItem(apiKey)
	return nil
}

func (m *AdminMemoryRepository) StoreUserPolicy(ctx context.Context, p entity.UserPolicy) error {
	ctx = rctx.BuildContext(ctx, rctx.Context(rctx.Component, "StoreUserPolicy"), rctx.Context(rctx.StartTime, time.Now()))
	defer logger.LogComponentTime(ctx)

	m.userPolicies.StoreItem(p.UserId, &p)
	return nil
}

func (m *AdminMemoryRepository) GetUserPolicy(ctx context.Context, userId string) (*entity.UserPolicy, error) {
	ctx = rctx.BuildContext(ctx, rctx.Context(rctx.Component, "GetUserPolicy"), rctx.Context(rctx.StartTime, time.Now()))
	defer logger.LogComponentTime(ctx)

	item, err := m.userPolicies.GetItem(userId)
	if err != nil {
		return nil, err
	}

	return item.(*entity.UserPolicy), nil
}

func (m *AdminMemoryRepository) StoreInstancePolicy(ctx context.Context, p entity.InstancePolicy) error {
	ctx = rctx.BuildContext(ctx, rctx.Context(rctx.Component, "StoreInstancePolicy"), rctx.Context(rctx.StartTime, time.Now()))
	defer logger.LogComponentTime(ctx)

	id := fmt.Sprintf("%s-%s", p.Project, p.Instance)
	m.instancePolicies.StoreItem(id, &p)
	return nil
}

func (m *AdminMemoryRepository) StoreInstancePolicies(ctx context.Context, ps []entity.InstancePolicy) error {
	ctx = rctx.BuildContext(ctx, rctx.Context(rctx.Component, "StoreInstancePolicies"), rctx.Context(rctx.StartTime, time.Now()))
	defer logger.LogComponentTime(ctx)

	for _, p := range ps {
		m.StoreInstancePolicy(ctx, p)
	}

	return nil
}

func (m *AdminMemoryRepository) GetInstanceMethodPolicy(ctx context.Context, project, instance, methodName string) (*entity.MethodPolicy, error) {

	ctx = rctx.BuildContext(ctx, rctx.Context(rctx.Component, "GetInstanceMethodPolicy"), rctx.Context(rctx.StartTime, time.Now()))
	defer logger.LogComponentTime(ctx)

	id := fmt.Sprintf("%s-%s", project, instance)
	item, err := m.instancePolicies.GetItem(id)
	if err != nil {
		return nil, err
	}

	ip := item.(*entity.InstancePolicy)
	return ip.GetMethodByName(methodName), nil
}

func (m *AdminMemoryRepository) GetInstancePolicy(ctx context.Context, project, instance string) (*entity.InstancePolicy, error) {

	ctx = rctx.BuildContext(ctx, rctx.Context(rctx.Component, "GetInstancePolicy"), rctx.Context(rctx.StartTime, time.Now()))
	defer logger.LogComponentTime(ctx)

	id := fmt.Sprintf("%s-%s", project, instance)
	item, err := m.instancePolicies.GetItem(id)
	if err != nil {
		return nil, err
	}

	return item.(*entity.InstancePolicy), nil
}

func (m *AdminMemoryRepository) GetInstanceMethodPolicies(ctx context.Context, project, instance, methodCategory string) ([]entity.MethodPolicy, error) {

	ctx = rctx.BuildContext(ctx, rctx.Context(rctx.Component, "GetInstanceMethodPolicies"), rctx.Context(rctx.StartTime, time.Now()))
	defer logger.LogComponentTime(ctx)

	id := fmt.Sprintf("%s-%s", project, instance)
	item, err := m.instancePolicies.GetItem(id)
	if err != nil {
		return nil, err
	}

	ip := item.(*entity.InstancePolicy)
	return ip.GetMethodsByCategory(methodCategory), nil
}

func (m *AdminMemoryRepository) StoreAPIKeyPolicy(ctx context.Context, p entity.APIKeyPolicy) error {
	ctx = rctx.BuildContext(ctx, rctx.Context(rctx.Component, "StoreAPIKeyPolicy"), rctx.Context(rctx.StartTime, time.Now()))
	defer logger.LogComponentTime(ctx)

	m.apikeyPolicies.StoreItem(p.Key, &p)
	return nil
}

func (m *AdminMemoryRepository) GetAPIKeyPolicy(ctx context.Context, key string) (*entity.APIKeyPolicy, error) {
	ctx = rctx.BuildContext(ctx, rctx.Context(rctx.Component, "GetAPIKeyPolicy"), rctx.Context(rctx.StartTime, time.Now()))
	defer logger.LogComponentTime(ctx)

	item, err := m.apikeyPolicies.GetItem(key)
	if err != nil {
		return nil, err
	}

	return item.(*entity.APIKeyPolicy), nil
}
