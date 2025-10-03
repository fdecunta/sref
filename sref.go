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

    "sref/crossref"
    "sref/db"
)

var d *db.DataBase

func main() {
    var file string
    var input string
    var doi string
    var add bool
    var del bool
    var read bool
//    var edit bool
    var toJson bool

    flag.StringVar(&file, "file", "", "Path to the JSON database file")
    flag.StringVar(&input, "input", "", "Input value to use. Can be a DOI or the paper's title")
    flag.BoolVar(&add, "a", false, "Add reference to the database")
    flag.BoolVar(&del, "d", false, "Delete reference from the database")
    flag.BoolVar(&read, "r", false, "Read reference from the database")
    flag.BoolVar(&toJson, "json", false, "Print reference(s) in JSON format")

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

    if toJson {
        for _, i := range d.Table {
            fmt.Println(i.ToJson())
        }
        return
    }

    doi, err = assertDoi(input)
    if err != nil {
        fmt.Println(err)
        flag.Usage()
        os.Exit(1)
    }

    r := d.Get(doi)

    if add {
        err := Add(d, r, doi)
        if err != nil {
            fmt.Println("Failed to store reference: %s\n", err)
            os.Exit(1)
        }
        return
    }

    // Next operations need r to exist:
    if r == nil {
       fmt.Println("ERROR: DOI not found")
       os.Exit(1)
    }

    if read {
        s, err := r.ToJson()
        if err != nil {
            fmt.Println("error: can't read reference \n%s\n", err)
            os.Exit(1)
        }
        fmt.Println(s)
    } else if del {
        err := d.Delete(doi)
        if err != nil {
            fmt.Fprintf(os.Stderr, "Failed fo delete DOI: %s\n", err)
            os.Exit(1)
        }
    } 
}


func Add(d *db.DataBase, r *crossref.Reference, doi string) error {
    if r != nil {
        fmt.Println("DOI already exists")
        return nil
    }

    r, err := crossref.SearchDoi(doi)
    if err != nil {       
        return err
    }
    
    err = d.Set(doi, r)
    if err != nil {
        return err
    }

    return nil
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

    // Try to capture DOI
    doi, ok := CaptureDoi(s)
    if ok {
        return doi, nil
    }

    // If not a DOI, use it as a title and search for it in the database
    if r := QueryTitle(s); r != nil {
        return r.DOI, nil
    }

    // Fallback to title search in CrossRef
    r, err := crossref.SearchTitle(s)
    if err != nil {
        return "", nil
    }
    if r == nil {
        return "", errors.New("DOI not found")
    }
    return r.DOI, nil
}


func QueryTitle(title string) (*crossref.Reference) {
    for _, i := range d.Table {
        if strings.Contains(strings.ToLower(i.Title), strings.ToLower(title)) {
            ref := d.Table[i.DOI]
            return &ref
        }
    }
    
    return nil
}
