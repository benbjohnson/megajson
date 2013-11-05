package main

type MyStruct struct {
    StringX string
    IntX int
    Int64X int64
    UintX uint
    Uint64X uint64
    Float32X float32
    Float64X float64
    BoolX bool
    IgnoreString string `json:"-"`
}

