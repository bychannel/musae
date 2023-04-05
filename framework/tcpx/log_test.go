package tcpx

import (
	"log"
	"os"
	"testing"
)

var logger2 *Log

func InitLog() {
	logger2 = &Log{
		Logger: log.New(os.Stderr, "[tcpx] ", log.LstdFlags|log.Llongfile),
		Mode:   DEBUG,
	}
}
func TestLog_Println(t *testing.T) {
	InitLog()

	logger2.Println("test-case logger hello")
	logger2.Mode = RELEASE
	logger2.Println("test-case hello")

	logger2.SetLogMode(DEBUG)
	logger2.SetLogFlags(log.Llongfile)

    SetLogFlags(log.Llongfile|log.LUTC)
    SetLogMode(DEBUG)
    Logger.Println("global logger hello")
}

func TestPrintDepth(t *testing.T) {
	InitLog()

	logger2.Println("test-case logger hello")
}
