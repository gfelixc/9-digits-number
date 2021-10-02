# Getting started

- Make sure to have go 1.17 installed
- Download dependencies `go mod download`
- Run tests `go test -count 1 ./...`

# Building binaries

In order to build binaries just run `CGO_ENABLED=0 go build -o <output_filename> <main_go_path>`
There are two `main.go` files for this application

## Service

Main application on the [requirements](./challenge.md) provided.

Path: `cmd/service/main.go`

## FF-Performance

Spin up 5 concurrent clients during 30 seconds. Each client sends random numbers. 
After the 30 seconds, prints a report of numbers generated, writes failed and average of numbers in 10 secs.
In case of not met the requirements (2.000.000 in 10 secs), a message is printed in Stderr and exit code is 1.
In case of met requirements exit code is 0.

This build is intended to be integrated in the staging/preproduction pipelines 
to make sure performance requirement are met before releasing to prod.

Path: `cmd/ff-performance/main.go`

# Notes

- Due to the short time available I decide to hardcode all the configs (frequency of log writer, frequency of reporter, log file path, ...)    
- Log file path is not configurable, it is generated in the same directory of the binary
- Server is able to handle clients with "server-native newline sequence" from the most common platforms (CRLF/LF) 
- For writing log file the assumption for a "server-native" is a unix-like (LF). Using build tags + os specific newlineSequence would make it platform specific
