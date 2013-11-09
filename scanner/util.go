package scanner

import (
	"fmt"
)

func debugf(format string, a ...interface{}) {
	fmt.Printf(format+"\n", a...)
}

func debugln(a ...interface{}) {
	fmt.Println(a...)
}
