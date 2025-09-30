package crossref

import (
    "log"
    "strings"

    "github.com/caltechlibrary/crossrefapi"
)

// TODO: Add support for books!

// See the structure of Message here: https://github.com/caltechlibrary/crossrefapi/blob/main/works.go

type Person struct {
    Given  string `json:"given"`
    Family string `json:"family"`
}

type Reference struct {
    Type         string   `json:"type"`
    DOI          string   `json:"doi"`
    Title        string   `json:"title"`
    Author       []Person `json:"author"`
    Year         int      `json:"year"`
    Journal      string   `json:"journal"`
    JournalShort string   `json:"journal-short"`
    Page         string   `json:"page"`
    Volume       string   `json:"volume"`
    URL          string   `json:"url"`
    Abstract     string   `json:"abstract"`
}


func BuildReference(msg *crossrefapi.Message) *Reference {
    var authorsList []Person
    for _, a := range msg.Author {
        authorsList = append(authorsList, Person{a.Given, a.Family})
    }

    r := Reference{
        Type: msg.Type,
        DOI: msg.DOI,
        Title: strings.Join(msg.Title, " "),
        Author: authorsList,
        Year: msg.Issued.DateParts[0][0],
        Journal: strings.Join(msg.ContainerTitle, " "),
        JournalShort: strings.Join(msg.ShortContainerTitle, " "),
        Page: msg.Page,
        Volume: msg.Volume,
        URL: msg.URL,
        Abstract: msg.Abstract,
    }

    return &r
}


func SearchDoi(doi string) *Reference {
    // TODO: don't hardcode the email
    client, err := crossrefapi.NewCrossRefClient("sref", "fdecunta@agro.uba.ar")
    if err != nil {
        log.Fatal(err)
    }
    works, err := client.Works(doi)
   
    if err != nil {
        log.Fatal(err)
    }
    if works.Status != "ok" {
        log.Fatal("request is not ok")
    }
    
    msg := works.Message 

    return BuildReference(msg)
}


func SearchTitle(title string) *Reference {
    // TODO: don't hardcode the email
    client, err := crossrefapi.NewCrossRefClient("sref", "fdecunta@agro.uba.ar")
    if err != nil {
        return nil
    }

    query := crossrefapi.WorksQuery{
        Fields: &crossrefapi.WorksQueryFields{
            Title: title,
        },
    }
    works, err := client.QueryWorks(query)
   
    if err != nil || works == nil || works.Status != "ok" {
        return nil
    }
    if len(works.Message.Items) == 0 {
        return nil
    }

    msg := works.Message.Items[0] 

    return BuildReference(&msg)
}
