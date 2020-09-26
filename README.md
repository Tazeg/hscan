# HSCAN

Scans recursively a path to match given sha1 checksums.
Usefull to find duplicate files, or to find relevant/irrelevant files for which you already have the checksums.

## USAGE

```bash
hscan -d <PATH> -db <PATH>
-d string
      Directory to scan recursively
-db string
      Directory containing text files with sha1 to search (1 checksum by line)
```

## Examples

You have the file `dbpath/sha1.txt` :

```
fed5cdfb1c9b121ea6d042dd54842407df3b4a6b
64725786589f263f0ecc1da55c2bcac7eb18e681
12d81f50767d4e09aa7877da077ad9d1b915d75b
```

Searching for files having those checksums in the directory `test/` :

```bash
hscan -d test -db dbpath

# result :
Loading database file "dbpath/sha1.txt"... 3 uniq checksum found in "46.975µs"

Scanning path "tmp"...
  1964 files - 0 unreadable files - 492 dirs - 0 unreadable dirs - 3 matches

RESULT
  sha1tmp.txt                              : 3 matches
  Total                                    : 3 matches

Done in 292.09673ms
```

The files `hscan_match.log` and `hscan_error.log` are created :

```
# hscan_match.log
64725786589f263f0ecc1da55c2bcac7eb18e681 sha1.txt tmp/runTest.sh
fed5cdfb1c9b121ea6d042dd54842407df3b4a6b sha1.txt tmp/files/CONTRIBUTING.md
12d81f50767d4e09aa7877da077ad9d1b915d75b sha1.txt tmp/test/LICENSE
```

## INSTALL

```bash
git config --global --add url."git@github.com:".insteadOf "https://github.com/"
go get github.com/Tazeg/hscan
cd ~/go/src/github.com/Tazeg/hscan

# Linux
env GOOS=linux GOARCH=amd64 go build hscan.go

# Windows
env GOOS=windows GOARCH=amd64 go build -o hscan.exe hscan.go

# Raspberry Pi
env GOARM=7 GOARCH=arm go build hscan.go

go install
```

## TEST

```bash
go test
```
