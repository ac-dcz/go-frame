package codec

import (
	"bytes"
	"encoding/gob"
	"io"
)

type GobCodec struct {
	cc      io.ReadWriteCloser
	buff    *bytes.Buffer
	encoder *gob.Encoder
	decoder *gob.Decoder
}

func newGobCodec(cc io.ReadWriteCloser) Codec {
	buff := bytes.NewBuffer(nil)
	return &GobCodec{
		cc:      cc,
		buff:    buff,
		encoder: gob.NewEncoder(buff),
		decoder: gob.NewDecoder(cc),
	}
}

func (c *GobCodec) ReadHeader(h *Header) error {
	return c.decoder.Decode(h)
}

// ReadBody: body must be a pointer
func (c *GobCodec) ReadBody(body any) error {
	return c.decoder.Decode(body)
}

func (c *GobCodec) Write(h *Header, body any) (err error) {
	defer func() {
		if err == nil {
			c.cc.Write(c.buff.Bytes())
		}
		c.buff.Reset()
	}()

	if err = c.encoder.Encode(h); err != nil {
		return err
	}
	if err = c.encoder.Encode(body); err != nil {
		return err
	}
	return
}

func (c *GobCodec) Close() error {
	return c.cc.Close()
}
