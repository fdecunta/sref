package db

import (
    "encoding/json"
    "errors"
    "fmt"
    "io"
    "os"
    "strings"

    "github.com/caltechlibrary/crossrefapi"
)

// Uses Message struct from crossrefapi:
// https://github.com/caltechlibrary/crossrefapi/blob/main/works.go

type DataBase struct {
    Path string
    Table map[string]crossrefapi.Message
}


func Open(path string) (*DataBase, error) {
    d := DataBase{
        Path: path,
        Table: make(map[string]crossrefapi.Message),
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
    f, err := os.OpenFile(db.Path, os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ") 

    err = enc.Encode(db.Table)
    if err != nil {
        return err
    }

	return nil
}


func (db *DataBase) Store(r *crossrefapi.Message) error {
    // Check that don't exist
    if ref, ok := db.Table[r.DOI]; ok {
        return fmt.Errorf("reference already exists %s", ref.DOI)
    }

    // Clear Reference. If not, you got a huge list with all the papers that cite this one
    r.Reference = r.Reference[:0]

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


func (db * DataBase) QueryTitle(title string) (*crossrefapi.Message) {
    for _, i := range db.Table {
        if strings.Contains(strings.ToLower(i.Title[0]), strings.ToLower(title)) {
            ref := db.Table[i.DOI]
            return &ref
        }
    }
    return nil
}
