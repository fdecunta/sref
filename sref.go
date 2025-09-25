package main

import (
    "fmt"

    "sref/crossref"
)

func main() {
    newRef := crossref.SearchDoi("10.1111/j.1461-0248.2008.01192.x")

    fmt.Printf("authors. (%d). %s. %s, %s:%s. %s\n",
         newRef.Year, newRef.Title, newRef.Journal, newRef.Volume, newRef.Page, newRef.URL)

    for _, a := range newRef.Author {
        fmt.Printf("%s, %s\n", a.Given, a.Family)
    }

    st := newRef.ToJson()
    fmt.Println(st)
}
