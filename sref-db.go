package main

import (
    "errors"
    "flag"
    "fmt"
    "io"
    "os"
    "path/filepath"
    "regexp"
    "strings"
    "time"

    "sref/db"
    "sref/export"

    "github.com/caltechlibrary/crossrefapi"
)

type Cmd struct {
    file    string
    verb    string
    doi     string
    logFile string
}

type State struct {
    Db        *db.DataBase
    Msg       *crossrefapi.Message
    LogOutput io.Writer
}

func main() {
    var err error
    cmd := Cmd{}
    state := State{nil, nil, nil}
    flag.Usage = usage

    if len(os.Args) < 2 {
        flag.Usage()
        os.Exit(1)
    }

    cmd.verb = os.Args[1]
    fs := flag.NewFlagSet(cmd.verb, flag.ExitOnError)
    fs.StringVar(&cmd.file, "f", "", "Path to JSON file with references")
    fs.StringVar(&cmd.doi, "doi", "", "Paper DOI")
    fs.StringVar(&cmd.logFile, "o", "", "Output log file")
    fs.Parse(os.Args[2:])

    if err := assertFile(&cmd); err != nil {
        fmt.Fprintln(os.Stderr, "error with input file %s", err)
        os.Exit(1)
    }

    state.Db, err = db.Open(cmd.file)
    if err != nil {
        fmt.Fprintln(os.Stderr, "error:", err)
        os.Exit(1)
    }

    // Look for reference in database using -doi. Return nil if not in database
    state.Msg, err = lookupReference(cmd.doi, state.Db)
    if err != nil {
        fmt.Fprintln(os.Stderr, "error during lookupReference:", err)
        os.Exit(1)
    }

    state.LogOutput = os.Stdout
    if cmd.logFile != "" {
        f, err := os.OpenFile(cmd.logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
        if err != nil {
            fmt.Fprintln(os.Stderr, "Cannot open log file %s: %v", cmd.logFile, err)
            os.Exit(1)
        }
        state.LogOutput = io.MultiWriter(os.Stdout, f)
        defer f.Close()
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
    fmt.Println("Usage: sref-db <subcommand> [-f FILE] [-doi DOI] [-o LOGFILE]")
    fmt.Println("Subcommands:")
    fmt.Println("  add     Add a new entry")
    fmt.Println("  read    Read an entry")
    fmt.Println("  del     Delete an entry")
    fmt.Println("  edit    Edit an entry")
    fmt.Println("\nArgs:")
    fmt.Println("  -f      JSON database. Default to ~/.config/sref/references.json")
    fmt.Println("  -doi    Paper DOI")
    fmt.Println("  -o      Output file for logging when adding references")
    fmt.Println("\nExamples:")
    fmt.Println("Add new paper to foo.json:")
    fmt.Println("  sref-db add -f foo.json -doi \"10.1007/s11104-024-06671-1\"")
    fmt.Println("\nPrint all the references from foo.json:")
    fmt.Println("  sref-db read -f foo.json")
}


func Add(cmd *Cmd, state *State) {
    var err error

    email, err := getUserEmail()
    if err != nil {
        fmt.Fprintln(os.Stderr, "Error: %s\nTo configure edit ~/config/sref/email.conf", err)
        os.Exit(1)
    }

    if state.Msg != nil {
        errMsg := ""
        err := fmt.Errorf("Reference already exists")
        if cmd.logFile != "" {
            errMsg = fmt.Sprintf("%s", FormatLog("FAILED", state.Msg.DOI, err))
        } else {
            errMsg = fmt.Sprintf("%s %s", err, state.Msg.DOI)
        }
        fmt.Fprintln(state.LogOutput, errMsg)
        return 
    }

 
    if cmd.doi != "" {
        state.Msg, err = SearchDoi(cmd.doi, email)
    } else {
        fmt.Println("No input provided")
        os.Exit(1)
    }
    
    if err != nil {
        errMsg := ""
        if cmd.logFile != "" {
            var search string
            if cmd.doi != "" {
                search = cmd.doi
            }
            errMsg = fmt.Sprintf("%s", FormatLog("FAILED", search, err))
        } else {
            errMsg = fmt.Sprintf("Failed to find reference: %s", err)
        }

        fmt.Fprintln(state.LogOutput, errMsg)
        os.Exit(1)
    }
    
    err = state.Db.Store(state.Msg)
    if err != nil {
        errMsg := ""
        if cmd.logFile != "" {
            var search string
            if cmd.doi != "" {
                search = cmd.doi
            }
            errMsg = fmt.Sprintf("%s", FormatLog("FAILED", search, err))
        } else {
            errMsg = fmt.Sprintf("Failed to store reference: %s", err)
        }

        fmt.Fprintln(state.LogOutput, errMsg)
        os.Exit(1)
    }

    successMsg := ""
    if cmd.logFile != "" {
        successMsg = fmt.Sprintf("%s", FormatLog("ADDED", state.Msg.DOI, nil))
    } else {
        successMsg = fmt.Sprintf("Added %s", state.Msg.DOI)
    }
    fmt.Fprintln(state.LogOutput, successMsg)
    return
}


func Read(cmd *Cmd, st *State) {
    var toPrint []*crossrefapi.Message

    if cmd.doi != "" {
        if st.Msg == nil {
            fmt.Fprintln(os.Stderr, "reference not found")
            return
        } else {
            toPrint = append(toPrint, st.Msg)
        }
    } else {
        for _, r := range st.Db.Table {
            MsgPtr := &r
            toPrint = append(toPrint, MsgPtr)
        }
    } 

    for _, i := range toPrint {
        s, err := export.Json(i)
        if err != nil {
            fmt.Fprintln(os.Stderr, "can't format json: %v", err)
            continue
        }
        fmt.Println(s)
    }
}


func Delete(st *State) {
    err := st.Db.Delete(st.Msg.DOI)
    if err != nil {
        fmt.Fprintln(os.Stderr, "Failed to delete reference: %s", err)
        os.Exit(1)
    }
    fmt.Printf("Deleted %s\n", st.Msg.DOI)
}


func lookupReference(doi string, d *db.DataBase) (*crossrefapi.Message, error) {
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

    doi = strings.ToLower(doi)

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

    if works == nil {
        return nil, errors.New("NULL response by crossrefapi")
    }

    if works.Status != "ok" {
        return nil, errors.New("request is not ok")
    }
    
    return works.Message, nil
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


func FormatLog(status string, doi string, err error) string {
    timestamp := time.Now().Format("2006-01-02 15:04:05")
    line := fmt.Sprintf("[%s] %s %s", timestamp, status, doi)

    if err != nil {
        line = fmt.Sprintf("%s \"%v\"", line, err)
    }

    return line
}
