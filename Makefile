bench:
	go run main.go -- test/code.go
	go test -v -test.bench=. ./test

cloc:
	cloc --not-match-f=_test.go --not-match-d=test .
