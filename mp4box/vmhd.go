package mp4box

import (
	"github.com/dsoprea/go-logging"
)

// VmhdBox is the "Video Media Header" box.
type VmhdBox struct {
	Box

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

	data, err := b.readBoxData()
	log.PanicIf(err)

	b.version = data[0]
	b.flags = defaultEndianness.Uint32(data[0:4])
	b.graphicsMode = defaultEndianness.Uint16(data[4:6])
	b.opColor = defaultEndianness.Uint16(data[6:8])

	return nil
}

type vmhdBoxFactory struct {
}

// Name returns the name of the type.
func (vmhdBoxFactory) Name() string {
	return "vmhd"
}

// New returns a new value instance.
func (vmhdBoxFactory) New(box Box) (cb CommonBox, err error) {
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

	return vmhdBox, nil
}

var (
	_ boxFactory = vmhdBoxFactory{}
	_ CommonBox  = &VmhdBox{}
)

func init() {
	registerAtom(vmhdBoxFactory{})
}
