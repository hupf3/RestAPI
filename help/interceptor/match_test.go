package interceptor

import (
	"testing"

	"github.com/issue9/assert"
)

func TestMatchDigit(t *testing.T) {
	a := assert.New(t)

	a.True(MatchDigit("123456"))
	a.True(MatchDigit("0123456"))
}

func TestMatchWord(t *testing.T) {
	a := assert.New(t)

	a.True(MatchWord("123456"))
	a.True(MatchWord("a123456"))
}
