package utils

import (
	"fmt"
	"testing"
)

func Test_XOR(t *testing.T) {
	s := "abc123123"
	k := "123ABCabc"

	e := EncryptXORByBytes([]byte(s), []byte(k))

	d := EncryptXORByBytes(e, []byte(k))

	println(fmt.Sprintf("原文:%s, 密文:%s, 解密后:%s", s, e, d))
	println(fmt.Sprintf("原文:%v, 密文:%v, 解密后:%v", []byte(s), e, d))
}
