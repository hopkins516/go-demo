package common

import (
	"fmt"
	"strings"
)

var Delimiter = func(args ...string) {
	fmt.Print(strings.Repeat("-", 16))
	if len(args) > 0 {
		fmt.Printf("%+v", args)
	}
	fmt.Println(strings.Repeat("-", 16))
}
