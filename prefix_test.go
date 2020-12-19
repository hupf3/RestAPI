package mux

import (
	"net/http"
	"testing"

	"github.com/issue9/assert"
)

func (t *tester) prefix(p string) *Prefix {
	return t.mux.Prefix(p)
}

func TestPrefix(t *testing.T) {
	test := newTester(t, false, true, false)
	p := test.prefix("/hupf")

	p.Get("/t/1", buildHandler(201))
	test.matchTrue(http.MethodGet, "/hupf/t/1", 201)
	p.GetFunc("/h/1", buildFunc(201))
	test.matchTrue(http.MethodGet, "/hupf/h/1", 201)
}

func TestMux_Prefix(t *testing.T) {
	a := assert.New(t)
	srvmux := New(false, true, false, nil, nil)
	a.NotNil(srvmux)

	p := srvmux.Prefix("/hupf")
	a.Equal(p.prefix, "/hupf")
	a.Equal(p.Mux(), srvmux)
}

func TestPrefix_Prefix(t *testing.T) {
	a := assert.New(t)
	srvmux := New(false, true, false, nil, nil)
	a.NotNil(srvmux)

	p := srvmux.Prefix("/hupf")
	pp := p.Prefix("/t")
	a.Equal(pp.prefix, "/hupf/t")
	a.Equal(p.Mux(), srvmux)
}

func TestPrefix_Name_URL(t *testing.T) {
	a := assert.New(t)
	srvmux := New(false, true, false, nil, nil)
	a.NotNil(srvmux)

	p := srvmux.Prefix("/api")
	p.Any("/v1", nil)
	a.NotNil(p)
	url, err := p.URL("/v1", map[string]string{"id": "1"})
	a.NotError(err).Equal(url, "/api/v1")
}
