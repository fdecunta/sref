package db

import (
    "encoding/json"
    "errors"
    "log"
    "os"
    "strings"
    "sref/crossref"
)

var db = make(map[string]crossref.Reference)

func WriteDB(filename string, db map[string]crossref.Reference) error {
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


func Get(file string, doi string) (crossref.Reference, error) {
    var r crossref.Reference

    data, err := LoadDB(file)
    if err != nil {
        return r, err
    }

    ref, ok := data[doi]
    if !ok {
        return crossref.Reference{}, errors.New("DOI not found")
    } 
    r = ref

    return r, nil
}


func Set(file string, doi string, r crossref.Reference) error {
    data, err := LoadDB(file)
    if err != nil {
        return err
    }

    if _, ok := data[doi]; ok {
        return errors.New("DOI already exists.")
    } 

    if err := WriteDB(file, data); err != nil {
        return err
    }

    return nil
}


func AddReference(file string, doi string) error {
    data, err := LoadDB(file)
    if err != nil {
        return err
    }

    if _, ok := data[doi]; ok {
        return errors.New("DOI already exists.")
    } 

    r := crossref.SearchDoi(doi)
    data[doi] = *r

    if err := WriteDB(file, data); err != nil {
        return err
    }

    return nil
}


func DeleteReference(file string, doi string) (error) {
    data, err := LoadDB(file)
    if err != nil {
        return err
    }

    if _, ok := data[doi]; !ok {
        return errors.New("DOI not found.")
    } 
    delete(data, doi)

    if err := WriteDB(file, data); err != nil {
        return err
    }

    return nil
}


func SearchByTitle(file string, title string) (*crossref.Reference, error) {
    data, err := LoadDB(file)
    if err != nil {
        return nil, err
    }

    for _, d := range data {
        if strings.Contains(strings.ToLower(d.Title), strings.ToLower(title)) {
            ref := data[d.DOI]
            return &ref, nil
        }
    }

    r := crossref.SearchTitle(title)
    if r == nil {
        log.Fatal(err)
    }
    return r, nil
}
 
