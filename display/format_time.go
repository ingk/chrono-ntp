package display

import (
	"fmt"
	"time"
)

var AllowedDateFormats = [...]string{"YYYY-MM-DD", "DD/MM/YYYY", "MM/DD/YYYY", "DD.MM.YYYY"}
var AllowedTimeFormats = [...]string{"ISO8601", "12h", "12h_AM_PM", ".beat", "septimal", "mars", "lunar", "unix"}

func FormatDate(t time.Time, dateFormat *string) string {
	switch *dateFormat {
	case "YYYY-MM-DD":
		return t.Format("2006-01-02")
	case "DD/MM/YYYY":
		return t.Format("02/01/2006")
	case "MM/DD/YYYY":
		return t.Format("01/02/2006")
	case "DD.MM.YYYY":
		return t.Format("02.01.2006")
	default:
		return t.Format("2006-01-02") // fallback to ISO
	}
}

func FormatTime(t time.Time, timeFormat *string) string {
	switch *timeFormat {
	case ".beat":
		return formatBeatTime(t)
	case "septimal":
		return formatSeptimalTime(t)
	case "mars":
		return formatMarsTime(t)
	case "lunar":
		return formatLunarTime(t)
	case "unix":
		return fmt.Sprintf("%d", t.Unix())
	default:
		timeFormatMap := map[string]string{
			"ISO8601":   "15:04:05",
			"12h":       "03:04:05",
			"12h_AM_PM": "03:04:05 PM",
		}
		return t.Format(timeFormatMap[*timeFormat])
	}
}

// formatMarsTime returns Coordinated Mars Time (MTC)
// See: https://en.wikipedia.org/wiki/Timekeeping_on_Mars
func formatMarsTime(t time.Time) string {
	// Julian Date (UTC)
	y, m, d := t.UTC().Date()
	h, min, s := t.UTC().Clock()
	ms := t.UTC().Nanosecond() / 1e6

	if m <= 2 {
		y -= 1
		m += 12
	}
	A := y / 100
	B := 2 - A + A/4
	jd := 365.25*float64(y+4716) + 30.6001*float64(m+1) + float64(d) + float64(B) - 1524.5
	fracDay := (float64(h) + float64(min)/60 + float64(s)/3600 + float64(ms)/3600000) / 24.0
	JD := jd + fracDay

	// Mars Sol Date (MSD)
	MSD := (JD - 2405522.0028779) / 1.0274912517
	mtc := 24.0 * (MSD - float64(int(MSD)))
	hh := int(mtc)
	mm := int((mtc - float64(hh)) * 60)
	ss := int((((mtc - float64(hh)) * 60) - float64(mm)) * 60)
	return fmt.Sprintf("%02d:%02d:%02d", hh, mm, ss)
}

// formatLunarTime returns Coordinated Lunar Time (LTC)
// See: https://en.wikipedia.org/wiki/Timekeeping_on_the_Moon
func formatLunarTime(t time.Time) string {
	// Reference epoch: 2000-01-06 18:14 UTC (known new moon, J2000)
	reference := time.Date(2000, 1, 6, 18, 14, 0, 0, time.UTC)
	lunarDay := 29.530589 * 86400 // seconds in a lunar day
	delta := t.UTC().Sub(reference).Seconds()
	lunarDays := delta / lunarDay
	fraction := lunarDays - float64(int(lunarDays))
	if fraction < 0 {
		fraction += 1
	}
	totalSeconds := fraction * lunarDay
	hh := int(totalSeconds / 3600)
	mm := int((totalSeconds - float64(hh*3600)) / 60)
	ss := int(totalSeconds) % 60
	return fmt.Sprintf("%02d:%02d:%02d", hh, mm, ss)
}

// formatSeptimalTime returns the time in septimal format (base-7 pairs)
// See: http://the-light.com/cal/veseptimal.html
func formatSeptimalTime(t time.Time) string {
	// Get time since midnight in local time
	year, month, day := t.Date()
	midnight := time.Date(year, month, day, 0, 0, 0, 0, t.Location())
	msSinceMidnight := float64(t.Sub(midnight).Milliseconds())

	// Use the original JS divisors for accuracy
	sep1 := int(msSinceMidnight / 12342857.14)
	sep2 := int(msSinceMidnight/1763265.306) % 7
	sep3 := int(msSinceMidnight/251895.0437) % 7
	sep4 := int(msSinceMidnight/35985.00625) % 7
	sep5 := int(msSinceMidnight/5140.715178) % 7
	sep6 := int((msSinceMidnight / 734.3878826)) % 7

	return fmt.Sprintf("%d%d %d%d %d%d", sep1, sep2, sep3, sep4, sep5, sep6)
}

// formatBeatTime returns Swatch Internet Time (.beat)
// @see https://en.wikipedia.org/wiki/Swatch_Internet_Time
func formatBeatTime(t time.Time) string {
	// Convert time to UTC+1 (Biel Mean Time)
	bmt := t.UTC().Add(1 * time.Hour)
	seconds := bmt.Hour()*3600 + bmt.Minute()*60 + bmt.Second()
	beat := float64(seconds) / 86.4
	return fmt.Sprintf("@%s", formatBeat(beat))
}

func formatBeat(beat float64) string {
	// .beat is three digits; we also include centibeats (two digits)
	// Round down both parts (floor behavior)
	intPart := int(beat) % 1000
	frac := beat - float64(intPart)
	centi := int(frac * 100)
	// Guard against floating-point rounding producing 100
	if centi >= 100 {
		centi = 0
		intPart = (intPart + 1) % 1000
	}
	return fmt.Sprintf("%s.%s", leftPadInt(intPart, 3), leftPadInt(centi, 2))
}

func leftPadInt(n int, width int) string {
	s := fmt.Sprintf("%d", n)
	for len(s) < width {
		s = "0" + s
	}
	return s
}
