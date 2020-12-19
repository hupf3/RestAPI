// Package syntax 负责处理路由语法
package syntax

import (
	"errors"
	"fmt"
	"strings"
)

// 路由项的类型
type Type int8

const (
	String Type = iota // 字符串类型
	Regexp             // 正则表达式
	Named              // 明明参数
)

// 特殊字符
const (
	startByte     = '{' // 其实字符
	endByte       = '}' // 结束字符
	separatorByte = ':' // 分隔符
)

func (t Type) String() string {
	switch t {
	case Named:
		return "named"
	case Regexp:
		return "regexp"
	case String:
		return "string"
	default:
		panic("不存在的类型！")
	}
}

// 转化为正则表达式
var repl = strings.NewReplacer(string(startByte), "(?P<",
	string(separatorByte), ">",
	string(endByte), ")")

// 字符串转化为字符串数组
func Split(str string) ([]*Segment, error) {
	if str == "" {
		return nil, errors.New("参数不能为空！")
	}

	ss := make([]*Segment, 0, strings.Count(str, string(startByte))+1)
	s := newState()

	for i := 0; i < len(str); i++ {
		switch str[i] {
		case startByte:
			start := s.start
			s.setStart(i)

			if s.err == "" && i > 0 {
				ss = append(ss, NewSegment(str[start:i]))
			}
		case separatorByte:
			s.setSeparator(i)
		case endByte:
			s.setEnd(i)
		}

		if s.err != "" {
			return nil, errors.New(s.err)
		}
	}

	if s.start < len(str) {
		if s.state != endByte {
			return nil, fmt.Errorf("缺少 %s 字符！", string(endByte))
		}

		ss = append(ss, NewSegment(str[s.start:]))
	}

	return ss, nil
}
