package mongo

import (
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Option func(*MongoClient)

func ConnectionURI(uri string) Option {
	return func(mc *MongoClient) {
		mc.uri = uri
	}
}

func Creds(name, pass string) Option {
	return func(mc *MongoClient) {
		mc.creds = options.Credential{
			AuthSource: "admin",
			Username:   name,
			Password:   pass,
		}
	}
}

func Database(db string) Option {
	return func(mc *MongoClient) {
		mc.db = db
	}
}
