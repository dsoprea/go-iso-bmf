package bmftype

import (
	"testing"
	"time"

	"github.com/dsoprea/go-logging"
	"github.com/dsoprea/go-utility/filesystem"

	"github.com/dsoprea/go-iso-bmf/common"
)

func TestMdhdBox_Version(t *testing.T) {
	mb := MdhdBox{
		version: 11,
	}

	if mb.Version() != 11 {
		t.Fatalf("Version() not correct.")
	}
}

func TestMdhdBox_Flags(t *testing.T) {
	mb := MdhdBox{
		flags: 22,
	}

	if mb.Flags() != 22 {
		t.Fatalf("Flags() not correct.")
	}
}

func TestMdhdBox_HasCreationTime_False(t *testing.T) {
	mb := MdhdBox{}

	if mb.HasCreationTime() != false {
		t.Fatalf("HasCreationTime() should be false.")
	}
}

func TestMdhdBox_HasCreationTime_True(t *testing.T) {
	now := bmfcommon.NowTime()
	creationEpoch := bmfcommon.EpochDelta(now)

	mb := MdhdBox{
		Standard32TimeSupport: bmfcommon.NewStandard32TimeSupport(creationEpoch, 0, 0, 0),
	}

	if mb.HasCreationTime() != true {
		t.Fatalf("HasCreationTime() should be true.")
	}
}

func TestMdhdBox_CreationTime(t *testing.T) {
	_, sts := getTestStandard32Time()

	mb := MdhdBox{
		Standard32TimeSupport: sts,
	}

	if mb.CreationTime() != sts.CreationTime() {
		t.Fatalf("CreationTime() not correct: [%s] != [%s]", mb.CreationTime(), sts.CreationTime())
	}
}

func TestMdhdBox_HasModificationTime_False(t *testing.T) {
	mb := MdhdBox{}

	if mb.HasModificationTime() != false {
		t.Fatalf("HasModificationTime() should be false.")
	}
}

func TestMdhdBox_HasModificationTime_True(t *testing.T) {
	_, sts := getTestStandard32Time()

	mb := MdhdBox{
		Standard32TimeSupport: sts,
	}

	if mb.HasModificationTime() != true {
		t.Fatalf("HasModificationTime() should be true.")
	}
}

func TestMdhdBox_ModificationTime(t *testing.T) {
	_, sts := getTestStandard32Time()

	mb := MdhdBox{
		Standard32TimeSupport: sts,
	}

	if mb.ModificationTime() != sts.ModificationTime() {
		t.Fatalf("ModificationTime() not correct.")
	}
}

func TestMdhdBox_TimeScale(t *testing.T) {
	_, sts := getTestStandard32Time()

	mb := MdhdBox{
		Standard32TimeSupport: sts,
	}

	if mb.TimeScale() != sts.TimeScale() {
		t.Fatalf("TimeScale() not correct.")
	}
}

func TestMdhdBox_ScaledDuration(t *testing.T) {
	mb := MdhdBox{
		Standard32TimeSupport: bmfcommon.NewStandard32TimeSupport(0, 0, 10, 0),
	}

	if mb.ScaledDuration() != 10 {
		t.Fatalf("ScaledDuration() not correct.")
	}
}

func TestMdhdBox_HasDuration_False(t *testing.T) {
	mb := MdhdBox{}

	if mb.HasDuration() != false {
		t.Fatalf("HasDuration() not correct.")
	}
}

func TestMdhdBox_HasDuration_True(t *testing.T) {
	timeScale := uint32(60)

	mb := MdhdBox{
		Standard32TimeSupport: bmfcommon.NewStandard32TimeSupport(0, 0, timeScale*10, timeScale),
	}

	if mb.HasDuration() != true {
		t.Fatalf("HasDuration() not correct.")
	}
}

func TestMdhdBox_Duration(t *testing.T) {
	timeScale := uint32(60)

	mb := MdhdBox{
		Standard32TimeSupport: bmfcommon.NewStandard32TimeSupport(0, 0, timeScale*10, timeScale),
	}

	d := time.Second * 10
	if mb.Duration() != d {
		t.Fatalf("Duration() not correct: [%s] != [%s]", mb.Duration(), d)
	}
}

func TestMdhdBox_LanguageString(t *testing.T) {

	mb := MdhdBox{
		// 00100 00101 00110
		language: 0b001000010100110,
	}

	l := mb.Language()

	if l != "def" {
		t.Fatalf("Language() not correct: (%d) [%s]", len(l), l)
	}
}

func TestMdhdBox_getLanguageString(t *testing.T) {
	mb := MdhdBox{}
	l := mb.getLanguageString(0b001000010100110)

	if l != "def" {
		t.Fatalf("Language() not correct: (%d) [%s]", len(l), l)
	}
}

func TestMdhdBox_String(t *testing.T) {
	box := bmfcommon.NewBox("mdhd", 1234, 5678, nil)

	epoch := uint32(3677725917)

	timeScale := uint32(60)

	sts := bmfcommon.NewStandard32TimeSupport(
		epoch,
		epoch+1,
		timeScale*10,
		timeScale)

	mb := MdhdBox{
		Box:                   box,
		version:               11,
		flags:                 22,
		Standard32TimeSupport: sts,

		// 00100 00101 00110
		language: 0b001000010100110,
	}

	if mb.String() != "mdhd<NAME=[mdhd] START=(1234) SIZE=(5678) VER=(0x0b) FLAGS=(0x00000016) LANG=[def] DUR-S=[10.00] CTIME=[2020-07-16 06:31:57 +0000 UTC] MTIME=[2020-07-16 06:31:58 +0000 UTC]>" {
		t.Fatalf("String() not correct: [%s]", mb.String())
	}
}

func TestMdhdBox_InlineString(t *testing.T) {
	box := bmfcommon.NewBox("mdhd", 1234, 5678, nil)

	epoch := uint32(3677725917)

	timeScale := uint32(60)

	mb := MdhdBox{
		Box:                   box,
		version:               11,
		flags:                 22,
		Standard32TimeSupport: bmfcommon.NewStandard32TimeSupport(epoch, epoch+1, timeScale*10, timeScale),

		// 00100 00101 00110
		language: 0b001000010100110,
	}

	if mb.InlineString() != "NAME=[mdhd] START=(1234) SIZE=(5678) VER=(0x0b) FLAGS=(0x00000016) LANG=[def] DUR-S=[10.00] CTIME=[2020-07-16 06:31:57 +0000 UTC] MTIME=[2020-07-16 06:31:58 +0000 UTC]" {
		t.Fatalf("InlineString() not correct: [%s]", mb.InlineString())
	}
}

func TestMdhdBoxFactory_Name(t *testing.T) {
	name := mdhdBoxFactory{}.Name()
	if name != "mdhd" {
		t.Fatalf("Name() not correct.")
	}
}

func TestMdhdBoxFactory_New(t *testing.T) {
	data := []byte{}

	// flags
	bmfcommon.PushBytes(&data, uint32(0x11223344))

	// creation and modified epochs

	epoch := uint32(3677725917)
	baseTime := bmfcommon.EpochToTime(epoch)

	// creation epoch
	bmfcommon.PushBytes(&data, epoch)

	// modification epoch
	bmfcommon.PushBytes(&data, epoch+1)

	// TimeScale()
	bmfcommon.PushBytes(&data, uint32(30))

	// ScaledDuration()
	bmfcommon.PushBytes(&data, uint32(300))

	// language

	// 00100 00101 00110
	bmfcommon.PushBytes(&data, uint16(0b001000010100110))

	b := []byte{}
	bmfcommon.PushBox(&b, "mdhd", data)

	// Parse.

	sb := rifs.NewSeekableBufferWithBytes(b)

	file := bmfcommon.NewFile(sb, int64(len(b)))

	box, err := file.ReadBaseBox(0)
	log.PanicIf(err)

	cb, err := mdhdBoxFactory{}.New(box)
	log.PanicIf(err)

	mb := cb.(*MdhdBox)

	if mb.Version() != 0x11 {
		t.Fatalf("Version() not correct: (0x%02x)", mb.Version())
	}

	if mb.Flags() != 0x11223344 {
		t.Fatalf("Flags() not correct: (0x%08x)", mb.Flags())
	}

	if mb.CreationTime() != baseTime {
		t.Fatalf("CreationTime() not correct: [%s] != [%s]", mb.CreationTime(), baseTime)
	}

	if mb.ModificationTime() != baseTime.Add(1*time.Second) {
		t.Fatalf("ModificationTime() not correct: %s", mb.ModificationTime())
	}

	if mb.TimeScale() != 30 {
		t.Fatalf("TimeScale() not correct.")
	}

	if mb.ScaledDuration() != 300 {
		t.Fatalf("ScaledDuration() not correct.")
	}

	d := 10 * time.Second

	if mb.Duration() != d {
		t.Fatalf("Duration() not correct: [%s] != [%s]", mb.Duration(), d)
	}

	if mb.Language() != "def" {
		t.Fatalf("Language() not correct.")
	}
}
