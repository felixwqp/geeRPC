package codec

import "io"

type Header struct {
	ServiceMethod string // should follow "Service.Method"
	Seq           uint64 // sequence number chosen by client, reqId for differnet req
	Error         string
}

type Codec interface {
	io.Closer
	ReadHeader(*Header) error
	ReadBody(interface{}) error
	Write(*Header, interface{}) error
}

type Type string

const GobType Type = "application/gob"
const JsonType Type = "application/json"

type NewCodecFunc func(io.ReadWriteCloser) Codec

var NewCodecFuncMap map[Type]NewCodecFunc

func init() {
	NewCodecFuncMap = make(map[Type]NewCodecFunc)
	NewCodecFuncMap[GobType] = NewGobCodec
}
