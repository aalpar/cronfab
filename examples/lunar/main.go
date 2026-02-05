// Command lunar demonstrates building a custom calendar with cronfab.
//
// It defines a MoonPhaseUnit and constructs a CrontabConfig whose fields
// are hour, moon phase, and month. This is possible because FieldConfig
// fields are exported, so any package can build its own calendar without
// modifying cronfab.
//
// Usage:
//
//	go run ./examples/lunar
package main

import (
	"fmt"
	"math"
	"time"

	"github.com/aalpar/cronfab"
)

// lunarEpoch is a known new moon (January 6, 2000 at 18:14 UTC).
var lunarEpoch = time.Date(2000, 1, 6, 18, 14, 0, 0, time.UTC)

// synodicPeriod is the average length of a synodic month in days.
const synodicPeriod = 29.53059

// phasePeriod is the average length of a single moon phase in days (1/8 of a synodic month).
const phasePeriod = synodicPeriod / 8

// MoonPhaseUnit represents the phase of the moon within a synodic month.
// The synodic month is divided into 8 phases (0-7).
type MoonPhaseUnit struct{}

func (MoonPhaseUnit) String() string { return "moon phase" }

func (MoonPhaseUnit) Add(t time.Time, n int) time.Time {
	hours := phasePeriod * float64(n) * 24
	return t.Add(time.Duration(hours * float64(time.Hour)))
}

func (MoonPhaseUnit) Less(u cronfab.Unit) bool {
	switch u.(type) {
	case cronfab.SecondUnit, cronfab.MinuteUnit, cronfab.HourUnit,
		cronfab.DayUnit, cronfab.WeekOfMonth, MoonPhaseUnit:
		return false
	}
	return true
}

func (MoonPhaseUnit) Trunc(t time.Time) time.Time {
	diff := t.Sub(lunarEpoch).Hours() / 24
	phaseDays := math.Mod(diff, phasePeriod)
	if phaseDays < 0 {
		phaseDays += phasePeriod
	}
	return t.Add(-time.Duration(phaseDays * 24 * float64(time.Hour)))
}

// moonPhaseIndex returns the current moon phase (0-7) for a given time.
func moonPhaseIndex(t time.Time) int {
	diff := t.Sub(lunarEpoch).Hours() / 24
	cycleDays := math.Mod(diff, synodicPeriod)
	if cycleDays < 0 {
		cycleDays += synodicPeriod
	}
	return int(cycleDays / phasePeriod)
}

var lunarConfig = cronfab.MustCrontabConfig([]cronfab.FieldConfig{
	{
		Unit: cronfab.HourUnit{},
		Name: "hour",
		Min:  0,
		Max:  23,
		GetIndex: func(t time.Time) int {
			return t.Hour()
		},
	},
	{
		Unit: MoonPhaseUnit{},
		Name: "moon phase",
		RangeNames: []string{
			"new", "waxingcrescent", "firstquarter", "waxinggibbous",
			"full", "waninggibbous", "thirdquarter", "waningcrescent",
		},
		Min: 0,
		Max: 7,
		GetIndex: func(t time.Time) int {
			return moonPhaseIndex(t)
		},
	},
	{
		Unit:       cronfab.MonthUnit{},
		Name:       "month",
		RangeNames: []string{"january", "february", "march", "april", "may", "june", "july", "august", "september", "october", "november", "december"},
		Min:        1,
		Max:        12,
		GetIndex: func(t time.Time) int {
			return int(t.Month())
		},
	},
})

func main() {
	// midnight on every full moon
	markers, err := lunarConfig.ParseCronTab("0 full *")
	if err != nil {
		fmt.Printf("parse error: %v\n", err)
		return
	}
	fmt.Printf("expression: %v\n\n", markers)

	fmt.Println("next full moon midnights from 2025-01-01:")
	t0 := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	for i := 0; i < 6; i++ {
		t0, err = lunarConfig.Next(markers, t0)
		if err != nil {
			fmt.Printf("error: %v\n", err)
			return
		}
		fmt.Printf("  %s\n", t0.Format(time.RFC3339))
	}

	// noon during waxing phases (new through first quarter)
	markers, err = lunarConfig.ParseCronTab("12 new-firstquarter *")
	if err != nil {
		fmt.Printf("parse error: %v\n", err)
		return
	}
	fmt.Printf("\nexpression: %v\n\n", markers)

	fmt.Println("next waxing-phase noons from 2025-06-01:")
	t0 = time.Date(2025, 6, 1, 0, 0, 0, 0, time.UTC)
	for i := 0; i < 6; i++ {
		t0, err = lunarConfig.Next(markers, t0)
		if err != nil {
			fmt.Printf("error: %v\n", err)
			return
		}
		fmt.Printf("  %s\n", t0.Format(time.RFC3339))
	}
}
