package tcpx

import (
	"errors"
	"fmt"
)

// Message contains the necessary parts of tcpx protocol
// MessagID is defining a message routing flag.
// Header is an attachment of a message.
// Body is the message itself, it should be raw message not serialized yet, like "hello", not []byte("hello")
type Message struct {
	MessageID  int32       `json:"message_id"`
	ErrCode    int32       `json:"err_code"`
	SecretKeys []string    `json:"screct_keys"`
	Body       interface{} `json:"body"`
}

func NewMessage(messageID int32, src interface{}) Message {
	return Message{
		MessageID:  messageID,
		SecretKeys: make([]string, 1),
		Body:       src,
	}
}
func NewURLPatternMessage(urlPattern string, src interface{}) Message {
	return Message{
		MessageID:  0,
		SecretKeys: make([]string, 1),
		Body:       src,
	}
}

func (m Message) Pack(marshaller Marshaller) ([]byte, error) {
	return PackWithMarshaller(m, marshaller)
}

// Get value of message's header whose key is 'key'
// Get and Set don't have lock to ensure concurrently safe, which means
// if you should never operate the header in multiple goroutines, it's better to design a context yourself per request
// rather than straightly use message.Header.
func (msg Message) Get(key string) interface{} {
	return errors.New("Message.Get not implemented")
}

// Get and Set don't have lock to ensure concurrently safe, which means
// if you should never operate the header in multiple goroutines, it's better to design a context yourself per request
// rather than straightly use message.Header.
func (msg *Message) Set(k string, v interface{}) {
	fmt.Println(errors.New("Message.Get not implemented"))
}
