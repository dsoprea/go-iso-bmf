package bmfcommon

import (
	"testing"
)

func TestDump_Unexported_CommonBox(t *testing.T) {
	tb1 := new(testBox1)

	dump(tb1, 0)
}

func TestDump_Unexported_BoxChildIndexer(t *testing.T) {
	lbi := LoadedBoxIndex{
		"abc": []CommonBox{&testBox1{}},
	}

	tb4 := &testBox4{
		LoadedBoxIndex: lbi,
	}

	dump(tb4, 0)
}

func TestDump_Unexported_BoxChildIndexer_Root(t *testing.T) {
	lbi := LoadedBoxIndex{
		"abc": []CommonBox{&testBox1{}},
	}

	br := &BmfResource{
		LoadedBoxIndex: lbi,
	}

	dump(br, 0)
}

func TestDump_Exported(t *testing.T) {
	tb1 := new(testBox1)

	Dump(tb1)
}
