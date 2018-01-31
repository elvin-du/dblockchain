package main

import (
	"fmt"
)

func IntToHex(i int64) []byte {
	return []byte(fmt.Sprintf("%x", i))
}
