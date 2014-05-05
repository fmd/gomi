package gomi

import (
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"strconv"
	"time"
)

type Migrator struct {
	Migrations *mgo.Collection
	Structures *mgo.Collection
}

//Creates a migrator instance based on the two collections the instance needs.
//Returns a new *Migrator.
func NewMigrator(migrations *mgo.Collection, structures *mgo.Collection) *Migrator {
	return &Migrator{
		Migrations: migrations,
		Structures: structures,
	}
}

//Creates a migration. Uses a *Migrator and a *Structure to create and return a *Migration.
//Returns the newly created migration.
func (m *Migrator) CreateMigration(s *Structure) *Migration {
	g := &Migration{}
	g.Timestamp = time.Now().UTC().UnixNano()
	g.Id = strconv.FormatInt(g.Timestamp, 16)
	g.Structure = s

	return g
}

//IsApplied checks whether a migration has been applied.
//It uses a *Migrator and a *Migration to see if the migration has already been applied to Mongo.
//Returns a bool - true or false depending on whether the migration has been applied.
func (m *Migrator) IsApplied(g *Migration) bool {
	ts := strconv.FormatInt(g.Timestamp, 16)

	r := &Migration{}
	err := m.Migrations.Find(bson.M{"_id": ts}).One(&r)
	if err != nil {
		return false
	}

	if len(r.Id) > 0 {
		return true
	}

	return false
}

//Apply applies a migration using a *Migrator from a *Migration.
//Returns an error if unsuccessful, or nil otherwise.
func (m *Migrator) Apply(g *Migration) error {
	a := m.IsApplied(g)
	if a {
		return nil
	}

	err := m.Migrations.Insert(g)
	if err != nil {
		return err
	}

	return nil
}
