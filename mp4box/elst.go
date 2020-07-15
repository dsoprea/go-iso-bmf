package mp4box

import (
	"github.com/dsoprea/go-logging"
)

// ElstBox is a "Edit List" box.
type ElstBox struct {
	Box

	version uint32
	entries []elstEntry
}

// Version is the version of this box.
func (eb *ElstBox) Version() uint32 {
	return eb.version
}

// Entries returns the entries.
func (eb *ElstBox) Entries() []elstEntry {
	return eb.entries
}

type elstEntry struct {
	// segmentDuration is the duration of this edit segment.
	segmentDuration uint32

	// mediaTime is the starting time within the media of this edit segment.
	mediaTime uint32

	// mediaRate is the relative rate at which to play the media corresponding
	// to this segment.
	mediaRate uint16

	mediaRateFraction uint16
}

// SegmentDuration is the duration of this edit segment.
func (ee elstEntry) SegmentDuration() uint32 {
	return ee.segmentDuration
}

// MediaTime is the starting time within the media of this edit segment.
func (ee elstEntry) MediaTime() uint32 {
	return ee.mediaTime
}

// MediaRate is the relative rate at which to play the media corresponding to
// this segment.
func (ee elstEntry) MediaRate() uint16 {
	return ee.mediaRate
}

func (ee elstEntry) MediaRateFraction() uint16 {
	return ee.mediaRateFraction
}

func (b *ElstBox) parse() (err error) {
	defer func() {
		if errRaw := recover(); errRaw != nil {
			err = log.Wrap(errRaw.(error))
		}
	}()

	data, err := b.readBoxData()
	log.PanicIf(err)

	b.version = defaultEndianness.Uint32(data[0:4])

	entryCount := int(defaultEndianness.Uint32(data[4:8]))
	b.entries = make([]elstEntry, entryCount)

	for i := 0; i < entryCount; i++ {
		b.entries[i].segmentDuration = defaultEndianness.Uint32(data[(8 + 12*i):(12 + 12*i)])
		b.entries[i].mediaTime = defaultEndianness.Uint32(data[(12 + 12*i):(16 + 12*i)])
		b.entries[i].mediaRate = defaultEndianness.Uint16(data[(16 + 12*i):(18 + 12*i)])
		b.entries[i].mediaRateFraction = defaultEndianness.Uint16(data[(18 + 12*i):(20 + 12*i)])
	}

	return nil
}

type elstBoxFactory struct {
}

// Name returns the name of the type.
func (elstBoxFactory) Name() string {
	return "elst"
}

// New returns a new value instance.
func (elstBoxFactory) New(box Box) (cb CommonBox, err error) {
	defer func() {
		if errRaw := recover(); errRaw != nil {
			err = log.Wrap(errRaw.(error))
		}
	}()

	elstBox := &ElstBox{
		Box: box,
	}

	err = elstBox.parse()
	log.PanicIf(err)

	return elstBox, nil
}

var (
	_ boxFactory = elstBoxFactory{}
	_ CommonBox  = &ElstBox{}
)

func init() {
	registerAtom(elstBoxFactory{})
}
