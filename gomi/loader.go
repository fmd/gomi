package gomi

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

//These are the names for the migrations and structures folders/collections.
var migrateName = "migrations"
var structureName = "structures"

//MigrationIndex gets the current number of migrations in dirname and adds one to the figure.
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

//LoadStructure loads a structure from the structure's name.
//It returns a *Structure and nil if successful, or nil and an error if unsuccessful.
func LoadStructure(name string) (*Structure, error) {
	var err error

	content, err := ioutil.ReadFile(fmt.Sprintf("%s/%s.json", structureName, name))
	if err != nil {
		return nil, err
	}

	s := &Structure{}
	err = json.Unmarshal(content, s)
	if err != nil {
		return nil, err
	}

	return s, nil
}

//func LoadMigration() {

//}

//LoadMigrations loads all migrations.
//It returns a slice of *Migrations and a nil error if successful,
//or nil and an error otherwise.
func LoadMigrations() ([]*Migration, error) {
	files, err := ioutil.ReadDir(migrateName)
	if err != nil {
		return nil, err
	}

	m := []*Migration{}

	for _, fn := range files {
		content, err := ioutil.ReadFile(fmt.Sprintf("%s/%s", migrateName, fn.Name()))
		if err != nil {
			return nil, err
		}

		g := &Migration{}

		err = json.Unmarshal(content, g)
		if err != nil {
			return nil, err
		}

		m = append(m, g)
	}

	return m, nil
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
