package mux

import "net/http"

// 简化写法

// 将前缀一样的路由项统一起来
type Prefix struct {
	mux    *Mux
	prefix string
}

func (p *Prefix) SetAllow(pattern string, allow string) error {
	return p.mux.SetAllow(p.prefix+pattern, allow)
}

func (p *Prefix) Options(pattern string, allow string) *Prefix {
	if err := p.SetAllow(pattern, allow); err != nil {
		panic(err)
	}
	return p
}

func (p *Prefix) Get(pattern string, h http.Handler) *Prefix {
	return p.handle(pattern, h, http.MethodGet)
}

func (p *Prefix) Post(pattern string, h http.Handler) *Prefix {
	return p.handle(pattern, h, http.MethodPost)
}

func (p *Prefix) Delete(pattern string, h http.Handler) *Prefix {
	return p.handle(pattern, h, http.MethodDelete)
}

func (p *Prefix) Put(pattern string, h http.Handler) *Prefix {
	return p.handle(pattern, h, http.MethodPut)
}

func (p *Prefix) Patch(pattern string, h http.Handler) *Prefix {
	return p.handle(pattern, h, http.MethodPatch)
}

func (p *Prefix) Any(pattern string, h http.Handler) *Prefix {
	return p.handle(pattern, h)
}

func (p *Prefix) GetFunc(pattern string, fun http.HandlerFunc) *Prefix {
	return p.handleFunc(pattern, fun, http.MethodGet)
}

func (p *Prefix) PutFunc(pattern string, fun http.HandlerFunc) *Prefix {
	return p.handleFunc(pattern, fun, http.MethodPut)
}

func (p *Prefix) PostFunc(pattern string, fun http.HandlerFunc) *Prefix {
	return p.handleFunc(pattern, fun, http.MethodPost)
}

func (p *Prefix) DeleteFunc(pattern string, fun http.HandlerFunc) *Prefix {
	return p.handleFunc(pattern, fun, http.MethodDelete)
}

func (p *Prefix) PatchFunc(pattern string, fun http.HandlerFunc) *Prefix {
	return p.handleFunc(pattern, fun, http.MethodPatch)
}

func (p *Prefix) AnyFunc(pattern string, fun http.HandlerFunc) *Prefix {
	return p.handleFunc(pattern, fun)
}

func (p *Prefix) Remove(pattern string, methods ...string) *Prefix {
	p.mux.Remove(p.prefix+pattern, methods...)
	return p
}

func (p *Prefix) Clean() *Prefix {
	p.mux.tree.Clean(p.prefix)
	return p
}

func (p *Prefix) Handle(pattern string, h http.Handler, methods ...string) error {
	return p.mux.Handle(p.prefix+pattern, h, methods...)
}

func (p *Prefix) handle(pattern string, h http.Handler, methods ...string) *Prefix {
	if err := p.Handle(pattern, h, methods...); err != nil {
		panic(err)
	}

	return p
}

func (p *Prefix) HandleFunc(pattern string, fun http.HandlerFunc, methods ...string) error {
	return p.Handle(pattern, fun, methods...)
}

func (p *Prefix) handleFunc(pattern string, fun http.HandlerFunc, methods ...string) *Prefix {
	if err := p.HandleFunc(pattern, fun, methods...); err != nil {
		panic(err)
	}
	return p
}

// 命名
func (p *Prefix) Name(name, pattern string) error {
	return p.mux.Name(name, p.prefix+pattern)
}

// 生成URL
func (p *Prefix) URL(name string, params map[string]string) (string, error) {
	p.mux.namesMu.RLock()
	pattern, found := p.mux.names[name]
	p.mux.namesMu.RUnlock()

	if !found {
		pattern = p.prefix + name
	}

	return p.mux.tree.URL(pattern, params)
}

func (p *Prefix) Prefix(prefix string) *Prefix {
	return &Prefix{
		mux:    p.mux,
		prefix: p.prefix + prefix,
	}
}

func (mux *Mux) Prefix(prefix string) *Prefix {
	return &Prefix{
		mux:    mux,
		prefix: prefix,
	}
}

// 返回相关Mux
func (p *Prefix) Mux() *Mux {
	return p.mux
}
