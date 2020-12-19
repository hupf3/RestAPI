package syntax

import (
	"testing"

	"github.com/issue9/assert"

	"github.com/hupf3/RestAPI/help/params"
)

func TestNewSegment(t *testing.T) {
	a := assert.New(t)

	seg := NewSegment("/post/1")
	a.Equal(seg.Type, String).Equal(seg.Value, "/post/1")
}

func TestLongestPrefix(t *testing.T) {
	a := assert.New(t)

	test := func(s1, s2 string, len int) {
		a.Equal(longestPrefix(s1, s2), len)
	}

	test("", "", 0)
	test("/", "", 0)
	test("/test", "test", 0)
}

func TestSegment_Similarity(t *testing.T) {
	a := assert.New(t)

	seg := NewSegment("{id}/author")
	a.NotNil(seg)

	s1 := NewSegment("{id}/author")
	a.Equal(-1, seg.Similarity(s1))
}

func TestSegment_Split(t *testing.T) {
	a := assert.New(t)

	seg := NewSegment("{id}/author")
	a.NotNil(seg)

	segs := seg.Split(4)
	a.Equal(segs[0].Value, "{id}").
		Equal(segs[1].Value, "/author")
}

func TestSegment_Match(t *testing.T) {
	a := assert.New(t)

	seg := NewSegment("{id:any}/author")
	ps := params.Params{}
	a.NotNil(seg)
	path := "1/author"
	index := seg.Match(path, ps)
	a.Empty(path[index:]).
		Equal(ps, params.Params{"id": "1"})
}
