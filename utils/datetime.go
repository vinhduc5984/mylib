package utils

import (
	"time"
)

const (
	defaultTimezone = "Asia/Ho_Chi_Minh"
)

func GetCurrentDate() int64 {
	currentTime := time.Now()
	return currentTime.UnixNano() / int64(time.Millisecond)
}

func GetMinDate() int64 {
	minDate := time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC)
	return minDate.UnixNano() / int64(time.Millisecond)
}

func GetMaxDate() int64 {
	maxDate := time.Date(9999, time.December, 31, 23, 59, 59, 999999999, time.UTC)
	return maxDate.UnixNano() / int64(time.Millisecond)
}

func AddDate(milliseconds int64, days int32) int64 {
	t := time.Unix(milliseconds/1000, (milliseconds%1000)*int64(time.Millisecond))
	newTime := t.AddDate(0, 0, int(days))
	return newTime.UnixNano() / int64(time.Millisecond)
}

func AddDateMillisecond(milliseconds int64, interval int64) int64 {
	t := time.Unix(milliseconds/1000, (milliseconds%1000)*int64(time.Millisecond))
	newTime := t.Add(time.Duration(interval) * time.Millisecond)
	return newTime.UnixNano() / int64(time.Millisecond)
}

func GetStartDate(milliseconds int64) int64 {
	t := time.Unix(milliseconds/1000, (milliseconds%1000)*int64(time.Millisecond))
	startTime := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
	return startTime.UnixNano() / int64(time.Millisecond)
}

func GetEndDate(milliseconds int64) int64 {
	t := time.Unix(milliseconds/1000, (milliseconds%1000)*int64(time.Millisecond))
	endTime := time.Date(t.Year(), t.Month(), t.Day(), 23, 59, 59, 999999999, t.Location())
	return endTime.UnixNano() / int64(time.Millisecond)
}

func GetStartWeek(milliseconds int64) int64 {
	t := time.Unix(0, milliseconds*int64(time.Millisecond))
	year, week := t.ISOWeek()
	startOfWeek := time.Date(year, 0, 0, 0, 0, 0, 0, t.Location()).AddDate(0, 0, (week-1)*7)
	return startOfWeek.UnixNano() / int64(time.Millisecond)
}

func GetStartMonth(milliseconds int64) int64 {
	t := time.Unix(0, milliseconds*int64(time.Millisecond))
	startOfMonth := time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, t.Location())
	return startOfMonth.UnixNano() / int64(time.Millisecond)
}

func GetStartQuarter(milliseconds int64) int64 {
	t := time.Unix(0, milliseconds*int64(time.Millisecond))
	startOfQuarter := time.Date(t.Year(), ((t.Month()-1)/3)*3+1, 1, 0, 0, 0, 0, t.Location())
	return startOfQuarter.UnixNano() / int64(time.Millisecond)
}

func GetStartYear(milliseconds int64) int64 {
	t := time.Unix(0, milliseconds*int64(time.Millisecond))
	startOfYear := time.Date(t.Year(), 1, 1, 0, 0, 0, 0, t.Location())
	return startOfYear.UnixNano() / int64(time.Millisecond)
}

func GetEndWeek(milliseconds int64) int64 {
	t := time.Unix(0, milliseconds*int64(time.Millisecond))
	year, week := t.ISOWeek()
	startOfWeek := time.Date(year, 0, 0, 0, 0, 0, 0, t.Location()).AddDate(0, 0, week*7-1)
	endOfWeek := startOfWeek.AddDate(0, 0, 1)
	return endOfWeek.UnixNano() / int64(time.Millisecond)
}

func GetEndMonth(milliseconds int64) int64 {
	t := time.Unix(0, milliseconds*int64(time.Millisecond))
	nextMonth := t.AddDate(0, 1, 0)
	startOfNextMonth := time.Date(nextMonth.Year(), nextMonth.Month(), 1, 0, 0, 0, 0, nextMonth.Location())
	endOfMonth := startOfNextMonth.AddDate(0, 0, -1)
	return endOfMonth.UnixNano() / int64(time.Millisecond)
}

func GetEndQuarter(milliseconds int64) int64 {
	t := time.Unix(0, milliseconds*int64(time.Millisecond))
	nextQuarter := t.AddDate(0, 3, 0)
	startOfNextQuarter := time.Date(nextQuarter.Year(), ((nextQuarter.Month()-1)/3)*3+1, 1, 0, 0, 0, 0, nextQuarter.Location())
	endOfQuarter := startOfNextQuarter.AddDate(0, 0, -1)
	return endOfQuarter.UnixNano() / int64(time.Millisecond)
}

func GetEndYear(milliseconds int64) int64 {
	t := time.Unix(0, milliseconds*int64(time.Millisecond))
	nextYear := t.AddDate(1, 0, 0)
	startOfNextYear := time.Date(nextYear.Year(), 1, 1, 0, 0, 0, 0, nextYear.Location())
	endOfYear := startOfNextYear.AddDate(0, 0, -1)
	return endOfYear.UnixNano() / int64(time.Millisecond)
}

func GetStartDateWithTimeZone(milliseconds int64, timezone string) int64 {
	if timezone == "" {
		timezone = defaultTimezone
	}
	loc, _ := time.LoadLocation(timezone)

	t := time.Unix(0, milliseconds*int64(time.Millisecond)).In(loc)
	formattedDate := t.Format("2006-01-02")

	startDate, _ := time.ParseInLocation("2006-01-02", formattedDate, loc)
	return startDate.UnixNano() / int64(time.Millisecond)
}

func GetDateWithLayoutAndTimeZone(milliseconds int64, layout, timezone string) string {
	t := time.Unix(0, milliseconds*int64(time.Millisecond))
	loc, err := time.LoadLocation(timezone)
	if err != nil {
		return ""
	}
	t = t.In(loc)
	return t.Format(layout)
}

func SetDate(year, month, day int32) int64 {
	t := time.Date(int(year), time.Month(month), int(day), 0, 0, 0, 0, time.UTC)
	return t.UnixNano() / int64(time.Millisecond)
}

func SetDateTime(year, month, day, hour, minute, second int32) int64 {
	t := time.Date(int(year), time.Month(month), int(day), int(hour), int(minute), int(second), 0, time.UTC)
	return t.UnixNano() / int64(time.Millisecond)
}

func DateToMilliseconds(dateString string, diffHours float64) (int64, error) {
	// Define the layout for parsing the date string
	layout := "02/01/2006"
	// Parse the date string into a time.Time object
	parsedTime, err := time.Parse(layout, dateString)
	if err != nil {
		return 0, err
	}
	// Convert the parsed time to milliseconds since Unix epoch
	return parsedTime.UTC().UnixMilli() - int64(diffHours*60*60*1000), nil
}
