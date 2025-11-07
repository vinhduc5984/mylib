package utils

import (
	"time"
)

const YYYYMMDD = "20060102"
const YYYYMMDDHHMM = "200601020304"

// function FormatDateTime return "yyyyMMddHHMM"
func FormatDateTime(millisecond int64) string {
	return formatDateTime(millisecond, YYYYMMDDHHMM)
}

// function FormatDate return "yyyyMMdd"
func FormatDate(millisecond int64) string {
	return formatDateTime(millisecond, YYYYMMDD)
}

// function FormatDateWithLayout
func FormatDateWithLayout(millisecond int64, layout string) string {
	return formatDateTime(millisecond, layout)
}

// function StandardFormatDate return "dd/MM/yyyyy"
func StandardFormatDate(millisecond int64) string {
	return formatDateTime(millisecond, "02/01/2006")
}

func formatDateTime(millisecond int64, layout string) string {
	return time.Unix(0, millisecond*int64(time.Millisecond)).In(time.UTC).Format(layout)
}

func MakeTimeWithTimezone(val time.Time, diffHour float64) time.Time {
	if diffHour == DiffHourNil {
		return val
	}
	return val.Add(time.Hour * time.Duration(diffHour))
}

func MakeNowWithTimezone(diffHour float64) time.Time {
	return MakeTimeWithTimezone(time.Now(), diffHour)
}

func MakeDateTimeWithDiffHour(millisecond int64, diffHour float64) int64 {
	return MakeTimeWithTimezone(time.UnixMilli(millisecond), diffHour).In(time.UTC).UnixMilli()
}
