bench:
	go run main.go -- test/code.go
	go test -v -test.bench=. ./test
