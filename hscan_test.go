package main

import "encoding/hex"
import "testing"

func TestSha1sum(t *testing.T) {
  s := sha1sum("testfiles/go-logo-blue.svg")
  got := hex.EncodeToString(s[:])
  want := "274fe3dc04269ecb6b5e2a3b659779b8df4bbf07"
  if got != want {
    t.Errorf("got %q want %q", got, want)
  }
}

func TestGetMimeType(t *testing.T) {
  got, err := getMimeType("nope")
  want := ""
  if got != want {
    t.Errorf("got %q want %q", got, want)
  }
  if err == nil {
    t.Errorf("err should be : ERROR reading file nope")
  }  

  got, err = getMimeType("testfiles/go-logo-blue.svg")
  want = "image/svg+xml"
  if err != nil {
    t.Errorf("err %v", err)
  }  
  if got != want {
    t.Errorf("got %q want %q", got, want)
  }

  got, err = getMimeType("testfiles/sha1.txt")
  want = "text/plain; charset=utf-8"
  if err != nil {
    t.Errorf("err %v", err)
  }  
  if got != want {
    t.Errorf("got %q want %q", got, want)
  }

  got, err = getMimeType("/etc/shadow")
  want = ""
  if got != want {
    t.Errorf("got %q want %q", got, want)
  }
}

func TestDirExists(t *testing.T) {
  got := dirExists("nope")
  want := false
  if got != want {
      t.Errorf("got %t want %t", got, want)
  }

  got = dirExists("testfiles/go-logo-blue.svg")
  want = false
  if got != want {
      t.Errorf("got %t want %t", got, want)
  }
  
  got = dirExists("testfiles/")
  want = true
  if got != want {
      t.Errorf("got %t want %t", got, want)
  }
}

func TestLoadChecksumFile(t *testing.T) {
  // sha1
  loadChecksumFile("testfiles/sha1.txt", 3)
  got := arrSha1
  if len(got) != 5 {
    t.Errorf("got length %d instead of 4", len(got))
  }
  if got["64725786589f263f0ecc1da55c2bcac7eb18e681"] != 3 {
      t.Errorf("got %q want %q", got, "64725786589f263f0ecc1da55c2bcac7eb18e681")
  }
  if got["0741e65ae292d5a68c7c167f04d0538254da8e8b"] != 3 {
    t.Errorf("got %q want %q", got, "0741e65ae292d5a68c7c167f04d0538254da8e8b")
  }
  if got["274fe3dc04269ecb6b5e2a3b659779b8df4bbf07"] != 3 {
    t.Errorf("got %q want %q", got, "274fe3dc04269ecb6b5e2a3b659779b8df4bbf07")
  }
  if got["12d81f50767d4e09aa7877da077ad9d1b915d75b"] != 3 {
    t.Errorf("got %q want %q", got, "12d81f50767d4e09aa7877da077ad9d1b915d75b")
  }
  if got["894b7cbc31d7647667b11eb9efe0526d55252711"] != 3 {
    t.Errorf("got %q want %q", got, "894b7cbc31d7647667b11eb9efe0526d55252711")
  }
}
