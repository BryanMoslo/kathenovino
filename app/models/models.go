package models

import (
	"bytes"
	"kathenovino"
	"kathenovino/config"
	"log"

	"github.com/gobuffalo/envy"
	"github.com/gobuffalo/pop/v6"
	"github.com/wawandco/ox/pkg/buffalotools"
)

var (
	// DB returns the DB connection for a given connection name
	db = buffalotools.DatabaseProvider(config.FS())
)

var env = envy.Get("GO_ENV", "development")

// DB is a connection to your database to be used
// throughout your application.
var DB *pop.Connection

func init() {
	initDB()
}

// FindDB allows to pull a connection by the name of it
// this comes handy when you have multiple databases to read
// from.
func FindDB(connname string) *pop.Connection {
	return db(connname)
}

func initDB() {
	bf, err := kathenovino.Config.Find("database.yml")
	if err != nil {
		log.Fatal(err)
	}

	err = pop.LoadFrom(bytes.NewReader(bf))
	if err != nil {
		log.Fatal(err)
	}

	DB, err = pop.Connect(env)
	if err != nil {
		log.Fatal(err)
	}

	pop.Debug = env == "development"
}