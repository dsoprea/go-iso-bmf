package bmfcommon

import (
	"github.com/dsoprea/go-logging"
)

// testBox1 has no fields and does not have children.
type testBox1 struct {
	// Box is the base box.
	Box
}

func (*testBox1) InlineString() string {
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

// testBox2 has fields and does not have children.
type testBox2 struct {
	// Box is the base box.
	Box

	string1 string
	string2 string
}

func (*testBox2) InlineString() string {
	return "TestBox2"
}

func (tb2 *testBox2) String1() string {
	return tb2.string1
}

func (tb2 *testBox2) String2() string {
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

// testBox3 has no fields but does have children.
type testBox3 struct {
	// Box is the base box.
	Box

	// LoadedBoxIndex contains this boxes children.
	LoadedBoxIndex
}

// SetLoadedBoxIndex sets the child boxes after a box has been manufactured
// and the children have been parsed. This allows parent boxes to be
// registered before the child boxes can look for them.
func (tb3 *testBox3) SetLoadedBoxIndex(lbi LoadedBoxIndex) {

	// TODO(dustin): !! Add test

	tb3.LoadedBoxIndex = lbi
}

type testBox3Factory struct {
}

// Name returns the name of the type.
func (testBox3Factory) Name() string {
	return "tb3 "
}

// New returns a new value instance.
func (testBox3Factory) New(box Box) (cb CommonBox, childBoxSeriesOffset int, err error) {
	defer func() {
		if errRaw := recover(); errRaw != nil {
			err = log.Wrap(errRaw.(error))
		}
	}()

	tb3 := &testBox3{
		Box: box,
	}

	return tb3, 0, nil
}

func pushTestBox1(b *[]byte) {
	PushBox(b, "tb1 ", nil)
}

func pushTestBox2(b *[]byte, data []byte) {
	PushBox(b, "tb2 ", data)
}

func pushTestBox3(b *[]byte, encodedChildBoxes []byte) {
	PushBox(b, "tb3 ", encodedChildBoxes)
}

func pushUnknownBox(b *[]byte, data []byte) {
	PushBox(b, "wxyz", data)
}

// testBox4 has no fields but does have children.
type testBox4 struct {
	// Box is the base box.
	Box

	// LoadedBoxIndex provides a GetChildBoxes() method that returns a child box
	// if present or panics with a descriptive error.
	LoadedBoxIndex
}

func newTestBox4(box Box, lbi LoadedBoxIndex) *testBox4 {
	return &testBox4{
		Box:            box,
		LoadedBoxIndex: lbi,
	}
}

func (*testBox4) InlineString() string {
	return "TestBox4"
}
