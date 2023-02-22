# Makefile
.EXPORT_ALL_VARIABLES:

GO111MODULE=on
GOPROXY=direct
GOSUMDB=off
GOPRIVATE=github.com/jeffotoni/quick

update:
	@echo "########## Compilando nossa API ... "
	@rm -f go.*
	go mod init github.com/jeffotoni/quick
	go mod tidy
	@echo "fim"
test: 
	go test -race -v ./...
	go test -v -tags musl -covermode atomic -coverprofile=coverage.out ./...
