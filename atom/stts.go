package atom

import (
	"encoding/binary"

	"github.com/dsoprea/go-logging"
)

// SttsBox - Decoding Time to Sample Box
// Box Type: stts
// Container: Sample Table Box (stbl)
// Mandatory: Yes
// Quantity: Exactly one.
type SttsBox struct {
	*Box

	Version      byte
	Flags        uint32
	EntryCount   uint32
	SampleCounts []uint32
	SampleDeltas []uint32
}

func (b *SttsBox) parse() (err error) {
	defer func() {
		if errRaw := recover(); errRaw != nil {
			err = log.Wrap(errRaw.(error))
		}
	}()

	data, err := b.readBoxData()
	log.PanicIf(err)

	b.Version = data[0]
	b.Flags = binary.BigEndian.Uint32(data[0:4])

	count := binary.BigEndian.Uint32(data[4:8])
	b.SampleCounts = make([]uint32, count)
	b.SampleDeltas = make([]uint32, count)

	for i := 0; i < int(count); i++ {
		b.SampleCounts[i] = binary.BigEndian.Uint32(data[(8 + 8*i):(12 + 8*i)])
		b.SampleDeltas[i] = binary.BigEndian.Uint32(data[(12 + 8*i):(16 + 8*i)])
	}

	return nil
}

type sttsBoxFactory struct {
}

// Name returns the name of the type.
func (sttsBoxFactory) Name() string {
	return "stts"
}

// New returns a new value instance.
func (sttsBoxFactory) New(box *Box) (cb CommonBox, err error) {
	defer func() {
		if errRaw := recover(); errRaw != nil {
			err = log.Wrap(errRaw.(error))
		}
	}()

	sttsBox := &SttsBox{
		Box: box,
	}

	err = sttsBox.parse()
	log.PanicIf(err)

	return sttsBox, nil
}

var (
	_ boxFactory = sttsBoxFactory{}
	_ CommonBox  = SttsBox{}
)

func init() {
	registerAtom(sttsBoxFactory{})
}
