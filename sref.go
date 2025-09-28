package main

import (
    "errors"
    "io"
    "log"
    "flag"
    "net/http"
    "net/url"
    "fmt"
    "os"
    "path/filepath"
    "sref/db"
    "sref/format"
)

var file string
var doi string
var add bool
var del bool

func main() {
    flag.StringVar(&file, "file", "", "JSON file")
    flag.StringVar(&doi, "doi", "", "DOI to use")
    flag.BoolVar(&add, "a", false, "Add DOI's reference")
    flag.BoolVar(&del, "d", false, "Delete DOI's reference")
    flag.Parse()


    // TODO: unteresting that it fails in get this on top of the search!
    SearchPaper("Plant neighbourhood diversity effects on leaf traits: A meta‐analysis")

    if len(file) == 0 {
        var err error
        file, err = GetDefaultJSON()
        if err != nil {
            log.Fatal(err)
        }
    }

    if err := isJsonFile(file); err != nil {
        fmt.Fprintf(os.Stderr, "Can't read file: %s\n", err)
        os.Exit(1)
    }

    if add {
        if len(doi) == 0 {
            // TODO: validate DOI. At least that looks like a DOI
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

    if del {
        if len(doi) == 0 {
            // TODO: validate DOI. At least that looks like a DOI
            fmt.Fprintf(os.Stderr, "missing DOI\n")
            os.Exit(1)
        }

        if err := db.DeleteReference(file, doi); err != nil {
            fmt.Fprintf(os.Stderr, "Failed fo delete DOI: %s\n", err)
            os.Exit(1)
        }
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
        format.PrintAbstract(r)
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


func SearchPaper(s string) {
    baseURL := fmt.Sprintf("https://api.crossref.org/works")
    params := url.Values{}
    params.Add("query.title", s) 

    requestURL := fmt.Sprintf("%s?%s", baseURL, params.Encode())

    resp, err := http.Get(requestURL)
    if err != nil {
        fmt.Printf("error making http request: %s\n", err)
        os.Exit(1)
     }
    defer resp.Body.Close()

    body, err := io.ReadAll(resp.Body)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println(string(body))
}
