package atom

import (
	"encoding/binary"

	"github.com/dsoprea/go-logging"
)

// StsdBox - Sample Description Box
// Box Type: stsd
// Container: Sample Table Box (stbl)
// Mandatory: Yes
// Quantity: Exactly one.
type StsdBox struct {
	*Box
	Version byte
	Flags   uint32
	Avc1    *Avc1Box
}

func (b *StsdBox) parse() (err error) {
	defer func() {
		if errRaw := recover(); errRaw != nil {
			err = log.Wrap(errRaw.(error))
		}
	}()

	data, err := b.readBoxData()
	log.PanicIf(err)

	b.Version = data[0]
	b.Flags = binary.BigEndian.Uint32(data[0:4])

	// Skip extra 8 bytes.
	boxes, err := b.Box.readBoxes(8)
	log.PanicIf(err)

	for _, box := range boxes {
		switch box.Name() {
		case "avc1":
			b.Avc1 = &Avc1Box{Box: box}

			err := b.Avc1.parse()
			log.PanicIf(err)
		}
	}

	return nil
}
