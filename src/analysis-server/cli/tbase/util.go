package tbase

import (
	"math"
	"regexp"
	"strings"

	"github.com/mattn/go-runewidth"
)

var ansi = regexp.MustCompile("\033\\[(?:[0-9]{1,3}(?:;[0-9]{1,3})*)?[m|K]")

func DisplayWidth(str string) int {
	return runewidth.StringWidth(ansi.ReplaceAllLiteralString(str, ""))
}

// Simple Condition for string
// Returns value based on condition
func ConditionString(cond bool, valid, inValid string) string {
	if cond {
		return valid
	}
	return inValid
}

// Pad String
// Attempts to play string in the center
func Pad(s, pad string, width int) string {
	gap := width - DisplayWidth(s)
	if gap > 0 {
		gapLeft := int(math.Ceil(float64(gap / 2)))
		gapRight := gap - gapLeft
		return strings.Repeat(string(pad), gapLeft) + s + strings.Repeat(string(pad), gapRight)
	}
	return s
}

// Pad String Right position
// This would pace string at the left side fo the screen
func PadRight(s, pad string, width int) string {
	gap := width - DisplayWidth(s)
	if gap > 0 {
		return s + strings.Repeat(string(pad), gap)
	}
	return s
}

// Pad String Left position
// This would pace string at the right side fo the screen
func PadLeft(s, pad string, width int) string {
	gap := width - DisplayWidth(s)
	if gap > 0 {
		return strings.Repeat(string(pad), gap) + s
	}
	return s
}
