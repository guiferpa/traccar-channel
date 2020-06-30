package timeutil

import "time"

type tutil struct {
	t time.Time
}

func Now() tutil {
	return tutil{time.Now()}
}
