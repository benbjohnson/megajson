# megajson ![status](https://img.shields.io/badge/status-unmaintained-red.svg)

> High performance Go JSON encoder and decoder.<br/>

> *Notice: This tool is unmaintained. Please use [ffjson](https://github.com/pquerna/ffjson) instead.*

## Overview

Go's builtin JSON support works great to provide simple, runtime JSON encoding and decoding.
However, it's based on reflection so it has a few drawbacks:

* *Performance* - The reflection library is slow and isn't optimized by the compiler at compile time.

* *Public Fields Only* - The reflection library can only reflect on exported fields.
  That means that you can't marshal private fields to JSON.

Megajson is built to get around some of these limitations.
It's a code generation tool that uses the `go/parser` and `go/ast` packages to write custom encoders and decoders for your types.
These encoders and decoders know your types so the reflection package is not necessary.


### Performance

Megajson encodes and decodes at approximately two times the speed of the `encoding/json` package using the built-in `encoding/json` test data in Go 1.2.
This is just a benchmark though.
Your mileage may vary.

Please test megajson encoders using real data for actual results.
This library is primarily focused on performance so performance improvement pull requests are very welcome.


## Installation

Installing megajson is easy.
Simply `go get` from the command line:

```sh
$ go get github.com/benbjohnson/megajson
```

And you're ready to go.


## Usage

Running megajson is simple.
Just provide the files or directories that you want to generate encoders and decoders for:

```sh
$ megajson mypkg/my_file.go
```

Two new files will be generated:

```
mypkg/my_file_encoder.go
mypkg/my_file_decoder.go
```

They live in the same package as your `my_file.go` code so they're ready to go.

Once your encoders and decoders are generated, you can use them just like the `json.Encoder` and `json.Decoder` except they're named after your types.
For a struct type inside `my_file.go` called `MyStruct`, the generated code can be used like this:

```go
err := NewMyStructEncoder(writer).Encode(val)
err := NewMyStructDecoder(reader).Decode(&val)
```


## Supported Types

The following struct field types are supported:

* `string`
* `int`, `int64`
* `uint`, `uint64`
* `float32`, `float64`
* `bool`
* Pointers to structs which have been megajsonified.
* Arrays of pointers to structs which have megajsonified.

If you have a type that you would like to see supported, please add an issue to the GitHub page.

