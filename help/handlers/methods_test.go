package handlers

import (
	"testing"

	"github.com/issue9/assert"
)

var max int

func init() {
	for _, v := range methodMap {
		max += v
	}
}

func TestMethods(t *testing.T) {
	a := assert.New(t)

	var val int
	for _, v := range methodMap {
		val += v
	}
	a.Equal(max, val, "methodMap 中的值与 max 不相同！")

	for _, m := range addAny {
		_, found := methodMap[m]
		a.True(found)
	}
}

func TestOptionsStrings(t *testing.T) {
	a := assert.New(t)

	for index, allow := range optionsStrings {
		if index == 0 {
			a.Empty(allow)
		} else {
			a.NotEmpty(allow, "索引 %d 的值为空", index)
		}
	}

	test := func(key int, str string) {
		a.Equal(optionsStrings[key], str, "key:%d,val:%s", key, optionsStrings[key])
	}

	test(0, "")
	test(1, "GET")
}
