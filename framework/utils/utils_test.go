package utils

import (
	"fmt"
	"testing"
)

func Test_CheckPath(t *testing.T) {
	bTrue := PathExists("log")
	if bTrue {
		fmt.Println("exist")
	} else {
		fmt.Println("not exist")

	}
}
