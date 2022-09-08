package mongo

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoClient struct {
	uri    string
	db     string
	creds  options.Credential
	client *mongo.Client
}

func NewClient() *MongoClient {
	mc := &MongoClient{}
	return mc
}

func (mc *MongoClient) SetOptions(opts ...Option) {
	for _, opt := range opts {
		opt(mc)
	}
}

func (mc *MongoClient) StartClient() error {
	client, err := mongo.NewClient(options.Client().ApplyURI(mc.uri).SetAuth(mc.creds))
	if err != nil {
		return err
	}
	mc.client = client
	return nil
}

func (mc *MongoClient) ConnectMongo(c context.Context) (*MongoConnection, error) {
	// declare context and try to set client connection
	ctx, cancel := context.WithTimeout(c, 10*time.Second)
	defer cancel()

	err := mc.client.Connect(ctx)
	if err != nil {
		return &MongoConnection{}, err
	}

	if err := mc.PingDB(ctx); err != nil {
		return nil, err
	}
	conn := mc.ConnectDB()

	// output
	log.Println("Connected to MongoDB")
	return conn, nil
}

func (mc *MongoClient) PingDB(ctx context.Context) error {
	// ping the database
	err := mc.client.Ping(ctx, nil)
	if err != nil {
		return err
	}
	return nil
}

func (mc *MongoClient) ConnectDB() *MongoConnection {
	conn := &MongoConnection{
		db: mc.client.Database(mc.db),
	}
	return conn
}
