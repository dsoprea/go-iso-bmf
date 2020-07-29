package bmfcommon

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/dsoprea/go-logging"
	"github.com/dsoprea/go-utility/filesystem"
)

func TestNewBmfResource(t *testing.T) {
	// Construct

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

	// Parse

	sb := rifs.NewSeekableBufferWithBytes(b)
	size := int64(len(b))

	resource := NewBmfResource(sb, size)

	// Validate

	if len(resource.LoadedBoxIndex) != 2 {
		t.Fatalf("Exactly two child boxes weren't found.")
	}
}

func TestBmfResource_Index(t *testing.T) {
	var b []byte

	sb := rifs.NewSeekableBufferWithBytes(b)
	size := int64(len(b))

	resource := NewBmfResource(sb, size)

	if reflect.DeepEqual(resource.fullBoxIndex, resource.Index()) != true {
		t.Fatalf("Index() does not return inner index field.")
	}
}

func TestBmfResource_readBytesAt_Front(t *testing.T) {
	data := []byte{
		0, 0, 0, 0,
		1, 2, 3, 4, 5,
		0, 0, 0, 0, 0,
	}

	sb := rifs.NewSeekableBufferWithBytes(data)

	resource := NewBmfResource(sb, 0)

	recovered, err := resource.readBytesAt(4, 5)
	log.PanicIf(err)

	if bytes.Equal(recovered, data[4:9]) != true {
		t.Fatalf("Read bytes not correct.")
	}
}

func TestBmfResource_readBytesAt_Middle(t *testing.T) {
	data := []byte{
		0, 0, 0, 0,
		1, 2, 3, 4, 5,
		0, 0, 0, 0, 0,
	}

	sb := rifs.NewSeekableBufferWithBytes(data)

	resource := NewBmfResource(sb, 0)

	recovered, err := resource.readBytesAt(5, 5)
	log.PanicIf(err)

	if bytes.Equal(recovered, data[5:10]) != true {
		t.Fatalf("Read bytes not correct.")
	}
}

func TestBmfResource_readBytesAt_MiddleToEnd(t *testing.T) {
	data := []byte{
		0, 0, 0, 0,
		1, 2, 3, 4, 5,
		6, 7, 8, 9, 10,
	}

	sb := rifs.NewSeekableBufferWithBytes(data)

	resource := NewBmfResource(sb, 0)

	recovered, err := resource.readBytesAt(4, 10)
	log.PanicIf(err)

	if bytes.Equal(recovered, data[4:14]) != true {
		t.Fatalf("Read bytes not correct.")
	}
}

func TestBmfResource_readBaseBox_32(t *testing.T) {
	var buffer []byte
	PushBox(&buffer, "abcd", []byte{6, 7, 8, 9})

	sb := rifs.NewSeekableBufferWithBytes(buffer)
	resource := NewBmfResource(sb, int64(len(buffer)))

	box, err := resource.readBaseBox(0)
	log.PanicIf(err)

	if box.Size() != int64(12) {
		t.Fatalf("Size not correct: (%d)", box.Size())
	} else if box.Name() != "abcd" {
		t.Fatalf("Type not correct: [%s]", box.Name())
	}
}

func TestBmfResource_readBaseBox_64(t *testing.T) {
	var buffer []byte
	PushBox(&buffer, "abcd", Data64BitDescribed{6, 7, 8, 9})

	sb := rifs.NewSeekableBufferWithBytes(buffer)
	resource := NewBmfResource(sb, int64(len(buffer)))

	box, err := resource.readBaseBox(0)
	log.PanicIf(err)

	if box.Size() != int64(20) {
		t.Fatalf("Size not correct: (%d)", box.Size())
	} else if box.Name() != "abcd" {
		t.Fatalf("Type not correct: [%s]", box.Name())
	}
}

func TestBmfResource_readBaseBox_Front(t *testing.T) {
	data := []byte{
		0, 0, 0, 12,
		'a', 'b', 'c', 'd',
		6, 7, 8, 9,
	}

	sb := rifs.NewSeekableBufferWithBytes(data)

	resource := NewBmfResource(sb, int64(len(data)))

	box, err := resource.readBaseBox(0)
	log.PanicIf(err)

	if box.Size() != int64(12) {
		t.Fatalf("Size not correct: (%d)", box.Size())
	} else if box.Name() != "abcd" {
		t.Fatalf("Type not correct: [%s]", box.Name())
	}
}

func TestBmfResource_readBaseBox_Middle(t *testing.T) {
	data := []byte{
		0, 0, 0, 12,
		'a', 'b', 'c', 'd',
		6, 7, 8, 9,

		0, 0, 0, 16,
		'e', 'f', 'g', 'h',
		10, 11, 12, 13, 14, 15, 16, 17,
	}

	sb := rifs.NewSeekableBufferWithBytes(data)

	resource := NewBmfResource(sb, int64(len(data)))

	box, err := resource.readBaseBox(12)
	log.PanicIf(err)

	if box.Size() != int64(16) {
		t.Fatalf("Size not correct: (%d)", box.Size())
	} else if box.Name() != "efgh" {
		t.Fatalf("Type not correct: [%s]", box.Name())
	}
}

func TestBmfResource_ReadBaseBox_32(t *testing.T) {
	var buffer []byte
	PushBox(&buffer, "abcd", []byte{6, 7, 8, 9})

	sb := rifs.NewSeekableBufferWithBytes(buffer)
	resource := NewBmfResource(sb, int64(len(buffer)))

	box, err := resource.ReadBaseBox(0)
	log.PanicIf(err)

	if box.Size() != int64(12) {
		t.Fatalf("Size not correct: (%d)", box.Size())
	} else if box.Name() != "abcd" {
		t.Fatalf("Type not correct: [%s]", box.Name())
	}
}

func TestBmfResource_ReadBaseBox_64(t *testing.T) {
	var buffer []byte
	PushBox(&buffer, "abcd", Data64BitDescribed{6, 7, 8, 9})

	sb := rifs.NewSeekableBufferWithBytes(buffer)
	resource := NewBmfResource(sb, int64(len(buffer)))

	box, err := resource.ReadBaseBox(0)
	log.PanicIf(err)

	if box.Size() != int64(20) {
		t.Fatalf("Size not correct: (%d)", box.Size())
	} else if box.Name() != "abcd" {
		t.Fatalf("Type not correct: [%s]", box.Name())
	}
}

func TestBmfResource_ReadBaseBox_Front(t *testing.T) {
	data := []byte{
		0, 0, 0, 12,
		'a', 'b', 'c', 'd',
		6, 7, 8, 9,
	}

	sb := rifs.NewSeekableBufferWithBytes(data)

	resource := NewBmfResource(sb, int64(len(data)))

	box, err := resource.ReadBaseBox(0)
	log.PanicIf(err)

	if box.Size() != int64(12) {
		t.Fatalf("Size not correct: (%d)", box.Size())
	} else if box.Name() != "abcd" {
		t.Fatalf("Type not correct: [%s]", box.Name())
	}
}

func TestBmfResource_ReadBaseBox_Middle(t *testing.T) {
	data := []byte{
		0, 0, 0, 12,
		'a', 'b', 'c', 'd',
		6, 7, 8, 9,

		0, 0, 0, 16,
		'e', 'f', 'g', 'h',
		10, 11, 12, 13, 14, 15, 16, 17,
	}

	sb := rifs.NewSeekableBufferWithBytes(data)

	resource := NewBmfResource(sb, int64(len(data)))

	box, err := resource.ReadBaseBox(12)
	log.PanicIf(err)

	if box.Size() != int64(16) {
		t.Fatalf("Size not correct: (%d)", box.Size())
	} else if box.Name() != "efgh" {
		t.Fatalf("Type not correct: [%s]", box.Name())
	}
}

func TestReadBox(t *testing.T) {
	// Construct

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
	pushUnknownBox(&b, nil)

	// Parse

	sb := rifs.NewSeekableBufferWithBytes(b)
	size := int64(len(b))

	resource := NewBmfResource(sb, size)

	cb1, known1, err := readBox(resource, nil, 0)
	log.PanicIf(err)

	if known1 != true {
		t.Fatalf("First box was not known.")
	} else if cb1.Name() != "tb2 " {
		t.Fatalf("First box name not correct.")
	}

	cb2, known2, err := readBox(resource, nil, 16)
	log.PanicIf(err)

	if known2 != true {
		t.Fatalf("Second box was not known.")
	} else if cb2.Name() != "tb1 " {
		t.Fatalf("Second box name not correct.")
	}

	cb3, known3, err := readBox(resource, nil, 24)
	log.PanicIf(err)

	if known3 != false {
		t.Fatalf("Third box should not have been known.")
	} else if cb3.Name() != "wxyz" {
		t.Fatalf("Third box name not correct.")
	}
}

func TestReadBox_InvalidBoxName(t *testing.T) {
	// Construct

	var b []byte
	PushBox(&b, "\001\002\003\004", nil)

	// Parse

	sb := rifs.NewSeekableBufferWithBytes(b)

	resource := NewBmfResource(sb, 0)

	_, _, err := readBox(resource, nil, 0)
	if err == nil {
		t.Fatalf("Expected error.")
	} else if err.Error() != "box starting at offset (0x0000000000000000) looks like garbage" {
		log.Panic(err)
	}
}

func TestReadBox_WithChildBoxes(t *testing.T) {
	// Construct

	ClearRegistrations()
	defer ClearRegistrations()

	RegisterBoxType(testBox1Factory{})
	RegisterBoxType(testBox2Factory{})
	RegisterBoxType(testBox3Factory{})

	var encodedChildBoxes []byte
	pushTestBox1(&encodedChildBoxes)

	var b []byte
	pushTestBox3(&b, encodedChildBoxes)

	// Parse

	sb := rifs.NewSeekableBufferWithBytes(b)
	size := int64(len(b))

	resource := NewBmfResource(sb, size)

	cb, _, err := readBox(resource, nil, 0)
	log.PanicIf(err)

	if cb.Name() != "tb3 " {
		t.Fatalf("Outer box not correct.")
	}

	tb1 := cb.(*testBox3)

	if len(tb1.LoadedBoxIndex) != 1 {
		t.Fatalf("Expected LBI to have one entry.")
	} else if _, found := tb1.LoadedBoxIndex["tb1 "]; found != true {
		t.Fatalf("Child box not correct.")
	}
}

func TestReadBoxes(t *testing.T) {
	defer func() {
		if errRaw := recover(); errRaw != nil {
			err := errRaw.(error)
			log.PrintError(err)

			t.Fatalf("Test failed.")
		}
	}()

	// Construct

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

	// Parse

	sb := rifs.NewSeekableBufferWithBytes(b)
	size := int64(len(b))

	resource := NewBmfResource(sb, size)

	boxes, err := readBoxes(resource, nil, 0, size)
	log.PanicIf(err)

	// Validate

	if len(boxes) != 2 {
		t.Fatalf("Expected two boxes: (%d)", len(boxes))
	}

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
