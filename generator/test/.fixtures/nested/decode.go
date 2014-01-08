package main

import (
	"fmt"
	"log"
	"strings"
)

const DATA = `{"StringX":"foo","BX":{"Name":"John","Age":20},"BY":null,"Bn":[{"Name":"Jane","Age":60},{"Name":"Jack","Age":-13}]}`

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
	fmt.Printf("%v|", obj.BY)
	fmt.Printf("%v|", len(obj.Bn))
	fmt.Printf("%v|", obj.Bn[0].Name)
	fmt.Printf("%v|", obj.Bn[0].Age)
	fmt.Printf("%v|", obj.Bn[1].Name)
	fmt.Printf("%v|", obj.Bn[1].Age)
}
