package mux

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/issue9/assert"
)

func (t *tester) resource(p string) *Resource {
	return t.mux.Resource(p)
}

func TestResource(t *testing.T) {
	a := assert.New(t)
	test := newTester(t, false, true, false)
	h := test.resource("/hupf/1")
	a.NotNil(h)

	h.Get(buildHandler(201))
	test.matchTrue(http.MethodGet, "/hupf/1", 201)
}

func TestMux_Resource(t *testing.T) {
	a := assert.New(t)
	srvmux := New(false, true, false, nil, nil)
	a.NotNil(srvmux)

	r1 := srvmux.Resource("/hupf/1")
	a.NotNil(r1)
	a.Equal(r1.Mux(), srvmux)
	a.Equal(r1.pattern, "/hupf/1")
}

func TestPrefix_Resource(t *testing.T) {
	a := assert.New(t)

	srvmux := Default()
	a.NotNil(srvmux)
	p := srvmux.Prefix("/p")

	r1 := p.Resource("/hupf/1")
	a.NotNil(r1)

	r1.Delete(buildHandler(201))
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodDelete, "/p/hupf/1", nil)
	srvmux.ServeHTTP(w, r)
	a.Equal(w.Result().StatusCode, 201)
}

func TestResource_Name_URL(t *testing.T) {
	a := assert.New(t)
	srvmux := New(false, true, false, nil, nil)
	a.NotNil(srvmux)

	res := srvmux.Resource("/api/v1")
	a.NotNil(res)
	url, err := res.URL(map[string]string{"id": "1"})
	a.NotError(err).Equal(url, "/api/v1")
}
