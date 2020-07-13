package atom

import (
	"github.com/dsoprea/go-logging"
)

// MdiaBox - Media Box
// Box Type: mdia
// Container: Track Box (trak)
// Mandatory: Yes
// Quantity: Exactly one.
// The mediaa declaration container contains all the objects that declare information
// about the media data within a track.
type MdiaBox struct {
	*Box

	Hdlr *HdlrBox
	Mdhd *MdhdBox
	Minf *MinfBox

	LoadedBoxIndex
}

func (b *MdiaBox) parse() (err error) {
	defer func() {
		if errRaw := recover(); errRaw != nil {
			err = log.Wrap(errRaw.(error))
		}
	}()

	boxes, err := b.Box.readBoxes(0)
	log.PanicIf(err)

	// for _, box := range boxes {
	// 	switch box.Name() {
	// 	case "hdlr":
	// 		b.Hdlr = &HdlrBox{Box: box}

	// 		err := b.Hdlr.parse()
	// 		log.PanicIf(err)

	// 	case "mdhd":
	// 		b.Mdhd = &MdhdBox{Box: box}

	// 		err := b.Mdhd.parse()
	// 		log.PanicIf(err)

	// 	case "minf":
	// 		b.Minf = &MinfBox{Box: box}

	// 		err := b.Minf.parse()
	// 		log.PanicIf(err)
	// 	}
	// }

	b.LoadedBoxIndex = boxes.Index()

	return nil
}

type mdiaBoxFactory struct {
}

// Name returns the name of the type.
func (mdiaBoxFactory) Name() string {
	return "mdia"
}

// New returns a new value instance.
func (mdiaBoxFactory) New(box *Box) (cb CommonBox, err error) {
	defer func() {
		if errRaw := recover(); errRaw != nil {
			err = log.Wrap(errRaw.(error))
		}
	}()

	mdiaBox := &MdiaBox{
		Box: box,
	}

	err = mdiaBox.parse()
	log.PanicIf(err)

	return mdiaBox, nil
}

var (
	_ boxFactory = mdiaBoxFactory{}
	_ CommonBox  = MdiaBox{}
)

func init() {
	registerAtom(mdiaBoxFactory{})
}
