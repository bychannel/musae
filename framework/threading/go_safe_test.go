package threading

import (
	"io/ioutil"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRunSafe(t *testing.T) {
	log.SetOutput(ioutil.Discard)

	i := 0

	defer func() {
		assert.Equal(t, 1, i)
	}()

	ch := make(chan int)
	GoSafe(func() {
		defer func() {
			ch <- 111
		}()

		panic("panic")
	})

	<-ch
	i++
}
