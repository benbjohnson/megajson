COVERPROFILE=/tmp/c.out

default:
	@echo "usage:"
	@echo "    make bench"
	@echo "    make cloc"
	@echo "    make PKG=./pkgname cover"
	@echo "    make fmt"
	@echo

assets: generator/decoder_tmpl.go

generator/decoder_tmpl.go: generator/tmpl/decoder.tmpl
	cat $< | go-bindata -f decoder_tmpl -p generator | gofmt > $@


.PHONY: assets

bench: benchpreq
	go test -v -test.bench=. ./test

bench-cpuprofile: benchpreq
	go test -v -test.bench=. -test.cpuprofile=test/cpu.out ./test
	go tool pprof test/cpu.out

benchpreq:
	go run main.go -- test/code.go

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

test:
	go test -v ./...

.PHONY: test
