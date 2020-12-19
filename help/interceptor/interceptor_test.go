package interceptor

import (
	"testing"

	"github.com/issue9/assert"
)

func TestInterceptors(t *testing.T) {
	a := assert.New(t)

	a.Panic(func() {
		Register(MatchAny)
	})

	a.NotError(Register(MatchWord, "[a-zA-Z0-9]+", "hupf"))
	a.Error(Register(MatchWord, "[a-zA-Z0-9]+"))
	_, found := Get("hupf")
	a.True(found)
	_, found = Get("[a-zA-Z0-9]+")
	a.True(found)
}
