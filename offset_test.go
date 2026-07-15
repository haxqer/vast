package vast

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOffsetMarshaler(t *testing.T) {
	b, err := Offset{}.MarshalText()
	if assert.NoError(t, err) {
		assert.Equal(t, "0%", string(b))
	}
	b, err = Offset{Percent: .1}.MarshalText()
	if assert.NoError(t, err) {
		assert.Equal(t, "10%", string(b))
	}
	b, err = Offset{Percent: .125}.MarshalText()
	if assert.NoError(t, err) {
		assert.Equal(t, "12.5%", string(b))
	}
	d := Duration(0)
	b, err = Offset{Duration: &d}.MarshalText()
	if assert.NoError(t, err) {
		assert.Equal(t, "00:00:00", string(b))
	}
	_, err = Offset{Percent: 1.01}.MarshalText()
	assert.EqualError(t, err, "invalid offset: 101%")
	for _, percent := range []float32{float32(math.NaN()), float32(math.Inf(1)), float32(math.Inf(-1))} {
		_, err = Offset{Percent: percent}.MarshalText()
		assert.Error(t, err)
	}
}

func TestOffsetUnmarshaler(t *testing.T) {
	var o Offset
	if assert.NoError(t, o.UnmarshalText([]byte("0%"))) {
		assert.Nil(t, o.Duration)
		assert.Equal(t, float32(0.0), o.Percent)
	}
	o = Offset{}
	if assert.NoError(t, o.UnmarshalText([]byte("10%"))) {
		assert.Nil(t, o.Duration)
		assert.Equal(t, float32(0.1), o.Percent)
	}
	o = Offset{}
	if assert.NoError(t, o.UnmarshalText([]byte(" 12.5% "))) {
		assert.Nil(t, o.Duration)
		assert.Equal(t, float32(0.125), o.Percent)
	}
	o = Offset{}
	if assert.NoError(t, o.UnmarshalText([]byte("00:00:00"))) {
		if assert.NotNil(t, o.Duration) {
			assert.Equal(t, Duration(0), *o.Duration)
		}
		assert.Equal(t, float32(0), o.Percent)
	}
	o = Offset{}
	assert.EqualError(t, o.UnmarshalText([]byte("abc%")), "invalid offset: abc%")
	assert.EqualError(t, o.UnmarshalText([]byte("-1%")), "invalid offset: -1%")
	assert.EqualError(t, o.UnmarshalText([]byte("101%")), "invalid offset: 101%")
	assert.EqualError(t, o.UnmarshalText([]byte("not-a-duration")), "invalid duration: not-a-duration")

	for _, value := range []string{
		"1e2%", "NaN%", "Inf%", "+10%", "-0%", ".5%", "1.%", "1.2.3%", "12.5 %", "1234%",
	} {
		t.Run(value, func(t *testing.T) {
			var offset Offset
			assert.EqualError(t, offset.UnmarshalText([]byte(value)), "invalid offset: "+value)
		})
	}
}
