package ziputil

import (
	"bytes"
	"compress/gzip"
	"io"
)

type Gzip struct {
}

func (g *Gzip) ZipEncode(input []byte) ([]byte, error) {
	var (
		err error
		buf bytes.Buffer // 创建一个新的byte输出流
	)
	//gzipWriter := gzip.NewWriter(&buf) // 创建一个新的gzip输出流
	//defer gzipWriter.Close()
	//
	//_, err := gzipWriter.Write(input) // 将 input byte 数组写入到此输出流中
	//gzipWriter.Flush()
	//if err != nil {
	//	return nil, err
	//}
	//gzipWriter.Close()

	w := gzip.NewWriter(&buf)
	_, err = w.Write(input)
	if err != nil {
		return nil, err
	}
	err = w.Close()
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (g *Gzip) ZipDecode(input []byte) ([]byte, error) {
	var (
		err error
		buf = new(bytes.Buffer)
	)

	// 创建一个新的 gzip.Reader
	//bytesReader := bytes.NewReader(input)
	buf.Write(input)
	gzipReader, err := gzip.NewReader(buf)
	if err != nil {
		return nil, err
	}
	defer gzipReader.Close()

	// 从 Reader 中读取出数据
	return io.ReadAll(gzipReader)
}
