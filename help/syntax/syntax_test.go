package syntax

import (
	"testing"

	"github.com/issue9/assert"
)

func TestRegexp(t *testing.T) {
	a := assert.New(t)

	a.Equal(repl.Replace("{id:\\d+}"), "(?P<id>\\d+)")
}

func TestType_String(t *testing.T) {
	a := assert.New(t)

	a.Equal(Named.String(), "named")
	a.Equal(Regexp.String(), "regexp")
	a.Equal(String.String(), "string")
}

func TestSplit(t *testing.T) {
	a := assert.New(t)
	test := func(str string, isError bool, ss ...*Segment) {
		s, err := Split(str)

		if isError {
			a.Error(err)
			return
		}

		a.NotError(err).Equal(len(s), len(ss))
		for index, seg := range ss {
			item := s[index]
			a.Equal(seg.Value, item.Value).
				Equal(seg.Name, item.Name).
				Equal(seg.Endpoint, item.Endpoint).
				Equal(seg.Suffix, item.Suffix)
		}
	}

	test("/", false, NewSegment("/"))

	test("/posts/1", false, NewSegment("/posts/1"))
}
