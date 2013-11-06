package main

import (
	"fmt"
	"log"
	"strings"
)

// const DATA = `{"StringX":"foo","BX":{"Name":"John","Age":20},"BY":null,"Bn":[{"Name":"Jane","Age":60}],"Bn2":[]}`
const DATA = `{"StringX":"foo","BX":{"Name":"John","Age":20}}`

func main() {
	var obj *A
	d := NewAJSONDecoder(strings.NewReader(DATA))
	if err := d.Decode(&obj); err != nil {
		log.Fatalln("Decoding error: ", err.Error())
	}

	fmt.Print("|")
	fmt.Printf("%v|", obj.StringX)
	fmt.Printf("%v|", obj.BX.Name)
	fmt.Printf("%v|", obj.BX.Age)
}
