package formatters

import "time"

func FormatDate(t time.Time) string {
	return t.Format("January 2, 2006 at 3:04 PM")
}

func BitwiseAnd(a int64, b int64) int32 {
	return int32(a & b)
}
