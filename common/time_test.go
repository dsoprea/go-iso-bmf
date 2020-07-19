package bmfcommon

import (
	"testing"
	"time"
)

func TestStandard32TimeSupport_HasCreationTime_False(t *testing.T) {
	sts := Standard32TimeSupport{}

	if sts.HasCreationTime() != false {
		t.Fatalf("HasCreationTime() should be false.")
	}
}

func TestStandard32TimeSupport_HasCreationTime_True(t *testing.T) {
	now := NowTime()
	creationEpoch := EpochDelta(now)

	sts := Standard32TimeSupport{
		creationEpoch: creationEpoch,
	}

	if sts.HasCreationTime() != true {
		t.Fatalf("HasCreationTime() should be true.")
	}
}

func TestStandard32TimeSupport_CreationTime(t *testing.T) {
	now := NowTime()
	creationEpoch := EpochDelta(now)

	sts := Standard32TimeSupport{
		creationEpoch: creationEpoch,
	}

	if sts.CreationTime() != now {
		t.Fatalf("CreationTime() not correct: [%s] != [%s]", sts.CreationTime(), now)
	}
}

func TestStandard32TimeSupport_HasModificationTime_False(t *testing.T) {
	sts := Standard32TimeSupport{}

	if sts.HasModificationTime() != false {
		t.Fatalf("HasModificationTime() should be false.")
	}
}

func TestStandard32TimeSupport_HasModificationTime_True(t *testing.T) {
	now := NowTime()
	modificationEpoch := EpochDelta(now)

	sts := Standard32TimeSupport{
		modificationEpoch: modificationEpoch,
	}

	if sts.HasModificationTime() != true {
		t.Fatalf("HasModificationTime() should be true.")
	}
}

func TestStandard32TimeSupport_ModificationTime(t *testing.T) {
	now := NowTime()
	modificationEpoch := EpochDelta(now)

	sts := Standard32TimeSupport{
		modificationEpoch: modificationEpoch,
	}

	if sts.ModificationTime() != now {
		t.Fatalf("ModificationTime() not correct.")
	}
}

func TestStandard32TimeSupport_TimeScale(t *testing.T) {
	sts := Standard32TimeSupport{
		timeScale: 55,
	}

	if sts.TimeScale() != 55 {
		t.Fatalf("TimeScale() not correct.")
	}
}

func TestStandard32TimeSupport_ScaledDuration(t *testing.T) {
	sts := Standard32TimeSupport{
		scaledDuration: 10,
	}

	if sts.ScaledDuration() != 10 {
		t.Fatalf("ScaledDuration() not correct.")
	}
}

func TestStandard32TimeSupport_HasDuration_False(t *testing.T) {
	sts := Standard32TimeSupport{}

	if sts.HasDuration() != false {
		t.Fatalf("HasDuration() not correct.")
	}
}

func TestStandard32TimeSupport_HasDuration_True(t *testing.T) {
	timeScale := uint32(60)

	sts := Standard32TimeSupport{
		timeScale:      timeScale,
		scaledDuration: timeScale * 10,
	}

	if sts.HasDuration() != true {
		t.Fatalf("HasDuration() not correct.")
	}
}

func TestStandard32TimeSupport_Duration(t *testing.T) {
	timeScale := uint32(60)

	sts := Standard32TimeSupport{
		timeScale:      timeScale,
		scaledDuration: timeScale * 10,
	}

	d := time.Second * 10
	if sts.Duration() != d {
		t.Fatalf("Duration() not correct: [%s] != [%s]", sts.Duration(), d)
	}
}
