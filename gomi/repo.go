package gomi

import (
	"errors"
	"fmt"
	"labix.org/v2/mgo"
	"os"
)

//A Repo is a directory which contains at least:
//The directories for migrations and structures,
//A mongo database with the migrations and structures collections.
//We store the current mgo session and database in order to perform operations on a Repo.
type Repo struct {
	Migrator *Migrator
	Session  *mgo.Session
	Db       *mgo.Database
}

//These are the names for the migrations and structures folders/collections.
var migrateName = "migrations"
var structureName = "structures"

//NewSession creates a new Mgo session.
//BUG(Needs to include other credentials from the Repo session)
//It returns the session and a nil error if successful, or nil and an error otherwise.
func NewSession(hostname string) (*mgo.Session, error) {
	s, err := mgo.Dial(hostname)
	if err != nil {
		return nil, err
	}

	s.SetMode(mgo.Monotonic, true)
	return s, nil
}

//NewRepo creates and returns a *Repo.
//It does NOT initialize a repo's folders and collections. See (r *Repo) Init().
//It takes Mongo credentials and a db in string format.
//It returns a new Repo and a nil error if successful, or nil and an error otherwise.
func NewRepo(hostname string, db string) (*Repo, error) {
	var err error
	r := &Repo{}
	r.Session, err = NewSession(hostname)
	if err != nil {
		return nil, err
	}
	r.Db = r.Session.DB(db)
	r.Migrator = NewMigrator(r.Db.C(migrateName), r.Db.C(structureName))

	return r, nil
}

//Creates a directory in the current directory.
//Take a name string to determine the name to give the new directory.
//Returns an error if unsuccessful, nil otherwise.
func CreateDir(name string) error {

	//Ensure nothing with this name already exists here.
	if _, err := os.Stat(name); err != nil {
		if !os.IsNotExist(err) {
			return err
		}
	} else {
		return errors.New(fmt.Sprintf("File named '%s' already exists in this directory.", name))
	}

	//Create the directory.
	if err := os.Mkdir(name, 0755); err != nil {
		return err
	}

	return nil
}

//Creates the necessary directories for a new Repo.
//Returns an error if unsuccessful, or nil otherwise.
func MakeDirs() error {
	var err error

	//Make the project directory.
	if err = CreateDir(migrateName); err != nil {
		return err
	}

	//Make the project directory.
	if err = CreateDir(structureName); err != nil {
		return err
	}

	return nil
}

//Creates the collections for a new Repo.
//Returns an error if unsuccessful, or nil otherwise.
func (r *Repo) MakeCollections() error {
	c := r.Db.C(migrateName)
	err := c.Create(&mgo.CollectionInfo{})

	if err != nil {
		return err
	}

	c = r.Db.C(structureName)
	err = c.Create(&mgo.CollectionInfo{})

	if err != nil {
		return err
	}

	return nil
}

//Initializes a new Repo.
//Returns an error if unsuccessful, or nil otherwise.
func (r *Repo) Init() error {
	var err error
	err = MakeDirs()
	if err != nil {
		return err
	}

	err = r.MakeCollections()
	if err != nil {
		return err
	}

	return nil
}
