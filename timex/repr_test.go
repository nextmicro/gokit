package timex

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestReprOfDuration(t *testing.T) {
	assert.Equal(t, "1000.000ms", ReprOfDuration(time.Second))
	assert.Equal(t, "1111.555ms", ReprOfDuration(
		time.Second+time.Millisecond*111+time.Microsecond*555))
}

func TestDuration(t *testing.T) {
	assert.Equal(t, 1000.0, Duration(time.Second))
}
