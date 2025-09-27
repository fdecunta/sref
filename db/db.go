package db

import (
    "encoding/json"
    "errors"
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


func QueryDOI(file string, doi string) (crossref.Reference, error) {
    var r crossref.Reference

	data, err := LoadDB(file)
	if err != nil {
		return r, err
	}

    if ref, ok := data[doi]; ok {
        r = ref
    } else {
        r = crossref.SearchDoi(doi)
    }

    return r, nil
}

func AddReference(file string, doi string) (crossref.Reference, error) {
    var r crossref.Reference

   	data, err := LoadDB(file)
	if err != nil {
		return r, err
	}

    if _, ok := data[doi]; ok {
        return r, errors.New("DOI already exists.")
    } 

    r = crossref.SearchDoi(doi)
    data[r.DOI] = r

	if err := SaveDB(file, data); err != nil {
		return r, err
	}

    return r, nil
}


