package main

import (
    "encoding/json"
    "fmt"
    "log"
    "os"
    "sref/crossref"
)

func main() {
    if len(os.Args) == 1 {
        fmt.Println("need some doi") 
        os.Exit(1)
    }

    switch os.Args[1] {
        case "-q":
      //    doi := os.Args[2]
            doi := "10.1111/j.1461-0248.2008.01192.x"
            newRef := crossref.SearchDoi(doi)
            st := newRef.ToJson()
            fmt.Println(st)
        case "-a":
           addReference(os.Args[2])
        case "-r":
           readJson("data.json")
    }
}


func addReference(doi string) error {
    r := crossref.SearchDoi(doi)

    jsonData := r.ToJson()

    // Open file in append mode
    f, err := os.OpenFile("data.json", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        log.Fatal(err)
    }
    defer f.Close()

    // Write JSON + newline
    if _, err := f.Write(append(jsonData, '\n')); err != nil {
        log.Fatal(err)
    }

    log.Printf("ADD %s\n", r.DOI)

    return nil
}


func readJson(path string) {
    bytes, err := os.ReadFile(path)
    if err != nil {
        panic(err)
    }

    var ref []crossref.Reference
    json.Unmarshal(bytes, &ref)

    fmt.Println(ref)

    for _, p := range ref {
        fmt.Println(p.Title)
    }
}
