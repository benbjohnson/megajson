COVERPROFILE=/tmp/c.out

default:
	@echo "usage:"
	@echo "    make bindata"
	@echo "    make bench"
	@echo "    make cloc"
	@echo "    make PKG=./pkgname cover"
	@echo "    make fmt"
	@echo

bindata: generator/encoder_tmpl.go generator/decoder_tmpl.go

generator/encoder_tmpl.go: generator/tmpl/encoder.tmpl
	cat $< | go-bindata -f encoder_tmpl -p generator | gofmt > $@

generator/decoder_tmpl.go: generator/tmpl/decoder.tmpl
	cat $< | go-bindata -f decoder_tmpl -p generator | gofmt > $@


bench: benchpreq
	go test -v -test.bench=. ./bench

bench-cpuprofile: benchpreq
	go test -v -test.bench=. -test.cpuprofile=test/cpu.out ./bench
	go tool pprof test/cpu.out

benchpreq:
	go run main.go -- bench/code.go

cloc:
	cloc --not-match-f=_test.go --not-match-d=test --not-match-d=bench .

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

test:
	go test -i ./...
	go test -v ./...


.PHONY: assets test
