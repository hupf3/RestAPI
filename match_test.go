package mux

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/issue9/assert"
)

func TestHosts_Match(t *testing.T) {
	a := assert.New(t)

	h := NewHosts("caixw.io", "caixw.oi", "*.example.com")
	a.NotNil(h)
	a.Equal(len(h.domains), 2).
		Equal(len(h.wildcards), 1)

	r := httptest.NewRequest(http.MethodGet, "http://caixw.io/test", nil)
	a.True(h.Match(r))

	r = httptest.NewRequest(http.MethodGet, "https://caixw.io/test", nil)
	a.True(h.Match(r))

	// 泛域名
	r = httptest.NewRequest(http.MethodGet, "https://xx.example.com/test", nil)
	a.True(h.Match(r))

	// 带端口
	r = httptest.NewRequest(http.MethodGet, "http://caixw.io:88/test", nil)
	a.True(h.Match(r))

	// 访问不允许的域名
	r = httptest.NewRequest(http.MethodGet, "http://sub.caixw.io/test", nil)
	a.False(h.Match(r))

	// 访问不允许的域名
	r = httptest.NewRequest(http.MethodGet, "http://sub.1example.com/test", nil)
	a.False(h.Match(r))
}

func TestHosts_Add_Delete(t *testing.T) {
	a := assert.New(t)

	h := NewHosts()

	h.Add("xx.example.com")
	h.Add("xx.example.com")
	h.Add("xx.example.com")
	h.Add("*.example.com")
	h.Add("*.example.com")
	h.Add("*.example.com")
	a.Equal(1, len(h.domains)).
		Equal(1, len(h.wildcards))

	h.Delete("*.example.com")
	a.Equal(1, len(h.domains)).
		Equal(0, len(h.wildcards))

	h.Delete("*.example.com")
	a.Equal(1, len(h.domains)).
		Equal(0, len(h.wildcards))

	h.Delete("xx.example.com")
	a.Equal(0, len(h.domains)).
		Equal(0, len(h.wildcards))
}

func TestMux_NewMux(t *testing.T) {
	a := assert.New(t)

	m := Default()
	r, ok := m.NewMux("host", NewHosts())
	a.True(ok).NotNil(r)
	a.Equal(r.name, "host").Equal(r.disableHead, m.disableHead)

	r, ok = m.NewMux("host", NewHosts())
	a.False(ok).Nil(r)

	r, ok = m.NewMux("host-2", NewHosts())
	a.True(ok).NotNil(r)
	a.Equal(r.name, "host-2").Equal(r.disableHead, m.disableHead)
}

// 带域名的路由项
func TestHosts(t *testing.T) {
	a := assert.New(t)

	m := Default()
	router, ok := m.NewMux("host", NewHosts("localhost"))
	a.True(ok).NotNil(router)
	w := httptest.NewRecorder()
	router.Get("/t1", buildHandler(201))
	r := httptest.NewRequest(http.MethodGet, "/t1", nil)
	m.ServeHTTP(w, r)
	a.Equal(w.Result().StatusCode, 404)

	w = httptest.NewRecorder()
	r = httptest.NewRequest(http.MethodGet, "https://localhost/t1", nil)
	m.ServeHTTP(w, r)
	a.Equal(w.Result().StatusCode, 201)

	w = httptest.NewRecorder()
	r = httptest.NewRequest(http.MethodGet, "https://localhost/t1", nil)
	router.ServeHTTP(w, r) // 由 h 直接访问
	a.Equal(w.Result().StatusCode, 201)

	w = httptest.NewRecorder()
	r = httptest.NewRequest(http.MethodGet, "/t1", nil)
	router.ServeHTTP(w, r) // 由 h 直接访问
	a.Equal(w.Result().StatusCode, 404)

	// resource
	m = Default()
	router, ok = m.NewMux("host", NewHosts("localhost"))
	a.True(ok).NotNil(router)
	res := router.Resource("/r1")
	res.Get(buildHandler(202))
	w = httptest.NewRecorder()
	r = httptest.NewRequest(http.MethodGet, "/r1", nil)
	m.ServeHTTP(w, r)
	a.Equal(w.Result().StatusCode, 404)

	w = httptest.NewRecorder()
	r = httptest.NewRequest(http.MethodGet, "http://localhost/r1", nil)
	m.ServeHTTP(w, r)
	a.Equal(w.Result().StatusCode, 202)

	// prefix
	m = Default()
	router, ok = m.NewMux("host", NewHosts("localhost"))
	a.True(ok).NotNil(router)
	p := router.Prefix("/prefix1")
	p.Get("/p1", buildHandler(203))
	w = httptest.NewRecorder()
	r = httptest.NewRequest(http.MethodGet, "/prefix1/p1", nil)
	m.ServeHTTP(w, r)
	a.Equal(w.Result().StatusCode, 404)

	w = httptest.NewRecorder()
	r = httptest.NewRequest(http.MethodGet, "http://localhost:88/prefix1/p1", nil)
	m.ServeHTTP(w, r)
	a.Equal(w.Result().StatusCode, 203)

	// prefix prefix
	m = Default()
	router, ok = m.NewMux("host", NewHosts("localhost"))
	a.True(ok).NotNil(router)
	p1 := router.Prefix("/prefix1")
	p2 := p1.Prefix("/prefix2")
	p2.GetFunc("/p2", buildFunc(204))
	w = httptest.NewRecorder()
	r = httptest.NewRequest(http.MethodGet, "/prefix1/prefix2/p2", nil)
	m.ServeHTTP(w, r)
	a.Equal(w.Result().StatusCode, 404)

	w = httptest.NewRecorder()
	r = httptest.NewRequest(http.MethodGet, "http://localhost/prefix1/prefix2/p2", nil)
	m.ServeHTTP(w, r)
	a.Equal(w.Result().StatusCode, 204)

	// 第二个 Prefix 为域名
	m = Default()
	p1 = m.Prefix("/prefix1")
	p2 = p1.Prefix("example.com")
	p2.GetFunc("/p2", buildFunc(205))
	w = httptest.NewRecorder()
	r = httptest.NewRequest(http.MethodGet, "/prefix1example.com/p2", nil)
	m.ServeHTTP(w, r)
	a.Equal(w.Result().StatusCode, 205)
}