package atom

import (
	"fmt"
	"strings"

	"encoding/binary"

	"github.com/dsoprea/go-logging"
)

// FtypBox is a file-type box.
type FtypBox struct {
	Box

	// MajorBrand is a brand identifer.
	MajorBrand string

	// MinorVersion is an informative integer for the minor version of the
	// major brand.
	MinorVersion uint32

	// CompatibleBrands is a list of brands.
	CompatibleBrands []string
}

// String returns a descriptive string.
func (fb *FtypBox) String() string {
	return fmt.Sprintf("ftyp<%s>", fb.InlineString())
}

// InlineString returns an undecorated string of field names and values.
func (fb *FtypBox) InlineString() string {
	return fmt.Sprintf("%s MAJOR-BRAND=[%s] MINOR-VER=(%d) COMPAT-BRANDS=[%s]", fb.Box.InlineString(), fb.MajorBrand, fb.MinorVersion, strings.Join(fb.CompatibleBrands, ","))
}

func (fb *FtypBox) parse() (err error) {
	defer func() {
		if errRaw := recover(); errRaw != nil {
			err = log.Wrap(errRaw.(error))
		}
	}()

	data, err := fb.readBoxData()
	log.PanicIf(err)

	fb.MajorBrand = string(data[0:4])
	fb.MinorVersion = binary.BigEndian.Uint32(data[4:8])

	if len(data) > 8 {
		for i := 8; i < len(data); i += 4 {
			fb.CompatibleBrands = append(fb.CompatibleBrands, string(data[i:i+4]))
		}
	}

	return nil
}

type ftypBoxFactory struct {
}

// Name returns the name of the type.
func (ftypBoxFactory) Name() string {
	return "ftyp"
}

// New returns a new value instance.
func (ftypBoxFactory) New(box Box) (cb CommonBox, err error) {
	defer func() {
		if errRaw := recover(); errRaw != nil {
			err = log.Wrap(errRaw.(error))
		}
	}()

	ftypBox := &FtypBox{
		Box: box,
	}

	err = ftypBox.parse()
	log.PanicIf(err)

	return ftypBox, nil
}

var (
	_ boxFactory = ftypBoxFactory{}
	_ CommonBox  = &FtypBox{}
)

func init() {
	registerAtom(ftypBoxFactory{})
}
