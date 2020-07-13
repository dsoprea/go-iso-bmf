package atom

import (
	"encoding/binary"

	"github.com/dsoprea/go-logging"
)

// HdlrBox is a "Handler Reference" box.
type HdlrBox struct {
	Box

	version byte
	flags   uint32
	handler string

	hdlrName string
}

func (hb *HdlrBox) Version() byte {
	return hb.version
}

func (hb *HdlrBox) Flags() uint32 {
	return hb.flags
}

func (hb *HdlrBox) Handler() string {
	return hb.handler
}

func (hb *HdlrBox) HdlrName() string {
	return hb.hdlrName
}

func (b *HdlrBox) parse() (err error) {
	defer func() {
		if errRaw := recover(); errRaw != nil {
			err = log.Wrap(errRaw.(error))
		}
	}()

	data, err := b.readBoxData()
	log.PanicIf(err)

	b.version = data[0]
	b.flags = binary.BigEndian.Uint32(data[0:4])
	b.handler = string(data[8:12])
	b.hdlrName = string(data[24 : b.Size()-boxHeaderSize])

	return nil
}

type hdlrBoxFactory struct {
}

// Name returns the name of the type.
func (hdlrBoxFactory) Name() string {
	return "hdlr"
}

// New returns a new value instance.
func (hdlrBoxFactory) New(box Box) (cb CommonBox, err error) {
	defer func() {
		if errRaw := recover(); errRaw != nil {
			err = log.Wrap(errRaw.(error))
		}
	}()

	hdlrBox := &HdlrBox{
		Box: box,
	}

	err = hdlrBox.parse()
	log.PanicIf(err)

	return hdlrBox, nil
}

var (
	_ boxFactory = hdlrBoxFactory{}
	_ CommonBox  = &HdlrBox{}
)

func init() {
	registerAtom(hdlrBoxFactory{})
}
