package main

import (
    "errors"
    "flag"
    "fmt"
    "os"
    "path/filepath"
    "regexp"
    "strings"

    "sref/db"
    "sref/export"

    "github.com/caltechlibrary/crossrefapi"
)

type Cmd struct {
    file  string
    verb  string
    doi   string
    title string
}

type State struct {
    db  *db.DataBase
    msg *crossrefapi.Message
}

func main() {
    var err error
    cmd := Cmd{}
    state := State{nil, nil}
    flag.Usage = usage

    if len(os.Args) < 2 {
        flag.Usage()
        os.Exit(1)
    }

    cmd.verb = os.Args[1]
    fs := flag.NewFlagSet(cmd.verb, flag.ExitOnError)
    fs.StringVar(&cmd.file, "f", "", "Path to JSON file with references")
    fs.StringVar(&cmd.doi, "doi", "", "Paper DOI")
    fs.StringVar(&cmd.title, "title", "", "Paper title")
    fs.Parse(os.Args[2:])

    if err := assertFile(&cmd); err != nil {
        fmt.Fprintf(os.Stderr, "error with input file %s\n", err)
        os.Exit(1)
    }

    state.db, err = db.Open(cmd.file)
    if err != nil {
        fmt.Fprintln(os.Stderr, "error:", err)
        os.Exit(1)
    }

    // Look for reference in database using -doi or -title Return nil is not in database
    state.msg, err = lookupReference(cmd.doi, cmd.title, state.db)
    if err != nil {
        fmt.Fprintln(os.Stderr, "error during lookupReference:", err)
        os.Exit(1)
    }

    switch cmd.verb {
    case "add":
        Add(&cmd, &state)
    case "read":
        Read(&cmd, &state)
    case "del":
        Delete(&state)
    case "edit":
        fmt.Println("not implemented yet")
    default:
        fmt.Println("unknown subcommand:", os.Args[1])
        flag.Usage()
        os.Exit(1)
    }
}


func usage() {
    fmt.Println("Usage: sref <subcommand> [options]")
    fmt.Println("Subcommands:")
    fmt.Println("  add     Add a new entry")
    fmt.Println("  read    Read an entry")
    fmt.Println("  del     Delete an entry")
    fmt.Println("  edit    Edit an entry")
    fmt.Println("\nGlobal options:")
    flag.PrintDefaults()
}


func Add(cmd *Cmd, state *State) {
    var err error

    if state.msg != nil {
        fmt.Fprintf(os.Stderr, "Reference already exists: %s\n", state.msg.DOI)
        return 
    }

    email, err := getUserEmail()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error: %s\nTo configure edit ~/config/sref/email.conf\n", err)
        os.Exit(1)
    }
 
    if cmd.doi != "" {
        state.msg, err = SearchDoi(cmd.doi, email)
    } else if cmd.title != "" {
        state.msg, err = SearchTitle(cmd.title, email)
    } else {
        fmt.Println("No input provided")
        os.Exit(1)
    }
    
    if err != nil {
        fmt.Printf("Failed to find reference: %s\n", err)
        os.Exit(1)
    }
    
    err = state.db.Store(state.msg)
    if err != nil {
        fmt.Printf("Failed to store reference: %s\n", err)
        os.Exit(1)
    }
    fmt.Printf("Added %s\n", state.msg.DOI)
}


func Read(cmd *Cmd, st *State) {
    var toPrint []*crossrefapi.Message

    if cmd.doi != "" || cmd.title != "" {
        if st.msg == nil {
            fmt.Fprintf(os.Stderr, "reference not found\n")
            return
        } else {
            toPrint = append(toPrint, st.msg)
        }
    } else {
        for _, r := range st.db.Table {
            msgPtr := &r
            toPrint = append(toPrint, msgPtr)
        }
    } 

    for _, i := range toPrint {
        s, err := export.Json(i)
        if err != nil {
            fmt.Fprintf(os.Stderr, "can't format json: %v\n", err)
            continue
        }
        fmt.Println(s)
    }
}


func Delete(st *State) {
    err := st.db.Delete(st.msg.DOI)
    if err != nil {
        fmt.Fprintf(os.Stderr, "Failed to delete reference: %s\n", err)
        os.Exit(1)
    }
    fmt.Printf("Deleted %s\n", st.msg.DOI)
}


func lookupReference(doi string, title string, d *db.DataBase) (*crossrefapi.Message, error) {
    // Check if reference already exists
    if doi != "" {
        doi, err := CaptureDoi(doi)
        if err != nil {
            return nil, err
        }

        ref, ok := d.Table[doi]
        if ok {
            return &ref, nil
        } 
    } 

    if title != "" {
        return d.QueryTitle(title), nil
    } 

    return nil, nil
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

    configFile := filepath.Join(configDir, "references.json")
    file, err := os.OpenFile(configFile, os.O_RDONLY|os.O_CREATE, 0644)
    if err != nil {
        return "", err
    }
    file.Close()

    return configFile, nil
}


func assertFile(c *Cmd) error {
    if c.file == "" {
        var err error
        c.file, err = GetDefaultJson()
        if err != nil {
            return err
        }
    }

    if _, err := os.Stat(c.file); err != nil {
        if os.IsNotExist(err) {
            return errors.New("File does not exist")
        }
        return err
    } 

    if filepath.Ext(c.file) != ".json" {
        return errors.New("File is not JSON")
    }

    return nil
}


func CaptureDoi(s string) (string, error) {
    s = strings.TrimSpace(s)

    re := regexp.MustCompile(`10\.\d{4,}/\S+`)
    doi := re.FindString(s)
    if doi == "" {
        return "", errors.New("Not a valid DOI")
    }

    return doi, nil
}


func SearchDoi(doi string, email string) (*crossrefapi.Message, error) {
    client, err := crossrefapi.NewCrossRefClient("sref", email)
    if err != nil {
        return nil, err
    }

    works, err := client.Works(doi)
    if err != nil {
        return nil, err
    }
    if works.Status != "ok" {
        return nil, errors.New("request is not ok")
    }
    
    return works.Message, nil
}


func SearchTitle(title string, email string) (*crossrefapi.Message, error) {
    client, err := crossrefapi.NewCrossRefClient("sref", email)
    if err != nil {
        return nil, err
    }

    query := crossrefapi.WorksQuery{
        Fields: &crossrefapi.WorksQueryFields{
            Title: title,
        },
    }
    works, err := client.QueryWorks(query)
   
    if err != nil {
        return nil, err
    }

    if works.Status != "ok" {
        return nil, errors.New("can't reach CrossrefAPI. Status not ok")
    }
    if len(works.Message.Items) == 0 {
        return nil, errors.New("no results for title")
    }

    return &works.Message.Items[0], nil
}


func getUserEmail() (string, error) {
    homeDir, err := os.UserHomeDir()
    if err != nil {
        return "", err
    }

    emailFile := filepath.Join(homeDir, ".config", "sref", "email.conf")
    emailBytes, err := os.ReadFile(emailFile)
    if err != nil {
        return "", err
    }

    email := strings.TrimSpace(string(emailBytes))

    // regexp to catch email
    re := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
    if !re.MatchString(email) {
        return "", errors.New("invalid email")
    }

    return email, nil
}

func IsEmail(s string) bool {
    // regexp to catch email
    re := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
    return re.MatchString(s) 
}
