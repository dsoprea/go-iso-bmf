package atom

import (
	"encoding/binary"

	"github.com/dsoprea/go-logging"
)

// HdlrBox is a "Handler Reference" box.
type HdlrBox struct {
	*Box

	Version byte
	Flags   uint32
	Handler string

	name string
}

func (b *HdlrBox) HdlrName() string {
	return b.name
}

func (b *HdlrBox) parse() (err error) {
	defer func() {
		if errRaw := recover(); errRaw != nil {
			err = log.Wrap(errRaw.(error))
		}
	}()

	data, err := b.readBoxData()
	log.PanicIf(err)

	b.Version = data[0]
	b.Flags = binary.BigEndian.Uint32(data[0:4])
	b.Handler = string(data[8:12])
	b.name = string(data[24 : b.Size()-boxHeaderSize])

	return nil
}

type hdlrBoxFactory struct {
}

// Name returns the name of the type.
func (hdlrBoxFactory) Name() string {
	return "hdlr"
}

// New returns a new value instance.
func (hdlrBoxFactory) New(box *Box) (cb CommonBox, err error) {
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
	_ CommonBox  = HdlrBox{}
)

func init() {
	registerAtom(hdlrBoxFactory{})
}
