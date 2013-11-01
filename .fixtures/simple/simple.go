package main

import (
	"log"
	"os"
)

func main() {
	s := &MyStruct{}
	e := NewMyStructJSONEncoder(os.Stdout)
	if err := e.Encode(s); err != nil {
		log.Fatalln("MyStruct encoding error: ", err.Error())
	}
}
