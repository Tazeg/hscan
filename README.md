# HSCAN

Scans recursively a path to match given sha1 checksums.
Usefull to find duplicate files, or to find relevant/irrelevant/unknown files.

## USAGE

```bash
hscan -d <PATH> -db <PATH>
-d string
      Directory to scan recursively
-db string
      Directory containing text files with sha1 to search (1 checksum by line)
```

## EXAMPLE

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

Matching files, unknown files, and errors are written in real time into `result.csv` :

```csv
# sha1,dbfile,filename,error
dff8a1731f59ccad056b346102d1e1d014b843f3,nsrl_uniq.txt,/home/jeff/tmp/.vscode/settings.json,
0841f15b7436126cb2877b094d632dbc2707eda0,,/home/jeff/tmp/img_20190502_175115.jpg,
98fb7452234c1d7666a54a53eb7340e501d8c173,sha1test.txt,/home/jeff/tmp/602352874.jpg,
,,/home/jeff/tmp/mysqltmp/undo_001,open /home/jeff/tmp/mysqltmp/undo_001: permission denied
```

A SQLite3 database named `result.db` with the same data as the CSV is created at the end of the process.

## INSTALL

Get the [latest release](https://github.com/Tazeg/hscan/releases) or download and install from source :

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

## BENCHMARKS

Tried on :

- OS : Linux
- HDD : 128 Gb SSD + 2 Tb HDD
- CPU: Intel(R) Xeon(R) CPU E5-1660 v3 @ 3.00GHz
- Memory: 32 Gb

Loading a NIST/NSRL file of 1,2Gb containing 29,459,433 took 22.14s.
Scanning 2Tb and 128 Gb of data took 1h32m34s. This depends on the data stored and the free space on the drive. Further tests will be done shortly.

```bash
$> hscan -d / -db bases_hash/
Loading database file "bases_hash/nsrl_sha1_uniq.txt"... 29459433 uniq checksum found in "22.146464941s"

Scanning path "/"...
  2012574 files - 12091 unreadable files - 274715 dirs - 2510 unreadable dirs - 287870 matches

RESULT
  nsrl_sha1_uniq.txt                       : 287870 matches
  Total                                    : 287870 matches

Done in 1h32m34.505006098s
```
