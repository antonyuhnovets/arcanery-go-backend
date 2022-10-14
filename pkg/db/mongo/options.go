// Methods to set, get mongo connection options

package mongo

import "go.mongodb.org/mongo-driver/mongo/options"

type Options struct {
	db    string             // db name
	uri   string             // connection uri
	creds options.Credential // authorisation creds
}

// Set DB connection options
func SetOptions(opts map[string]string) Options {
	return Options{
		db:    opts["name"],
		uri:   opts["uri"],
		creds: SetCreds(opts["username"], opts["password"]),
	}
}

// Set DB authorisation creds
func SetCreds(username, pass string) options.Credential {
	return options.Credential{
		Username:   username,
		Password:   pass,
		AuthSource: "admin",
	}
}

func (o Options) GetDB() string {
	return o.db
}

func (o Options) GetURI() string {
	return o.uri
}

func (o Options) GetCreds() options.Credential {
	return o.creds
}
