package format

import (
    "fmt"
    "strings"
    "sref/crossref"
)


func FormatCite(r crossref.Reference) string {
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
