package db

import (
    "encoding/json"
    "errors"
    "io"
    "log"
    "os"
    "strings"
    "sref/crossref"
)


type DataBase struct {
    Path string
    Table map[string]crossref.Reference
}


func Open(path string) (*DataBase, error) {
    d := DataBase{
        Path: path,
        Table: make(map[string]crossref.Reference),
    }

	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	err = json.NewDecoder(f).Decode(&d.Table)
    if err == io.EOF {
        return &d, nil
    }
	if err != nil {
		return nil, err
	}

	return &d, nil
}


func (db *DataBase) Write() error {
	f, err := os.Create(db.Path)
	if err != nil {
		return err
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ") 

	return enc.Encode(db.Table)
}


func (db *DataBase) Get(doi string) *crossref.Reference {
    r, ok := db.Table[doi]
    if !ok {
        return nil
    } 
    return &r
}


func (db *DataBase) Set(doi string, r *crossref.Reference) error {
    db.Table[doi] = *r

    if err := db.Write(); err != nil {
        return err
    }
    return nil
}


func (db *DataBase) Delete(doi string) error {
    if _, ok := db.Table[doi]; !ok {
        return errors.New("DOI not found.")
    } 

    delete(db.Table, doi)

    if err := db.Write(); err != nil {
        return err
    }

    return nil
}


func SearchByTitle(file string, title string) (*crossref.Reference, error) {
    data, err := Open(file)
    if err != nil {
        return nil, err
    }

    for _, d := range data.Table {
        if strings.Contains(strings.ToLower(d.Title), strings.ToLower(title)) {
            ref := data.Table[d.DOI]
            return &ref, nil
        }
    }

    r := crossref.SearchTitle(title)
    if r == nil {
        log.Fatal(err)
    }
    return r, nil
}
 
