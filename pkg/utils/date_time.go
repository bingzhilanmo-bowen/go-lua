package utils

import "time"

const (
	BaseDateTime = "2006-01-02 15:04:05"
	BaseDate        = "2006-01-02"
	ShortDate   = "06-01-02"
	BaseTimes       = "15:04:05"
	ShortTime   = "1504"
)

func String2Date(d string) *time.Time {
	date, _ := time.Parse(BaseDateTime, d)
    return &date
}

func Unix2ShortTime(ux int64)string {
	ti := time.Unix(ux, 0)
	st := ti.Format(ShortTime)
	return st
}

func IsUnixTime(ux int64) bool {
	return ux >= 1136185445
}
