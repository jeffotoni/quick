#!/bin/bash
echo -ne "\ncoverage starting\n"
go test -v -count=2 -cover -failfast -coverprofile coverage.out ./
go tool cover -html=coverage.out -o coverage.html
echo -ne "\ncoverage completed\n"
