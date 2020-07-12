package atom

import (
	"encoding/binary"

	"github.com/dsoprea/go-logging"
)

// ElstBox - Edit List Box
// Box Type: elst
// Container: Edit Box (edts)
// Mandatory: No
// Quantity: Zero or one
type ElstBox struct {
	*Box

	Version    uint32 // Version of this box.
	EntryCount uint32 // Integer that gives the number of entries.
	Entries    []elstEntry
}

type elstEntry struct {
	SegmentDuration   uint32 // Duration of this edit segment.
	MediaTime         uint32 // Starting time within the media of this edit segment.
	MediaRate         uint16 // Relative rate at which to play the media corresponding to this segment.
	MediaRateFraction uint16
}

func (b *ElstBox) parse() (err error) {
	defer func() {
		if errRaw := recover(); errRaw != nil {
			err = log.Wrap(errRaw.(error))
		}
	}()

	data, err := b.readBoxData()
	log.PanicIf(err)

	b.Version = binary.BigEndian.Uint32(data[0:4])
	b.EntryCount = binary.BigEndian.Uint32(data[4:8])
	b.Entries = make([]elstEntry, b.EntryCount)

	for i := 0; i < len(b.Entries); i++ {
		b.Entries[i].SegmentDuration = binary.BigEndian.Uint32(data[(8 + 12*i):(12 + 12*i)])
		b.Entries[i].MediaTime = binary.BigEndian.Uint32(data[(12 + 12*i):(16 + 12*i)])
		b.Entries[i].MediaRate = binary.BigEndian.Uint16(data[(16 + 12*i):(18 + 12*i)])
		b.Entries[i].MediaRateFraction = binary.BigEndian.Uint16(data[(18 + 12*i):(20 + 12*i)])
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
func (elstBoxFactory) New(box *Box) (cb CommonBox, err error) {
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
	_ CommonBox  = ElstBox{}
)

func init() {
	registerAtom(elstBoxFactory{})
}
