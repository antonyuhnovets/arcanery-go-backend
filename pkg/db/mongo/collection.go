package mongo

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoConnection struct {
	db *mongo.Database
}

type MongoCollection struct {
	coll *mongo.Collection
}

func (mg *MongoConnection) GetCollection(collectionName string) *MongoCollection {
	collection := mg.db.Collection(collectionName)
	coll := &MongoCollection{
		coll: collection,
	}
	return coll
}

func (mc *MongoCollection) Create(doc interface{}, ctx context.Context) error {
	_, err := mc.coll.InsertOne(ctx, doc)
	return err
}

func (mc *MongoCollection) ReadById(doc interface{}, ctx context.Context, filter int64) error {
	return mc.coll.FindOne(ctx, bson.M{"_id": filter}).Decode(doc)
}

func (mc *MongoCollection) ReadAll(doc interface{}, ctx context.Context) error {
	cursor, err := mc.coll.Aggregate(ctx, bson.D{})
	if err != nil {
		return err
	}
	return cursor.All(ctx, doc)
}

func (mc *MongoCollection) UpdateById(doc interface{}, ctx context.Context, filter int64) error {
	_, err := mc.coll.UpdateOne(ctx, bson.M{"_id": filter}, bson.M{"$set": doc})
	return err
}

func (mc *MongoCollection) DeleteById(ctx context.Context, filter int64) error {
	_, err := mc.coll.DeleteOne(ctx, bson.M{"_id": filter})
	return err
}

func (mc *MongoCollection) DeleteAll(ctx context.Context) (int64, error) {
	result, err := mc.coll.DeleteMany(ctx, bson.M{})
	return result.DeletedCount, err
}
