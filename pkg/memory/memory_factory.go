package memory

import (
	"context"

	"github.com/zbitech/common/interfaces"
)

type MemoryRepositoryFactory struct {
	project interfaces.ProjectRepositoryIF
	//	team    interfaces.TeamRepositoryIF
	admin interfaces.AdminRepositoryIF
}

func NewMemoryRepositoryFactory() interfaces.RepositoryFactoryIF {
	return &MemoryRepositoryFactory{}
}

func (r *MemoryRepositoryFactory) Init(ctx context.Context, create_db, load_db bool) error {

	return nil
}

func (r *MemoryRepositoryFactory) OpenConnection(ctx context.Context) error {
	return nil
}

func (r *MemoryRepositoryFactory) CloseConnection(ctx context.Context) error {
	return nil
}

func (r *MemoryRepositoryFactory) GetProjectRepository() interfaces.ProjectRepositoryIF {
	return nil
}

//func (r *MemoryRepositoryFactory) GetTeamRepository() interfaces.TeamRepositoryIF {
//	return nil
//}

func (r *MemoryRepositoryFactory) GetAdminRepository() interfaces.AdminRepositoryIF {
	return nil
}

func (r *MemoryRepositoryFactory) CreateDatabase(ctx context.Context, purge, load bool) error {

	return nil
}
