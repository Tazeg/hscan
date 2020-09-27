package main

import (
  "bufio"
  "crypto/sha1"
  "flag"
  "fmt"
  "github.com/gabriel-vasile/mimetype"
  "github.com/gammazero/workerpool"
  "github.com/saracen/walker"
  "io/ioutil"
  "os"
  "path"
  "strings"
  "sync"
  "time"
)

//-----------------------------------------------------------------------------
// global vars
//-----------------------------------------------------------------------------

// path to scan
var argDir = flag.String("d", "", "Directory to scan recursively")
var argDbSha1 = flag.String("db", "", "Directory containing text files with sha1 to search (1 checksum by line)")

// stats
var nbFiles = 0
var nbDirs = 0
var nbUnreadableDir = 0
var nbUnreadableFile = 0
var nbSha1Match = map[uint8]int{} // nbSha1Match[filename_index] = count matches
var nbTotalSha1Match = 0

// maps of relevant sha1 to look for
var arrSha1 = map[string]uint8{} // arrSha1[strSha1] = db index filename

// log files
var logSha1 *os.File
var logError *os.File

// limited worker pool to calculate hash files to avoid "too many open files"
var wp = workerpool.New(5)

// const
const strVersion = "1.0.1"

// hash databases file names, i.e. checksumFilenames[0]="nsrl.txt"
var checksumFilenames []string

// concurrency to update map
var l = sync.Mutex{}


//-----------------------------------------------------------------------------
// main
//-----------------------------------------------------------------------------

func main() {
  checkArgs()
  var err error
  
  // walk function called for every path found, see https://golang.org/pkg/os/#FileInfo
  walkFn := func(path string, info os.FileInfo) error {
    if info.IsDir() {
      nbDirs++
      return nil
    } 
		
		// skip symbolic links and 0 size files (i.e. /dev/dri/card0)
    if !(info.Mode() & os.ModeSymlink == os.ModeSymlink) && info.Size() > 0 {
      // fmt.Printf("path:%q name:%q size:%d\n", path, info.Name(), info.Size())

      wp.Submit(func() {
        // arbitraty skip files > 238,41 Mb (250000000 b)
        if info.Size() > 250000000 {
          writeErrorLog(fmt.Sprintf("skip file size > 238 Mb : %q\n", path))
          return
        }
        workerPoolAction(path)
      })

      nbFiles++
      showInfos()
    }
    return nil
  }

  // error function called for every error encountered
  errorCallbackOption := walker.WithErrorCallback(func(path string, err error) error {
    writeErrorLog(fmt.Sprintf("could not read file %q: %v\n", path, err))
    nbUnreadableDir++
    return nil
  })

  // create log files
  logSha1, err = os.OpenFile("hscan_match.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
  if err != nil {
    panic(err)
  }
  defer logSha1.Close()
  logError, err = os.OpenFile("hscan_error.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
  if err != nil {
    panic(err)
  }
  defer logError.Close()
  
  // load sha1 files
  loadDbFiles(*argDbSha1)

  // start scan
  startTime := time.Now()
  fmt.Println()
  fmt.Println("Scanning path \"" + *argDir + "\"...")
  walker.Walk(*argDir, walkFn, errorCallbackOption)  

  wp.StopWait()
  endTime := time.Now()
  showInfos() // last refresh
  fmt.Println()
  fmt.Println()
  fmt.Println("RESULT")
  for k, v := range nbSha1Match { 
    fmt.Printf("  %-40s : %d matches\n", checksumFilenames[k], v)
  }
  fmt.Printf("  %-40s : %d matches\n", "Total", nbTotalSha1Match)
  fmt.Println()
  fmt.Println("Done in", endTime.Sub(startTime))
}

func checkArgs() {
  flag.Parse()
  if *argDir == "" {
    showUsage()
    os.Exit(0)
  }
  if *argDbSha1 == "" {
    showUsage()
    os.Exit(0)
  }
  if *argDbSha1 != "" && !dirExists(*argDbSha1) {
    fmt.Printf("ERROR loading databases: %q does not exists or is not a directory\n", *argDbSha1)
    os.Exit(0)
  }
}

// Returns the mime type of a file
// @param {string} filename, full path of a file, ex: "/home/user/file.txt"
// @returns {string} "text/plain" or "" if unknown
func getMimeType(filename string) (string, error) {
  mime, err := mimetype.DetectFile(filename)
  if err != nil {
    return "", err
  }
  return mime.String(), nil
}

func showUsage() {
  fmt.Println("NAME")
  fmt.Println("  hscan")
  fmt.Println("  Look for files recursively matching a list of checksums (sha-1 20 bytes base 16)")
  fmt.Println()
  fmt.Println("VERSION")
  fmt.Printf("  v%s\n", strVersion)
  fmt.Println()
  fmt.Println("USAGE")
  fmt.Println("  hscan -d <PATH> -db <PATH>")
  flag.PrintDefaults()
  fmt.Println("  Results are written in log files in current directory")
  fmt.Println()
  fmt.Println("EXAMPLES")
  fmt.Println("  hscan -db /home/user/sha1files/ -d /mnt/dir/")
  fmt.Println("    Loads text files containing sha1 checksums from the directory /home/user/sha1files.")
  fmt.Println("    Those files must have one checksum per line.")
  fmt.Println("    The path /mnt/dir is scanned recursively to look for matches.")
  fmt.Println()
  fmt.Println("AUTHOR")
  fmt.Println("  Written by Twitter:@JeffProd")
  fmt.Println()
  fmt.Println("LICENCE")
  fmt.Println("  MIT License - Copyright (c) 2020 JeffProd.com")
}

// load txt files from the given path
// each line must contain 1 sha1sum
// @param {string} rootpath "/home/user/dir" or relative path
func loadDbFiles(rootpath string) {
  files, err := ioutil.ReadDir(rootpath)
  if err != nil {
    fmt.Printf("ERROR accessing path %q: %v\n", rootpath, err)
    os.Exit(1)
  }
  cpt := uint8(0)
  for _, file := range files {
    if file.IsDir() { continue } // skip dirs
    strFilename := path.Join(rootpath, file.Name())
    strType, err := getMimeType(strFilename)
    if err != nil {
      fmt.Println(err)
      os.Exit(1)
    }
    if !strings.HasPrefix(strType, "text/plain") { continue } // skip non txt files
    checksumFilenames = append(checksumFilenames, file.Name())
    loadChecksumFile(path.Join(rootpath, file.Name()), cpt)
    cpt++
  }
  if cpt == 0 {
    fmt.Printf("No database text file found in %q\n", rootpath)
    os.Exit(0)
  }
} // loadDbFiles

// progress information
func showInfos() {
  fmt.Printf("\r  %d files - %d unreadable files - %d dirs - %d unreadable dirs - %d matches", nbFiles, nbUnreadableFile, nbDirs, nbUnreadableDir, nbTotalSha1Match)
  // if range nbSha1Match here: concurrent map iteration and map write, so we use nbTotalSha1Match
}

// Calculate SHA1 of a file
// filename, ex: "/home/user/file.txt"
func sha1sum(filename string) [20]byte {
  data, err := ioutil.ReadFile(filename)
  if err != nil {
    writeErrorLog(fmt.Sprintf("could not sha1sum file %q: %v\n", filename, err))
    nbUnreadableFile++
    return [20]byte{}
  }
  return sha1.Sum(data)
}

// dirExists checks if a directory exists
func dirExists(path string) bool {
  info, err := os.Stat(path)
  if os.IsNotExist(err) {
    return false
  }
  if err != nil { return false }
  return info.IsDir()
}

// loads in memory the sha1 checksums from a txt file
// using a map for fastest check if hash key exists instead of parsing a []string
// @params {string} filename "/home/user/toto.txt"
// @params {uint8} index of the file in checksumFilenames
func loadChecksumFile(filename string, idx uint8) {
  startTime := time.Now()
  fmt.Printf("Loading database file %q... ", filename)

  file, err := os.Open(filename)
  if err != nil {
    fmt.Printf("\nError reading the SHA1 checksums file %q: %v\n", filename, err)
    os.Exit(0)
  }
  defer file.Close()

  l := ""
  cpt := 0
  scanner := bufio.NewScanner(file)
  for scanner.Scan() {
    l = strings.TrimSpace(scanner.Text())
    if len(l) > 40 { l = l[0:40] }    
    l = strings.ToLower(l)    
    if l == "" { continue }
    arrSha1[l] = idx
    cpt++
  }

  // init hash match count for this db text file
  nbSha1Match[idx] = 0

  endTime := time.Now()
  fmt.Printf("%d lines, %d uniq checksum found in %q\n", cpt, len(arrSha1), endTime.Sub(startTime))
}

func writeSha1Log(s string) {
  if _, err := logSha1.WriteString(s); err != nil { panic(err) }
}

func writeErrorLog(s string) {
  if _, err := logError.WriteString(s); err != nil { panic(err) }
}

// Action on a file within the worker pool.
// <!> This is a file, not a directory. <!>
// @params {string} filename, i.e. "/home/user/file.txt"
func workerPoolAction(filename string) {
  bSha1 := sha1sum(filename) // binary
  sSha1 := fmt.Sprintf("%x", bSha1) // string

  // sha1 exists in arrSha1 ?
  // info : arrSha1[sSha1] is the db txt file name
  if _, ok := arrSha1[sSha1]; ok {
    // we have a sha1 match
    writeSha1Log(sSha1 + " " + checksumFilenames[arrSha1[sSha1]] + " " + filename + "\n")

		// can't update a map concurrently otherwhite we get "fatal error: concurrent map writes"
    l.Lock()
    nbSha1Match[arrSha1[sSha1]]++
    l.Unlock()

    nbTotalSha1Match++
    showInfos()
    return
  }
}
