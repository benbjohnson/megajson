PKG=./...
TEST=.
BENCH=.
TMPL=$(wildcard generator/**/*.tmpl)
TMPLBIN=$(patsubst %,%.go,${TMPL})
COVERPROFILE=/tmp/c.out

default:
	@echo "usage:"
	@echo "    make bindata"
	@echo "    make bench"
	@echo "    make cloc"
	@echo "    make PKG=./pkgname cover"
	@echo "    make fmt"
	@echo


bindata: $(TMPLBIN)
	echo $(TMPLBIN)

generator/encoder/encoder.tmpl.go: generator/encoder/encoder.tmpl
	cat $< | go-bindata -func tmplsrc -pkg encoder | gofmt > $@

generator/decoder/decoder.tmpl.go: generator/decoder/decoder.tmpl
	cat $< | go-bindata -func tmplsrc -pkg decoder | gofmt > $@


bench: benchpreq
	go test -v -test.bench=$(BENCH) ./.bench

bench-cpu: benchpreq
	cd .bench; \
	rm ./.bench.test; \
	go test -c; \
	./.bench.test -test.v -test.bench=$(BENCH) -test.benchtime=30s -test.cpuprofile=cpu.prof; \
	go tool pprof .bench.test cpu.prof

benchpreq: bindata
	go run main.go -- .bench/code.go

cloc:
	cloc --not-match-f=_test.go --not-match-d=test --not-match-d=.bench .

cover: coverpreq fmt
	go test -coverprofile=$(COVERPROFILE) $(PKG)
	go tool cover -html=$(COVERPROFILE)
	rm $(COVERPROFILE)

coverpreq:
	@if [[ -z "$(PKG)" ]]; then \
		echo "usage: make PKG=./mypkg cover"; \
		exit 1; \
	fi

fmt:
	go fmt ./...

test: bindata
	go test -i -test.run=$(TEST) $(PKG)
	go test -v -test.run=$(TEST) $(PKG)

goveralls: bindata
	goveralls -package=./... $(COVERALLS_TOKEN)

.PHONY: assets test
