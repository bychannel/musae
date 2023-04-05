package utils

import (
	"fmt"
	"github.com/google/uuid"
	"testing"
)

func TestUUID(t *testing.T) {

	ui := uuid.New()

	for i := 0; i < 10; i++ {
		fmt.Println(ui.ID())
	}
	for i := 0; i < 10; i++ {
		fmt.Println(GenStrUUID())
	}

	for i := 0; i < 10; i++ {
		fmt.Println(GenIntUUID())
	}
}
