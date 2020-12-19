package mux

import "net/http"

// 简化各个函数方法的写法

// 资源地址的路由配置
type Resource struct {
	mux     *Mux
	pattern string
}

func (r *Resource) SetAllow(allow string) error {
	return r.mux.SetAllow(r.pattern, allow)
}

func (r *Resource) Options(allow string) *Resource {
	if err := r.SetAllow(allow); err != nil {
		panic(err)
	}
	return r
}

func (r *Resource) Get(h http.Handler) *Resource {
	return r.handle(h, http.MethodGet)
}

func (r *Resource) Post(h http.Handler) *Resource {
	return r.handle(h, http.MethodPost)
}

func (r *Resource) Delete(h http.Handler) *Resource {
	return r.handle(h, http.MethodDelete)
}

func (r *Resource) Put(h http.Handler) *Resource {
	return r.handle(h, http.MethodPut)
}

func (r *Resource) Patch(h http.Handler) *Resource {
	return r.handle(h, http.MethodPatch)
}

func (r *Resource) Any(h http.Handler) *Resource {
	return r.handle(h)
}

func (r *Resource) GetFunc(fun http.HandlerFunc) *Resource {
	return r.handleFunc(fun, http.MethodGet)
}

func (r *Resource) PutFunc(fun http.HandlerFunc) *Resource {
	return r.handleFunc(fun, http.MethodPut)
}

func (r *Resource) PostFunc(fun http.HandlerFunc) *Resource {
	return r.handleFunc(fun, http.MethodPost)
}

func (r *Resource) DeleteFunc(fun http.HandlerFunc) *Resource {
	return r.handleFunc(fun, http.MethodDelete)
}

func (r *Resource) PatchFunc(fun http.HandlerFunc) *Resource {
	return r.handleFunc(fun, http.MethodPatch)
}

func (r *Resource) AnyFunc(fun http.HandlerFunc) *Resource {
	return r.handleFunc(fun)
}

func (r *Resource) Remove(methods ...string) *Resource {
	r.mux.Remove(r.pattern, methods...)
	return r
}

func (r *Resource) Clean() *Resource {
	r.mux.Remove(r.pattern)
	return r
}

func (r *Resource) Handle(h http.Handler, methods ...string) error {
	return r.mux.Handle(r.pattern, h, methods...)
}

func (r *Resource) handle(h http.Handler, methods ...string) *Resource {
	if err := r.Handle(h, methods...); err != nil {
		panic(err)
	}
	return r
}

func (r *Resource) HandleFunc(fun http.HandlerFunc, methods ...string) error {
	return r.Handle(fun, methods...)
}

func (r *Resource) handleFunc(fun http.HandlerFunc, methods ...string) *Resource {
	if err := r.HandleFunc(fun, methods...); err != nil {
		panic(err)
	}

	return r
}

func (r *Resource) Name(name string) error {
	return r.mux.Name(name, r.pattern)
}

func (r *Resource) URL(params map[string]string) (string, error) {
	return r.mux.URL(r.pattern, params)
}

func (mux *Mux) Resource(pattern string) *Resource {
	return &Resource{
		mux:     mux,
		pattern: pattern,
	}
}

func (p *Prefix) Resource(pattern string) *Resource {
	return &Resource{
		mux:     p.mux,
		pattern: p.prefix + pattern,
	}
}

// 返回 Mux
func (r *Resource) Mux() *Mux {
	return r.mux
}
