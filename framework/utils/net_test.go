package utils

import (
	"fmt"
	"testing"
)

func Test_GetIp(t *testing.T) {
	fmt.Println(ExternalIP())
}
