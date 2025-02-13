# Makefile
.EXPORT_ALL_VARIABLES:

GO111MODULE=on
GOPROXY=direct
GOSUMDB=off
GOPRIVATE=github.com/jeffotoni/quick

update:
	@echo "########## Compilando nossa API ... "
	@rm -f go.*
	go mod init github.com/jeffotoni/quick/examples
	go mod tidy
	@echo "fim"
