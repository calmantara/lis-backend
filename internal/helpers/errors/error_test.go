package errors

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	err := New("some error")
	assert.EqualError(t, err, "some error")
}

func TestNewWithWrap(t *testing.T) {
	err1 := New("some error")
	err := Wrapf(err1, "other error")
	assert.EqualError(t, err, "other error: some error")

	err2 := Wrapf(err, "other error")
	assert.EqualError(t, err2, "other error: other error: some error")
}

func TestWrapError(t *testing.T) {
	err1 := New("some error")
	err := Wrap(err1, "other error")
	assert.EqualError(t, err, "other error: some error")

	err2 := Wrap(err, "other error 2")
	assert.EqualError(t, err2, "other error 2: other error: some error")

	err3 := Wrap(err2, "other error 3")
	assert.True(t, Is(err3, err1))
}
