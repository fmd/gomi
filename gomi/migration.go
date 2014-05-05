package gomi

import (
	"labix.org/v2/mgo"
)

var migrationsDir string = "migrations"

type Migration struct {
	Id        string      `bson:"_id"`
	Timestamp int64       `bson:"timestamp"`
	Structure interface{} `bson:"structure"`
}

type Migrator struct {
	Migrations *mgo.Collection
	Structures *mgo.Collection
}

func NewMigrator(migrations *mgo.Collection, structures *mgo.Collection) *Migrator {
	return &Migrator{
		Migrations: migrations,
		Structures: structures,
	}
}
