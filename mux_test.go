package mux

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/issue9/assert"
	"github.com/issue9/assert/rest"

	"github.com/hupf3/RestAPI/help/handlers"
)

func buildHandler(code int) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(code)
	})
}

func buildFunc(code int) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(code)
	})
}

type tester struct {
	mux *Mux
	srv *rest.Server
}

func newTester(t testing.TB, disableOptions, disableHead, skipClean bool) *tester {
	mux := New(disableOptions, disableHead, skipClean, nil, nil)
	return &tester{
		mux: mux,
		srv: rest.NewServer(t, mux, nil),
	}
}

func (t *tester) matchTrue(method, path string, code int) {
	t.srv.NewRequest(method, path).Do().Status(code)
}

func (t *tester) matchContent(method, path string, code int, content string) {
	t.srv.NewRequest(method, path).Do().Status(code).StringBody(content)
}

func (t *tester) optionsTrue(path string, code int, allow string) {
	t.srv.NewRequest(http.MethodOptions, path).Do().Status(code).Header("Allow", allow)
}

func TestMux(t *testing.T) {
	test := newTester(t, false, true, false)

	test.mux.Get("/", buildHandler(201))
	test.matchTrue(http.MethodGet, "", 201)
	test.matchTrue(http.MethodGet, "/", 201)
	test.matchTrue(http.MethodHead, "/", http.StatusMethodNotAllowed)
	test.matchTrue(http.MethodGet, "/hupf", http.StatusNotFound)
}

func TestMux_All(t *testing.T) {
	a := assert.New(t)

	m := Default()
	a.NotNil(m)

	m.Get("/hupf", buildHandler(1))
	m.Post("/hupf", buildHandler(1))
	a.Equal(m.All(false, false), []*Router{
		{
			Name: "",
			Routes: map[string][]string{
				"/hupf": {"GET", "HEAD", "OPTIONS", "POST"},
			},
		},
	})

}

func TestMux_Head(t *testing.T) {
	test := newTester(t, false, false, false)

	test.mux.Get("/", buildHandler(201))
	test.matchTrue(http.MethodGet, "", 201)
	test.matchTrue(http.MethodGet, "/", 201)
	test.matchTrue(http.MethodHead, "", 201)
	test.matchTrue(http.MethodHead, "/", 201)
	test.matchContent(http.MethodHead, "/", 201, "")
}

func TestMux_Handle_Remove(t *testing.T) {
	a := assert.New(t)
	test := newTester(t, false, true, false)

	a.NotError(test.mux.HandleFunc("/api/1", buildFunc(201), http.MethodGet))
	a.NotError(test.mux.HandleFunc("/api/1", buildFunc(201), http.MethodPut))
	a.NotError(test.mux.HandleFunc("/api/2", buildFunc(202), http.MethodGet))
	test.matchTrue(http.MethodGet, "/api/1", 201)
	test.matchTrue(http.MethodPut, "/api/1", 201)
	test.matchTrue(http.MethodGet, "/api/2", 202)
}

func TestMux_Options(t *testing.T) {
	a := assert.New(t)
	test := newTester(t, false, true, false)

	a.NotError(test.mux.Handle("/api/1", buildHandler(201), http.MethodGet))
	test.optionsTrue("/api/1", http.StatusOK, "GET, OPTIONS")
}

func TestMux_Params(t *testing.T) {
	a := assert.New(t)
	srvmux := Default()
	a.NotNil(srvmux)
	params := map[string]string{}

	buildParamsHandler := func() http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ps := Params(r)
			a.NotNil(ps)
			params = ps
		})
	}

	requestParams := func(method, url string, status int, ps map[string]string) {
		w := httptest.NewRecorder()
		a.NotNil(w)

		r, err := http.NewRequest(method, url, nil)
		a.NotError(err).NotNil(r)

		srvmux.ServeHTTP(w, r)

		a.Equal(w.Code, status)
		if ps != nil {
			a.Equal(params, ps)
		}
		params = nil
	}

	a.NotError(srvmux.Patch("/api/{version:\\d+}", buildParamsHandler()))
	requestParams(http.MethodPatch, "/api/2", http.StatusOK, map[string]string{"version": "2"})
	requestParams(http.MethodPatch, "/api/1", http.StatusOK, map[string]string{"version": "1"})
	requestParams(http.MethodGet, "/api/1", http.StatusMethodNotAllowed, nil)

}

func TestMux_Clean(t *testing.T) {
	a := assert.New(t)

	m := New(false, false, false, nil, nil)
	m.Get("/m", buildHandler(200)).
		Post("/m", buildHandler(201))
	router, ok := m.NewMux("host", NewHosts("hupf.com"))
	a.True(ok).NotNil(router)
	router.Get("/m", buildHandler(202)).
		Post("/m", buildHandler(203))

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/m", nil)
	m.ServeHTTP(w, r)
	a.Equal(w.Result().StatusCode, 200)
}

func TestMux_ServeHTTP(t *testing.T) {
	test := newTester(t, false, true, false)

	test.mux.Handle("/posts/{path}.html", buildHandler(201))
	test.matchTrue(http.MethodGet, "/posts/2020/1.html", 201)
}

func TestMux_ServeHTTP_Order(t *testing.T) {
	a := assert.New(t)
	test := newTester(t, false, true, false)

	a.NotError(test.mux.GetFunc("/posts/1", buildFunc(201)))
	test.matchTrue(http.MethodGet, "/posts/1", 201)

}

func TestMethods(t *testing.T) {
	a := assert.New(t)
	a.Equal(Methods(), handlers.Methods)
}

func TestIsWell(t *testing.T) {
	a := assert.New(t)

	a.Error(IsWell("/{path"))
}

func TestClearPath(t *testing.T) {
	a := assert.New(t)

	a.Equal(cleanPath(""), "/")
}
