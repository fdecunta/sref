package db

import (
    "log"
    "os"
    "path/filepath"
    "testing"
    "github.com/caltechlibrary/crossrefapi"
)


// Helper to create a test database with a temp file
func newTestDB(t *testing.T) (*DataBase, string) {
    tmpDir := t.TempDir()
    path := filepath.Join(tmpDir, "test.json")

    d := DataBase{
        Path: path,
        Table: make(map[string]crossrefapi.Message),
    }

    f, err := os.Open(d.Path)
    if err != nil {
        log.Panic(err)
    }
    defer f.Close()

    return &d, path
}


func TestOpen(t *testing.T) {
    tmpDir := t.TempDir()
    path := filepath.Join(tmpDir, "db.json")      

    // File don't exist, must be error
    _, err := Open(path)
    if err == nil {
        t.Errorf("db.Open don't fail with non-existing file")
    } 
   
}
