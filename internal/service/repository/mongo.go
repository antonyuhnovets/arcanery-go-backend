// Abstract implementation of mongo repository.
// Manage connecttion and CRUD methods.

package repository

import (
	"context"
	"log"

	"github.com/hetonei/arcanery-go-backend/internal/service"
	"github.com/hetonei/arcanery-go-backend/pkg/db/mongo"
)

// Connection to db
type MongoConnection struct {
	db   string
	ctx  context.Context
	conn *mongo.MongoConnection
}

// Concrete collection
// Implements CRUD methods of repository interface from interactor
type MongoRepository struct {
	usecase string
	ctx     context.Context
	coll    *mongo.MongoCollection
}

func (mr *MongoRepository) Create(result interface{}) error {
	if err := mr.coll.Create(result, mr.ctx); err != nil {
		return err
	}

	return nil
}

func (mr *MongoRepository) ReadById(result interface{}, filter int64) error {
	if err := mr.coll.ReadById(result, mr.ctx, filter); err != nil {
		return err
	}

	return nil
}

func (mr *MongoRepository) ReadAll(result interface{}) error {
	if err := mr.coll.ReadAll(result, mr.ctx); err != nil {
		return err
	}

	return nil
}

func (mr *MongoRepository) UpdateById(result interface{}, filter int64) error {
	if err := mr.coll.UpdateById(result, mr.ctx, filter); err != nil {
		return err
	}

	return nil
}

func (mr *MongoRepository) DeleteById(filter int64) error {
	if err := mr.coll.DeleteById(mr.ctx, filter); err != nil {
		return err
	}

	return nil
}

func (mr *MongoRepository) DeleteAll() (int64, error) {
	count, err := mr.coll.DeleteAll(mr.ctx)
	if err != nil {
		return 0, err
	}

	return count, nil
}

// Get repository interface with it's methods
func (mc *MongoConnection) GetRepoService(collection string) service.RepositoryService {
	return &MongoRepository{
		usecase: collection,
		ctx:     mc.ctx,
		coll:    mc.conn.GetCollection(collection),
	}
}

// Start client with options, connect to db
func GetConnectionMongo(ctx context.Context, opts map[string]string) *MongoConnection {
	mc := &mongo.MongoClient{
		Opts: mongo.SetOptions(opts),
	}

	if err := mc.StartClient(); err != nil {
		log.Println(err)
	}

	conn, err := mc.ConnectMongo(ctx)
	if err != nil {
		log.Println(err)
	}

	return &MongoConnection{
		db:   mc.Opts.GetDB(),
		ctx:  ctx,
		conn: conn,
	}
}
