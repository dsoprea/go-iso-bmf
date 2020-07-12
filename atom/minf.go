package atom

import (
	"github.com/dsoprea/go-logging"
)

// MinfBox - Media Information Box
// Box Type: minf
// Container: Media Box (mdia)
// Mandatory: Yes
// Quantity: Exactly one.
//
// This box contains all the objects that declare characteristics information of the
// media in the track.
type MinfBox struct {
	*Box

	Vmhd *VmhdBox
	Stbl *StblBox
	Hmhd *HmhdBox
}

func (b *MinfBox) parse() (err error) {
	defer func() {
		if errRaw := recover(); errRaw != nil {
			err = log.Wrap(errRaw.(error))
		}
	}()

	boxes, err := readBoxes(b.File(), b.Start()+BoxHeaderSize, b.Size()-BoxHeaderSize)
	log.PanicIf(err)

	for _, box := range boxes {
		switch box.Name() {
		case "vmhd":
			b.Vmhd = &VmhdBox{Box: box}

			err := b.Vmhd.parse()
			log.PanicIf(err)

		case "stbl":
			b.Stbl = &StblBox{Box: box}

			err := b.Stbl.parse()
			log.PanicIf(err)

		case "hmhd":
			b.Hmhd = &HmhdBox{Box: box}

			err := b.Hmhd.parse()
			log.PanicIf(err)
		}
	}
	return nil
}
