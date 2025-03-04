cover:
	@bash ./coverage.sh;
	@rm -f ./cover.out;
	@rm -f ./coverage.out;

bench:
	go test -bench=. -benchtime=1s -benchmem

gosec:
	gosec ./...

ineffassign:
	ineffassign ./...

staticcheck:
	staticcheck ./...
