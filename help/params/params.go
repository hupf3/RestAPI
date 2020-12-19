// Package params 获取和转换路由中的参数信息
package params

import (
	"errors"
	"net/http"
	"strconv"
)

type contextKey int

// 存取路由参数的关键字
const ContextKeyParams contextKey = 0

// 报错信息
var ErrParamNotExists = errors.New("不存在该参数")

// 获取路由中的参数信息
type Params map[string]string

// 获取一个Params
func Get(r *http.Request) Params {
	if params := r.Context().Value(ContextKeyParams); params != nil {
		return params.(Params)
	}

	return nil
}

// 查找制定的参数是否存在
func (p Params) Exists(key string) bool {
	_, found := p[key]
	return found
}

// 将变量转化为int
func (p Params) Int(key string) (int64, error) {
	str, found := p[key]
	if !found {
		return 0, ErrParamNotExists
	}

	return strconv.ParseInt(str, 10, 64)
}

// 将变量转化为uint
func (p Params) Uint(key string) (uint64, error) {
	str, found := p[key]
	if !found {
		return 0, ErrParamNotExists
	}

	return strconv.ParseUint(str, 10, 64)
}

// 将变量转化为float
func (p Params) Float(key string) (float64, error) {
	str, found := p[key]
	if !found {
		return 0, ErrParamNotExists
	}

	return strconv.ParseFloat(str, 64)
}

// 将变量转化为bool
func (p Params) Bool(key string) (bool, error) {
	str, found := p[key]
	if !found {
		return false, ErrParamNotExists
	}

	return strconv.ParseBool(str)
}

// 将变量转化为string
func (p Params) String(key string) (string, error) {
	v, found := p[key]
	if !found {
		return "", ErrParamNotExists
	}

	return v, nil
}

// 将变量转化为int
func (p Params) MustInt(key string, def int64) int64 {
	str, found := p[key]
	if !found {
		return def
	}

	if val, err := strconv.ParseInt(str, 10, 64); err == nil {
		return val
	}

	return def
}

// 将变量转化为uint
func (p Params) MustUint(key string, def uint64) uint64 {
	str, found := p[key]
	if !found {
		return def
	}

	if val, err := strconv.ParseUint(str, 10, 64); err == nil {
		return val
	}

	return def
}

// 将变量转化为float
func (p Params) MustFloat(key string, def float64) float64 {
	str, found := p[key]
	if !found {
		return def
	}

	if val, err := strconv.ParseFloat(str, 64); err == nil {
		return val
	}

	return def
}

// 将变量转化为bool
func (p Params) MustBool(key string, def bool) bool {
	str, found := p[key]
	if !found {
		return def
	}

	if val, err := strconv.ParseBool(str); err == nil {
		return val
	}

	return def
}

// 将变量转化为string
func (p Params) MustString(key, def string) string {
	v, found := p[key]
	if !found {
		return def
	}

	return v
}
