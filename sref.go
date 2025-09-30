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
)

var file string
var input string
var doi string
var add bool
var del bool
var read bool

func main() {
    flag.StringVar(&file, "file", "", "JSON file")
    flag.StringVar(&input, "input", "", "INPUT to use. Must be DOI or TITLE")

    flag.BoolVar(&add, "a", false, "Add reference")
    flag.BoolVar(&del, "d", false, "Delete reference")
    flag.BoolVar(&read, "r", false, "Read reference")
    flag.Parse()

    file, err := assertFile(file)
    if err != nil {
        fmt.Println(err)
        flag.Usage()
        os.Exit(1)
    }

    doi, err := assertDoi(input)
    if err != nil {
        fmt.Println(err)
        flag.Usage()
        os.Exit(1)
    }

    if read {
        r, err := db.Get(file, doi)
        if err != nil {
           log.Fatal(err)
        }
           
        fmt.Println(r)
        return
    }

    if add {
        err := db.AddReference(file, doi) 
        if err != nil {
            fmt.Fprintf(os.Stderr, "Failed fo store reference: %s\n", err)
            os.Exit(1)
        }
        return 
    }

    if del {
        if err := db.DeleteReference(file, doi); err != nil {
            fmt.Fprintf(os.Stderr, "Failed fo delete DOI: %s\n", err)
            os.Exit(1)
        }
        return 
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


func IsDoi(s string) bool {
   re := regexp.MustCompile(`10\.\d{4,}/\S+`)
   match := re.FindString(s)
   if match != "" {
       return true
   } else {
       return false
   }
}


func CaptureDoi(s string) (string, bool) {
    re := regexp.MustCompile(`10\.\d{4,}/\S+`)
    match := re.FindString(s)
    if match != "" {
        return match, true
    }
    return "", false
}


func assertFile(file string) (string, error) {
    if len(file) == 0 {
        var err error
        file, err = GetDefaultJSON()
        if err != nil {
            return "", err
        }
    }

    if err := isJsonFile(file); err != nil {
        fmt.Fprintf(os.Stderr, "Can't read file: %s\n", err)
        os.Exit(1)
    }
    return file, nil
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
