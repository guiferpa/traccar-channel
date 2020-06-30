package timeutil

import (
	"fmt"
)

func (tu tutil) UTCClockString() (string, string, string) {
	hour, minute, second := tu.t.UTC().Clock()

	nhour := fmt.Sprintf("%v", hour)
	if hour < 10 {
		nhour = fmt.Sprintf("0%v", hour)
	}

	nminute := fmt.Sprintf("%v", minute)
	if minute < 10 {
		nminute = fmt.Sprintf("0%v", minute)
	}

	nsecond := fmt.Sprintf("%v", second)
	if second < 10 {
		nsecond = fmt.Sprintf("0%v", second)
	}

	return nhour, nminute, nsecond
}
