cover:
	@bash ./coverage.sh;
	@rm -f ./cover.out;

bench:
	go test -bench=. -benchtime=1s -benchmem

gosec:
	gosec ./...

ineffassign:
	ineffassign ./...

staticcheck:
	staticcheck ./...
