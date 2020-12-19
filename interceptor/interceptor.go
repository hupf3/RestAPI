package interceptor

import (
	"fmt"
	"sync"
)

var interceptors = &sync.Map{}

// 查找制定的处理函数
func Get(v string) (MatchFunc, bool) {
	f, found := interceptors.Load(v)
	if !found {
		return nil, false
	}
	return f.(MatchFunc), true
}

// 对于注册的处理
func Register(f MatchFunc, val ...string) error {
	if len(val) == 0 {
		panic("参数不能为空！")
	}

	for _, v := range val {
		if _, exists := interceptors.Load(v); exists {
			return fmt.Errorf("%s 已经存在！", v)
		}

		interceptors.Store(v, f)
	}

	return nil
}

// 对于注销的处理
func Deregister(val ...string) {
	for _, v := range val {
		interceptors.Delete(v)
	}
}
