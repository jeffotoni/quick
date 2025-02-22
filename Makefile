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
	@bash ./scripts/test.sh;
	
cover:
	@bash ./scripts/coverage.sh;
	@cd http/client && make cover
	#@rm -f ./coverage.out;
	#@rm -f ./cover.out;

#### install
### go install github.com/securego/gosec/v2/cmd/gosec@latest
gosec:
	gosec ./...

#### install
#### go install github.com/gordonklaus/ineffassign@latest
ineffassign:
	ineffassign ./...


#### install
#### go install honnef.co/go/tools/cmd/staticcheck@latest
staticcheck:
	staticcheck ./...
