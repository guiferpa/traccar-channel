package timeutil

import (
	"fmt"
)

func (tu tutil) DateString() (string, string, string) {
	year, month, day := tu.t.Date()

	nyear := fmt.Sprintf("%v", year)

	nmonth := fmt.Sprintf("%v", int(month))
	if int(month) < 10 {
		nmonth = fmt.Sprintf("0%v", int(month))
	}

	nday := fmt.Sprintf("%v", day)
	if day < 10 {
		nday = fmt.Sprintf("0%v", day)
	}

	return nyear, nmonth, nday
}
