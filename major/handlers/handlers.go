package handlers

import (
	"fmt"
	"net/http"
	"sort"
)

// 请求
type response struct {
	http.ResponseWriter
}

type optionsState int8
type headState int8

// 处理方式
const (
	optionsStateDefault      optionsState = iota // 默认情况
	optionsStateFixedString                      // 字符串
	optionsStateFixedHandler                     // Handler
	optionsStateDisable                          // 禁用
)

// 是否自动生成
const (
	headStateDefault headState = iota // 无额外操作
	headStateAuto                     // 自动生成
	headStateFixed                    // 固定的
)

// HTTP的所有请求方法
var Methods = []string{
	http.MethodGet,
	http.MethodPost,
	http.MethodDelete,
	http.MethodPut,
	http.MethodPatch,
	http.MethodConnect,
	http.MethodTrace,
	http.MethodOptions,
	http.MethodHead,
}

var addAny = Methods[:len(Methods)-2]

// 对应于每个请求方法的处理函数
type Handlers struct {
	handlers     map[string]http.Handler
	optionsAllow string       // 报头内容。
	optionsState optionsState // 请求的处理方式
	headState    headState
}

func (resp *response) Write(data []byte) (int, error) {
	return 0, nil
}

func (hs *Handlers) getOptionsAllow() string {
	var index int
	for method := range hs.handlers {
		index += methodMap[method]
	}
	return optionsStrings[index]
}

// 声明一个新的 Handler
func New(disableOptions, disableHead bool) *Handlers {
	ret := &Handlers{
		handlers:     make(map[string]http.Handler, 4),
		optionsState: optionsStateDefault,
		headState:    headStateDefault,
	}

	if !disableHead {
		ret.headState = headStateAuto
	}

	if disableOptions {
		ret.optionsState = optionsStateDisable
	} else {
		ret.handlers[http.MethodOptions] = http.HandlerFunc(ret.optionsServeHTTP)
		ret.optionsAllow = ret.getOptionsAllow()
	}

	return ret
}

// 添加一个处理函数
func (hs *Handlers) Add(h http.Handler, methods ...string) error {
	if len(methods) == 0 {
		methods = addAny
	}

	for _, m := range methods {
		if err := hs.addSingle(h, m); err != nil {
			return err
		}
	}

	return nil
}

// 方法是否存在
func methodExists(m string) bool {
	for _, mm := range Methods {
		if mm == m {
			return true
		}
	}
	return false
}

// 添加单个
func (hs *Handlers) addSingle(h http.Handler, m string) error {
	switch m {
	case http.MethodOptions:
		if hs.optionsState == optionsStateFixedHandler {
			return fmt.Errorf("该请求方法 %s 已经存在！", m)
		}

		hs.handlers[m] = h
		hs.optionsState = optionsStateFixedHandler
	case http.MethodHead:
		if hs.headState == headStateFixed {
			return fmt.Errorf("该请求方法 %s 已经存在！", m)
		}

		hs.handlers[m] = h
		hs.headState = headStateFixed
	default:
		if !methodExists(m) {
			return fmt.Errorf("该请求方法 %s 不被支持！", m)
		}

		if _, found := hs.handlers[m]; found {
			return fmt.Errorf("该请求方法 %s 已经存在！", m)
		}
		hs.handlers[m] = h

		if m == http.MethodGet && hs.headState == headStateAuto {
			hs.handlers[http.MethodHead] = hs.headServeHTTP(h)
		}

		if hs.optionsState == optionsStateDefault {
			hs.optionsAllow = hs.getOptionsAllow()
		}
	}
	return nil
}

func (hs *Handlers) optionsServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Allow", hs.optionsAllow)
}

func (hs *Handlers) headServeHTTP(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h.ServeHTTP(&response{ResponseWriter: w}, r)
	})
}

// 去除某个请求方法的处理函数
func (hs *Handlers) Remove(methods ...string) bool {
	if len(methods) == 0 {
		hs.handlers = make(map[string]http.Handler, 8)
		hs.optionsAllow = ""
		return true
	}

	for _, m := range methods {
		delete(hs.handlers, m)

		if m == http.MethodOptions {
			hs.optionsState = optionsStateDisable
		} else if m == http.MethodGet && hs.headState == headStateAuto {
			delete(hs.handlers, http.MethodHead)
		}
	}

	if hs.Len() == 0 {
		hs.optionsAllow = ""
		return true
	}

	if hs.Len() == 1 &&
		hs.handlers[http.MethodOptions] != nil &&
		hs.optionsState == optionsStateDefault {
		delete(hs.handlers, http.MethodOptions)
		hs.optionsAllow = ""
		return true
	}

	if hs.optionsState == optionsStateDefault {
		hs.optionsAllow = hs.getOptionsAllow()
	}
	return false
}

// 设置请求报头
func (hs *Handlers) SetAllow(optionsAllow string) {
	if hs.optionsState == optionsStateDisable {
		hs.handlers[http.MethodOptions] = http.HandlerFunc(hs.optionsServeHTTP)
	}
	hs.optionsAllow = optionsAllow
	hs.optionsState = optionsStateFixedString
}

// 获取指定请求方法对应的处理函数
func (hs *Handlers) Handler(method string) http.Handler {
	return hs.handlers[method]
}

// 获取指定请求方法对应的列表字符串
func (hs *Handlers) Options() string {
	return hs.optionsAllow
}

// 获取指定请求方法对应的数量
func (hs *Handlers) Len() int {
	return len(hs.handlers)
}

// 获取该节点的请求方法
func (hs *Handlers) Methods(ignoreHead, ignoreOptions bool) []string {
	methods := make([]string, 0, len(hs.handlers))

	for key := range hs.handlers {
		if (key == http.MethodOptions && ignoreOptions && hs.optionsState == optionsStateDefault) ||
			key == http.MethodHead && ignoreHead && hs.headState == headStateAuto {
			continue
		}

		methods = append(methods, key)
	}

	sort.Strings(methods)
	return methods
}
