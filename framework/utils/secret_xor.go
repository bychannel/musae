package utils

//
// EncryptXOR 异或加密
//  @param message 明文
//  @param keywords 秘钥
//  @return string 密文
//
func EncryptXOR(message string, keywords string) string {
	return strByXOR(message, keywords)
}

//
// DecryptXOR 异或解密
//  @param enStr 密文
//  @param keywords 秘钥
//  @return string 明文
//
func DecryptXOR(enStr string, keywords string) string {
	return strByXOR(enStr, keywords)
}

// 字符串异或加密
func strByXOR(message string, keywords string) string {
	if keywords == "" {
		return message
	}

	messageLen := len(message)
	keywordsLen := len(keywords)

	result := ""

	for i := 0; i < messageLen; i++ {
		result += string(message[i] ^ keywords[i%keywordsLen])
	}
	return result
}

// EncryptXORByBytes 字节数组异或加密
func EncryptXORByBytes(message, keywords []byte) []byte {
	return bytesByXOR(message, keywords)
}

// 字节数组异或加密
func bytesByXOR(message, keywords []byte) []byte {
	if keywords == nil || len(keywords) <= 0 {
		return message
	}

	messageLen := len(message)
	keywordsLen := len(keywords)

	result := make([]byte, messageLen)

	for i := 0; i < messageLen; i++ {
		result[i] = message[i] ^ keywords[i%keywordsLen]
	}
	return result
}
