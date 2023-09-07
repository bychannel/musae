package tcpx

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"gitlab.musadisca-games.com/wangxw/musae/framework/baseconf"
	"gitlab.musadisca-games.com/wangxw/musae/framework/errorx"
	"gitlab.musadisca-games.com/wangxw/musae/framework/logger"
	"gitlab.musadisca-games.com/wangxw/musae/framework/utils"
	"io"
	"reflect"
)

var KEEP_PACK_SIZE int32 = 4

// tcpx's tool to help build expected stream for communicating
type Packx struct {
	Marshaller Marshaller
}

// a package scoped packx instance
var packx = NewPackx(nil)
var PackJSON = NewPackx(JsonMarshaller{})
var PackTOML = NewPackx(TomlMarshaller{})
var PackXML = NewPackx(XmlMarshaller{})
var PackYAML = NewPackx(YamlMarshaller{})
var PackProtobuf = NewPackx(ProtobufMarshaller{})

// New a packx instance, specific a marshaller for communication.
// If marshaller is nil, official jsonMarshaller is put to used.
func NewPackx(marshaller Marshaller) *Packx {

	if marshaller == nil {
		marshaller = JsonMarshaller{}
	}
	return &Packx{
		Marshaller: marshaller,
	}
}

// Pack src with specific messageID and optional headers
// Src has not been marshaled yet.Whatever you put as src, it will be marshaled by packx.Marshaller.
func (packx Packx) Pack(messageID int32, errCode int32, src interface{}, keys ...string) ([]byte, error) {
	var (
		e error
	)
	if keys == nil || len(keys) == 0 {
		keys = make([]string, 0)
		//return PackWithMarshaller(Message{MessageID: messageID, SecretKeys: make([]string, 0), Body: src}, packx.Marshaller)
	}

	body := make([]byte, 0)
	if src != nil {
		//if packx.Marshaller == nil {
		//	packx.Marshaller = JsonMarshaller{}
		//}

		body, e = packx.Marshaller.Marshal(src)
		if e != nil {
			return nil, e
		}
	}

	return packx.PackWithBody(messageID, errCode, body, keys...)
	//return PackWithMarshallerAndBody(Message{MessageID: messageID, SecretKeys: keys, Body: nil}, body)
	//return PackWithMarshaller(Message{MessageID: messageID, SecretKeys: keys, Body: src}, packx.Marshaller)
}

// PackWithBody is used for self design protocol
func (packx Packx) PackWithBody(messageID int32, errCode int32, body []byte, keys ...string) ([]byte, error) {
	if keys == nil || len(keys) == 0 {
		//return PackWithMarshallerAndBody(Message{MessageID: messageID, SecretKeys: make([]string, 0), Body: nil}, body)
		keys = make([]string, 0)
	}

	return PackWithMarshallerAndBody(
		Message{
			MessageID:  messageID,
			ErrCode:    errCode,
			SecretKeys: keys,
			Body:       nil,
		}, body)
}

// Unpack
// Stream is a block of length,messageID,headerLength,bodyLength,header,body.
// Dest refers to the body, it can be dynamic by messageID.
//
// Before use this, users should be aware of which struct used as `dest`.
// You can use stream's messageID for judgement like:
// messageID,_:= packx.MessageIDOf(stream)
//
//	switch messageID {
//	    case 1:
//	      packx.Unpack(stream, &struct1)
//	    case 2:
//	      packx.Unpack(stream, &struct2)
//	    ...
//	}
func (packx Packx) Unpack(stream []byte, dest interface{}) (Message, error) {
	return UnpackWithMarshaller(stream, dest, packx.Marshaller)
}

// a stream from a reader can be apart by protocol.
// FirstBlockOf helps tear apart the first block []byte from reader
func (packx Packx) FirstBlockOf(r io.Reader) ([]byte, error) {
	return FirstBlockOf(r)
}
func (packx Packx) FirstBlockOfLimitMaxByte(r io.Reader, maxByte int32) ([]byte, error) {
	if maxByte <= 0 {
		return FirstBlockOf(r)
	}
	return FirstBlockOfLimitMaxByte(r, maxByte)
}

// returns the first block's messageID, header, body marshalled stream, error.
func UnPackFromReader(r io.Reader) (int32, map[string]interface{}, []byte, error) {
	buf, e := UnpackToBlockFromReader(r)
	if e != nil {
		return 0, nil, nil, e
	}

	messageID, e := MessageIDOf(buf)
	if e != nil {
		return 0, nil, nil, e
	}

	header, e := HeaderOf(buf)
	if e != nil {
		return 0, nil, nil, e
	}

	body, e := BodyBytesOf(buf)
	if e != nil {
		return 0, nil, nil, e
	}
	return messageID, header, body, nil
}

// Since FirstBlockOf has nothing to do with packx instance, so make it alone,
// for old usage remaining useful, old packx.FirstBlockOf is still useful
func FirstBlockOf(r io.Reader) ([]byte, error) {
	return UnpackToBlockFromReader(r)
}

func FirstBlockOfLimitMaxByte(r io.Reader, maxByte int32) ([]byte, error) {
	if maxByte <= 0 {
		return UnpackToBlockFromReader(r)
	}
	return UnpackToBlockFromReaderLimitMaxLengthOfByte(r, int(maxByte))
}

// a stream from a buffer which can be apart by protocol.
// FirstBlockOfBytes helps tear apart the first block []byte from a []byte buffer
func (packx Packx) FirstBlockOfBytes(buffer []byte) ([]byte, error) {
	return FirstBlockOfBytes(buffer)
}
func FirstBlockOfBytes(buffer []byte) ([]byte, error) {
	if len(buffer) < 16 {
		return nil, errors.New(fmt.Sprintf("require buffer length more than 16 but got %d", len(buffer)))
	}
	var length = binary.BigEndian.Uint32(buffer[0:4])
	if len(buffer) < 4+int(length) {
		return nil, errors.New(fmt.Sprintf("require buffer length more than %d but got %d", 4+int(length), len(buffer)))

	}
	return buffer[:4+int(length)], nil
}

// messageID of a stream.
// Use this to choose which struct for unpacking.
func (packx Packx) MessageIDOf(stream []byte) (int32, error) {
	return MessageIDOf(stream)
}

// messageID of a stream.
// Use this to choose which struct for unpacking.
func MessageIDOf(stream []byte) (int32, error) {
	if len(stream) < 8 {
		return 0, errors.New(fmt.Sprintf("MessageIDOf stream lenth should be bigger than 8"))
	}
	//if len(stream) >= 12 {
	//	cmdLen := binary.BigEndian.Uint32(stream[8:12])
	//	logger.Debugf("MessageIDOf cmdLen:%d", cmdLen)
	//}

	messageID := binary.BigEndian.Uint32(stream[4:8])
	return int32(messageID), nil
}

// Length of the stream starting validly.
// Length doesn't include length flag itself, it refers to a valid message length after it.
func (packx Packx) LengthOf(stream []byte) (int32, error) {
	return LengthOf(stream)
}

// Length of the stream starting validly.
// Length doesn't include length flag itself, it refers to a valid message length after it.
func LengthOf(stream []byte) (int32, error) {
	//if len(stream) >= 4 {
	//	allLen := binary.BigEndian.Uint32(stream[0:4])
	//	logger.Debugf("allLen:%d", allLen+4)
	//}
	//if len(stream) >= 8 {
	//	crcLen := binary.BigEndian.Uint32(stream[4:8])
	//	logger.Debugf("crcLen:%d", crcLen)
	//}

	if len(stream) < 4 {
		return 0, errors.New(fmt.Sprintf("LengthOf stream lenth should be bigger than 4"))
	}
	length := binary.BigEndian.Uint32(stream[0:4])
	return int32(length), nil
}

// Header length of a stream received
func (packx Packx) HeaderLengthOf(stream []byte) (int32, error) {
	return HeaderLengthOf(stream)
}

// Header length of a stream received
func HeaderLengthOf(stream []byte) (int32, error) {
	if len(stream) < 20 {
		return 0, errors.New(fmt.Sprintf("HeaderLengthOf stream lenth should be bigger than 20"))
	}

	headerLength := binary.BigEndian.Uint32(stream[16:20])
	return int32(headerLength), nil
}

// Body length of a stream received
func (packx Packx) BodyLengthOf(stream []byte) (int32, error) {
	return BodyLengthOf(stream)
}

// Body length of a stream received
func BodyLengthOf(stream []byte) (int32, error) {
	if len(stream) < 24 {
		return 0, errors.New(fmt.Sprintf("BodyLengthOf stream lenth should be bigger than %d", 24))
	}
	bodyLength := binary.BigEndian.Uint32(stream[20:24])
	return int32(bodyLength), nil
}

// request of order-index
func ReqIndexOf(stream []byte) (uint32, error) {
	if len(stream) < 12 {
		return 0, errors.New(fmt.Sprintf("ReqIndexOf stream lenth should be bigger than %d", 12))
	}
	reqIdx := binary.BigEndian.Uint32(stream[8:12])
	logger.Debugf("----- reqIdx -----%d", reqIdx)
	return uint32(reqIdx), nil
}

// order-index length of a stream received
func (packx Packx) ReqIndexLengthOf(stream []byte) (uint32, error) {
	return ReqIndexOf(stream)
}

// CRC length of a stream received
func CRCLengthOf(stream []byte) (int32, error) {
	if len(stream) < 16 {
		return 0, errors.New(fmt.Sprintf("CRCLengthOf stream lenth should be bigger than %d", 16))
	}
	bodyLength := binary.BigEndian.Uint32(stream[12:16])
	return int32(bodyLength), nil
}

// Header bytes of a block
func (packx Packx) HeaderBytesOf(stream []byte) ([]byte, error) {
	return HeaderBytesOf(stream)
}

// Header bytes of a block
func HeaderBytesOf(stream []byte) ([]byte, error) {
	headerLen, e := HeaderLengthOf(stream)
	if e != nil {
		return nil, e
	}
	if len(stream) < 16+int(headerLen) {
		return nil, errors.New(fmt.Sprintf("HeaderBytesOf stream lenth should be bigger than %d", 16+int(headerLen)))
	}
	header := stream[16 : 16+headerLen]
	return header, nil
}

// header of a block
func (packx Packx) HeaderOf(stream []byte) (map[string]interface{}, error) {
	return HeaderOf(stream)
}

// header of a block
func HeaderOf(stream []byte) (map[string]interface{}, error) {
	var header map[string]interface{}
	headerBytes, e := HeaderBytesOf(stream)
	if e != nil {
		return nil, errorx.Wrap(e)
	}
	//wangxw
	if len(headerBytes) == 0 {
		return nil, nil
	}
	e = json.Unmarshal(headerBytes, &header)
	if e != nil {
		return nil, errorx.Wrap(e)
	}
	return header, nil
}

// body bytes of a block
func (packx Packx) BodyBytesOf(stream []byte) ([]byte, error) {
	return BodyBytesOf(stream)
}

// decrpt client stream
func Decrypt(stream []byte, secretKey string) ([]byte, error) {
	//logger.Debugf("----- 收到的数据 -----%v, secretKey:%s", stream, secretKey)
	allLen, err := LengthOf(stream)
	//logger.Debugf("----- allLen -----%d", allLen)
	if err != nil {
		logger.Warnf(err.Error())
		return nil, err
	}

	// encrpt data
	crcLen, err := CRCLengthOf(stream)
	//logger.Debugf("----- crcLen -----%d", crcLen)
	if err != nil {
		logger.Warnf(err.Error())
		return nil, err
	}

	info := stream[:16]
	cryptData := stream[16 : allLen+KEEP_PACK_SIZE]
	//logger.Debugf("----- cryptData -----%d", len(cryptData))

	// crc
	myCRC := utils.GenerateCheckSum(cryptData)
	//logger.Debugf("----- myCRC -----%d", myCRC)
	if int32(myCRC) != crcLen {
		return nil, errors.New(fmt.Sprintf("Decrypt, CRC check faild, client crc:%d, server crc:%d", crcLen, myCRC))
	}

	messageID, err := MessageIDOf(stream)
	if err != nil {
		return nil, err
	}

	if baseconf.GetBaseConf() != nil && baseconf.GetBaseConf().UseEncrypt == 1 && NeedEncrypt(messageID) {
		// 解密
		//logger.Debugf("----- cryptData -----%v", cryptData)
		cryptData = utils.EncryptXORByBytes(cryptData, []byte(secretKey))
		//Logger.Println(fmt.Sprintf("解密后:%v", cryptData))
	}

	return append(info, cryptData...), nil
}

// body bytes of a block
func BodyBytesOf(stream []byte) ([]byte, error) {
	headerLen, e := HeaderLengthOf(stream)
	if e != nil {
		return nil, e
	}
	bodyLen, e := BodyLengthOf(stream)
	if e != nil {
		return nil, e
	}

	//if len(stream) >= 8 {
	//	cmdId := binary.BigEndian.Uint32(stream[4:8])
	//	logger.Debugf("cmdId:%d", cmdId)
	//}
	//if len(stream) >= 12 {
	//	cmdId := binary.BigEndian.Uint32(stream[8:12])
	//	logger.Debugf("reqIdx:%d", cmdId)
	//}
	//if len(stream) >= 16 {
	//	src := binary.BigEndian.Uint32(stream[12:16])
	//	logger.Debugf("clientSrc:%d", src)
	//}
	//if len(stream) >= 20 {
	//	headerLength := binary.BigEndian.Uint32(stream[16:20])
	//	logger.Debugf("headerLength:%d", headerLength)
	//}
	//if len(stream) >= 24 {
	//	bodyLength := binary.BigEndian.Uint32(stream[20:24])
	//	logger.Debugf("bodyLength:%d", bodyLength)
	//}
	//logger.Debugf("stream:%d, headerLen:%d, bodyLen:%d", len(stream), int(headerLen), int(bodyLen))

	if len(stream) < 24+int(headerLen)+int(bodyLen) {
		return nil, errors.New(fmt.Sprintf("BodyBytesOf stream lenth should be bigger than %d", 24+int(headerLen)+int(bodyLen)))
	}
	body := stream[24+headerLen : 24+headerLen+bodyLen]

	return body, nil
}

// PackWithMarshaller will encode message into blocks of length,messageID,headerLength,header,bodyLength,body.
// Users don't need to know how pack serializes itself if users use UnpackPWithMarshaller.
//
// If users want to use this protocol across languages, here are the protocol details:
// (they are ordered as list)
// [0 0 0 24 0 0 0 1 0 0 0 6 0 0 0 6 2 1 19 18 13 11 11 3 1 23 12 132]
// header: [0 0 0 24]
// mesageID: [0 0 0 1]
// headerLength, bodyLength [0 0 0 6]
// header: [2 1 19 18 13 11]
// body: [11 3 1 23 12 132]
// [4]byte -- length             fixed_size,binary big endian encode
// [4]byte -- messageID          fixed_size,binary big endian encode
// [4]byte -- headerLength       fixed_size,binary big endian encode
// [4]byte -- bodyLength         fixed_size,binary big endian encode
// []byte -- header              marshal by json
// []byte -- body                marshal by marshaller
func PackWithMarshaller(message Message, marshaller Marshaller) ([]byte, error) {
	if marshaller == nil {
		marshaller = JsonMarshaller{}
	}
	var e error
	var lengthBuf = make([]byte, 4)
	var messageIDBuf = make([]byte, 4)
	binary.BigEndian.PutUint32(messageIDBuf, uint32(message.MessageID))
	var headerLengthBuf = make([]byte, 4)
	var bodyLengthBuf = make([]byte, 4)
	var headerBuf []byte
	var bodyBuf []byte
	//headerBuf, e = json.Marshal(message.Header)
	if e != nil {
		return nil, e
	}
	binary.BigEndian.PutUint32(headerLengthBuf, uint32(len(headerBuf)))

	if message.Body != nil {
		if e != nil {
			return nil, e
		}
	}

	if message.Body != nil {
		bodyBuf, e = marshaller.Marshal(message.Body)
		if e != nil {
			return nil, e
		}

		if baseconf.GetBaseConf() != nil && baseconf.GetBaseConf().UseEncrypt == 1 {
			secretKeys := message.SecretKeys
			if len(secretKeys) > 0 {
				encryptKey := secretKeys[0]
				//if encryptKey != "" {
				//bodyBuf, e = openssl.AesCBCEncrypt(bodyBuf, []byte(encryptKey), make([]byte, 16), openssl.PKCS7_PADDING)
				bodyBuf = utils.EncryptXORByBytes(bodyBuf, []byte(encryptKey))
				if e != nil {
					return nil, e
				}
				//}
			}
		}
	}

	binary.BigEndian.PutUint32(bodyLengthBuf, uint32(len(bodyBuf)))
	var content = make([]byte, 0, 1024)

	content = append(content, messageIDBuf...)
	content = append(content, headerLengthBuf...)
	content = append(content, bodyLengthBuf...)
	content = append(content, headerBuf...)
	content = append(content, bodyBuf...)

	binary.BigEndian.PutUint32(lengthBuf, uint32(len(content)))

	var packet = make([]byte, 0, 1024)

	packet = append(packet, lengthBuf...)
	packet = append(packet, content...)

	//logger.Debug("PackWithMarshaller, msg:", message, lengthBuf, messageIDBuf, headerLengthBuf, bodyLengthBuf)
	return packet, nil
}

// same as above
func PackWithMarshallerName(message Message, marshallerName string) ([]byte, error) {
	var marshaller Marshaller
	switch marshallerName {
	case "json":
		marshaller = JsonMarshaller{}
	case "xml":
		marshaller = XmlMarshaller{}
	case "toml", "tml":
		marshaller = TomlMarshaller{}
	case "yaml", "yml":
		marshaller = YamlMarshaller{}
	case "protobuf", "proto":
		marshaller = ProtobufMarshaller{}
	default:
		return nil, errors.New("only accept ['json', 'xml', 'toml','yaml','protobuf']")
	}
	return PackWithMarshaller(message, marshaller)
}

// unpack stream from PackWithMarshaller
// If users want to use this protocol across languages, here are the protocol details:
// (they are ordered as list)
// [4]byte -- length             fixed_size,binary big endian encode
// [4]byte -- messageID          fixed_size,binary big endian encode
// [4]byte -- headerLength       fixed_size,binary big endian encode
// [4]byte -- bodyLength         fixed_size,binary big endian encode
// []byte -- header              marshal by json
// []byte -- body                marshal by marshaller
func UnpackWithMarshaller(stream []byte, dest interface{}, marshaller Marshaller) (Message, error) {
	if marshaller == nil {
		marshaller = JsonMarshaller{}
	}
	var e error

	logger.Debugf("stream====>>> %v", stream)
	// 包长
	length := binary.BigEndian.Uint32(stream[0:KEEP_PACK_SIZE])
	stream = stream[0 : length+uint32(KEEP_PACK_SIZE)]
	// messageID
	messageID := binary.BigEndian.Uint32(stream[4:8])
	// reqIdx
	_ = binary.BigEndian.Uint32(stream[8:12])
	// crc
	_ = binary.BigEndian.Uint32(stream[12:16])

	// header长度
	headerLength := binary.BigEndian.Uint32(stream[16:20])
	// body长度
	bodyLength := binary.BigEndian.Uint32(stream[20:24])

	// header
	//var header map[string]interface{}
	//if headerLength != 0 {
	//	e = json.Unmarshal(stream[16:(16+headerLength)], &header)
	//	if e != nil {
	//		return Message{}, e
	//	}
	//}

	// body
	if bodyLength != 0 {
		e = marshaller.Unmarshal(stream[24+headerLength:(24+headerLength+bodyLength)], dest)
		if e != nil {
			return Message{}, e
		}
	}

	return Message{
		MessageID:  int32(messageID),
		SecretKeys: make([]string, 0),
		Body:       reflect.Indirect(reflect.ValueOf(dest)).Interface(),
	}, nil
}

// same as above
func UnpackWithMarshallerName(stream []byte, dest interface{}, marshallerName string) (Message, error) {
	var marshaller Marshaller
	switch marshallerName {
	case "json":
		marshaller = JsonMarshaller{}
	case "xml":
		marshaller = XmlMarshaller{}
	case "toml", "tml":
		marshaller = TomlMarshaller{}
	case "yaml", "yml":
		marshaller = YamlMarshaller{}
	case "protobuf", "proto":
		marshaller = ProtobufMarshaller{}
	default:
		return Message{}, errors.New("only accept ['json', 'xml', 'toml','yaml','protobuf']")
	}
	return UnpackWithMarshaller(stream, dest, marshaller)
}

// unpack the first block from the reader.
// protocol is PackWithMarshaller().
// [4]byte -- length             fixed_size,binary big endian encode
// [4]byte -- messageID          fixed_size,binary big endian encode
// [4]byte -- headerLength       fixed_size,binary big endian encode
// [4]byte -- bodyLength         fixed_size,binary big endian encode
// []byte -- header              marshal by json
// []byte -- body                marshal by marshaller
// ussage:
//
//	for {
//	    blockBuf, e:= UnpackToBlockFromReader(reader)
//		   go func(buf []byte){
//	        // handle a message block apart
//	    }(blockBuf)
//	    continue
//	}
func UnpackToBlockFromReader(reader io.Reader) ([]byte, error) {
	if reader == nil {
		return nil, errors.New("reader is nil")
	}
	var info = make([]byte, 4, 4)
	if e := readUntil(reader, info); e != nil {
		if e == io.EOF {
			return nil, e
		}
		return nil, errorx.Wrap(e)
	}

	length, e := packx.LengthOf(info)
	logger.Debugf("UnpackToBlockFromReader ===> %d", length)
	if e != nil {
		return nil, e
	}
	var content = make([]byte, length, length)
	if e := readUntil(reader, content); e != nil {
		if e == io.EOF {
			return nil, e
		}
		return nil, errorx.Wrap(e)
	}

	return append(info, content...), nil
}

func UnpackToBlockFromReaderLimitMaxLengthOfByte(reader io.Reader, maxByTe int) ([]byte, error) {
	if reader == nil {
		return nil, errors.New("reader is nil")
	}
	var info = make([]byte, 4, 4)
	if e := readUntil(reader, info); e != nil {
		if e == io.EOF {
			return nil, e
		}
		return nil, errorx.Wrap(e)
	}

	length, e := packx.LengthOf(info)
	if e != nil {
		return nil, e
	}

	if length < 0 || length > int32(maxByTe) {
		return nil, errorx.NewFromStringf("recv message beyond max byte length limit(%d), got (%d)", maxByTe, length)
	}

	var content = make([]byte, length, length)
	if e := readUntil(reader, content); e != nil {
		if e == io.EOF {
			return nil, e
		}
		return nil, errorx.Wrap(e)
	}

	val := append(info, content...)
	return val, nil
}

func readUntil(reader io.Reader, buf []byte) error {
	if len(buf) == 0 {
		return nil
	}
	var offset int
	for {
		n, e := reader.Read(buf[offset:])
		if e != nil {
			if e == io.EOF {
				return e
			}
			return errorx.Wrap(e)
		}

		offset += n
		if offset >= len(buf) {
			break
		}
	}
	return nil
}

// This method is used to pack message whose body is well-marshaled.
func PackWithMarshallerAndBody(message Message, body []byte) ([]byte, error) {
	//var e error
	// 包体长度
	var lengthBuf = make([]byte, KEEP_PACK_SIZE)
	// 包头部分
	var messageIDBuf = make([]byte, 4)
	var errCodeBuf = make([]byte, 4)
	var crcBuf = make([]byte, 4)
	var headerLengthBuf = make([]byte, 4)
	var bodyLengthBuf = make([]byte, 4)

	// 包头数据
	var headerBuf []byte
	// 包体数据
	var bodyBuf []byte
	//headerBuf, e = json.Marshal(message.Header)
	//if e != nil {
	//	return nil, e
	//}

	binary.BigEndian.PutUint32(messageIDBuf, uint32(message.MessageID))
	binary.BigEndian.PutUint32(errCodeBuf, uint32(message.ErrCode))

	bodyBuf = body

	binary.BigEndian.PutUint32(headerLengthBuf, uint32(len(headerBuf)))
	binary.BigEndian.PutUint32(bodyLengthBuf, uint32(len(bodyBuf)))
	var content = make([]byte, 0, 1024)

	//content = append(content, messageIDBuf...)
	content = append(content, headerLengthBuf...)
	content = append(content, bodyLengthBuf...)
	//content = append(content, reqIdxBuf...)
	content = append(content, headerBuf...)
	content = append(content, bodyBuf...)

	if baseconf.GetBaseConf() != nil && NeedEncrypt(message.MessageID) {
		secretKey := ""
		secretKeys := message.SecretKeys
		if len(secretKeys) > 0 {
			secretKey = secretKeys[0]
			//logger.Debugf("加密--------, %s", secretKey)
			//if len(secretKey) <= 0 {
			//	return nil, errors.New(fmt.Sprintf("PackWithMarshallerAndBody, do not get secret key value, message.MessageID=%d", message.MessageID))
			//}
			//} else {
			//	return nil, errors.New(fmt.Sprintf("PackWithMarshallerAndBody, do not get secret key value, message.MessageID=%d", message.MessageID))
			content = utils.EncryptXORByBytes(content, []byte(secretKey))
		}

		//logger.Debugf(fmt.Sprintf("messageId=%d, 异或加密", message.MessageID))
		//logger.Debugf(fmt.Sprintf("原文:%v", body))
		//logger.Debugf(fmt.Sprintf("%s, 秘钥:%v", secretKey, []byte(secretKey)))
		//logger.Debugf(fmt.Sprintf("密文:%v", content))
	}
	//Logger.Println(fmt.Sprintf("content:%v", content))

	//Logger.Println(fmt.Sprintf("content: messageIDBuf=%d, headerLengthBuf=%d, bodyLengthBuf=%d",
	//	messageIDBuf, headerLengthBuf, bodyLengthBuf))

	crcLen := utils.GenerateCheckSum(content)
	//Logger.Println(fmt.Sprintf("发送, crc 长度 :%d", crcLen))
	binary.BigEndian.PutUint32(crcBuf, crcLen)

	totalLen :=
		//4 + // + totalLenSize	（记录包体长度的int，不记录到整个包体中）
		4 + //  + cmdIdSize
			4 + // + indexSize
			4 + // + crcLenSize
			uint32(len(content)) // + bodySize

	//Logger.Println(fmt.Sprintf("发送, 数据长度 :%d", totalLen))
	binary.BigEndian.PutUint32(lengthBuf, totalLen)

	var packet = make([]byte, 0, 1024)

	packet = append(packet, lengthBuf...)
	packet = append(packet, messageIDBuf...)
	packet = append(packet, errCodeBuf...)
	packet = append(packet, crcBuf...)
	//packet = append(packet, messageIDBuf...)
	packet = append(packet, content...)
	//logger.Debugf("发送数据[%d]: 数据包长度:%d, 数据长度:%d, crc值:%d, 包体长度:%d",
	//	message.MessageID, totalLen+4, totalLen, crcLen, len(content))
	//logger.Debug("PackWithMarshaller, msg",
	//	message, lengthBuf, messageIDBuf, headerLengthBuf, bodyLengthBuf, crcBuf, len(packet))
	//logger.Debugf("===>>>PackWithMarshallerAndBody len:[%v],data:%v", len(packet), packet)
	return packet, nil
}

// NeedEncrypt 协议数据是否需要加密
func NeedEncrypt(msgId int32) bool {
	conf := baseconf.GetBaseConf()

	if conf.UseEncrypt != 1 {
		return false
	}

	for _, each := range conf.IgnoreEncryptCmdIds {
		if each == msgId {
			return false
		}
	}

	return true
}

func PackHeartbeat() []byte {
	buf, e := PackWithMarshallerAndBody(Message{
		MessageID: DEFAULT_HEARTBEAT_MESSAGEID,
	}, nil)
	if e != nil {
		panic(e)
	}
	return buf
}

// pack short signal which only contains messageID
func PackStuff(messageID int32) []byte {
	buf, e := PackWithMarshallerAndBody(Message{
		MessageID: messageID,
	}, nil)
	if e != nil {
		panic(e)
	}
	return buf
}

func URLPatternOf(stream []byte) (string, error) {
	header, e := HeaderOf(stream)
	if e != nil {
		return "", errorx.Wrap(e)
	}
	str, _, e := headerGetString(header, HEADER_ROUTER_VALUE)
	if e != nil {
		return "", errorx.Wrap(e)
	}
	return str, nil
}

func RouteTypeOf(stream []byte) (string, error) {
	header, e := HeaderOf(stream)
	if e != nil {
		return "", errorx.Wrap(e)
	}
	str, _, e := headerGetString(header, HEADER_ROUTER_KEY)
	if e != nil {
		return "", errorx.Wrap(e)
	}

	return str, nil
}

// pack detail
func Pack(messageID int32, header map[string]interface{}, src interface{}, marshaller Marshaller) ([]byte, error) {
	return PackWithMarshaller(Message{
		MessageID:  messageID,
		SecretKeys: make([]string, 0),
		Body:       src,
	}, marshaller)
}
