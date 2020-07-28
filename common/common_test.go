package bmfcommon

import (
	"github.com/dsoprea/go-logging"
)

type testBox1 struct {
	Box
}

func (testBox1) InlineString() string {
	return "TestBox1"
}

type testBox1Factory struct {
}

// Name returns the name of the type.
func (testBox1Factory) Name() string {
	return "tb1 "
}

// New returns a new value instance.
func (testBox1Factory) New(box Box) (cb CommonBox, childBoxSeriesOffset int, err error) {
	defer func() {
		if errRaw := recover(); errRaw != nil {
			err = log.Wrap(errRaw.(error))
		}
	}()

	tb1 := &testBox1{
		Box: box,
	}

	return tb1, -1, nil
}

type testBox2 struct {
	Box

	string1 string
	string2 string
}

func (testBox2) InlineString() string {
	return "TestBox2"
}

func (tb2 testBox2) String1() string {
	return tb2.string1
}

func (tb2 testBox2) String2() string {
	return tb2.string2
}

type testBox2Factory struct {
}

// Name returns the name of the type.
func (testBox2Factory) Name() string {
	return "tb2 "
}

// New returns a new value instance.
func (testBox2Factory) New(box Box) (cb CommonBox, childBoxSeriesOffset int, err error) {
	defer func() {
		if errRaw := recover(); errRaw != nil {
			err = log.Wrap(errRaw.(error))
		}
	}()

	tb2 := &testBox2{
		Box: box,
	}

	data, err := tb2.ReadBoxData()
	log.PanicIf(err)

	tb2.string1 = string(data[0:4])
	tb2.string2 = string(data[4:8])

	return tb2, -1, nil
}
