package geerpc

import "geerpc/codec"

const MagicNum = 0xbeffffeb

type Option struct {
	SecretKey uint
	CodecType codec.CodecType
}

var DefaultOption = Option{
	SecretKey: MagicNum,
	CodecType: codec.GobType,
}

type Server struct {
}
