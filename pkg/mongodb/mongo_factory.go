package mongodb

import (
	"context"
	"time"

	"github.com/zbitech/common/interfaces"
	"github.com/zbitech/common/pkg/logger"
	"github.com/zbitech/common/pkg/rctx"
	"github.com/zbitech/common/pkg/vars"
)

type MongoRepositoryFactory struct {
	project interfaces.ProjectRepositoryIF
	admin   interfaces.AdminRepositoryIF
	conn    *MongoDBConnection
}

func NewMongoRepositoryFactory() interfaces.RepositoryFactoryIF {
	return &MongoRepositoryFactory{}
}

func (r *MongoRepositoryFactory) Init(ctx context.Context, create_db, load_db bool) error {

	ctx = rctx.BuildContext(ctx, rctx.Context(rctx.Component, "MongoRepositoryFactory"), rctx.Context(rctx.StartTime, time.Now()))
	defer logger.LogComponentTime(ctx)

	r.conn = NewMongoDBConnection(vars.AppConfig.Database.Mongodb.Url)
	r.conn.OpenConnection(ctx)
	if create_db {
		logger.Infof(ctx, "Creating Mongo DB documents")
		r.CreateDatabase(ctx, true, load_db)
	}

	r.project = NewProjectMongoRepository(r.conn)
	r.admin = NewAdminMongoRepository(r.conn)

	//bootstrap connection to load users. TODO - determine best way to bootstrap database pre-loading
	r.admin.GetUsers(ctx)

	return nil
}

func (r *MongoRepositoryFactory) OpenConnection(ctx context.Context) error {
	return r.conn.OpenConnection(ctx)
}

func (r *MongoRepositoryFactory) CloseConnection(ctx context.Context) error {
	r.conn.CloseConnection(ctx)
	return nil
}

func (r *MongoRepositoryFactory) GetProjectRepository() interfaces.ProjectRepositoryIF {
	return r.project
}

//func (r *MongoRepositoryFactory) GetTeamRepository() interfaces.TeamRepositoryIF {
//	return r.team
//}

func (r *MongoRepositoryFactory) GetAdminRepository() interfaces.AdminRepositoryIF {
	return r.admin
}

func (r *MongoRepositoryFactory) CreateDatabase(ctx context.Context, purge, load bool) error {
	var err error

	//	r.conn.OpenConnection(ctx)
	//	defer r.conn.CloseConnection(ctx)

	if purge {
		err = PurgeCollections(ctx, r.conn)
		if err != nil {
			return err
		}
	}

	err = CreateCollections(ctx, r.conn)
	if err != nil {
		return err
	}

	if load {
		LoadDatabase(ctx, r.conn, true)
	}
	return nil
}
