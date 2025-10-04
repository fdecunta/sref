package db

import (
    "encoding/json"
    "errors"
    "io"
    "os"
    "strings"
    "strconv"
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


func (db *DataBase) Store(r *crossref.Reference) error {
    id := db.assignId(r)
    if id == "" {
        return errors.New("can't assign new ID")
    }
    r.ID = id

    db.Table[r.DOI] = *r

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


func (db *DataBase) assignId(r *crossref.Reference) string {
    var family string
    var year string

    family = strings.TrimSpace(r.Author[0].Family)
    family = strings.ReplaceAll(family, " ", "")

    year = strconv.Itoa(r.Issued.DateParts[0][0])

    base := family + year
    if isIdFree(base, db) {
        return base
    }

    // try counters until one is free
    for i := 2; i < 100; i++ {
        candidate := base + "_" + strconv.Itoa(i)
        if isIdFree(candidate, db) {
            return candidate
        }
    }

    // if it reach this something bad happened
    return base + "_x"
}


func isIdFree(id string, db *DataBase) bool {
    for _, r := range db.Table {
        if (id == r.ID) {
            return false
        }
    }
    return true
}


func (db * DataBase) QueryTitle(title string) (*crossref.Reference) {
    for _, i := range db.Table {
        if strings.Contains(strings.ToLower(i.Title), strings.ToLower(title)) {
            ref := db.Table[i.DOI]
            return &ref
        }
    }
    return nil
}


func (db *DataBase) QueryId(id string) (*crossref.Reference) {
    for _, i := range db.Table {
        if strings.ToLower(i.ID) == strings.ToLower(id) {
            ref := db.Table[i.DOI]
            return &ref
        }
    }
    return nil
}
