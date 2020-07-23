package bmftype

import (
	"github.com/dsoprea/go-logging"

	"github.com/dsoprea/go-iso-bmf/common"
)

// TrakBox is the "Track" box.
type TrakBox struct {
	bmfcommon.Box

	// SamplesDuration
	// SamplesSize
	// SampleGroupsInfo

	// chunks []Chunk
	// samples []Sample

	// LoadedBoxIndex contains this boxes children.
	bmfcommon.LoadedBoxIndex
}

func (b *TrakBox) parse() (err error) {
	defer func() {
		if errRaw := recover(); errRaw != nil {
			err = log.Wrap(errRaw.(error))
		}
	}()

	boxes, err := b.Box.ReadBoxes(0)
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
func (trakBoxFactory) New(box bmfcommon.Box) (cb bmfcommon.CommonBox, err error) {
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
	_ bmfcommon.BoxFactory = trakBoxFactory{}
	_ bmfcommon.CommonBox  = &TrakBox{}
)

func init() {
	bmfcommon.RegisterBoxType(trakBoxFactory{})
}
