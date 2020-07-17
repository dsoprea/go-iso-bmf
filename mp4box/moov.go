package mp4box

import (
	"github.com/dsoprea/go-logging"
)

// MoovBox is a "Movie" box.
//
// The metadata for a presentation is stored in the single Movie Box which occurs
// at the top-level of a file. Normally this box is close to the beginning or end
// of the file, though this is not required.
type MoovBox struct {
	Box

	// TODO(dustin): Add test for this.
	isFragmented bool

	LoadedBoxIndex
}

// IsFragmented return true if mvex box is present.
func (mv MoovBox) IsFragmented() bool {
	return mv.isFragmented
}

func (b *MoovBox) parse() (err error) {
	defer func() {
		if errRaw := recover(); errRaw != nil {
			err = log.Wrap(errRaw.(error))
		}
	}()

	boxes, err := b.Box.readBoxes(0)
	log.PanicIf(err)

	b.LoadedBoxIndex = boxes.Index()

	_, b.isFragmented = b.LoadedBoxIndex["mvex"]

	return nil
}

type moovBoxFactory struct {
}

// Name returns the name of the type.
func (moovBoxFactory) Name() string {
	return "moov"
}

// New returns a new value instance.
func (moovBoxFactory) New(box Box) (cb CommonBox, err error) {
	defer func() {
		if errRaw := recover(); errRaw != nil {
			err = log.Wrap(errRaw.(error))
		}
	}()

	moovBox := &MoovBox{
		Box: box,
	}

	err = moovBox.parse()
	log.PanicIf(err)

	return moovBox, nil
}

var (
	_ boxFactory = moovBoxFactory{}
	_ CommonBox  = &MoovBox{}
)

func init() {
	RegisterBoxType(moovBoxFactory{})
}
