package main

import (
    "testing"
)

func TestIsEmail(t *testing.T) {
    goodEmail := "foo@gmail.com"
    if !IsEmail(goodEmail) {
        t.Errorf("IsEmail failed to detect correct email")
    }

    badEmail := "badstuff.com"
    if IsEmail(badEmail) {
        t.Errorf("IsEmail failed to detect bad email")
    }
}


func TestCaptureDoi(t *testing.T) {
    s := "https://doi.org/10.1038/nature01014"
    want := "10.1038/nature01014"

    doi, err := CaptureDoi(s)
    if err != nil {
        t.Errorf("CaptureDoi return error with valid DOI")
    }
    if doi != want {
        t.Errorf("CaptureDoi failed to capture valid DOI")
    }

    // ---
    badDoi := "https://www.jstor.org/stable/45135882"
    doi, err = CaptureDoi(badDoi)
    if err == nil {
        t.Errorf("CaptureDoi return ok with bad DOI")
    }
    if doi != "" {
        t.Errorf("CaptureDoi return non empty string with bad DOI")
    }
}

