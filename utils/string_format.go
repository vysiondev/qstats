package utils

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"
)

func GetKeymodeString(is7K bool) string {
	if is7K {
		return "7K"
	} else {
		return "4K"
	}
}

func GetKeymodeIntAsStr(is7K bool) string {
	if is7K {
		return "2"
	} else {
		return "1"
	}
}

func GetCountryStr(proposedCountry string) string {
	newC := strings.ToLower(proposedCountry)
	if newC == "xx" {
		newC = "aq"
	}
	return newC
}

// https://github.com/icza/gox
func AddCommas(n int64) string {
	in := strconv.FormatInt(n, 10)
	numOfDigits := len(in)
	if n < 0 {
		numOfDigits-- // First character is the - sign (not a digit)
	}
	numOfCommas := (numOfDigits - 1) / 3

	out := make([]byte, len(in)+numOfCommas)
	if n < 0 {
		in, out[0] = in[1:], '-'
	}

	for i, j, k := len(in)-1, len(out)-1, 0; ; i, j = i-1, j-1 {
		out[j] = in[i]
		if i == 0 {
			return string(out)
		}
		if k++; k == 3 {
			j, k = j-1, 0
			out[j] = ','
		}
	}
}

func s(x float64) string {
	if int(x) == 1 {
		return ""
	}
	return "s"
}

// https://www.socketloop.com/tutorials/golang-human-readable-time-elapsed-format-such-as-5-days-ago
func TimeElapsed(now time.Time, then time.Time, full bool) string {
	var parts []string
	var text string

	year2, month2, day2 := now.Date()
	hour2, minute2, second2 := now.Clock()

	year1, month1, day1 := then.Date()
	hour1, minute1, second1 := then.Clock()

	year := math.Abs(float64(year2 - year1))
	month := math.Abs(float64(month2 - month1))
	day := math.Abs(float64(day2 - day1))
	hour := math.Abs(float64(hour2 - hour1))
	minute := math.Abs(float64(minute2 - minute1))
	second := math.Abs(float64(second2 - second1))

	week := math.Floor(day / 7)

	if year > 0 {
		parts = append(parts, strconv.Itoa(int(year))+" year"+s(year))
	}

	if month > 0 {
		parts = append(parts, strconv.Itoa(int(month))+" month"+s(month))
	}

	if week > 0 {
		parts = append(parts, strconv.Itoa(int(week))+" week"+s(week))
	}

	if day > 0 {
		parts = append(parts, strconv.Itoa(int(day))+" day"+s(day))
	}

	if hour > 0 {
		parts = append(parts, strconv.Itoa(int(hour))+" hour"+s(hour))
	}

	if minute > 0 {
		parts = append(parts, strconv.Itoa(int(minute))+" minute"+s(minute))
	}

	if second > 0 {
		parts = append(parts, strconv.Itoa(int(second))+" second"+s(second))
	}

	if now.After(then) {
		text = " ago"
	} else {
		text = " after"
	}

	if len(parts) == 0 {
		return "just now"
	}

	if full {
		return strings.Join(parts, ", ") + text
	}
	return parts[0] + text
}

func FmtDuration(d time.Duration) string {
	d = d.Round(time.Second)
	m := d / time.Minute
	d -= m * time.Minute
	s := d / time.Second
	return fmt.Sprintf("%02d:%02d", m, s)
}
