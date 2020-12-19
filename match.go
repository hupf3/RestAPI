package mux

import (
	"net/http"
	"strings"

	"github.com/hupf3/RestAPI/major/handlers"
	"github.com/hupf3/RestAPI/params"
	"github.com/issue9/sliceutil"
)

// 验证请求是否符合要求
type Matcher interface {
	Match(*http.Request) bool
}

// 域名匹配
type Hosts struct {
	domains   []string // 域名
	wildcards []string // 泛域名
}

// 声明一个新的Hosts
func NewHosts(domain ...string) *Hosts {
	h := &Hosts{
		domains:   make([]string, 0, len(domain)),
		wildcards: make([]string, 0, len(domain)),
	}

	h.Add(domain...)

	return h
}

// 添加子路由组
func (mux *Mux) NewMux(name string, matcher Matcher) (*Mux, bool) {
	if mux.routers == nil {
		mux.routers = make([]*Mux, 0, 5)
	}

	if sliceutil.Count(mux.routers, func(i int) bool { return mux.routers[i].name == name }) > 0 {
		return nil, false
	}

	m := New(mux.disableOptions, mux.disableHead, mux.skipCleanPath, mux.notFound, mux.methodNotAllowed)
	m.name = name
	m.matcher = matcher
	mux.routers = append(mux.routers, m)
	return m, true
}

// 匹配
func (mux *Mux) match(r *http.Request) (*handlers.Handlers, params.Params) {
	path := r.URL.Path
	if !mux.skipCleanPath {
		path = cleanPath(path)
	}

	for _, m := range mux.routers {
		if m.matcher.Match(r) {
			if hs, ps := m.tree.Handler(path); hs != nil {
				return hs, ps
			}
		}
	}

	if mux.matcher == nil || mux.matcher.Match(r) {
		return mux.tree.Handler(path)
	}
	return nil, nil
}

func (hs *Hosts) Match(r *http.Request) bool {
	hostname := r.URL.Hostname()
	for _, domain := range hs.domains {
		if domain == hostname {
			return true
		}
	}

	for _, wildcard := range hs.wildcards {
		if strings.HasSuffix(hostname, wildcard) {
			return true
		}
	}

	return false
}

// 添加新的域名
func (hs *Hosts) Add(domain ...string) {
	for _, d := range domain {
		switch {
		case strings.HasPrefix(d, "*."):
			d = d[1:]
			if sliceutil.Count(hs.wildcards, func(i int) bool { return d == hs.wildcards[i] }) <= 0 {
				hs.wildcards = append(hs.wildcards, d)
			}
		default:
			if sliceutil.Count(hs.domains, func(i int) bool { return d == hs.domains[i] }) <= 0 {
				hs.domains = append(hs.domains, d)
			}
		}
	}
}

// 删除域名
func (hs *Hosts) Delete(domain string) {
	switch {
	case strings.HasPrefix(domain, "*."):
		size := sliceutil.Delete(hs.wildcards, func(i int) bool { return hs.wildcards[i] == domain[1:] })
		hs.wildcards = hs.wildcards[:size]
	default:
		size := sliceutil.Delete(hs.domains, func(i int) bool { return hs.domains[i] == domain })
		hs.domains = hs.domains[:size]
	}
}
