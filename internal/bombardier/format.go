package bombardier

import "fmt"

// The formatting below is copied from bombardier's own source code so the
// numbers in generated reports match its console output byte for byte.

// units describes a scale ladder, e.g. us -> ms -> s.
type units struct {
	scale uint64
	base  string
	units []string
}

var (
	binaryUnits = units{
		scale: 1024,
		base:  "",
		units: []string{"KB", "MB", "GB", "TB", "PB"},
	}
	timeUnitsUs = units{
		scale: 1000,
		base:  "us",
		units: []string{"ms", "s"},
	}
	timeUnitsS = units{
		scale: 60,
		base:  "s",
		units: []string{"m", "h"},
	}
)

func formatUnits(n float64, m units, prec int) string {
	amt := n
	unit := m.base

	scale := float64(m.scale) * 0.85

	for i := 0; i < len(m.units) && amt >= scale; i++ {
		amt /= float64(m.scale)
		unit = m.units[i]
	}

	return fmt.Sprintf("%.*f%s", prec, amt, unit)
}

// FormatBinary formats a byte quantity the way bombardier does,
// e.g. 11189196 -> "10.67MB".
func FormatBinary(n float64) string {
	return formatUnits(n, binaryUnits, 2)
}

// FormatTimeUs formats a duration given in microseconds the way
// bombardier does, e.g. 1560.13 -> "1.56ms".
func FormatTimeUs(n float64) string {
	units := timeUnitsUs
	if n >= 1000000.0 {
		n /= 1000000.0
		units = timeUnitsS
	}

	return formatUnits(n, units, 2)
}
