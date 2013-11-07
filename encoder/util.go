package encoder

import (
	"fmt"
)

func debugln(a ...interface{}) (n int, err error) {
	return fmt.Println(a...)
}
