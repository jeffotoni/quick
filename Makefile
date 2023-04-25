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
	@rm -f ./coverage.out;
	@rm -f ./cover.out;
