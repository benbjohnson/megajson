COVERPROFILE=/tmp/c.out

default:
	@echo "usage:"
	@echo "    make bench"
	@echo "    make cloc"
	@echo "    make PKG=./pkgname cover"
	@echo "    make fmt"
	@echo

bench:
	go run main.go -- test/code.go
	go test -v -test.bench=. ./test

cloc:
	cloc --not-match-f=_test.go --not-match-d=test .

cover: coverpreq fmt
	go test -v -coverprofile=$(COVERPROFILE) $(PKG)
	go tool cover -html=$(COVERPROFILE)
	rm $(COVERPROFILE)

coverpreq:
	@if [[ -z "$(PKG)" ]]; then \
		echo "usage: make PKG=./mypkg cover"; \
		exit 1; \
	fi

fmt:
	go fmt ./...
