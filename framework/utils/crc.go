package utils

// 生成校验和
func GenerateCheckSum(src []byte) uint32 {
	checksum := uint32(0)
	srcLen := len(src)
	wideLen := srcLen - srcLen%2

	for i := 0; i < wideLen; i += 2 {
		checksum = checksum + (uint32(src[i])<<8 | uint32(src[i+1]))
	}

	if srcLen%2 != 0 {
		checksum = checksum + uint32(src[srcLen-1])<<8
	}

	checksum = (checksum >> 16) | (checksum & 0xffff)
	checksum += checksum >> 16
	checksum = 0xffff - checksum

	return checksum
}
