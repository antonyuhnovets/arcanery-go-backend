// Package for mongo lib.
// Contains tools for set up client, connect and communicate with db and collection.

package mongo

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Client with data, needed for connecting
type MongoClient struct {
	Opts   Options
	client *mongo.Client // client
}

// New empty MongoClient struct
func NewClient(opts Options) *MongoClient {
	mc := &MongoClient{
		Opts: opts,
	}

	return mc
}

// Add pointer of mongo client to struct
func (mc *MongoClient) StartClient() error {
	client, err := mongo.NewClient(options.Client().ApplyURI(mc.Opts.uri).SetAuth(mc.Opts.creds))
	if err != nil {
		return err
	}
	mc.client = client

	return nil
}

// Get connection to database.
func (mc *MongoClient) ConnectMongo(c context.Context) (*MongoConnection, error) {
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
	log.Println("Connected to MongoDB")

	return conn, nil
}

// Ping the database
func (mc *MongoClient) PingDB(ctx context.Context) error {
	err := mc.client.Ping(ctx, nil)
	if err != nil {
		return err
	}

	return nil
}

// Connect to database by name
func (mc *MongoClient) ConnectDB() *MongoConnection {
	conn := &MongoConnection{
		conn: mc.client.Database(mc.Opts.db),
	}

	return conn
}

// Connection to database
type MongoConnection struct {
	conn *mongo.Database
}

func (mc *MongoConnection) GetCollection(collectionName string) *MongoCollection {
	collection := mc.conn.Collection(collectionName)
	coll := &MongoCollection{
		coll: collection,
	}

	return coll
}
