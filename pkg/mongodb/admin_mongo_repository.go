package mongodb

import (
	"context"
	"github.com/zbitech/common/pkg/id"
	"github.com/zbitech/common/pkg/utils"
	"time"

	"github.com/zbitech/common/pkg/errs"
	"github.com/zbitech/common/pkg/logger"
	"github.com/zbitech/common/pkg/model/entity"
	"github.com/zbitech/common/pkg/rctx"
	"github.com/zbitech/common/pkg/vars"
	"github.com/zbitech/repo/internal/helper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type AdminMongoRepository struct {
	conn *MongoDBConnection
}

func NewAdminMongoRepository(conn *MongoDBConnection) *AdminMongoRepository {
	return &AdminMongoRepository{conn: conn}
}

func (m *AdminMongoRepository) RegisterUser(ctx context.Context, user *entity.User, pass *entity.UserPassword) error {
	ctx = rctx.BuildContext(ctx, rctx.Context(rctx.Component, "RegisterUser"), rctx.Context(rctx.StartTime, time.Now()))
	defer logger.LogComponentTime(ctx)

	userColl := m.conn.GetCollection(vars.MONGODB_NAME, MONGODB_COLL_USERS)

	filter := bson.M{"userid": user.UserId}
	userResult := userColl.FindOne(ctx, filter)

	userErr := userResult.Err()
	if userErr != nil {

		if userErr == mongo.ErrNoDocuments {
			user.Created = time.Now()
			user.Active = true
			user.LastUpdate = time.Now()

			iResult, err := userColl.InsertOne(ctx, user)
			if err != nil {
				logger.Errorf(ctx, "Error inserting user info - %s", err)
				return errs.ErrDBItemInsertFailed
			}

			if pass != nil {
				passColl := m.conn.GetCollection(vars.MONGODB_NAME, MONGODB_COLL_PASSWORD)
				_, err = passColl.InsertOne(ctx, pass)
				if err != nil {
					logger.Errorf(ctx, "Error setting user password - %s", err)
					return errs.ErrDBItemInsertFailed
				}
			}

			logger.Infof(ctx, "Inserted user with id %s", iResult.InsertedID)
			return nil

		}

		return handleMongoError(ctx, userErr)
	} else {
		return errs.ErrUserAlreadyExists
	}
}

func (m *AdminMongoRepository) DeactivateUser(ctx context.Context, userid string) error {
	ctx = rctx.BuildContext(ctx, rctx.Context(rctx.Component, "DeactivateUser"), rctx.Context(rctx.StartTime, time.Now()))
	defer logger.LogComponentTime(ctx)

	userColl := m.conn.GetCollection(vars.MONGODB_NAME, MONGODB_COLL_USERS)
	filter := bson.M{"userid": userid}
	updateFilter := bson.M{"$set": bson.M{"active": false, "lastupdate": time.Now()}}

	result := userColl.FindOneAndUpdate(ctx, filter, updateFilter)
	err := result.Err()

	if err != nil {
		return handleMongoError(ctx, err)
	}

	return nil
}

func (m *AdminMongoRepository) ReactivateUser(ctx context.Context, userid string) error {
	ctx = rctx.BuildContext(ctx, rctx.Context(rctx.Component, "ReactivateUser"), rctx.Context(rctx.StartTime, time.Now()))
	defer logger.LogComponentTime(ctx)

	userColl := m.conn.GetCollection(vars.MONGODB_NAME, MONGODB_COLL_USERS)
	filter := bson.M{"userid": userid}
	updateFilter := bson.M{"$set": bson.M{"active": true, "lastupdate": time.Now()}}

	result := userColl.FindOneAndUpdate(ctx, filter, updateFilter)
	err := result.Err()

	if err != nil {
		return handleMongoError(ctx, err)
	}

	return nil
}

//func (m *AdminMongoRepository) NewUser(ctx context.Context, userid, name, email, password string, role ztypes.Roles) (*entity.User, error) {
//	ctx = rctx.BuildContext(ctx, rctx.Context(rctx.Component, "NewUser"), rctx.Context(rctx.StartTime, time.Now()))
//	defer logger.LogComponentTime(ctx)
//
//	userColl := m.conn.GetCollection(vars.MONGODB_NAME, MONGODB_COLL_USERS)
//	passColl := m.conn.GetCollection(vars.MONGODB_NAME, MONGODB_COLL_PASSWORD)
//
//	filter := bson.M{"userid": userid}
//	userResult := userColl.FindOne(ctx, filter)
//
//	userErr := userResult.Err()
//	if userErr != nil {
//
//		if userErr == mongo.ErrNoDocuments {
//			hashedPass, err := fn.HashAndSaltPassword([]byte(password))
//			if err != nil {
//				logger.Errorf(ctx, "Error hashing password - %s", err)
//				return nil, errs.ErrMarshalFailed
//			}
//
//			user := entity.NewUser(userid, name, email, role)
//			userPass := entity.NewUserPassword(userid, hashedPass)
//
//			iResult, err := userColl.InsertOne(ctx, user)
//			if err != nil {
//				logger.Errorf(ctx, "Error inserting user info - %s", err)
//				return nil, errs.ErrDBItemInsertFailed
//			}
//
//			_, err = passColl.InsertOne(ctx, userPass)
//			if err != nil {
//				logger.Errorf(ctx, "Error setting user password - %s", err)
//				return nil, errs.ErrDBItemInsertFailed
//			}
//
//			logger.Infof(ctx, "Inserted user with id %s", iResult.InsertedID)
//			return user, nil
//
//		}
//
//		return nil, handleMongoError(ctx, userErr)
//	} else {
//		return nil, errs.ErrUserAlreadyExists
//	}
//}

func (m *AdminMongoRepository) GetUsers(ctx context.Context) []entity.User {

	ctx = rctx.BuildContext(ctx, rctx.Context(rctx.Component, "GetUsers"), rctx.Context(rctx.StartTime, time.Now()))
	defer logger.LogComponentTime(ctx)

	collection := m.conn.GetCollection(vars.MONGODB_NAME, MONGODB_COLL_USERS)
	items, err := collection.Find(ctx, bson.D{})
	if err != nil {
		logger.Errorf(ctx, "Got an error instead of list of users - %s", err)
		return nil
	}

	users := make([]entity.User, 0)
	if err = items.All(ctx, &users); err != nil {
		return nil
	}

	return users
}

func (m *AdminMongoRepository) GetUser(ctx context.Context, userId string) (*entity.User, error) {

	ctx = rctx.BuildContext(ctx, rctx.Context(rctx.Component, "GetUser"), rctx.Context(rctx.StartTime, time.Now()))
	defer logger.LogComponentTime(ctx)

	coll := m.conn.GetCollection(vars.MONGODB_NAME, MONGODB_COLL_USERS)

	logger.Infof(ctx, "Searching for user %s in repository", userId)
	result := coll.FindOne(ctx, bson.M{"userid": userId})

	if result.Err() != nil {
		return nil, handleMongoError(ctx, result.Err())
	}

	user := entity.User{}
	if err := result.Decode(&user); err != nil {
		return nil, errs.ErrMarshalFailed
	}

	return &user, nil
}

func (m *AdminMongoRepository) UpdateUser(ctx context.Context, user *entity.User) error {
	ctx = rctx.BuildContext(ctx, rctx.Context(rctx.Component, "UpdateUser"), rctx.Context(rctx.StartTime, time.Now()))
	defer logger.LogComponentTime(ctx)

	coll := m.conn.GetCollection(vars.MONGODB_NAME, MONGODB_COLL_USERS)
	filter := bson.M{"userid": user.UserId}

	user.LastUpdate = time.Now()
	result := coll.FindOneAndUpdate(ctx, filter, user)
	err := result.Err()

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return errs.ErrDBItemNotFound
		}

		return errs.ErrDBItemUpdateFailed
	}

	return nil
}

func (m *AdminMongoRepository) UpdatePassword(ctx context.Context, userid string, password *entity.UserPassword) error {

	ctx = rctx.BuildContext(ctx, rctx.Context(rctx.Component, "UpdatePassword"), rctx.Context(rctx.StartTime, time.Now()))
	defer logger.LogComponentTime(ctx)

	//	hashedPass, err := fn.HashAndSaltPassword([]byte(password))
	//	if err != nil {
	//		logger.Errorf(ctx, "Error hashing password - %s", err)
	//		return errs.ErrMarshalFailed
	//	}

	//	userPass := entity.NewUserPassword(userid, hashedPass)

	coll := m.conn.GetCollection(vars.MONGODB_NAME, MONGODB_COLL_USERS)
	filter := bson.M{"userid": userid}
	//	updateFilter := bson.M{"$set": bson.M{"password": userPass}}
	result := coll.FindOneAndUpdate(ctx, filter, password)
	err := result.Err()

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return errs.ErrDBItemNotFound
		}

		return errs.ErrDBItemUpdateFailed
	}

	return nil
}

func (m *AdminMongoRepository) AuthenticateUser(ctx context.Context, userId, password string) (*string, error) {

	ctx = rctx.BuildContext(ctx, rctx.Context(rctx.Component, "AuthenticateUser"), rctx.Context(rctx.StartTime, time.Now()))
	defer logger.LogComponentTime(ctx)

	userColl := m.conn.GetCollection(vars.MONGODB_NAME, MONGODB_COLL_USERS)
	passColl := m.conn.GetCollection(vars.MONGODB_NAME, MONGODB_COLL_PASSWORD)

	userResult := userColl.FindOne(ctx, bson.M{"userid": userId})
	if userResult.Err() != nil {
		return nil, handleMongoError(ctx, userResult.Err())
	}

	passResult := passColl.FindOne(ctx, bson.M{"userid": userId})
	if passResult.Err() != nil {
		return nil, handleMongoError(ctx, passResult.Err())
	}

	user := entity.User{}
	pass := entity.UserPassword{}

	if err := userResult.Decode(&user); err != nil {
		return nil, errs.ErrMarshalFailed
	}

	if err := passResult.Decode(&pass); err != nil {
		return nil, errs.ErrMarshalFailed
	}

	if id.ValidatePassword(pass.Password, []byte(password)) {
		return helper.GenerateJwtToken(user)
	}

	return nil, errs.ErrAuthFailed
}

func (m *AdminMongoRepository) GetAPIKeys(ctx context.Context, userId string) ([]string, error) {

	ctx = rctx.BuildContext(ctx, rctx.Context(rctx.Component, "GetAPIKey"), rctx.Context(rctx.StartTime, time.Now()))
	defer logger.LogComponentTime(ctx)

	opts := options.Find().SetProjection(bson.M{"key": 1})
	coll := m.conn.GetCollection(vars.MONGODB_NAME, MONGODB_COLL_APIKEY)
	items, err := coll.Find(ctx, bson.M{"userid": userId}, opts)
	if err != nil {
		return nil, handleMongoError(ctx, err)
	}

	keys := make([]string, 0)
	for items.Next(ctx) {
		var item = bson.M{}
		err = items.Decode(&item)
		if err == nil {
			logger.Debugf(ctx, "Item - %s", item["key"])
			keys = append(keys, item["key"].(string))
		} else {
			logger.Errorf(ctx, "Failed to marshal key")
		}
	}

	logger.Infof(ctx, "Items - %v", utils.MarshalObject(keys))

	return keys, nil
}

func (m *AdminMongoRepository) GetAPIKey(ctx context.Context, apiKey string) (*entity.APIKey, error) {
	ctx = rctx.BuildContext(ctx, rctx.Context(rctx.Component, "AdminMongoRepository.GetAPIKey"),
		rctx.Context(rctx.StartTime, time.Now()))
	defer logger.LogComponentTime(ctx)

	coll := m.conn.GetCollection(vars.MONGODB_NAME, MONGODB_COLL_APIKEY)
	result := coll.FindOne(ctx, bson.M{"key": apiKey})
	if result.Err() != nil {
		return nil, handleMongoError(ctx, result.Err())
	}

	item := entity.APIKey{}
	if err := result.Decode(&item); err != nil {
		return nil, errs.ErrMarshalFailed
	}

	return &item, nil
}

func (m *AdminMongoRepository) CreateAPIKey(ctx context.Context, user_id string) (*entity.APIKey, error) {

	ctx = rctx.BuildContext(ctx, rctx.Context(rctx.Component, "AdminMongoRepository.CreateAPIKey"),
		rctx.Context(rctx.StartTime, time.Now()))
	defer logger.LogComponentTime(ctx)

	coll := m.conn.GetCollection(vars.MONGODB_NAME, MONGODB_COLL_APIKEY)

	apikey := entity.NewAPIKey(user_id, vars.HOURS_IN_YEAR)

	_, err := coll.InsertOne(ctx, apikey)
	if err != nil {
		logger.Errorf(ctx, "Unable to store API key - %s", err)
		return nil, errs.ErrDBItemInsertFailed
	}

	return &apikey, nil
}

func (m *AdminMongoRepository) DeleteAPIKey(ctx context.Context, apiKey string) error {
	ctx = rctx.BuildContext(ctx, rctx.Context(rctx.Component, "AdminMongoRepository.DeleteAPIKey"),
		rctx.Context(rctx.StartTime, time.Now()))
	defer logger.LogComponentTime(ctx)

	coll := m.conn.GetCollection(vars.MONGODB_NAME, MONGODB_COLL_USERS)
	_, err := coll.DeleteOne(ctx, bson.M{"key": apiKey})
	if err != nil {
		return handleMongoError(ctx, err)
	}
	return nil
}

func (m *AdminMongoRepository) StoreUserPolicy(ctx context.Context, p entity.UserPolicy) error {
	ctx = rctx.BuildContext(ctx, rctx.Context(rctx.Component, "AdminMongoRepository.StoreUserPolicy"),
		rctx.Context(rctx.StartTime, time.Now()))
	defer logger.LogComponentTime(ctx)

	p.Updated = time.Now()
	coll := m.conn.GetCollection(vars.MONGODB_NAME, MONGODB_COLL_USER_POLICY)
	result := coll.FindOneAndReplace(ctx, bson.M{"userid": p.UserId}, p)
	err := result.Err()

	if err != nil {
		if err == mongo.ErrNoDocuments {
			_, err := coll.InsertOne(ctx, p)
			if err != nil {
				return handleMongoError(ctx, err)
			}
		} else {
			return handleMongoError(ctx, err)
		}
	}

	return nil
}

func (m *AdminMongoRepository) GetUserPolicy(ctx context.Context, userId string) (*entity.UserPolicy, error) {
	ctx = rctx.BuildContext(ctx, rctx.Context(rctx.Component, "GetUserPolicy"), rctx.Context(rctx.StartTime, time.Now()))
	defer logger.LogComponentTime(ctx)

	coll := m.conn.GetCollection(vars.MONGODB_NAME, MONGODB_COLL_USER_POLICY)
	result := coll.FindOne(ctx, bson.M{"userid": userId})
	err := result.Err()
	if err != nil {
		return nil, handleMongoError(ctx, err)
	}

	item := entity.UserPolicy{}
	if err := result.Decode(&item); err != nil {
		return nil, handleMongoError(ctx, err)
	}

	return &item, nil
}

func (m *AdminMongoRepository) StoreInstancePolicy(ctx context.Context, p entity.InstancePolicy) error {
	ctx = rctx.BuildContext(ctx, rctx.Context(rctx.Component, "StoreInstancePolicy"), rctx.Context(rctx.StartTime, time.Now()))
	defer logger.LogComponentTime(ctx)

	coll := m.conn.GetCollection(vars.MONGODB_NAME, MONGODB_COLL_INSTANCE_POLICY)
	result := coll.FindOneAndReplace(ctx, bson.M{"project": p.Project, "instance": p.Instance}, p)
	err := result.Err()

	if err != nil {
		if result.Err() == mongo.ErrNoDocuments {
			_, err := coll.InsertOne(ctx, p)
			if err != nil {
				return handleMongoError(ctx, err)
			}
		} else {
			return handleMongoError(ctx, err)
		}
	}

	return nil
}

func (m *AdminMongoRepository) StoreInstancePolicies(ctx context.Context, p []entity.InstancePolicy) error {
	ctx = rctx.BuildContext(ctx, rctx.Context(rctx.Component, "StoreInstancePolicies"), rctx.Context(rctx.StartTime, time.Now()))
	defer logger.LogComponentTime(ctx)

	coll := m.conn.GetCollection(vars.MONGODB_NAME, MONGODB_COLL_INSTANCE_POLICY)
	for _, i := range p {

		result := coll.FindOneAndReplace(ctx, bson.M{"project": i.Project, "instance": i.Instance}, i)
		err := result.Err()
		if err != nil {
			if result.Err() == mongo.ErrNoDocuments {
				_, err := coll.InsertOne(ctx, i)
				if err != nil {
					return handleMongoError(ctx, err)
				}
			} else {
				return handleMongoError(ctx, err)
			}
		}
	}

	return nil
}

func (m *AdminMongoRepository) GetInstanceMethodPolicy(ctx context.Context, project, instance, methodName string) (*entity.MethodPolicy, error) {

	ctx = rctx.BuildContext(ctx, rctx.Context(rctx.Component, "GetInstanceMethodPolicy"), rctx.Context(rctx.StartTime, time.Now()))
	defer logger.LogComponentTime(ctx)

	coll := m.conn.GetCollection(vars.MONGODB_NAME, MONGODB_COLL_INSTANCE_POLICY)
	result := coll.FindOne(ctx, bson.M{"project": project, "instance": instance})
	err := result.Err()

	if err != nil {
		return nil, handleMongoError(ctx, err)
	}

	item := entity.InstancePolicy{}
	if err = result.Decode(&item); err != nil {
		return nil, errs.ErrMarshalFailed
	}

	return item.GetMethodByName(methodName), nil
}

func (m *AdminMongoRepository) GetInstancePolicy(ctx context.Context, project, instance string) (*entity.InstancePolicy, error) {

	ctx = rctx.BuildContext(ctx, rctx.Context(rctx.Component, "AdminMongoRepository.GetInstancePolicy"), rctx.Context(rctx.StartTime, time.Now()))
	defer logger.LogComponentTime(ctx)

	coll := m.conn.GetCollection(vars.MONGODB_NAME, MONGODB_COLL_INSTANCE_POLICY)
	result := coll.FindOne(ctx, bson.M{"project": project, "instance": instance})
	err := result.Err()

	if err != nil {
		return nil, handleMongoError(ctx, err)
	}

	item := entity.InstancePolicy{}
	if err = result.Decode(&item); err != nil {
		return nil, handleMongoError(ctx, err)
	}

	return &item, nil
}

func (m *AdminMongoRepository) GetInstanceMethodPolicies(ctx context.Context, project, instance, methodCategory string) ([]entity.MethodPolicy, error) {

	ctx = rctx.BuildContext(ctx, rctx.Context(rctx.Component, "GetInstanceMethodPolicies"), rctx.Context(rctx.StartTime, time.Now()))
	defer logger.LogComponentTime(ctx)

	coll := m.conn.GetCollection(vars.MONGODB_NAME, MONGODB_COLL_INSTANCE_POLICY)
	results := coll.FindOne(ctx, bson.M{"project": project, "instance": instance})
	err := results.Err()

	if err != nil {
		return nil, handleMongoError(ctx, err)
	}

	item := entity.InstancePolicy{}
	if err = results.Decode(&item); err != nil {
		return nil, handleMongoError(ctx, err)
	}

	return item.GetMethodsByCategory(methodCategory), nil
}

func (m *AdminMongoRepository) StoreAPIKeyPolicy(ctx context.Context, p entity.APIKeyPolicy) error {
	ctx = rctx.BuildContext(ctx, rctx.Context(rctx.Component, "StoreAPIKeyPolicy"), rctx.Context(rctx.StartTime, time.Now()))
	defer logger.LogComponentTime(ctx)

	coll := m.conn.GetCollection(vars.MONGODB_NAME, MONGODB_COLL_APIKEY_POLICY)

	p.Updated = time.Now()
	result := coll.FindOneAndReplace(ctx, bson.M{"key": p.Key}, p)
	err := result.Err()

	if err != nil {
		if result.Err() == mongo.ErrNoDocuments {
			_, err := coll.InsertOne(ctx, p)
			if err != nil {
				return handleMongoError(ctx, err)
			}
		} else {
			return handleMongoError(ctx, err)
		}
	}

	return nil
}

func (m *AdminMongoRepository) GetAPIKeyPolicy(ctx context.Context, key string) (*entity.APIKeyPolicy, error) {
	ctx = rctx.BuildContext(ctx, rctx.Context(rctx.Component, "GetAPIKeyPolicy"), rctx.Context(rctx.StartTime, time.Now()))
	defer logger.LogComponentTime(ctx)

	coll := m.conn.GetCollection(vars.MONGODB_NAME, MONGODB_COLL_APIKEY_POLICY)
	result := coll.FindOne(ctx, bson.M{"key": key})
	err := result.Err()

	if err != nil {
		return nil, handleMongoError(ctx, err)
	}

	item := entity.APIKeyPolicy{}
	if err = result.Decode(&item); err != nil {
		return nil, errs.ErrMarshalFailed
	}

	return &item, nil
}
