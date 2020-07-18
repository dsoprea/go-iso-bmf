package boxtype

import (
	"fmt"
	"strings"

	"github.com/dsoprea/go-logging"
)

// FtypBox is a file-type box.
type FtypBox struct {
	Box

	// MajorBrand is a brand identifer.
	majorBrand string

	// MinorVersion is an informative integer for the minor version of the
	// major brand.
	minorVersion uint32

	// CompatibleBrands is a list of brands.
	compatibleBrands []string
}

// MajorBrand is a brand identifer.
func (fb *FtypBox) MajorBrand() string {
	return fb.majorBrand
}

// MinorVersion is an informative integer for the minor version of the
// major brand.
func (fb *FtypBox) MinorVersion() uint32 {
	return fb.minorVersion
}

// CompatibleBrands is a list of brands.
func (fb *FtypBox) CompatibleBrands() []string {
	return fb.compatibleBrands
}

// String returns a descriptive string.
func (fb *FtypBox) String() string {
	return fmt.Sprintf("ftyp<%s>", fb.InlineString())
}

// InlineString returns an undecorated string of field names and values.
func (fb *FtypBox) InlineString() string {
	return fmt.Sprintf(
		"%s MAJOR-BRAND=[%s] MINOR-VER=(0x%08x) COMPAT-BRANDS=[%s]",
		fb.Box.InlineString(), fb.majorBrand, fb.minorVersion, strings.Join(fb.compatibleBrands, ","))
}

func (fb *FtypBox) parse() (err error) {
	defer func() {
		if errRaw := recover(); errRaw != nil {
			err = log.Wrap(errRaw.(error))
		}
	}()

	data, err := fb.readBoxData()
	log.PanicIf(err)

	fb.majorBrand = string(data[0:4])
	fb.minorVersion = defaultEndianness.Uint32(data[4:8])

	if len(data) > 8 {
		for i := 8; i < len(data); i += 4 {
			fb.compatibleBrands = append(fb.compatibleBrands, string(data[i:i+4]))
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
	RegisterBoxType(ftypBoxFactory{})
}