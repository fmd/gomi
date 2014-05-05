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
func (m *Migrator) NewMigration(s *Structure) *Migration {
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

//IsStructured checks whether a structure already exists within Mongo.
//If uses a *Migrator and a *Structure and returns whether the structure exists.
func (m *Migrator) HasStructure(s *Structure) bool {
	t := &Structure{}
	err := m.Structures.Find(bson.M{"_id": s.Id}).One(&t)
	if err != nil {
		return false
	}

	if len(t.Id) > 0 {
		return true
	}

	return false
}

//Apply applies a migration using a *Migrator from a *Migration.
//BUG(Needs to actually find a way to determine what changes are actually necessary and then do generate datamigrations, etc.)
//Returns an error if unsuccessful, or nil otherwise.
func (m *Migrator) Apply(g *Migration) error {
	var err error
	a := m.IsApplied(g)
	if a {
		return nil
	}

	if !m.HasStructure(g.Structure) {
		err = m.Structures.Insert(g.Structure)
		if err != nil {
			return err
		}

		return nil
	}

	err = m.Structures.Update(bson.M{"_id": g.Structure.Id}, g.Structure)
	if err != nil {
		return err
	}

	err = m.Migrations.Insert(g)
	if err != nil {
		return err
	}

	return nil
}

//Rollback uses a *Migrator to roll back to the previous migration.
//BUG(Again, needs to actually perform data migrations, otherwise it's kind of pointless having this structure.)
//Returns an error if unsuccessful, or nil otherwise.
func (m *Migrator) Rollback() error {
    var err error

    g := &Migration{}
    err = m.Migrations.Find(nil).Sort("-timestamp").One(&g)
    if err != nil {
        return err
    }

    err = m.Migrations.Remove(g)
    if err != nil {
        return err
    }

    err = m.Migrations.Find(nil).Sort("-timestamp").One(&g)
    if err != nil {
        s := &Structure{}
        err = m.Structures.Find("_id":bson.M{g.Structure.Id}).One(&s)
        if err != nil {
            return err
        }

        err = m.Structures.Remove(s)
        if err != nil {
            return err
        }

        return nil
    }

    err = m.Apply(g)
    if err != nil {
        return err
    }

    return nil
}
