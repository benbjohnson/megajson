package main

import (
	"fmt"
	"log"
	"strings"
)

const DATA = `{"StringX":"foo","IntX":200,"Int64X":189273,"myuint":2392,"Uint64X":172389984,"Float32X":182.23,"Float64X":19380.1312,"BoolX":true}`

func main() {
	s := &MyStruct{}

	d := NewMyStructJSONDecoder(strings.NewReader(DATA))
	if err := d.Decode(&s); err != nil {
		log.Fatalln("MyStruct decoding error: ", err.Error())
	}

	fmt.Print("|")
	fmt.Printf("%v|", s.StringX)
	fmt.Printf("%v|", s.IntX)
	fmt.Printf("%v|", s.Int64X)
	fmt.Printf("%v|", s.UintX)
	fmt.Printf("%v|", s.Uint64X)
	fmt.Printf("%v|", s.Float32X)
	fmt.Printf("%v|", s.Float64X)
	fmt.Printf("%v|", s.BoolX)
}
