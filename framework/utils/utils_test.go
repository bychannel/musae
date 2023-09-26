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

func TestPrettyJsonLimit(t *testing.T) {
	tests := []struct {
		Name string `json:"name"`
		Args string `json:"args"`
	}{
		{Name: "jack1", Args: "参数111"},
		{Name: "jack2", Args: "参数222"},
		{Name: "jack3", Args: "参数333"},
		{Name: "jack4", Args: "参数444"},
		{Name: "jack5", Args: "参数555"},
		{Name: "jack6", Args: "参数666"},
	}
	fmt.Println(PrettyJsonLimit(tests))
}
