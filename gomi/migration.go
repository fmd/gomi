package gomi

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"labix.org/v2/mgo"
	"strconv"
	"strings"
	"time"
)

var migrationsDir string = "migrations"

type Structure struct {
	Id     string      `bson:"_id"    json:"_id"`
	Schema interface{} `bson:"schema" json:"schema"`
}

type Migration struct {
	Id        string     `bson:"_id"       json:"_id"`
	Timestamp int64      `bson:"timestamp" json:"timestamp"`
	Structure *Structure `bson:"structure" json:"structure"`
}

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

//Saves a new migration based on a structure file.
//Returns an error if unsuccessful, or nil otherwise.
func (m *Migrator) Structure(name string) error {
	var err error

	content, err := ioutil.ReadFile(fmt.Sprintf("%s/%s.json", structureName, name))
	if err != nil {
		return err
	}

	s := &Structure{}
	err = json.Unmarshal(content, s)
	if err != nil {
		return err
	}

	fmt.Println(s)

	g := m.CreateMigration(s)
	err = g.Save()
	if err != nil {
		return err
	}

	return nil
}

//MigrationIndex gets the current number of migrations in ./migrations and adds one to the figure.
//Returns a string padded up to five chars with zeroes and a nil error if successful,
//or a blank string and an error if unsuccessful.
func MigrationIndex() (string, error) {
	files, err := ioutil.ReadDir(migrateName)
	if err != nil {
		return "", err
	}

	num := strconv.Itoa(len(files) + 1)
	return fmt.Sprintf("%s%s", strings.Repeat("0", 6-len(num)), num), nil
}

//Creates a migration. Uses a migrator and a structure to create and return a Migration.
//Returns the newly created migration.
func (m *Migrator) CreateMigration(s *Structure) *Migration {
	g := &Migration{}
	g.Timestamp = time.Now().UTC().UnixNano()
	g.Id = strconv.FormatInt(g.Timestamp, 16)
	g.Structure = s

	return g
}

//Serializes the migration to JSON for saving.
//Returns an empty byte slice and an error if it fails,
//or a byte slice of JSON characters and a nil error if successful.
func (g *Migration) Serialize() ([]byte, error) {
	j, err := json.MarshalIndent(g, "", "    ")
	if err != nil {
		return []byte{}, err
	}

	return j, nil
}

//Gets a migration's filename.
//Returns an empty string and and error if unsucessful, or the filename and nil otherwise.
func (g *Migration) GetFilename() (string, error) {
	idx, err := MigrationIndex()
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s_%s_%s.json", idx, g.Id, strings.ToLower(g.Structure.Id)), nil
}

//Saves the migration to a file.
//Returns an error if unsucessful, or nil otherwise.
func (g *Migration) Save() error {
	filename, err := g.GetFilename()
	if err != nil {
		return err
	}

	migration, err := g.Serialize()
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(fmt.Sprintf("%s/%s", migrateName, filename), migration, 0755)
	if err != nil {
		return err
	}

	return nil
}
