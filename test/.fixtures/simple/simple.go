package main

import (
	"log"
	"os"
)

func main() {
	s := &MyStruct{
		StringX: "foo",
		IntX: 200,
		Int64X: 189273,
		UintX: 2392,
		Uint64X: 172389984,
		Float32X: 182.23,
		Float64X: 19380.1312,
		BoolX: true,
	}
	e := NewMyStructJSONEncoder(os.Stdout)
	if err := e.Encode(s); err != nil {
		log.Fatalln("MyStruct encoding error: ", err.Error())
	}
}
