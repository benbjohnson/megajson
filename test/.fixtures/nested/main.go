package main

import (
	"log"
	"os"
)

func main() {
	obj := &A{
		StringX: "foo",
		BX: &B{
			Name: "John",
			Age: 20,
		},
		BY: nil,
	}
	e := NewAJSONEncoder(os.Stdout)
	if err := e.Encode(obj); err != nil {
		log.Fatalln("Encoding error: ", err.Error())
	}
}
