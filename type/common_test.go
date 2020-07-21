package bmftype

import (
	"time"

	"github.com/dsoprea/go-iso-bmf/common"
)

func getTestStandard32Time() (now time.Time, sts bmfcommon.Standard32TimeSupport) {
	now = bmfcommon.NowTime()
	epoch := bmfcommon.TimeToEpoch(now)

	sts = bmfcommon.NewStandard32TimeSupport(
		epoch,
		epoch+1,
		10,
		55)

	return now, sts
}
