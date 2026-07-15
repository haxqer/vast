package vast

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// Duration is a VAST duration expressed a hh:mm:ss
type Duration time.Duration

// MarshalText implements the encoding.TextMarshaler interface.
func (dur Duration) MarshalText() ([]byte, error) {
	if dur < 0 {
		return nil, fmt.Errorf("invalid duration: %s", time.Duration(dur))
	}
	h := dur / Duration(time.Hour)
	m := dur % Duration(time.Hour) / Duration(time.Minute)
	s := dur % Duration(time.Minute) / Duration(time.Second)
	ms := dur % Duration(time.Second) / Duration(time.Millisecond)
	if ms == 0 {
		return []byte(fmt.Sprintf("%02d:%02d:%02d", h, m, s)), nil
	}
	return []byte(fmt.Sprintf("%02d:%02d:%02d.%03d", h, m, s, ms)), nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface.
func (dur *Duration) UnmarshalText(data []byte) (err error) {
	s := string(data)
	s = strings.TrimSpace(s)
	if s == "" || strings.ToLower(s) == "undefined" {
		*dur = 0
		return nil
	}
	parts := strings.SplitN(s, ":", 3)
	if len(parts) != 3 {
		return fmt.Errorf("invalid duration: %s", data)
	}
	var parsed Duration
	if i := strings.IndexByte(parts[2], '.'); i > 0 {
		fraction := parts[2][i+1:]
		if len(fraction) != 3 {
			return fmt.Errorf("invalid duration: %s", data)
		}
		ms, err := strconv.ParseInt(fraction, 10, 32)
		if err != nil {
			return fmt.Errorf("invalid duration: %s", data)
		}
		parts[2] = parts[2][:i]
		parsed += Duration(ms) * Duration(time.Millisecond)
	}
	f := Duration(time.Second)
	for i := 2; i >= 0; i-- {
		n, err := strconv.ParseInt(parts[i], 10, 32)
		if err != nil || n < 0 || n > 59 {
			return fmt.Errorf("invalid duration: %s", data)
		}
		parsed += Duration(n) * f
		f *= 60
	}
	*dur = parsed
	return nil
}
