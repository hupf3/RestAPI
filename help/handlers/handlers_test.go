package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/issue9/assert"
)

var getHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	w.Write([]byte("hupf"))
})

var optionsHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Allow", "options")
})

func TestNew(t *testing.T) {
	a := assert.New(t)

	hs := New(true, true)
	a.NotNil(hs)
	a.NotError(hs.Add(getHandler, http.MethodGet))
	a.Equal(hs.Len(), 1)
}

func TestHandlers_Add(t *testing.T) {
	a := assert.New(t)

	hs := New(false, true)
	a.NotNil(hs)
	a.NotError(hs.Add(getHandler))
	a.Equal(hs.Len(), len(addAny)+1)
}

func TestHandlers_Add_Remove(t *testing.T) {
	a := assert.New(t)

	hs := New(false, false)
	a.NotNil(hs)

	a.NotError(hs.Add(getHandler, http.MethodDelete, http.MethodPost))
	a.Error(hs.Add(getHandler, http.MethodPost))
	a.False(hs.Remove(http.MethodPost))
	a.True(hs.Remove(http.MethodDelete))
	a.True(hs.Remove(http.MethodDelete))
}

func TestHandlers_optionsAllow(t *testing.T) {
	a := assert.New(t)

	hs := New(false, true)
	a.NotNil(hs)

	test := func(allow string) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("OPTIONS", "/empty", nil)
		h := hs.Handler(http.MethodOptions)
		a.NotNil(h)
		h.ServeHTTP(w, r)
		a.Equal(w.Header().Get("Allow"), allow)
	}
	a.Equal(hs.Options(), "OPTIONS")

	a.NotError(hs.Add(getHandler, http.MethodGet))
	test("GET, OPTIONS")
	a.Equal(hs.Options(), "GET, OPTIONS")
}

func TestHandlers_head(t *testing.T) {
	a := assert.New(t)

	hs := New(false, false)
	a.NotNil(hs)

	test := func(val string) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("HEAD", "/empty", nil)
		h := hs.Handler(http.MethodHead)
		a.NotNil(h)
		h.ServeHTTP(w, r)
		a.Equal(w.Body.String(), val)
	}

	// 通过 Get 获取的 Head
	a.NotError(hs.Add(getHandler, http.MethodGet))
	test("")

	// 主动添加
	a.NotError(hs.Add(getHandler, http.MethodHead))
	test("hupf")

}

func TestHandlers_Methods(t *testing.T) {
	a := assert.New(t)

	hs := New(false, false)
	a.NotNil(hs)
	hs.Add(getHandler, http.MethodPut)
	a.Equal(hs.Methods(true, true), []string{http.MethodPut})
}
