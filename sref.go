package main

import (
    "errors"
    "log"
    "flag"
    "fmt"
    "os"
    "path/filepath"
    "regexp"
    "strings"

    "sref/db"
    "sref/crossref"
)

var file string
var input string
var doi string
var add bool
var del bool
var read bool
var edit bool

var d *db.DataBase

func main() {
    flag.StringVar(&file, "file", "", "JSON file")
    flag.StringVar(&input, "input", "", "INPUT to use. Must be DOI or TITLE")
    flag.BoolVar(&add, "a", false, "Add reference")
    flag.BoolVar(&del, "d", false, "Delete reference")
    flag.BoolVar(&read, "r", false, "Read reference")
    flag.BoolVar(&edit, "e", false, "Edit reference")
    flag.Parse()

    file, err := assertFile(file)
    if err != nil {
        fmt.Println(err)
        flag.Usage()
        os.Exit(1)
    }

    d, err = db.Open(file)
    if err != nil {
        log.Fatal(err)
    }

    doi, err := assertDoi(input)
    if err != nil {
        fmt.Println(err)
        flag.Usage()
        os.Exit(1)
    }

    r := d.Get(doi)

    if read {
        if r == nil {
           fmt.Println("DOI not found")
           return
        }
        fmt.Println(*r)
    } else if add {
        if r != nil {
            fmt.Println("DOI already exists")
            return
        }
        err := d.Set(doi, crossref.SearchDoi(doi))
        if err != nil {
            fmt.Fprintf(os.Stderr, "Failed fo store reference: %s\n", err)
            os.Exit(1)
        }
    } else if del {
        if r == nil {
            fmt.Println("DOI not found")
            return
        }
        err := d.Delete(doi)
        if err != nil {
            fmt.Fprintf(os.Stderr, "Failed fo delete DOI: %s\n", err)
            os.Exit(1)
        }
    } else if edit {
        if r == nil {
            fmt.Println("DOI not found")
            return
        }
        // TODO:
        fmt.Println("not implemented")
    } 
}


func GetDefaultJSON() (string, error) {
    homeDir, err := os.UserHomeDir()
    if err != nil {
        return "", err
    }

    configDir := filepath.Join(homeDir, ".config/sref")
    if err := os.MkdirAll(configDir, os.ModePerm); err != nil {
        return "", err
    }

    configFile := filepath.Join(configDir, "sref.json")
    file, err := os.OpenFile(configFile, os.O_RDONLY|os.O_CREATE, 0644)
    if err != nil {
        return "", err
    }
    file.Close()

    return configFile, nil
}


func assertFile(file string) (string, error) {
    if file == "" {
        var err error
        file, err = GetDefaultJSON()
        if err != nil {
            return "", err
        }
    }

    if _, err := os.Stat(file); err != nil {
        if os.IsNotExist(err) {
            return "", errors.New("File does not exist")
        }
        return "", err
    } 

    if filepath.Ext(file) != ".json" {
        return "", errors.New("File is not JSON")
    }

    return file, nil
}


func CaptureDoi(s string) (string, bool) {
    re := regexp.MustCompile(`10\.\d{4,}/\S+`)
    match := re.FindString(s)
    if match != "" {
        return match, true
    }
    return "", false
}


func assertDoi(s string) (string, error) {
    s = strings.TrimSpace(s)
    if len(s) == 0 {
        return "", errors.New("empty input")
    }

    doi, ok := CaptureDoi(s)
    if !ok {
        r, err := db.SearchByTitle(file, input)
        if err != nil {
            return "", err
        }
        doi = r.DOI
    }

    return doi, nil
}
