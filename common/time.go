package bmfcommon

import (
	"fmt"
	"time"

	"github.com/dsoprea/go-logging"
)

// Standard32TimeSupport packages some common time-handling fields and
// functionality.
type Standard32TimeSupport struct {
	// creationEpoch is the creation time expressed as an MP4 epoch.
	creationEpoch uint32

	// modificationEpoch is the modification time expressed as an MP4 epoch.
	modificationEpoch uint32

	// scaledDuration is the duration expressed as a number of ticks (scaled
	// per timeScale).
	scaledDuration uint32

	// timeScale is the number of ticks per second.
	timeScale uint32
}

func NewStandard32TimeSupport(creationEpoch, modificationEpoch uint32, scaledDuration, timeScale uint32) Standard32TimeSupport {

	// TODO(dustin): Add test.

	return Standard32TimeSupport{
		creationEpoch:     creationEpoch,
		modificationEpoch: modificationEpoch,
		scaledDuration:    scaledDuration,
		timeScale:         timeScale,
	}
}

// CreationTime returns the creation time.
func (sts Standard32TimeSupport) CreationTime() time.Time {
	t := EpochToTime(sts.creationEpoch)
	return t
}

// HasCreationTime returns true if the creation-time looks present.
func (sts Standard32TimeSupport) HasCreationTime() bool {
	return sts.creationEpoch > 0
}

// ModificationTime returns the modification time.
func (sts Standard32TimeSupport) ModificationTime() time.Time {
	t := EpochToTime(sts.modificationEpoch)
	return t
}

// HasModificationTime returns true if the modification-time looks present.
func (sts Standard32TimeSupport) HasModificationTime() bool {
	return sts.modificationEpoch > 0
}

// timeScale returns the time-scale.
func (sts Standard32TimeSupport) TimeScale() uint32 {
	return sts.timeScale
}

// HasDuration returns true if the duration has a meaningful value.
func (sts Standard32TimeSupport) HasDuration() bool {
	return sts.scaledDuration > 0
}

// scaledDuration returns the duration in timescale units (divide this number by
// the time-scale to get the number of seconds).
func (sts Standard32TimeSupport) ScaledDuration() uint32 {
	if sts.HasDuration() == false {
		log.Panicf("duration not set (scaled-duration)")
	}

	return sts.scaledDuration
}

// Duration returns the duration as a `time.Duration`.
func (sts Standard32TimeSupport) Duration() time.Duration {
	if sts.HasDuration() == false {
		log.Panicf("duration not set (duration)")
	}

	durationSeconds := float64(sts.scaledDuration) / float64(sts.timeScale)

	return time.Duration(durationSeconds * float64(time.Second))
}

// InlineString returns an undecorated string of field names and values.
func (sts Standard32TimeSupport) InlineString() string {

	// TODO(dustin): Add test.

	optional := ""

	if sts.HasCreationTime() == true {
		optional = fmt.Sprintf("%s CTIME=[%s]", optional, sts.CreationTime())
	}

	if sts.HasModificationTime() == true {
		optional = fmt.Sprintf("%s MTIME=[%s]", optional, sts.ModificationTime())
	}

	return fmt.Sprintf(
		"DUR-S=[%.02f]%s",
		float64(sts.Duration())/float64(time.Second), optional)
}