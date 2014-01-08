package main

type A struct {
    StringX string
    BX *B
    BY *B
    Bn []*B
    Bn2 []*B
}

type B struct {
    Name string
    Age int
}
