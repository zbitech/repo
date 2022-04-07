package memory

import (
	"context"
	"fmt"
	"github.com/zbitech/common/pkg/errs"
	mem "github.com/zbitech/common/pkg/memory"
	"time"

	"github.com/zbitech/common/pkg/logger"
	"github.com/zbitech/common/pkg/model/entity"
	"github.com/zbitech/common/pkg/rctx"
)

type ProjectMemoryRepository struct {
	projects  *mem.MemoryStore
	instances *mem.MemoryStore
	resources *mem.MemoryStore
	summaries *mem.MemoryStore
}

func NewProjectMemoryRepository() *ProjectMemoryRepository {
	return &ProjectMemoryRepository{
		projects:  mem.NewMemoryStore(),
		instances: mem.NewMemoryStore(),
		resources: mem.NewMemoryStore(),
		summaries: mem.NewMemoryStore(),
	}
}

func (m *ProjectMemoryRepository) GetProjects(ctx context.Context) ([]entity.Project, error) {

	ctx = rctx.BuildContext(ctx, rctx.Context(rctx.Component, "GetProjects"), rctx.Context(rctx.StartTime, time.Now()))
	defer logger.LogComponentTime(ctx)

	items := m.projects.GetItems()
	projects := make([]entity.Project, 0)
	for _, item := range items {
		project := item.(*entity.Project)
		projects = append(projects, *project)
	}

	return projects, nil
}

func (m *ProjectMemoryRepository) CreateProject(ctx context.Context, project *entity.Project) error {

	ctx = rctx.BuildContext(ctx, rctx.Context(rctx.Component, "CreateProject"), rctx.Context(rctx.StartTime, time.Now()))
	defer logger.LogComponentTime(ctx)

	_, err := m.projects.GetItem(project.Name)
	if err != nil {
		m.projects.StoreItem(project.Name, project)
		m.resources.StoreItem(project.Name, make([]*entity.KubernetesResource, 0))
		return nil
	}

	return errs.ErrDBKeyAlreadyExists
}

func (m *ProjectMemoryRepository) UpdateProject(ctx context.Context, project *entity.Project) error {

	ctx = rctx.BuildContext(ctx, rctx.Context(rctx.Component, "UpdateProject"), rctx.Context(rctx.StartTime, time.Now()))
	defer logger.LogComponentTime(ctx)

	m.projects.StoreItem(project.Name, project)
	return nil
}

func (m *ProjectMemoryRepository) GetProject(ctx context.Context, name string) (*entity.Project, error) {

	ctx = rctx.BuildContext(ctx, rctx.Context(rctx.Component, "GetProjectByName"), rctx.Context(rctx.StartTime, time.Now()))
	defer logger.LogComponentTime(ctx)

	item, err := m.projects.GetItem(name)
	if err != nil {
		return nil, err
	}

	return item.(*entity.Project), nil
}

func (m *ProjectMemoryRepository) GetProjectsByOwner(ctx context.Context, owner string) ([]entity.Project, error) {
	ctx = rctx.BuildContext(ctx, rctx.Context(rctx.Component, "GetProjectsByOwner"), rctx.Context(rctx.StartTime, time.Now()))
	defer logger.LogComponentTime(ctx)

	items := m.projects.GetItems()
	projects := make([]entity.Project, 0)
	for _, item := range items {
		project := item.(*entity.Project)
		if project.Owner == owner {
			projects = append(projects, *project)
		}
	}

	return projects, nil
}

func (m *ProjectMemoryRepository) GetProjectsByTeam(ctx context.Context, team string) ([]entity.Project, error) {
	ctx = rctx.BuildContext(ctx, rctx.Context(rctx.Component, "GetProjectsByTeam"), rctx.Context(rctx.StartTime, time.Now()))
	defer logger.LogComponentTime(ctx)

	items := m.projects.GetItems()
	projects := make([]entity.Project, 0)
	for _, item := range items {
		project := item.(*entity.Project)
		if project.TeamId == team {
			projects = append(projects, *project)
		}
	}

	return projects, nil
}

func (m *ProjectMemoryRepository) UpdateProjectStatus(ctx context.Context, name string, status string) error {

	ctx = rctx.BuildContext(ctx, rctx.Context(rctx.Component, "UpdateProjectStatus"), rctx.Context(rctx.StartTime, time.Now()))
	defer logger.LogComponentTime(ctx)

	item, err := m.projects.GetItem(name)
	if err != nil {
		return err
	}

	project := item.(*entity.Project)
	project.Status = status
	project.Timestamp = time.Now()

	m.projects.StoreItem(project.Name, project)

	return nil
}

func (m *ProjectMemoryRepository) GetInstances(ctx context.Context) ([]entity.Instance, error) {

	ctx = rctx.BuildContext(ctx, rctx.Context(rctx.Component, "GetInstances"), rctx.Context(rctx.StartTime, time.Now()))
	defer logger.LogComponentTime(ctx)

	items := m.instances.GetItems()
	instances := make([]entity.Instance, 0)
	for _, item := range items {
		instance := item.(*entity.Instance)
		instances = append(instances, *instance)
	}

	return instances, nil
}

func (m *ProjectMemoryRepository) CreateInstance(ctx context.Context, instance *entity.Instance) error {

	ctx = rctx.BuildContext(ctx, rctx.Context(rctx.Component, "CreateInstance"), rctx.Context(rctx.StartTime, time.Now()))
	defer logger.LogComponentTime(ctx)

	key := fmt.Sprintf("%s-%s", instance.Project, instance.Name)
	_, err := m.instances.GetItem(key)
	if err != nil {
		m.instances.StoreItem(key, instance)
		m.resources.StoreItem(key, make([]*entity.KubernetesResource, 0))
		return nil
	}

	return errs.ErrDBKeyAlreadyExists
}

func (m *ProjectMemoryRepository) GetInstance(ctx context.Context, project, name string) (*entity.Instance, error) {

	ctx = rctx.BuildContext(ctx, rctx.Context(rctx.Component, "GetInstanceByName"), rctx.Context(rctx.StartTime, time.Now()))
	defer logger.LogComponentTime(ctx)

	key := fmt.Sprintf("%s-%s", project, name)
	item, err := m.instances.GetItem(key)
	if err != nil {
		return nil, err
	}

	return item.(*entity.Instance), nil
}

func (m *ProjectMemoryRepository) GetInstancesByProject(ctx context.Context, project string) ([]entity.Instance, error) {
	ctx = rctx.BuildContext(ctx, rctx.Context(rctx.Component, "GetInstancesByProject"), rctx.Context(rctx.StartTime, time.Now()))
	defer logger.LogComponentTime(ctx)

	items := m.instances.GetItems()
	instances := make([]entity.Instance, 0)
	for _, item := range items {
		instance := item.(*entity.Instance)
		if instance.Project == project {
			instances = append(instances, *instance)
		}
	}

	return instances, nil
}

func (m *ProjectMemoryRepository) GetInstancesByOwner(ctx context.Context, owner string) ([]entity.Instance, error) {

	ctx = rctx.BuildContext(ctx, rctx.Context(rctx.Component, "GetInstancesByOwner"), rctx.Context(rctx.StartTime, time.Now()))
	defer logger.LogComponentTime(ctx)

	items := m.instances.GetItems()
	instances := make([]entity.Instance, 0)
	for _, item := range items {
		instance := item.(*entity.Instance)
		if instance.Owner == owner {
			instances = append(instances, *instance)
		}
	}

	return instances, nil
}

func (m *ProjectMemoryRepository) UpdateInstanceStatus(ctx context.Context, project, name string, status string) error {

	ctx = rctx.BuildContext(ctx, rctx.Context(rctx.Component, "UpdateInstanceStatus"), rctx.Context(rctx.StartTime, time.Now()))
	defer logger.LogComponentTime(ctx)

	key := fmt.Sprintf("%s-%s", project, name)
	item, err := m.instances.GetItem(key)
	if err != nil {
		return err
	}

	instance := item.(*entity.Instance)
	instance.Status = status
	instance.Timestamp = time.Now()

	m.projects.StoreItem(key, project)

	return nil
}

func (m *ProjectMemoryRepository) GetProjectResources(ctx context.Context, project string) ([]entity.KubernetesResource, error) {

	ctx = rctx.BuildContext(ctx, rctx.Context(rctx.Component, "GetProjectResources"), rctx.Context(rctx.StartTime, time.Now()))
	defer logger.LogComponentTime(ctx)

	items, err := m.resources.GetItem(project)
	if err != nil {
		logger.Errorf(ctx, "No resources found - %s", err)
		return []entity.KubernetesResource{}, nil
	}

	array := items.([]interface{})
	projRsc := make([]entity.KubernetesResource, 0)
	for _, item := range array {
		rsc := item.(*entity.KubernetesResource)
		projRsc = append(projRsc, *rsc)
	}

	return projRsc, nil
}

func (m *ProjectMemoryRepository) SaveProjectResource(ctx context.Context, project string, resource *entity.KubernetesResource) error {

	ctx = rctx.BuildContext(ctx, rctx.Context(rctx.Component, "SaveProjectResource"), rctx.Context(rctx.StartTime, time.Now()))
	defer logger.LogComponentTime(ctx)

	items, err := m.resources.GetItem(project)
	if err != nil {
		logger.Errorf(ctx, "No resources found - %s", err)
		return nil
	}

	array := items.([]interface{})
	for _, item := range array {
		rsc := item.(*entity.KubernetesResource)
		if rsc.Id == resource.Id {
			rsc.State = resource.State
			rsc.Timestamp = resource.Timestamp
			return nil
		}
	}

	m.resources.AddItem(project, resource)
	return nil
}

func (m *ProjectMemoryRepository) SaveProjectResources(ctx context.Context, project string, resources []entity.KubernetesResource) error {

	ctx = rctx.BuildContext(ctx, rctx.Context(rctx.Component, "SaveProjectResources"), rctx.Context(rctx.StartTime, time.Now()))
	defer logger.LogComponentTime(ctx)

	for _, resource := range resources {
		if err := m.SaveProjectResource(ctx, project, &resource); err != nil {
			return err
		}
	}

	return nil
}

func (m *ProjectMemoryRepository) GetInstanceResources(ctx context.Context, project, instance string) ([]entity.KubernetesResource, error) {

	ctx = rctx.BuildContext(ctx, rctx.Context(rctx.Component, "GetInstanceResources"), rctx.Context(rctx.StartTime, time.Now()))
	defer logger.LogComponentTime(ctx)

	key := fmt.Sprintf("%s-%s", project, instance)
	items, err := m.resources.GetItem(key)
	if err != nil {
		logger.Errorf(ctx, "No resources found - %s", err)
		return []entity.KubernetesResource{}, nil
	}

	array := items.([]interface{})
	resources := make([]entity.KubernetesResource, 0)
	for _, item := range array {
		rsc := item.(*entity.KubernetesResource)
		resources = append(resources, *rsc)
	}

	return resources, nil
}

func (m *ProjectMemoryRepository) SaveInstanceResource(ctx context.Context, project, instance string, resource *entity.KubernetesResource) error {

	ctx = rctx.BuildContext(ctx, rctx.Context(rctx.Component, "SaveInstanceResource"), rctx.Context(rctx.StartTime, time.Now()))
	defer logger.LogComponentTime(ctx)

	key := fmt.Sprintf("%s-%s", project, instance)
	items, err := m.resources.GetItem(key)
	if err != nil {
		logger.Errorf(ctx, "No resources found - %s", err)
		return nil
	}

	array := items.([]interface{})
	for _, item := range array {
		rsc := item.(*entity.KubernetesResource)
		if rsc.Id == resource.Id {
			rsc.State = resource.State
			rsc.Timestamp = resource.Timestamp
			return nil
		}
	}

	m.resources.AddItem(project, key)
	return nil
}

func (m *ProjectMemoryRepository) SaveInstanceResources(ctx context.Context, project, instance string, resources []entity.KubernetesResource) error {

	ctx = rctx.BuildContext(ctx, rctx.Context(rctx.Component, "SaveInstanceResources"), rctx.Context(rctx.StartTime, time.Now()))
	defer logger.LogComponentTime(ctx)

	for _, resource := range resources {
		if err := m.SaveInstanceResource(ctx, project, instance, &resource); err != nil {
			return err
		}
	}

	return nil
}

func (m *ProjectMemoryRepository) GetUserSummary(ctx context.Context, userId string) (*entity.ResourceSummary, error) {

	ctx = rctx.BuildContext(ctx, rctx.Context(rctx.Component, "GetUserSummary"), rctx.Context(rctx.StartTime, time.Now()))
	defer logger.LogComponentTime(ctx)

	return nil, nil
}
