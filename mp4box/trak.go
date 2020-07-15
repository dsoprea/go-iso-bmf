package mp4box

import (
	"github.com/dsoprea/go-logging"
)

// TrakBox is a "Track" box.
type TrakBox struct {
	Box

	// SamplesDuration
	// SamplesSize
	// SampleGroupsInfo

	// chunks []Chunk
	// samples []Sample

	LoadedBoxIndex
}

func (b *TrakBox) parse() (err error) {
	defer func() {
		if errRaw := recover(); errRaw != nil {
			err = log.Wrap(errRaw.(error))
		}
	}()

	boxes, err := b.Box.readBoxes(0)
	log.PanicIf(err)

	b.LoadedBoxIndex = boxes.Index()

	return nil
}

type trakBoxFactory struct {
}

// Name returns the name of the type.
func (trakBoxFactory) Name() string {
	return "trak"
}

// New returns a new value instance.
func (trakBoxFactory) New(box Box) (cb CommonBox, err error) {
	defer func() {
		if errRaw := recover(); errRaw != nil {
			err = log.Wrap(errRaw.(error))
		}
	}()

	trakBox := &TrakBox{
		Box: box,
	}

	err = trakBox.parse()
	log.PanicIf(err)

	return trakBox, nil
}

var (
	_ boxFactory = trakBoxFactory{}
	_ CommonBox  = &TrakBox{}
)

func init() {
	RegisterBoxType(trakBoxFactory{})
}
