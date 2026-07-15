package vast

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

// Offset represents either a vast.Duration or a percentage of the video duration.
type Offset struct {
	// If not nil, the Offset is duration based
	Duration *Duration
	// If Duration is nil, the Offset is percent based
	Percent float32
}

// MarshalText implements the encoding.TextMarshaler interface.
func (o Offset) MarshalText() ([]byte, error) {
	if o.Duration != nil {
		return o.Duration.MarshalText()
	}
	if math.IsNaN(float64(o.Percent)) || math.IsInf(float64(o.Percent), 0) || o.Percent < 0 || o.Percent > 1 {
		return nil, fmt.Errorf("invalid offset: %s", strconv.FormatFloat(float64(o.Percent*100), 'f', -1, 32)+"%")
	}
	percent := strconv.FormatFloat(float64(o.Percent*100), 'f', -1, 32)
	return []byte(percent + "%"), nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface.
func (o *Offset) UnmarshalText(data []byte) error {
	value := strings.TrimSpace(string(data))
	if strings.HasSuffix(value, "%") {
		percentage := value[:len(value)-1]
		if !validPercentage(percentage) {
			return fmt.Errorf("invalid offset: %s", data)
		}
		p, err := strconv.ParseFloat(percentage, 32)
		if err != nil || p < 0 || p > 100 {
			return fmt.Errorf("invalid offset: %s", data)
		}
		o.Duration = nil
		o.Percent = float32(p / 100)
		return nil
	}
	var d Duration
	if err := d.UnmarshalText([]byte(value)); err != nil {
		return err
	}
	o.Duration = &d
	o.Percent = 0
	return nil
}

func validPercentage(value string) bool {
	integerDigits := 0
	fractionDigits := 0
	decimalPoint := false
	for _, char := range value {
		switch {
		case char >= '0' && char <= '9':
			if decimalPoint {
				fractionDigits++
			} else {
				integerDigits++
			}
		case char == '.' && !decimalPoint && integerDigits > 0:
			decimalPoint = true
		default:
			return false
		}
	}
	return integerDigits >= 1 && integerDigits <= 3 && (!decimalPoint || fractionDigits > 0)
}
