package atom

import (
	"encoding/binary"

	"github.com/dsoprea/go-logging"
)

// VmhdBox - Video Media Header Box
// Box Type: vmhd
// Container: Media Information Box (minf)
// Mandatory: Yes
// Quantity: Exactly one specific media header shall be present.
type VmhdBox struct {
	*Box

	Version      byte
	Flags        uint32
	GraphicsMode uint16
	OpColor      uint16
}

func (b *VmhdBox) parse() (err error) {
	defer func() {
		if errRaw := recover(); errRaw != nil {
			err = log.Wrap(errRaw.(error))
		}
	}()

	data, err := b.readBoxData()
	log.PanicIf(err)

	b.Version = data[0]
	b.Flags = binary.BigEndian.Uint32(data[0:4])
	b.GraphicsMode = binary.BigEndian.Uint16(data[4:6])
	b.OpColor = binary.BigEndian.Uint16(data[6:8])

	return nil
}
