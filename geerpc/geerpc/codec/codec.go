package codec

import "io"

type CodecType int

const (
	GobType CodecType = iota
	JsonType
)

type Header struct {
	ServiceMethod string
	Error         string
	Seq           int
}

type Codec interface {
	io.Closer
	ReadHeader(*Header) error
	ReadBody(any) error
	Write(*Header, any) error
}

type newCodecFunc func(cc io.ReadWriteCloser) Codec

var defaultCodecFuncMap map[CodecType]newCodecFunc

func init() {
	defaultCodecFuncMap = make(map[CodecType]newCodecFunc)
	defaultCodecFuncMap[GobType] = newGobCodec
}

func DefaultCodecFuncMap(typ CodecType) newCodecFunc {
	return defaultCodecFuncMap[typ]
}
