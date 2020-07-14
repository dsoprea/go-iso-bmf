package atom

import (
	"github.com/dsoprea/go-logging"
)

// TkhdBox is a "Track Header" box.
type TkhdBox struct {
	Box

	version          byte
	flags            uint32
	creationTime     uint32
	modificationTime uint32
	trackID          uint32
	duration         uint32
	layer            uint16
	alternateGroup   uint16
	volume           Fixed16
	matrix           []byte
	width            Fixed32
	height           Fixed32
}

func (tb *TkhdBox) Version() byte {
	return tb.version
}

func (tb *TkhdBox) Flags() uint32 {
	return tb.flags
}

func (tb *TkhdBox) CreationTime() uint32 {
	return tb.creationTime
}

func (tb *TkhdBox) ModificationTime() uint32 {
	return tb.modificationTime
}

func (tb *TkhdBox) TrackID() uint32 {
	return tb.trackID
}

func (tb *TkhdBox) Duration() uint32 {
	return tb.duration
}

func (tb *TkhdBox) Layer() uint16 {
	return tb.layer
}

func (tb *TkhdBox) AlternateGroup() uint16 {
	return tb.alternateGroup
}

func (tb *TkhdBox) Volume() Fixed16 {
	return tb.volume
}

func (tb *TkhdBox) Matrix() []byte {
	return tb.matrix
}

// Width returns a calculated tkhd width.
func (b *TkhdBox) Width() Fixed32 {
	return b.width / (1 << 16)
}

// Height returns a calculated tkhd height.
func (b *TkhdBox) Height() Fixed32 {
	return b.height / (1 << 16)
}

func (b *TkhdBox) parse() (err error) {
	defer func() {
		if errRaw := recover(); errRaw != nil {
			err = log.Wrap(errRaw.(error))
		}
	}()

	data, err := b.readBoxData()
	log.PanicIf(err)

	b.version = data[0]
	b.flags = defaultEndianness.Uint32(data[0:4])
	b.creationTime = defaultEndianness.Uint32(data[4:8])
	b.modificationTime = defaultEndianness.Uint32(data[8:12])
	b.trackID = defaultEndianness.Uint32(data[12:16])
	b.duration = defaultEndianness.Uint32(data[20:24])
	b.layer = defaultEndianness.Uint16(data[32:34])
	b.alternateGroup = defaultEndianness.Uint16(data[34:36])
	b.volume = fixed16(data[36:38])
	b.matrix = data[40:76]
	b.width = fixed32(data[76:80])
	b.height = fixed32(data[80:84])

	return nil
}

type tkhdBoxFactory struct {
}

// Name returns the name of the type.
func (tkhdBoxFactory) Name() string {
	return "tkhd"
}

// New returns a new value instance.
func (tkhdBoxFactory) New(box Box) (cb CommonBox, err error) {
	defer func() {
		if errRaw := recover(); errRaw != nil {
			err = log.Wrap(errRaw.(error))
		}
	}()

	tkhdBox := &TkhdBox{
		Box: box,
	}

	err = tkhdBox.parse()
	log.PanicIf(err)

	return tkhdBox, nil
}

var (
	_ boxFactory = tkhdBoxFactory{}
	_ CommonBox  = &TkhdBox{}
)

func init() {
	registerAtom(tkhdBoxFactory{})
}
