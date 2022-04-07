package repo

import (
	"context"
	"fmt"
	"github.com/zbitech/repo/pkg/iam"
	"github.com/zbitech/repo/pkg/memory"
	"github.com/zbitech/repo/pkg/mongodb"

	"github.com/zbitech/common/interfaces"
	"github.com/zbitech/common/pkg/logger"
	"github.com/zbitech/common/pkg/vars"
)

func NewRepositoryFactory(ctx context.Context) interfaces.RepositoryFactoryIF {

	dbFactory := vars.AppConfig.Database.Factory

	logger.Infof(ctx, "Creating db factory for %s", dbFactory)

	switch dbFactory {
	case "mongo":
		return mongodb.NewMongoRepositoryFactory()
	case "memory":
		return memory.NewMemoryRepositoryFactory()
	}

	panic(fmt.Errorf("Unknown database type %s", dbFactory))
}

func NewAuthorizationFactory(ctx context.Context) interfaces.AuthorizationFactoryIF {
	return iam.NewBasicAuthorizationFactory()
}
