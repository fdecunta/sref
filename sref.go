package main

import (
    "errors"
    "log"
    "flag"
    "fmt"
    "os"
    "path/filepath"
    "sref/db"
    "sref/format"
)

var file string
var doi string
var add bool

func main() {
    flag.StringVar(&file, "file", "data.json", "JSON file")
    flag.StringVar(&doi, "doi", "", "DOI to use")
    flag.BoolVar(&add, "add", false, "Option - Add DOI's reference into FILE")

    flag.Parse()

    // Check file 
    if err := isJsonFile(file); err != nil {
        fmt.Fprintf(os.Stderr, "Can't read file: %s\n", err)
        os.Exit(1)
    }

    if add {
        if len(doi) == 0 {
            // TODO: validate DOI. At least that looks like a DOOI
            fmt.Fprintf(os.Stderr, "missing DOI\n")
            os.Exit(1)
        }

        r, err := db.AddReference(file, doi) 
        if err != nil {
            fmt.Fprintf(os.Stderr, "Failed fo save DOI: %s\n", err)
            os.Exit(1)
        }
        fmt.Println(format.FormatCite(r))
        return 
    }

    data, err := db.LoadDB(file)
    if err != nil {
        log.Panicln(err)
    }

    if len(doi) == 0 {
        for _, r := range data {
        fmt.Println(format.FormatCite(r))
        }

    } else {
        r, err := db.QueryDOI(file, doi)
        if err != nil {
            log.Fatal(err)
        }
        fmt.Println(format.FormatCite(r))
    }
}


func isJsonFile(path string) error {
    if _, err := os.Stat(path); os.IsNotExist(err) {
        return errors.New("File does not exist")
    } 

    if filepath.Ext(path) != ".json" {
        return errors.New("File is not JSON")
    }
    return nil
}
