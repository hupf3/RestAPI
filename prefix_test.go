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
	p := test.prefix("/p")

	p.Get("/h/1", buildHandler(201))
	test.matchTrue(http.MethodGet, "/p/h/1", 201)
	p.GetFunc("/f/1", buildFunc(201))
	test.matchTrue(http.MethodGet, "/p/f/1", 201)

	p.Post("/h/1", buildHandler(202))
	test.matchTrue(http.MethodPost, "/p/h/1", 202)
	p.PostFunc("/f/1", buildFunc(202))
	test.matchTrue(http.MethodPost, "/p/f/1", 202)

	p.Put("/h/1", buildHandler(203))
	test.matchTrue(http.MethodPut, "/p/h/1", 203)
	p.PutFunc("/f/1", buildFunc(203))
	test.matchTrue(http.MethodPut, "/p/f/1", 203)

	p.Patch("/h/1", buildHandler(204))
	test.matchTrue(http.MethodPatch, "/p/h/1", 204)
	p.PatchFunc("/f/1", buildFunc(204))
	test.matchTrue(http.MethodPatch, "/p/f/1", 204)

	p.Delete("/h/1", buildHandler(205))
	test.matchTrue(http.MethodDelete, "/p/h/1", 205)
	p.DeleteFunc("/f/1", buildFunc(205))
	test.matchTrue(http.MethodDelete, "/p/f/1", 205)

	// Any
	p.Any("/h/any", buildHandler(206))
	test.matchTrue(http.MethodGet, "/p/h/any", 206)
	test.matchTrue(http.MethodPost, "/p/h/any", 206)
	test.matchTrue(http.MethodPut, "/p/h/any", 206)
	test.matchTrue(http.MethodPatch, "/p/h/any", 206)
	test.matchTrue(http.MethodDelete, "/p/h/any", 206)
	test.matchTrue(http.MethodTrace, "/p/h/any", 206)

	p.AnyFunc("/f/any", buildFunc(206))
	test.matchTrue(http.MethodGet, "/p/f/any", 206)
	test.matchTrue(http.MethodPost, "/p/f/any", 206)
	test.matchTrue(http.MethodPut, "/p/f/any", 206)
	test.matchTrue(http.MethodPatch, "/p/f/any", 206)
	test.matchTrue(http.MethodDelete, "/p/f/any", 206)
	test.matchTrue(http.MethodTrace, "/p/f/any", 206)

	p.Options("/h/1", "ABC")
	test.optionsTrue("/p/h/1", 200, "ABC")

	a := assert.New(t)
	a.Panic(func() {
		p.Options("/h/{1", "ABC")
	})

	// remove
	p.Remove("/f/any", http.MethodDelete, http.MethodGet)
	test.matchTrue(http.MethodGet, "/p/f/any", 405)   // 已经删除
	test.matchTrue(http.MethodTrace, "/p/f/any", 206) // 未删除

	// clean
	p.Clean()
	test.matchTrue(http.MethodTrace, "/p/f/any", 404)
	test.optionsTrue("/p/h/1", 404, "")
	test.matchTrue(http.MethodDelete, "/p/f/1", 404)
}

func TestMux_Prefix(t *testing.T) {
	a := assert.New(t)
	srvmux := New(false, true, false, nil, nil)
	a.NotNil(srvmux)

	p := srvmux.Prefix("/abc")
	a.Equal(p.prefix, "/abc")
	a.Equal(p.Mux(), srvmux)

	p = srvmux.Prefix("")
	a.Equal(p.prefix, "")
}

func TestPrefix_Prefix(t *testing.T) {
	a := assert.New(t)
	srvmux := New(false, true, false, nil, nil)
	a.NotNil(srvmux)

	p := srvmux.Prefix("/abc")
	pp := p.Prefix("/def")
	a.Equal(pp.prefix, "/abc/def")
	a.Equal(p.Mux(), srvmux)

	p = srvmux.Prefix("")
	pp = p.Prefix("/abc")
	a.Equal(pp.prefix, "/abc")
}

func TestPrefix_Name_URL(t *testing.T) {
	a := assert.New(t)
	srvmux := New(false, true, false, nil, nil)
	a.NotNil(srvmux)

	// 非正则
	p := srvmux.Prefix("/api")
	p.Any("/v1", nil)
	a.NotNil(p)
	url, err := p.URL("/v1", map[string]string{"id": "1"})
	a.NotError(err).Equal(url, "/api/v1")

	// 正常的单个参数
	p = srvmux.Prefix("/api")
	p.Any("/{id:\\d+}/{path}", nil)
	a.NotNil(p)
	url, err = p.URL("/{id:\\d+}/{path}", map[string]string{"id": "1", "path": "p1"})
	a.NotError(err).Equal(url, "/api/1/p1")

	// 多个参数
	p = srvmux.Prefix("/api")
	p.Any("/{action}/{id:\\d+}", nil)
	a.NotNil(p)
	url, err = p.URL("/{action}/{id:\\d+}", map[string]string{"id": "1", "action": "blog"})
	a.NotError(err).Equal(url, "/api/blog/1")
	// 缺少参数
	url, err = p.URL("/{action}/{id:\\d+}", map[string]string{"id": "1"})
	a.Error(err).Equal(url, "")

	a.NotError(p.Name("action", "/{action}/{id:\\d+}"))
	url, err = p.Mux().URL("action", map[string]string{"id": "1", "action": "blog"})
	a.NotError(err).Equal(url, "/api/blog/1")

	url, err = p.URL("action", map[string]string{"id": "1", "action": "blog"})
	a.NotError(err).Equal(url, "/api/blog/1")
}