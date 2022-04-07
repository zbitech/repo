package mongodb

import (
	"context"
	"fmt"
	"github.com/zbitech/common/pkg/model/config"
	"github.com/zbitech/common/pkg/utils"
	"time"

	"github.com/zbitech/common/pkg/errs"
	"github.com/zbitech/common/pkg/logger"
	"github.com/zbitech/common/pkg/model/entity"
	"github.com/zbitech/common/pkg/rctx"
	"github.com/zbitech/common/pkg/vars"

	//	"github.com/zbi/utils/internal/helper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	MONGODB_COLL_PROJECTS        = "projects"
	MONGODB_COLL_INSTANCES       = "instances"
	MONGODB_COLL_RESOURCES       = "k8s_resources"
	MONGODB_COLL_USERS           = "users"
	MONGODB_COLL_PASSWORD        = "password"
	MONGODB_COLL_INSTANCE_POLICY = "instance_policy"
	MONGODB_COLL_USER_POLICY     = "user_policy"
	MONGODB_COLL_APIKEY          = "apikeys"
	MONGODB_COLL_APIKEY_POLICY   = "apikey_policy"
	MONGODB_COLL_TEAMS           = "teams"
	MONGODB_COLL_TEAM_MEMBERS    = "team_members"
	//	MONGODB_COLL_ROLE            = "roles"
	//	MONGODB_COLL_GLOBAL_POLICY   = "global_policy"
	//	MONGODB_COLL_SUMMARY         = "summary"

	USER_INDEXES            = []MongoIndex{{Name: "userid", Order: 1, Unique: true}, {Name: "email", Order: 1, Unique: true}}
	PASS_INDEXES            = []MongoIndex{{Name: "userid", Order: 1, Unique: true}}
	APIKEY_INDEXES          = []MongoIndex{{Name: "userid", Order: 1, Unique: true}, {Name: "key", Order: 1, Unique: true}, {Name: "expires", Order: 1, Unique: false}}
	INSTANCE_POLICY_INDEXES = []MongoIndex{{Name: "project", Order: 1, Unique: false}, {Name: "instance", Order: 1, Unique: false}}
	USER_POLICY_INDEXES     = []MongoIndex{{Name: "userid", Order: 1, Unique: true}}
	APIKEY_POLICY_INDEXES   = []MongoIndex{{Name: "key", Order: 1, Unique: false}}
	PROJECT_INDEXES         = []MongoIndex{{Name: "name", Order: 1, Unique: true}, {Name: "owner", Order: 1, Unique: false}, {Name: "team", Order: 1, Unique: false}}
	INSTANCE_INDEXES        = []MongoIndex{{Name: "project", Order: 1, Unique: false}, {Name: "name", Order: 1, Unique: false}, {Name: "owner", Order: 1, Unique: false}, {Name: "type", Order: 1, Unique: false}}
	RESOURCE_INDEXES        = []MongoIndex{{Name: "project", Order: 1, Unique: false}, {Name: "instance", Order: 1, Unique: false}, {Name: "type", Order: 1, Unique: false}}
	TEAM_INDEXES            = []MongoIndex{{Name: "name", Order: 1, Unique: false}, {Name: "owner", Order: 1, Unique: false}}
	TEAM_MEMBER_INDEXES     = []MongoIndex{{Name: "team", Order: 1, Unique: false}, {Name: "email", Order: 1, Unique: false}, {Name: "key", Order: 1, Unique: true}}
	//	ROLE_INDEXES            = []MongoIndex{{Name: "name", Order: 1, Unique: true}}
	//	ORGANIZATION_INDEXES    = []MongoIndex{{Name: "name", Order: 1, Unique: true}}
	//	SUMMARY_INDEXES         = []MongoIndex{{Name: "type", Order: 1, Unique: false}}
)

type MongoIndex struct {
	//	Fields []struct {
	//		Name  string
	//		Order int
	//	}
	Name   string
	Order  int
	Unique bool
}

type MongoIndexFields struct {
	Name  string
	Order int
}

func handleMongoError(ctx context.Context, err error) error {
	logger.Errorf(ctx, "Database error - %s", err)
	if err == mongo.ErrNoDocuments {
		return errs.ErrDBItemNotFound
	}

	return errs.ErrDBError
}

func CreateCollection(ctx context.Context, db *mongo.Database, collName string, indexes []MongoIndex, errs []error) {

	opt := options.CreateCollection()
	err := db.CreateCollection(ctx, collName, opt)
	if err != nil {
		logger.Errorf(ctx, "Unable to Create collection %s - %s", collName, err)
		errs = append(errs, err)
	}

	coll := db.Collection(collName)

	for _, index := range indexes {
		m_index := mongo.IndexModel{
			Keys:    bson.M{index.Name: index.Order},
			Options: options.Index().SetUnique(index.Unique).SetName(index.Name)}
		_, err := coll.Indexes().CreateOne(ctx, m_index)
		if err != nil {
			logger.Errorf(ctx, "Unable to create index %s.%s - %s", collName, m_index.Keys, err)
			errs = append(errs, err)
		}
	}

}

func PurgeCollection(ctx context.Context, db *mongo.Database, collName string, errs []error) (int64, error) {

	coll := db.Collection(collName)
	result, err := coll.DeleteMany(ctx, bson.D{})
	if err != nil {
		logger.Errorf(ctx, "Unable to Purge collection %s - %s", collName, err)
		errs = append(errs, err)
		return 0, err
	}

	logger.Infof(ctx, "Deleted %d rows from collection %s", result.DeletedCount, collName)
	return result.DeletedCount, nil
}

func DropCollection(ctx context.Context, db *mongo.Database, collName string, errs []error) {
	coll := db.Collection(collName)
	err := coll.Drop(ctx)

	if err != nil {
		logger.Errorf(ctx, "Unable to Drop collection %s - %s", collName, err)
		errs = append(errs, err)
	} else {
		logger.Infof(ctx, "Dropped collection %s", collName)
	}
}

func CreateCollections(ctx context.Context, conn *MongoDBConnection) error {

	db := conn.GetDatabase(vars.MONGODB_NAME)

	errs := make([]error, 0)

	CreateCollection(ctx, db, MONGODB_COLL_USERS, USER_INDEXES, errs)
	CreateCollection(ctx, db, MONGODB_COLL_PASSWORD, PASS_INDEXES, errs)
	CreateCollection(ctx, db, MONGODB_COLL_APIKEY, APIKEY_INDEXES, errs)

	CreateCollection(ctx, db, MONGODB_COLL_PROJECTS, PROJECT_INDEXES, errs)
	CreateCollection(ctx, db, MONGODB_COLL_INSTANCES, INSTANCE_INDEXES, errs)
	CreateCollection(ctx, db, MONGODB_COLL_RESOURCES, RESOURCE_INDEXES, errs)

	CreateCollection(ctx, db, MONGODB_COLL_INSTANCE_POLICY, INSTANCE_POLICY_INDEXES, errs)
	CreateCollection(ctx, db, MONGODB_COLL_USER_POLICY, USER_POLICY_INDEXES, errs)
	CreateCollection(ctx, db, MONGODB_COLL_APIKEY_POLICY, APIKEY_POLICY_INDEXES, errs)

	CreateCollection(ctx, db, MONGODB_COLL_TEAMS, TEAM_INDEXES, errs)
	CreateCollection(ctx, db, MONGODB_COLL_TEAM_MEMBERS, TEAM_MEMBER_INDEXES, errs)

	if len(errs) > 0 {
		return fmt.Errorf("Create Error: %s", errs)
	}
	return nil
}

func PurgeCollections(ctx context.Context, conn *MongoDBConnection) error {

	db := conn.GetDatabase(vars.MONGODB_NAME)

	errs := make([]error, 0)

	PurgeCollection(ctx, db, MONGODB_COLL_USERS, errs)
	PurgeCollection(ctx, db, MONGODB_COLL_PASSWORD, errs)
	PurgeCollection(ctx, db, MONGODB_COLL_APIKEY, errs)

	PurgeCollection(ctx, db, MONGODB_COLL_PROJECTS, errs)
	PurgeCollection(ctx, db, MONGODB_COLL_INSTANCES, errs)
	PurgeCollection(ctx, db, MONGODB_COLL_RESOURCES, errs)

	PurgeCollection(ctx, db, MONGODB_COLL_INSTANCE_POLICY, errs)
	PurgeCollection(ctx, db, MONGODB_COLL_USER_POLICY, errs)
	PurgeCollection(ctx, db, MONGODB_COLL_APIKEY_POLICY, errs)

	PurgeCollection(ctx, db, MONGODB_COLL_TEAMS, errs)
	PurgeCollection(ctx, db, MONGODB_COLL_TEAM_MEMBERS, errs)

	if len(errs) > 0 {
		return fmt.Errorf("Purge Error: %s", errs)
	}
	return nil
}

func DropCollections(ctx context.Context, conn *MongoDBConnection) error {

	errs := make([]error, 0)

	db := conn.GetDatabase(vars.MONGODB_NAME)

	DropCollection(ctx, db, MONGODB_COLL_USERS, errs)
	DropCollection(ctx, db, MONGODB_COLL_PASSWORD, errs)
	DropCollection(ctx, db, MONGODB_COLL_APIKEY, errs)

	DropCollection(ctx, db, MONGODB_COLL_PROJECTS, errs)
	DropCollection(ctx, db, MONGODB_COLL_INSTANCES, errs)
	DropCollection(ctx, db, MONGODB_COLL_RESOURCES, errs)

	DropCollection(ctx, db, MONGODB_COLL_INSTANCE_POLICY, errs)
	DropCollection(ctx, db, MONGODB_COLL_USER_POLICY, errs)
	DropCollection(ctx, db, MONGODB_COLL_APIKEY_POLICY, errs)

	DropCollection(ctx, db, MONGODB_COLL_TEAMS, errs)
	DropCollection(ctx, db, MONGODB_COLL_TEAM_MEMBERS, errs)

	if len(errs) > 0 {
		return fmt.Errorf("Drop Error: %s", errs)
	}
	return nil
}

func LoadDatabase(ctx context.Context, conn *MongoDBConnection, create_users bool) {
	ctx = rctx.BuildContext(ctx, rctx.Context(rctx.Component, "LoadDatabase"),
		rctx.Context(rctx.StartTime, time.Now()))
	defer logger.LogComponentTime(ctx)

	if create_users {
		var adminConfig config.AdminConfig
		config_path := fmt.Sprintf("%s/admin.yaml", vars.ASSET_PATH_DIRECTORY)
		err := utils.ReadConfig(config_path, nil, &adminConfig)
		if err != nil {
			logger.Errorf(ctx, "Unable to load users & keys: %s", err)
		} else {

			userColl := conn.GetCollection(vars.MONGODB_NAME, MONGODB_COLL_USERS)
			polColl := conn.GetCollection(vars.MONGODB_NAME, MONGODB_COLL_USER_POLICY)
			passColl := conn.GetCollection(vars.MONGODB_NAME, MONGODB_COLL_PASSWORD)

			for _, user := range adminConfig.Users {
				if user.Memberships == nil {
					user.Memberships = make([]entity.UserTeam, 0)
				}
				user.Created = time.Now()
				user.LastUpdate = time.Now()
				user.Active = true

				userid := user.UserId
				password := adminConfig.Passwords[userid]
				userPass := entity.NewUserPassword(userid, password)
				userPolicy := entity.NewUserPolicy(user.UserId)

				userColl.InsertOne(ctx, user)
				passColl.InsertOne(ctx, userPass)
				polColl.InsertOne(ctx, userPolicy)
			}

			apikeyColl := conn.GetCollection(vars.MONGODB_NAME, MONGODB_COLL_APIKEY)
			keyPolColl := conn.GetCollection(vars.MONGODB_NAME, MONGODB_COLL_APIKEY_POLICY)

			for _, apikey := range adminConfig.Keys {
				apikey.Created = time.Now()
				apikey.Expires = apikey.Created.Add(time.Hour * 8760)
				keyPolicy := entity.NewAPIKeyPolicy(apikey.Key, true)

				apikeyColl.InsertOne(ctx, apikey)
				keyPolColl.InsertOne(ctx, keyPolicy)
			}

			teamColl := conn.GetCollection(vars.MONGODB_NAME, MONGODB_COLL_TEAMS)
			for _, team := range adminConfig.Teams {
				team.Created = time.Now()
				team.LastUpdate = time.Now()
				teamColl.InsertOne(ctx, team)
			}

			mbrColl := conn.GetCollection(vars.MONGODB_NAME, MONGODB_COLL_TEAM_MEMBERS)
			for _, mbr := range adminConfig.Members {
				mbr.CreatedOn = time.Now()
				mbr.LastUpdate = time.Now()
				mbrColl.InsertOne(ctx, mbr)
			}
		}
	}

	//	globalPolicy, err := helper.GetGlobalPolicy(vars.ASSET_PATH_DIRECTORY)
	//	if err != nil {
	//		logger.Fatalf(ctx, "Unable to load global policy - %s", err)
	//	}

	//	policy_coll := conn.GetCollection(vars.MONGODB_NAME, MONGODB_COLL_GLOBAL_POLICY)
	//	policy_coll.InsertOne(ctx, globalPolicy)
}

func CountDocuments(ctx context.Context, conn *MongoDBConnection, database, collectionName string, filter bson.M) int {

	coll := conn.GetCollection(database, collectionName)
	count, err := coll.CountDocuments(ctx, filter)
	if err != nil {
		logger.Errorf(ctx, "Count of %s.%s with %s failed - %s", database, collectionName, filter, err)
		//		return -1
	}

	return int(count)
}
