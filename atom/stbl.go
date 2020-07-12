package atom

import (
	"github.com/dsoprea/go-logging"
)

// StblBox - Sample Table Box
// Box Type: stbl
// Container: Media Information Box (minf)
// Mandatory: Yes
// Quantity: Exactly one.
type StblBox struct {
	*Box

	Stts *SttsBox
	Stsd *StsdBox
}

func (b *StblBox) parse() (err error) {
	defer func() {
		if errRaw := recover(); errRaw != nil {
			err = log.Wrap(errRaw.(error))
		}
	}()

	boxes, err := readBoxes(b.File, b.Start+BoxHeaderSize, b.Size-BoxHeaderSize)
	log.PanicIf(err)

	for _, box := range boxes {
		switch box.Name {
		case "stts":
			b.Stts = &SttsBox{Box: box}

			err := b.Stts.parse()
			log.PanicIf(err)

		case "stsd":
			b.Stsd = &StsdBox{Box: box}

			err := b.Stsd.parse()
			log.PanicIf(err)
		}
	}
	return nil
}
