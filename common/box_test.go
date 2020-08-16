package bmfcommon

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/dsoprea/go-logging"
	"github.com/dsoprea/go-utility/filesystem"
)

func TestNewBox(t *testing.T) {
	resource := &BmfResource{}
	box := NewBox("name", 1, 2, 3, resource)

	if box.name != "name" {
		t.Fatalf("name field not correct.")
	} else if box.start != 1 {
		t.Fatalf("start field not correct.")
	} else if box.size != 2 {
		t.Fatalf("size field not correct.")
	} else if box.headerSize != 3 {
		t.Fatalf("headerSize field not correct.")
	} else if box.resource != resource {
		t.Fatalf("resource field not correct.")
	}
}

func TestBox_InlineString(t *testing.T) {
	resource := &BmfResource{}
	box := NewBox("name", 1, 2, 3, resource)

	if box.InlineString() != "NAME=[name] PARENT=[ROOT] START=(0x0000000000000001) SIZE=(2)" {
		t.Fatalf("InlineString() not correct: [%s]", box.InlineString())
	}
}

func TestBox_Name(t *testing.T) {
	resource := &BmfResource{}
	box := NewBox("name", 1, 2, 3, resource)

	if box.Name() != "name" {
		t.Fatalf("Name() not correct: [%s]", box.Name())
	}
}

func TestBox_Start(t *testing.T) {
	resource := &BmfResource{}
	box := NewBox("name", 1, 2, 3, resource)

	if box.Start() != 1 {
		t.Fatalf("Start() not correct: (%d)", box.Start())
	}
}

func TestBox_Size(t *testing.T) {
	resource := &BmfResource{}
	box := NewBox("name", 1, 2, 3, resource)

	if box.Size() != 2 {
		t.Fatalf("Size() not correct: (%d)", box.Size())
	}
}

func TestBox_HeaderSize(t *testing.T) {
	resource := &BmfResource{}
	box := NewBox("name", 1, 2, 3, resource)

	if box.HeaderSize() != 3 {
		t.Fatalf("HeaderSize() not correct: (%d)", box.HeaderSize())
	}
}

func TestBox_Parent(t *testing.T) {
	resource := &BmfResource{}
	parentBox := NewBox("parent", 1, 2, 3, resource)

	tb := &testBox1{
		Box: parentBox,
	}

	childBox := NewBox("child", 1, 2, 3, resource)
	childBox.parent = tb

	if childBox.Parent() != tb {
		t.Fatalf("Parent() not correct: [%s]", childBox.Parent())
	}
}

func TestBox_Index(t *testing.T) {
	resource := &BmfResource{}
	resource.fullBoxIndex = make(FullBoxIndex)

	box := NewBox("name", 1, 2, 3, resource)

	if reflect.DeepEqual(box.Index(), resource.fullBoxIndex) != true {
		t.Fatalf("Index() not correct.")
	}
}

func TestBox_ReadBytesAt(t *testing.T) {
	data := []byte{
		0, 0, 0, 0,
		1, 2, 3, 4, 5,
		0, 0, 0, 0, 0,
	}

	sb := rifs.NewSeekableBufferWithBytes(data)

	resource, err := NewBmfResource(sb, 0)
	log.PanicIf(err)

	box := NewBox("name", 1, 2, 3, resource)

	extracted, err := box.ReadBytesAt(4, 5)
	log.PanicIf(err)

	if bytes.Equal(extracted, data[4:9]) != true {
		t.Fatalf("Extracted data not correct.")
	}
}

func TestBox_CopyBytesAt(t *testing.T) {
	data := []byte{
		0, 0, 0, 0,
		1, 2, 3, 4, 5,
		0, 0, 0, 0, 0,
	}

	sb := rifs.NewSeekableBufferWithBytes(data)

	resource, err := NewBmfResource(sb, 0)
	log.PanicIf(err)

	box := NewBox("name", 1, 2, 3, resource)

	b := new(bytes.Buffer)

	err = box.CopyBytesAt(4, 5, b)
	log.PanicIf(err)

	if bytes.Equal(b.Bytes(), data[4:9]) != true {
		t.Fatalf("Extracted data not correct.")
	}
}

func TestBox_ReadBoxes(t *testing.T) {
	ClearRegistrations()
	defer ClearRegistrations()

	RegisterBoxType(testBox1Factory{})
	RegisterBoxType(testBox2Factory{})

	var b []byte

	var data2 []byte
	data2 = append(data2, 'a', 'b', 'c', 'd')
	data2 = append(data2, 'e', 'f', 'g', 'h')

	pushTestBox2(&b, data2)
	pushTestBox1(&b)

	sb := rifs.NewSeekableBufferWithBytes(b)

	resource, err := NewBmfResource(sb, 0)
	log.PanicIf(err)

	box := NewBox("root", 0, int64(len(b)), 0, resource)

	boxes, err := box.ReadBoxes(0, nil)
	log.PanicIf(err)

	if len(boxes) != 2 {
		t.Fatalf("The number of boxes is not correct.")
	} else if boxes[0].Name() != "tb2 " {
		t.Fatalf("The first box is not correct.")
	} else if boxes[1].Name() != "tb1 " {
		t.Fatalf("The second box is not correct.")
	}

	tb2 := boxes[0].(*testBox2)

	if tb2.String1() != "abcd" {
		t.Fatalf("The first string is not correct.")
	} else if tb2.String2() != "efgh" {
		t.Fatalf("The second string is not correct.")
	}
}

func TestBox_ReadBoxData(t *testing.T) {
	// Load stream.

	ClearRegistrations()
	defer ClearRegistrations()

	RegisterBoxType(testBox1Factory{})
	RegisterBoxType(testBox2Factory{})

	var b []byte

	var data2 []byte
	data2 = append(data2, 'a', 'b', 'c', 'd')
	data2 = append(data2, 'e', 'f', 'g', 'h')
	PushBox(&b, "tb2 ", data2)

	PushBox(&b, "tb1 ", nil)

	// Parse.

	sb := rifs.NewSeekableBufferWithBytes(b)

	resource, err := NewBmfResource(sb, int64(len(b)))
	log.PanicIf(err)

	entries1 := resource.LoadedBoxIndex["tb2 "]
	recovered1, err := entries1[0].ReadBoxData()
	log.PanicIf(err)

	if bytes.Equal(recovered1, data2) != true {
		t.Fatalf("Recovered data 1 not correct.")
	}

	entries2 := resource.LoadedBoxIndex["tb1 "]
	recovered2, err := entries2[0].ReadBoxData()
	log.PanicIf(err)

	if len(recovered2) != 0 {
		t.Fatalf("Recovered data 2 not correct.")
	}
}

func TestBox_ReadBoxData_TooSmall(t *testing.T) {
	box := Box{
		name:       "test",
		size:       4,
		headerSize: 8,
	}

	_, err := box.ReadBoxData()
	if err == nil {
		t.Fatalf("Expected error.")
	} else if err.Error() != "box [test] total-size (4) is smaller then box header-size (8)" {
		t.Fatalf("Error not correct: [%s]", err.Error())
	}
}

func TestBoxes_Index(t *testing.T) {
	resource := new(BmfResource)

	box1 := NewBox("box1", 0, 0, 0, resource)
	tb1 := &testBox1{
		Box: box1,
	}

	box2a := NewBox("box2", 0, 0, 0, resource)
	tb2a := &testBox1{
		Box: box2a,
	}

	box2b := NewBox("box2", 0, 0, 0, resource)
	tb2b := &testBox1{
		Box: box2b,
	}

	boxes := Boxes{
		tb1,
		tb2a,
		tb2b,
	}

	index := boxes.Index()

	if len(index) != 2 {
		t.Fatalf("Index not correct size.")
	} else if reflect.DeepEqual(index["box1"], []CommonBox{tb1}) != true {
		t.Fatalf("First indexed entry not correct: %v", index["box1"])
	} else if reflect.DeepEqual(index["box2"], []CommonBox{tb2a, tb2b}) != true {
		t.Fatalf("Second indexed entry not correct: %v", index["box2"])
	}
}
