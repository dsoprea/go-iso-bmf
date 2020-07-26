package bmftype

import (
	"github.com/dsoprea/go-logging"

	"github.com/dsoprea/go-iso-bmf/common"
)

// VmhdBox is the "Video Media Header" box.
type VmhdBox struct {
	bmfcommon.Box

	version      byte
	flags        uint32
	graphicsMode uint16
	opColor      uint16
}

func (vb *VmhdBox) Version() byte {
	return vb.version
}

func (vb *VmhdBox) Flags() uint32 {
	return vb.flags
}

func (vb *VmhdBox) GraphicsMode() uint16 {
	return vb.graphicsMode
}

func (vb *VmhdBox) OpColor() uint16 {
	return vb.opColor
}

func (b *VmhdBox) parse() (err error) {
	defer func() {
		if errRaw := recover(); errRaw != nil {
			err = log.Wrap(errRaw.(error))
		}
	}()

	data, err := b.ReadBoxData()
	log.PanicIf(err)

	b.version = data[0]
	b.flags = bmfcommon.DefaultEndianness.Uint32(data[0:4])
	b.graphicsMode = bmfcommon.DefaultEndianness.Uint16(data[4:6])
	b.opColor = bmfcommon.DefaultEndianness.Uint16(data[6:8])

	return nil
}

type vmhdBoxFactory struct {
}

// Name returns the name of the type.
func (vmhdBoxFactory) Name() string {
	return "vmhd"
}

// New returns a new value instance.
func (vmhdBoxFactory) New(box bmfcommon.Box) (cb bmfcommon.CommonBox, childBoxSeriesOffset int, err error) {
	defer func() {
		if errRaw := recover(); errRaw != nil {
			err = log.Wrap(errRaw.(error))
		}
	}()

	vmhdBox := &VmhdBox{
		Box: box,
	}

	err = vmhdBox.parse()
	log.PanicIf(err)

	return vmhdBox, childBoxSeriesOffset, nil
}

var (
	_ bmfcommon.BoxFactory = vmhdBoxFactory{}
	_ bmfcommon.CommonBox  = &VmhdBox{}
)

func init() {
	bmfcommon.RegisterBoxType(vmhdBoxFactory{})
}
