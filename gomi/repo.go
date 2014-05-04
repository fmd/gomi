package gomi

import (
	"fmt"
	"labix.org/v2/mgo"
)

type Repo struct {
	Session *mgo.Session
	Db      *mgo.Database
}

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
//It takes Mongo credentials and a db in string format.
//It returns a new Repo and a nil error if successful, or nil and an error otherwise.
func NewRepo(hostname string, db string) (*Repo, error) {
	r := &Repo{}
	r.Session, err = NewSession(hostname)
	if err != nil {
		return nil, err
	}

	r.Db = r.Session.DB(db)
	return r, nil
}

func MakeCollections() error {
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

func MakeDir(name string) error {

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

func MakeDirs() {
	var err error

	//Make the project directory.
	if err = Make(migrateName); err != nil {
		return err
	}

	//Make the project directory.
	if err = Make(structureName); err != nil {
		return err
	}
}
