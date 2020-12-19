package mux

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/issue9/assert"
)

func TestHosts_Add_Delete(t *testing.T) {
	a := assert.New(t)
	h := NewHosts()

	h.Add("xx.example.com")
	h.Add("*.example.com")
	a.Equal(1, len(h.domains)).
		Equal(1, len(h.wildcards))
}

func TestMux_NewMux(t *testing.T) {
	a := assert.New(t)

	m := Default()
	r, ok := m.NewMux("host", NewHosts())
	a.True(ok).NotNil(r)
	a.Equal(r.name, "host").Equal(r.disableHead, m.disableHead)
}

func TestHosts(t *testing.T) {
	a := assert.New(t)

	m := Default()
	router, ok := m.NewMux("host", NewHosts("localhost"))
	a.True(ok).NotNil(router)
	w := httptest.NewRecorder()
	router.Get("/h", buildHandler(201))
	r := httptest.NewRequest(http.MethodGet, "/h", nil)
	m.ServeHTTP(w, r)
	a.Equal(w.Result().StatusCode, 404)
}
