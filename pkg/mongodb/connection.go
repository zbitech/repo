package mongodb

import (
	"context"
	"github.com/zbitech/common/pkg/rctx"
	"log"
	"time"

	"github.com/zbitech/common/pkg/logger"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//var SERVER = "mongodb+srv://admin:5ysaDYYS7vbs5r4@cluster0.hi6d5.mongodb.net/zbiRepo?retryWrites=true&w=majority"

type MongoDBConnection struct {
	client *mongo.Client
	uri    string
	cancel context.CancelFunc
}

func NewMongoDBConnection(conn_url string) *MongoDBConnection {
	return &MongoDBConnection{
		client: nil,
		uri:    conn_url,
	}
}

func (m *MongoDBConnection) OpenConnection(ctx context.Context) error {

	ctx = rctx.BuildContext(ctx, rctx.Context(rctx.Component, "OpenConnection"), rctx.Context(rctx.StartTime, time.Now()))
	defer logger.LogComponentTime(ctx)

	logger.Infof(ctx, "Opening connection to Mongo-DB @ %s", m.uri)
	clientOptions := options.Client().ApplyURI(m.uri)
	//ctx, cancel := context.WithTimeout(ctx, 10*time.Second)

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return err
	}

	m.client = client

	return nil
}

func (m *MongoDBConnection) CloseConnection(ctx context.Context) {
	logger.Infof(ctx, "Closing connection to Mongo-DB")

	if m.client != nil {
		err := m.client.Disconnect(ctx)
		if err != nil {
			log.Printf("Error while closing connection - %s", err)
		}
	}
	m.client = nil
}

func (m *MongoDBConnection) GetDatabase(database string) *mongo.Database {
	return m.client.Database(database)
}

func (m *MongoDBConnection) GetCollection(database, document string) *mongo.Collection {
	if m.client == nil {
		logger.Debug(context.Background(), "MongoDB Client is nil")
	}
	return m.client.Database(database).Collection(document)
}
