package bmfcommon

import (
	"testing"
	"time"
)

func TestNewStandard32TimeSupport(t *testing.T) {
	creationEpoch := uint64(1)
	modificationEpoch := uint64(2)
	scaledDuration := uint64(3)
	timeScale := uint64(4)

	sts := NewStandard32TimeSupport(creationEpoch, modificationEpoch, scaledDuration, timeScale)

	if sts.creationEpoch != creationEpoch {
		t.Fatalf("createdEpoch is not correct.")
	} else if sts.modificationEpoch != modificationEpoch {
		t.Fatalf("modificationEpoch is not correct.")
	} else if sts.scaledDuration != scaledDuration {
		t.Fatalf("scaledDuration is not correct.")
	} else if sts.timeScale != timeScale {
		t.Fatalf("timeScale is not correct.")
	}
}

func TestStandard32TimeSupport_HasCreationTime_False(t *testing.T) {
	sts := Standard32TimeSupport{}

	if sts.HasCreationTime() != false {
		t.Fatalf("HasCreationTime() should be false.")
	}
}

func TestStandard32TimeSupport_HasCreationTime_True(t *testing.T) {
	now := NowTime()
	creationEpoch := TimeToEpoch(now)

	sts := Standard32TimeSupport{
		creationEpoch: creationEpoch,
	}

	if sts.HasCreationTime() != true {
		t.Fatalf("HasCreationTime() should be true.")
	}
}

func TestStandard32TimeSupport_CreationTime(t *testing.T) {
	now := NowTime()
	creationEpoch := TimeToEpoch(now)

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
	modificationEpoch := TimeToEpoch(now)

	sts := Standard32TimeSupport{
		modificationEpoch: modificationEpoch,
	}

	if sts.HasModificationTime() != true {
		t.Fatalf("HasModificationTime() should be true.")
	}
}

func TestStandard32TimeSupport_ModificationTime(t *testing.T) {
	now := NowTime()
	modificationEpoch := TimeToEpoch(now)

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
	timeScale := uint64(60)

	sts := Standard32TimeSupport{
		timeScale:      timeScale,
		scaledDuration: timeScale * 10,
	}

	if sts.HasDuration() != true {
		t.Fatalf("HasDuration() not correct.")
	}
}

func TestStandard32TimeSupport_Duration(t *testing.T) {
	timeScale := uint64(60)

	sts := Standard32TimeSupport{
		timeScale:      timeScale,
		scaledDuration: timeScale * 10,
	}

	d := time.Second * 10
	if sts.Duration() != d {
		t.Fatalf("Duration() not correct: [%s] != [%s]", sts.Duration(), d)
	}
}

func TestStandard32TimeSupport_InlineString(t *testing.T) {
	creationEpoch := uint64(1)
	modificationEpoch := creationEpoch + 1

	timeScale := uint64(60)

	sts := Standard32TimeSupport{
		creationEpoch:     creationEpoch,
		modificationEpoch: modificationEpoch,
		timeScale:         timeScale,
		scaledDuration:    timeScale * 10,
	}

	if sts.InlineString() != "DUR-S=[10.00] CTIME=[1904-01-01 00:00:01 +0000 UTC] MTIME=[1904-01-01 00:00:02 +0000 UTC]" {
		t.Fatalf("InlineString() not correct: [%s]", sts.InlineString())
	}
}

func TestTimeToEpoch(t *testing.T) {
	originalTime := time.Unix(1, 0).UTC()
	epoch := TimeToEpoch(originalTime)
	if epoch != 2082844801 {
		t.Fatalf("Epoch not correct.")
	}

	recovered := EpochToTime(epoch)

	if recovered != originalTime {
		t.Fatalf("Recovered time not correct.")
	}
}

func TestEpochToTime(t *testing.T) {
	originalTime := time.Unix(1, 0).UTC()

	epoch := uint64(2082844801)

	recovered := EpochToTime(epoch)
	if recovered != originalTime {
		t.Fatalf("Time not correct: %s != %s", recovered, originalTime)
	}
}

func TestNowTime(t *testing.T) {
	actual := NowTime()
	expected := time.Now().UTC().Round(time.Second)

	if actual != expected {
		t.Fatalf("Time not expected: %s != %s", actual, expected)
	}
}

func TestGetDurationString(t *testing.T) {
	s := GetDurationString(123456789, 12345)

	if s != "02:46:40:549" {
		t.Fatalf("Duration string not correct.")
	}
}
