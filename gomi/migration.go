package gomi

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"
)

type Structure struct {
	Id     string      `bson:"_id"    json:"_id"`
	Schema interface{} `bson:"schema" json:"schema"`
}

type Migration struct {
	Id        string     `bson:"_id"       json:"_id"`
	Timestamp int64      `bson:"timestamp" json:"timestamp"`
	Structure *Structure `bson:"structure" json:"structure"`
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
