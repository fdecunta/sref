package main

import (
    "errors"
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
    var doi, title, id string
    var add, read, del, toJson, toBib bool

    flag.StringVar(&file, "file", "", "Path to the JSON database file")
    flag.StringVar(&doi, "doi", "", "Paper DOI")
    flag.StringVar(&title, "title", "", "Paper title")
    flag.StringVar(&id, "id", "", "Paper id")
    flag.BoolVar(&add, "a", false, "Add reference to the database")
    flag.BoolVar(&read, "r", false, "Read reference from the database")
    flag.BoolVar(&del, "d", false, "Delete reference from the database")
    flag.BoolVar(&toJson, "json", false, "Print reference(s) in JSON format")
    flag.BoolVar(&toBib, "bib", false, "Print reference(s) in BibTeX format")

    flag.Parse()

    file, err := assertFile(file)
    if err != nil {
        fmt.Println(err)
        flag.Usage()
        os.Exit(1)
    }

    d, err = db.Open(file)
    if err != nil {
        fmt.Fprintln(os.Stderr, "error:", err)
        os.Exit(1)
    }


    if !add && !del && !read && !toJson && !toBib {
        flag.Usage()
        os.Exit(1)
    }

    if toJson {
        if doi != "" || title != "" || id != "" {
            flag.Usage()
            os.Exit(1)
        }

        for _, i := range d.Table {
            s, err := i.ToJson()
            if err != nil {
                fmt.Println("can't format json")
            }
            fmt.Println(s)
        }
        return
    }


    if toBib {
        for _, i := range d.Table {
           fmt.Println(i.ToBibTeX())
        }
        return
    }

    // Accept the input variable and check if already exists
    var r *crossref.Reference
    if doi != "" {
        doi, err = assertDoi(doi, d)
        if err != nil {
            fmt.Println(err)
            os.Exit(1)
        }
        ref, ok := d.Table[doi]
        if ok {
            r = &ref
        } else {
            r = nil
        }
    } else if title != "" {
        r = d.QueryTitle(title)
    } else if id != "" {
        r = d.QueryId(id)
    } 


    if add {
        if r != nil {
            fmt.Println("Reference already exists")
            return 
        }

        if doi != "" {
            r, err = crossref.SearchDoi(doi)
        } else if title != "" {
            r, err = crossref.SearchTitle(title)
        } else {
            fmt.Println("No input provided")
            os.Exit(1)
        }

        if err != nil {
            fmt.Printf("Failed to store reference: %s\n", err)
            os.Exit(1)
        }

        err = d.Store(r)
        if err != nil {
            fmt.Printf("Failed to store reference: %s\n", err)
        }
    
        fmt.Printf("Added %s\n", r.ID)
        return
    }


    // Next operations need r to exist:
    if r == nil {
       fmt.Println("ERROR: reference not found")
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


func GetDefaultJson() (string, error) {
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
        file, err = GetDefaultJson()
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


func assertDoi(s string, d *db.DataBase) (string, error) {
    s = strings.TrimSpace(s)

    doi, ok := CaptureDoi(s)
    if !ok {
        return "", errors.New("Not a valid DOI")
    } 

    return doi, nil
}
