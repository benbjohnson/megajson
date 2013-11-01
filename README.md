### (This project is currently under development)

# megajson

> High performance Go JSON encoder and decoder.


## Overview

Go's builtin JSON support works great to provide simple, runtime JSON encoding and decoding.
However, it's based on reflection so it has a few drawbacks:

* *Performance* - The reflection library is slow and cannot be optimized by the compiler at compile time.
* *Public Fields Only* - The reflection library can only reflect on exported fields. That means that you can't marshal private fields to JSON.

Megajson is built to get around some of these limitations.
It's a code generation tool that uses the `go/parser` and `go/ast` packages to write custom encoders and decoders for your types.
These encoders and decoders know your types so the reflection package is not necessary.


## Installation

Installing megajson is easy.
Simply `go get` from the command line:

```sh
$ go get github.com/benbjohnson/megajson
```

And you're ready to go.


## Generation

Running megajson is simple.
Just provide the files or directories that you want to generate encoders and decoders for:

```sh
$ megajson mypkg/my_prog.go
```

Two new fields will be generated:

```
mypkg/my_prog_encoder.go
mypkg/my_prog_decoder.go
```

They live in the same package as your `my_prog.go` code so they're ready to go.


## Using

Once your encoders and decoders are generated, you can use them just like the `json.Encoder` and `json.Decoder` except they're named after your types.
For a struct type called `MyStruct`, the generated code can be used like this:

```go
err := NewMyStructEncoder(writer).Encode(val)
err := NewMyStructDecoder(reader).Decode(&val)
```
