package ziputil

type ZipInterface interface {
	ZipEncode(input []byte) ([]byte, error)
	ZipDecode(input []byte) ([]byte, error)
}
