package atom

import (
	"encoding/binary"

	"github.com/dsoprea/go-logging"
)

// VmhdBox is the "Video Media Header" box.
type VmhdBox struct {
	*Box

	Version      byte
	Flags        uint32
	GraphicsMode uint16
	OpColor      uint16
}

func (b *VmhdBox) parse() (err error) {
	defer func() {
		if errRaw := recover(); errRaw != nil {
			err = log.Wrap(errRaw.(error))
		}
	}()

	data, err := b.readBoxData()
	log.PanicIf(err)

	b.Version = data[0]
	b.Flags = binary.BigEndian.Uint32(data[0:4])
	b.GraphicsMode = binary.BigEndian.Uint16(data[4:6])
	b.OpColor = binary.BigEndian.Uint16(data[6:8])

	return nil
}

type vmhdBoxFactory struct {
}

// Name returns the name of the type.
func (vmhdBoxFactory) Name() string {
	return "vmhd"
}

// New returns a new value instance.
func (vmhdBoxFactory) New(box *Box) (cb CommonBox, err error) {
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
	_ CommonBox  = VmhdBox{}
)

func init() {
	registerAtom(vmhdBoxFactory{})
}
