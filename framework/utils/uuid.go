package utils

import (
	"github.com/google/uuid"
)

func GenStrUUID() string {
	return uuid.New().String()
}

func GenIntUUID() uint32 {
	return uuid.New().ID()
}
