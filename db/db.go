package db

import (
    "encoding/json"
    "os"
    "sref/crossref"
)

var db = make(map[string]crossref.Reference)

func SaveDB(filename string, db map[string]crossref.Reference) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ") // pretty print
	return enc.Encode(db)
}


func LoadDB(filename string) (map[string]crossref.Reference, error) {
	f, err := os.Open(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return make(map[string]crossref.Reference), nil // start empty
		}
		return nil, err
	}
	defer f.Close()
	var db map[string]crossref.Reference
	err = json.NewDecoder(f).Decode(&db)
	if err != nil {
		return make(map[string]crossref.Reference), nil
	}
	return db, nil
}

