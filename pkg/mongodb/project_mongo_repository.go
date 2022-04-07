package mongodb

import (
	"context"
	"fmt"
	"github.com/zbitech/common/pkg/utils"
	"time"

	"github.com/zbitech/common/pkg/errs"
	"github.com/zbitech/common/pkg/logger"
	"github.com/zbitech/common/pkg/model/entity"
	"github.com/zbitech/common/pkg/rctx"
	"github.com/zbitech/common/pkg/vars"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type ProjectMongoRepository struct {
	conn *MongoDBConnection
}

func NewProjectMongoRepository(conn *MongoDBConnection) *ProjectMongoRepository {
	return &ProjectMongoRepository{conn: conn}
}

func (m *ProjectMongoRepository) GetProjects(ctx context.Context) ([]entity.Project, error) {

	ctx = rctx.BuildContext(ctx, rctx.Context(rctx.Component, "GetProjects"), rctx.Context(rctx.StartTime, time.Now()))
	defer logger.LogComponentTime(ctx)

	collection := m.conn.GetCollection(vars.MONGODB_NAME, MONGODB_COLL_PROJECTS)

	items, err := collection.Find(ctx, bson.D{})
	if err != nil {
		logger.Errorf(ctx, "Got an error instead of list - %s", err)
		return nil, fmt.Errorf("Unable to get projects")
	}

	var projects []entity.Project
	if err := items.All(ctx, &projects); err != nil {
		return nil, errs.ErrMarshalFailed
	}

	return projects, nil
}

func (m *ProjectMongoRepository) CreateProject(ctx context.Context, project *entity.Project) error {

	ctx = rctx.BuildContext(ctx, rctx.Context(rctx.Component, "CreateProject"), rctx.Context(rctx.StartTime, time.Now()))
	defer logger.LogComponentTime(ctx)

	collection := m.conn.GetCollection(vars.MONGODB_NAME, MONGODB_COLL_PROJECTS)

	logger.Infof(ctx, "Saving project %s - %s", project.GetName(), utils.MarshalObject(project))

	result, err := collection.InsertOne(ctx, project)
	if err != nil {
		return errs.ErrDBItemInsertFailed
	}

	logger.Infof(ctx, "Saved project %s - %s", result.InsertedID, utils.MarshalObject(result))
	return nil
}

func (m *ProjectMongoRepository) UpdateProject(ctx context.Context, project *entity.Project) error {

	ctx = rctx.BuildContext(ctx, rctx.Context(rctx.Component, "CreateProject"), rctx.Context(rctx.StartTime, time.Now()))
	defer logger.LogComponentTime(ctx)

	collection := m.conn.GetCollection(vars.MONGODB_NAME, MONGODB_COLL_PROJECTS)

	logger.Infof(ctx, "Updating project %s - %s", project.GetName(), utils.MarshalObject(project))
	filter := bson.M{"name": project.Name}
	result, err := collection.UpdateOne(ctx, filter, project)
	if err != nil {
		logger.Errorf(ctx, "Project update failed - %s", err)
		return errs.ErrDBItemUpdateFailed
	}

	logger.Infof(ctx, "Updated %d project", result.MatchedCount)

	if result.MatchedCount != 1 {
		return errs.ErrDBItemUpdateFailed
	}

	return nil
}

func (m *ProjectMongoRepository) GetProject(ctx context.Context, name string) (*entity.Project, error) {

	ctx = rctx.BuildContext(ctx, rctx.Context(rctx.Component, "GetProjectByName"), rctx.Context(rctx.StartTime, time.Now()))
	defer logger.LogComponentTime(ctx)

	collection := m.conn.GetCollection(vars.MONGODB_NAME, MONGODB_COLL_PROJECTS)
	var project entity.Project
	single_result := collection.FindOne(ctx, bson.M{"name": name})
	if single_result.Err() != nil {
		return nil, single_result.Err()
	}

	single_result.Decode(&project)

	return &project, nil
}

func (m *ProjectMongoRepository) GetProjectsByOwner(ctx context.Context, owner string) ([]entity.Project, error) {
	ctx = rctx.BuildContext(ctx, rctx.Context(rctx.Component, "GetProjectsByOwner"), rctx.Context(rctx.StartTime, time.Now()))
	defer logger.LogComponentTime(ctx)

	collection := m.conn.GetCollection(vars.MONGODB_NAME, MONGODB_COLL_PROJECTS)

	projects := []entity.Project{}
	items, err := collection.Find(ctx, bson.M{"owner": owner})
	if err != nil {
		return nil, err
	}

	if err = items.All(ctx, &projects); err != nil {
		return nil, errs.ErrMarshalFailed
	}

	return projects, nil
}

func (m *ProjectMongoRepository) GetProjectsByTeam(ctx context.Context, team string) ([]entity.Project, error) {
	ctx = rctx.BuildContext(ctx, rctx.Context(rctx.Component, "GetProjectsByTeam"), rctx.Context(rctx.StartTime, time.Now()))
	defer logger.LogComponentTime(ctx)

	filter := bson.M{"team": team}
	coll := m.conn.GetCollection(vars.MONGODB_NAME, MONGODB_COLL_PROJECTS)
	items, err := coll.Find(ctx, filter)
	if err != nil {
		return nil, handleMongoError(ctx, err)
	}

	projects := []entity.Project{}
	if err = items.All(ctx, &projects); err != nil {
		return nil, errs.ErrMarshalFailed
	}

	return projects, nil
}

func (m *ProjectMongoRepository) UpdateProjectStatus(ctx context.Context, project string, status string) error {

	ctx = rctx.BuildContext(ctx, rctx.Context(rctx.Component, "ProjectMongoRepository.UpdateProjectStatus"), rctx.Context(rctx.StartTime, time.Now()))
	defer logger.LogComponentTime(ctx)

	collection := m.conn.GetCollection(vars.MONGODB_NAME, MONGODB_COLL_PROJECTS)

	filter := bson.M{"project": project}
	updateFilter := bson.M{"$set": bson.M{"status": status, "timestamp": time.Now()}}
	result := collection.FindOneAndUpdate(ctx, filter, updateFilter)
	err := result.Err()

	if err != nil {
		return err
	}

	return nil
}

func (m *ProjectMongoRepository) GetInstances(ctx context.Context) ([]entity.Instance, error) {

	ctx = rctx.BuildContext(ctx, rctx.Context(rctx.Component, "GetInstances"), rctx.Context(rctx.StartTime, time.Now()))
	defer logger.LogComponentTime(ctx)

	collection := m.conn.GetCollection(vars.MONGODB_NAME, MONGODB_COLL_INSTANCES)
	items, err := collection.Find(context.Background(), bson.D{})
	if err != nil {
		return nil, err
	}

	instances := []entity.Instance{}

	for items.Next(ctx) {

		var rawData bson.Raw
		err := items.Decode(&rawData)
		if err != nil {
			return nil, err
		}

		instance, err := vars.ManagerFactory.GetProjectDataManager(ctx).UnmarshalBSONInstance(ctx, rawData)
		if err != nil {
			return nil, err
		}
		instances = append(instances, *instance)
	}

	logger.Infof(ctx, "Returning Instances: %s", utils.MarshalObject(instances))
	return instances, nil
}

func (m *ProjectMongoRepository) CreateInstance(ctx context.Context, instance *entity.Instance) error {

	ctx = rctx.BuildContext(ctx, rctx.Context(rctx.Component, "CreateInstance"), rctx.Context(rctx.StartTime, time.Now()))
	defer logger.LogComponentTime(ctx)

	collection := m.conn.GetCollection(vars.MONGODB_NAME, MONGODB_COLL_INSTANCES)
	result, err := collection.InsertOne(ctx, instance)
	if err != nil {
		return err
	}

	logger.Infof(ctx, "Saved instance %s - %s", result.InsertedID, utils.MarshalObject(result))

	return nil
}

func (m *ProjectMongoRepository) GetInstance(ctx context.Context, project, name string) (*entity.Instance, error) {

	ctx = rctx.BuildContext(ctx, rctx.Context(rctx.Component, "GetInstanceByName"), rctx.Context(rctx.StartTime, time.Now()))
	defer logger.LogComponentTime(ctx)

	collection := m.conn.GetCollection(vars.MONGODB_NAME, MONGODB_COLL_INSTANCES)

	var rawData bson.Raw
	result := collection.FindOne(ctx, bson.M{"name": name, "project": project})
	if result.Err() != nil {
		logger.Errorf(ctx, "Unable to find instance %s in project %s - %s", name, project, result.Err())
		return nil, errs.ErrDBItemNotFound
	}

	result.Decode(&rawData)
	instance, err := vars.ManagerFactory.GetProjectDataManager(ctx).UnmarshalBSONInstance(ctx, rawData)
	if err != nil {
		return nil, err
	}

	return instance, nil
}

func (m *ProjectMongoRepository) GetInstancesByProject(ctx context.Context, project string) ([]entity.Instance, error) {
	ctx = rctx.BuildContext(ctx, rctx.Context(rctx.Component, "GetInstancesByProject"), rctx.Context(rctx.StartTime, time.Now()))
	defer logger.LogComponentTime(ctx)

	instances := []entity.Instance{}

	collection := m.conn.GetCollection(vars.MONGODB_NAME, MONGODB_COLL_INSTANCES)
	items, err := collection.Find(ctx, bson.M{"project": project})
	if err != nil {
		return nil, err
	}

	for items.Next(ctx) {
		var rawData bson.Raw
		err := items.Decode(&rawData)
		if err != nil {
			return nil, err
		}

		instance, err := vars.ManagerFactory.GetProjectDataManager(ctx).UnmarshalBSONInstance(ctx, rawData)
		if err != nil {
			return nil, err
		}

		instances = append(instances, *instance)
	}

	return instances, nil
}

func (m *ProjectMongoRepository) GetInstancesByOwner(ctx context.Context, owner string) ([]entity.Instance, error) {

	ctx = rctx.BuildContext(ctx, rctx.Context(rctx.Component, "GetInstancesByOwner"), rctx.Context(rctx.StartTime, time.Now()))
	defer logger.LogComponentTime(ctx)

	collection := m.conn.GetCollection(vars.MONGODB_NAME, MONGODB_COLL_INSTANCES)
	items, err := collection.Find(ctx, bson.M{"owner": owner})
	if err != nil {
		return nil, err
	}

	var instances []entity.Instance
	for items.Next(ctx) {
		var rawData bson.Raw
		err := items.Decode(&rawData)
		if err != nil {
			return nil, err
		}

		instance, err := vars.ManagerFactory.GetProjectDataManager(ctx).UnmarshalBSONInstance(ctx, rawData)
		if err != nil {
			return nil, err
		}

		instances = append(instances, *instance)
	}

	return instances, nil
}

func (m *ProjectMongoRepository) UpdateInstanceStatus(ctx context.Context, project, name, status string) error {

	ctx = rctx.BuildContext(ctx, rctx.Context(rctx.Component, "UpdateInstanceStatus"), rctx.Context(rctx.StartTime, time.Now()))
	defer logger.LogComponentTime(ctx)

	collection := m.conn.GetCollection(vars.MONGODB_NAME, MONGODB_COLL_INSTANCES)
	filter := bson.M{"project": project, "instance": name}
	updateFilter := bson.M{"$set": bson.M{"status": status, "timestamp": time.Now()}}
	result := collection.FindOneAndUpdate(ctx, filter, updateFilter)
	err := result.Err()
	if err != nil {
		return err
	}

	return nil
}

func (m *ProjectMongoRepository) GetProjectResources(ctx context.Context, project string) ([]entity.KubernetesResource, error) {

	ctx = rctx.BuildContext(ctx, rctx.Context(rctx.Component, "GetProjectResources"), rctx.Context(rctx.StartTime, time.Now()))
	defer logger.LogComponentTime(ctx)

	collection := m.conn.GetCollection(vars.MONGODB_NAME, MONGODB_COLL_RESOURCES)
	result, err := collection.Find(ctx, bson.M{"project": project, "level": "project"})

	if err != nil || result.Err() != nil {
		return nil, handleMongoError(ctx, err)
	}

	var resources []entity.KubernetesResource

	for result.Next(ctx) {

		var resource entity.KubernetesResource
		if err = result.Current.Lookup("resource").Unmarshal(&resource); err != nil {
			return nil, errs.ErrMarshalFailed
		}

		resources = append(resources, resource)
	}

	return resources, nil
}

func (m *ProjectMongoRepository) SaveProjectResource(ctx context.Context, project string, resource *entity.KubernetesResource) error {

	ctx = rctx.BuildContext(ctx, rctx.Context(rctx.Component, "SaveProjectResource"), rctx.Context(rctx.StartTime, time.Now()))
	defer logger.LogComponentTime(ctx)

	collection := m.conn.GetCollection(vars.MONGODB_NAME, MONGODB_COLL_RESOURCES)

	filter := bson.M{"_id": resource.Id}
	updateFilter := bson.M{"$set": bson.M{"resource.state": resource.State, "resource.timestamp": resource.Timestamp}}

	result := collection.FindOneAndUpdate(ctx, filter, updateFilter)
	err := result.Err()

	if err != nil {
		if err == mongo.ErrNoDocuments {
			data := bson.M{"_id": resource.Id, "project": project, "level": "project", "resource": resource}
			_, err := collection.InsertOne(ctx, data)

			if err != nil {
				logger.Errorf(ctx, "Insert of kubernetes resource %s failed - %s", utils.MarshalObject(resource), err)
				return handleMongoError(ctx, err)
			}
		} else {
			logger.Errorf(ctx, "Update of kubernetes resource %s failed - %s", utils.MarshalObject(resource), err)
			return handleMongoError(ctx, err)
		}
	}

	return nil
}

func (m *ProjectMongoRepository) SaveProjectResources(ctx context.Context, project string, resources []entity.KubernetesResource) error {

	ctx = rctx.BuildContext(ctx, rctx.Context(rctx.Component, "SaveProjectResources"), rctx.Context(rctx.StartTime, time.Now()))
	defer logger.LogComponentTime(ctx)

	for _, resource := range resources {
		err := m.SaveProjectResource(ctx, project, &resource)
		if err != nil {
			return err
		}
	}

	return nil

}

func (m *ProjectMongoRepository) GetInstanceResources(ctx context.Context, project, instance string) ([]entity.KubernetesResource, error) {

	ctx = rctx.BuildContext(ctx, rctx.Context(rctx.Component, "GetInstanceResources"), rctx.Context(rctx.StartTime, time.Now()))
	defer logger.LogComponentTime(ctx)

	collection := m.conn.GetCollection(vars.MONGODB_NAME, MONGODB_COLL_RESOURCES)
	result, err := collection.Find(ctx, bson.M{"project": project, "instance": instance, "level": "instance"})

	if err != nil || result.Err() != nil {
		return nil, handleMongoError(ctx, err)
	}

	var resources []entity.KubernetesResource
	for result.Next(ctx) {
		var resource entity.KubernetesResource
		if err = result.Current.Lookup("resource").Unmarshal(&resource); err != nil {
			return nil, errs.ErrMarshalFailed
		}

		resources = append(resources, resource)
	}

	return resources, nil
}

func (m *ProjectMongoRepository) SaveInstanceResource(ctx context.Context, project, instance string, resource *entity.KubernetesResource) error {

	ctx = rctx.BuildContext(ctx, rctx.Context(rctx.Component, "SaveInstanceResource"), rctx.Context(rctx.StartTime, time.Now()))
	defer logger.LogComponentTime(ctx)

	coll := m.conn.GetCollection(vars.MONGODB_NAME, MONGODB_COLL_RESOURCES)

	filter := bson.M{"_id": resource.Id}
	updateFilter := bson.M{"$set": bson.M{"resource.state": resource.State, "resource.timestamp": resource.Timestamp}}

	result := coll.FindOneAndUpdate(ctx, filter, updateFilter)
	err := result.Err()

	if err != nil {

		if err == mongo.ErrNoDocuments {

			data := bson.M{"_id": resource.Id, "project": project, "instance": instance, "level": "instance", "resource": resource}
			_, err := coll.InsertOne(ctx, data)

			if err != nil {
				logger.Errorf(ctx, "insert of kubernetes resource %s failed - %s", utils.MarshalObject(resource), err)
				return handleMongoError(ctx, err)
			}
		} else {
			logger.Errorf(ctx, "Update of kubernetes resource %s failed - %s", utils.MarshalObject(resource), err)
			return handleMongoError(ctx, err)
		}
	}

	return nil
}

func (m *ProjectMongoRepository) SaveInstanceResources(ctx context.Context, project, instance string, resources []entity.KubernetesResource) error {

	ctx = rctx.BuildContext(ctx, rctx.Context(rctx.Component, "SaveInstanceResources"), rctx.Context(rctx.StartTime, time.Now()))
	defer logger.LogComponentTime(ctx)

	for _, resource := range resources {
		err := m.SaveInstanceResource(ctx, project, instance, &resource)
		if err != nil {
			return err
		}
	}

	return nil

}

func (m *ProjectMongoRepository) GetUserSummary(ctx context.Context, userId string) (*entity.ResourceSummary, error) {
	ctx = rctx.BuildContext(ctx, rctx.Context(rctx.Component, "GetUserSummary"), rctx.Context(rctx.StartTime, time.Now()))
	defer logger.LogComponentTime(ctx)

	useridFilter := bson.M{"userid": userId}
	ownerFilter := bson.M{"owner": userId}

	keyCount := CountDocuments(ctx, m.conn, vars.MONGODB_NAME, MONGODB_COLL_APIKEY, useridFilter)
	projectCount := CountDocuments(ctx, m.conn, vars.MONGODB_NAME, MONGODB_COLL_PROJECTS, ownerFilter)
	instanceCount := CountDocuments(ctx, m.conn, vars.MONGODB_NAME, MONGODB_COLL_INSTANCES, ownerFilter)
	teamCount := CountDocuments(ctx, m.conn, vars.MONGODB_NAME, MONGODB_COLL_TEAMS, ownerFilter)

	if keyCount == -1 || projectCount == -1 || instanceCount == -1 || teamCount == -1 {
		return nil, errs.ErrDBError
	}

	return &entity.ResourceSummary{
		TotalAPIKeys:   keyCount,
		TotalProjects:  projectCount,
		TotalInstances: instanceCount,
		TotalTeams:     teamCount,
	}, nil
}
