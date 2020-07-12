package atom

import (
	"encoding/binary"

	"github.com/dsoprea/go-logging"
)

// HdlrBox - Handler Reference Box
// Box Type: hdlr
// Container: Media Box (mdia) or Meta Box (meta)
// Mandatory: Yes
// Quantity: Exactly one
type HdlrBox struct {
	*Box

	Version byte
	Flags   uint32
	Handler string
	Name    string
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
	b.Name = string(data[24 : b.Size-BoxHeaderSize])

	return nil
}
