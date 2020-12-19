package mux

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"sync"

	"github.com/issue9/errwrap"

	"github.com/hupf3/RestAPI/help/handlers"
	"github.com/hupf3/RestAPI/help/syntax"
	"github.com/hupf3/RestAPI/help/tree"
	"github.com/hupf3/RestAPI/help/params"
)

// 报错信息
var ErrNameExists = errors.New("存在相同名称的路由项！")

// 返回参数
type Router struct {
	Name   string
	Routes map[string][]string
}

// 路由匹配
type Mux struct {
	name             string     // 路由名称
	routers          []*Mux     // 子路由
	matcher          Matcher    // 路由匹配条件
	tree             *tree.Tree // 路由项
	notFound         http.HandlerFunc
	methodNotAllowed http.HandlerFunc

	disableOptions, disableHead, skipCleanPath bool

	names   map[string]string // 路由项和名称对应关系
	namesMu sync.RWMutex
}

var (
	defaultNotFound = func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
	}

	defaultMethodNotAllowed = func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	}
)

// 默认参数
func Default() *Mux {
	return New(false, false, false, nil, nil)
}

// 声明一个新的Mux
func New(disableOptions, disableHead, skipCleanPath bool, notFound, methodNotAllowed http.HandlerFunc) *Mux {
	if notFound == nil {
		notFound = defaultNotFound
	}
	if methodNotAllowed == nil {
		methodNotAllowed = defaultMethodNotAllowed
	}

	mux := &Mux{
		tree: tree.New(disableOptions, disableHead),

		disableOptions: disableOptions,
		disableHead:    disableHead,
		skipCleanPath:  skipCleanPath,

		names:            make(map[string]string, 50),
		notFound:         notFound,
		methodNotAllowed: methodNotAllowed,
	}

	return mux
}

// 清除路由项
func (mux *Mux) Clean() *Mux {
	for _, m := range mux.routers {
		m.Clean()
	}
	mux.tree.Clean("")

	return mux
}

// 返回所有路由项
func (mux *Mux) All(ignoreHead, ignoreOptions bool) []*Router {
	routers := make([]*Router, 0, len(mux.routers)+1)

	for _, router := range mux.routers {
		routers = append(routers, router.All(ignoreHead, ignoreOptions)...)
	}

	return append(routers, &Router{
		Name:   mux.name,
		Routes: mux.tree.All(ignoreHead, ignoreOptions),
	})
}

// 移除指定路由
func (mux *Mux) Remove(pattern string, methods ...string) *Mux {
	mux.tree.Remove(pattern, methods...)
	return mux
}

// 添加一个新的路由
func (mux *Mux) Handle(pattern string, h http.Handler, methods ...string) error {
	return mux.tree.Add(pattern, h, methods...)
}

// 设置报头的值
func (mux *Mux) SetAllow(pattern string, allow string) error {
	return mux.tree.SetAllow(pattern, allow)
}

// 设置报头的值
func (mux *Mux) Options(pattern string, allow string) *Mux {
	if err := mux.SetAllow(pattern, allow); err != nil {
		panic(err)
	}
	return mux
}

func (mux *Mux) handle(pattern string, h http.Handler, methods ...string) *Mux {
	if err := mux.Handle(pattern, h, methods...); err != nil {
		panic(err)
	}
	return mux
}

// 清除路径中重复字符
func cleanPath(p string) string {
	if p == "" {
		return "/"
	}

	var b errwrap.StringBuilder
	b.Grow(len(p))

	if p[0] != '/' {
		b.WByte('/')
	}

	index := strings.Index(p, "//")
	if index == -1 {
		b.WString(p)
		if b.Err != nil {
			panic(b.Err)
		}
		return b.String()
	}

	b.WString(p[:index+1])

	slash := true
	for i := index + 2; i < len(p); i++ {
		if p[i] == '/' {
			if slash {
				continue
			}
			slash = true
		} else {
			slash = false
		}
		b.WByte(p[i])
	}

	if b.Err != nil {
		panic(b.Err)
	}
	return b.String()
}

// 简化写法

func (mux *Mux) Get(pattern string, h http.Handler) *Mux {
	return mux.handle(pattern, h, http.MethodGet)
}

func (mux *Mux) Post(pattern string, h http.Handler) *Mux {
	return mux.handle(pattern, h, http.MethodPost)
}

func (mux *Mux) Delete(pattern string, h http.Handler) *Mux {
	return mux.handle(pattern, h, http.MethodDelete)
}

func (mux *Mux) Put(pattern string, h http.Handler) *Mux {
	return mux.handle(pattern, h, http.MethodPut)
}

func (mux *Mux) Patch(pattern string, h http.Handler) *Mux {
	return mux.handle(pattern, h, http.MethodPatch)
}

func (mux *Mux) Any(pattern string, h http.Handler) *Mux {
	return mux.handle(pattern, h)
}

func (mux *Mux) HandleFunc(pattern string, fun http.HandlerFunc, methods ...string) error {
	return mux.Handle(pattern, fun, methods...)
}

func (mux *Mux) handleFunc(pattern string, fun http.HandlerFunc, methods ...string) *Mux {
	return mux.handle(pattern, fun, methods...)
}

func (mux *Mux) GetFunc(pattern string, fun http.HandlerFunc) *Mux {
	return mux.handleFunc(pattern, fun, http.MethodGet)
}

func (mux *Mux) PutFunc(pattern string, fun http.HandlerFunc) *Mux {
	return mux.handleFunc(pattern, fun, http.MethodPut)
}

func (mux *Mux) PostFunc(pattern string, fun http.HandlerFunc) *Mux {
	return mux.handleFunc(pattern, fun, http.MethodPost)
}

func (mux *Mux) DeleteFunc(pattern string, fun http.HandlerFunc) *Mux {
	return mux.handleFunc(pattern, fun, http.MethodDelete)
}

func (mux *Mux) PatchFunc(pattern string, fun http.HandlerFunc) *Mux {
	return mux.handleFunc(pattern, fun, http.MethodPatch)
}

func (mux *Mux) AnyFunc(pattern string, fun http.HandlerFunc) *Mux {
	return mux.handleFunc(pattern, fun)
}

func (mux *Mux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	hs, ps := mux.match(r)
	if hs == nil {
		mux.notFound(w, r)
		return
	}

	h := hs.Handler(r.Method)
	if h == nil {
		w.Header().Set("Allow", hs.Options())
		mux.methodNotAllowed(w, r)
		return
	}

	if len(ps) > 0 {
		ctx := context.WithValue(r.Context(), params.ContextKeyParams, ps)
		r = r.WithContext(ctx)
	}

	h.ServeHTTP(w, r)
}

// 命名路由项
func (mux *Mux) Name(name, pattern string) error {
	mux.namesMu.Lock()
	defer mux.namesMu.Unlock()

	if _, found := mux.names[name]; found {
		return ErrNameExists
	}

	mux.names[name] = pattern
	return nil
}

// 根据参数生成URL地址
func (mux *Mux) URL(name string, params map[string]string) (string, error) {
	mux.namesMu.RLock()
	pattern, found := mux.names[name]
	mux.namesMu.RUnlock()

	if !found {
		pattern = name
	}

	return mux.tree.URL(pattern, params)
}

// 路由的参数集合
func Params(r *http.Request) params.Params {
	return params.Get(r)
}

// 判断语法是否正确
func IsWell(pattern string) error {
	_, err := syntax.Split(pattern)
	return err
}

// 返回所有请求方法
func Methods() []string {
	methods := make([]string, len(handlers.Methods))
	copy(methods, handlers.Methods)
	return methods
}
