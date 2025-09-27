package main

import (
    "errors"
    "fmt"
    "os"
    "strings"
    "sref/crossref"
    "sref/db"
)

func main() {
    if len(os.Args) == 1 {
        fmt.Println("need some doi") 
        os.Exit(1)
    }

    switch os.Args[1] {
        case "-q":
            file := os.Args[2]
            doi := os.Args[3]
            queryDOI(file, doi)
            
        case "-a":
            file := os.Args[2]
            doi := os.Args[3]
            if err := addReference(file, doi); err != nil {
                fmt.Fprintf(os.Stderr, "%s\n", err)
            }
    }
}


func queryDOI(file string, doi string) {
	data, err := db.LoadDB(file)
	if err != nil {
		panic(err)
	}

    var r crossref.Reference
    if ref, ok := data[doi]; ok {
        r = ref
    } else {
        r = crossref.SearchDoi(doi)
    }

    fmt.Println(formatCite(r))
    return
}

func addReference(file string, doi string) error {
   	data, err := db.LoadDB(file)
	if err != nil {
		panic(err)
	}

    if _, ok := data[doi]; ok {
        return errors.New("DOI already exists.")
    } 

    r := crossref.SearchDoi(doi)
    data[doi] = r

	if err := db.SaveDB(file, data); err != nil {
		panic(err)
	}

    fmt.Println(formatCite(r))

    return nil
}


func formatCite(r crossref.Reference) string {
    author := formatAuthor(r)

    return fmt.Sprintf("%s(%d). %s. %s %s(%s). %s",
         author, r.Year, r.Title, r.Journal, r.Volume, r.Page, r.URL)
}


func formatAuthor(r crossref.Reference) string {
    s := ""

    for i, a := range r.Author {
        if len(r.Author) == 1 {
            return fmt.Sprintf("%s, %s.", a.Family, toInitials(a.Given))
        }

        prefix := ""
        suffix := ""
        if i == len(r.Author)-1 {
            prefix += " & "
            suffix += ""
        } else if i == 0 {
            prefix += ""
            suffix += ","
        } else {
            prefix += " "
            suffix += ","
        }
        s += fmt.Sprintf("%s%s, %s%s", prefix, a.Family, toInitials(a.Given), suffix)
    }

    return s
}


func toInitials(n string) string {
    n = strings.Replace(n, ".", "", -1)
    
    s := ""
    for i, name := range strings.Split(n, " ") {
        s += fmt.Sprintf("%s.", name[:1])
        // if not last
        if i != len(strings.Split(n, " "))-1 {
            s += " "
        }
    }

    return s
}
