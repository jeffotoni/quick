#!/bin/bash
echo -ne "\ncoverage starting\n"
go test -v -count=1 -cover -failfast -coverprofile cover.out ./...
go tool cover -html=cover.out -o coverage.html
echo -ne "\ncoverage completed\n"
