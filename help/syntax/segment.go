package syntax

import (
	"regexp"
	"strings"

	"github.com/hupf3/RestAPI/help/interceptor"
	"github.com/hupf3/RestAPI/help/params"
)

// 路由项被拆分之后的内容
type Segment struct {
	Value string // 数值
	Type  Type   // 类型

	Endpoint bool // 是否为终点

	Name   string // 参数名称
	Suffix string // 后缀字符串

	expr *regexp.Regexp // 正则表达式参数

	matcher interceptor.MatchFunc // 处理函数
}

// 生命新的Segment结构体
func NewSegment(val string) *Segment {
	seg := &Segment{
		Value: val,
		Type:  String,
	}

	start := strings.IndexByte(val, startByte)
	if start < 0 {
		return seg
	}
	separator := strings.IndexByte(val, separatorByte)
	end := strings.IndexByte(val, endByte)

	if separator < 0 {
		seg.Type = Named
		seg.Name = val[start+1 : end]
		seg.Suffix = val[end+1:]
		seg.Endpoint = val[len(val)-1] == endByte
		seg.matcher = interceptor.MatchAny
		return seg
	}

	matcher, found := interceptor.Get(val[separator+1 : end])
	if found {
		seg.Type = Named
		seg.Name = val[start+1 : separator]
		seg.Suffix = val[end+1:]
		seg.Endpoint = val[len(val)-1] == endByte
		seg.matcher = matcher
		return seg
	}

	seg.Type = Regexp
	seg.expr = regexp.MustCompile(repl.Replace(val))
	seg.Name = val[start+1 : separator]
	seg.Suffix = val[end+1:]
	return seg
}

// 表示相似度
func (seg *Segment) Similarity(s1 *Segment) int {
	if s1.Value == seg.Value {
		return -1
	}

	return longestPrefix(s1.Value, seg.Value)
}

// 拆分
func (seg *Segment) Split(pos int) []*Segment {
	return []*Segment{
		NewSegment(seg.Value[:pos]),
		NewSegment(seg.Value[pos:]),
	}
}

// 路径与当前节点是否匹配
func (seg *Segment) Match(path string, params params.Params) int {
	switch seg.Type {
	case String:
		if strings.HasPrefix(path, seg.Value) {
			return len(seg.Value)
		}
	case Named:
		if seg.Endpoint {
			if seg.matcher(path) {
				params[seg.Name] = path
				return len(path)
			}
		} else if index := strings.Index(path, seg.Suffix); index >= 0 {
			for {
				if val := path[:index]; seg.matcher(val) {
					params[seg.Name] = val
					return index + len(seg.Suffix)
				}

				i := strings.Index(path[index+len(seg.Suffix):], seg.Suffix)
				if i < 0 {
					return -1
				}
				index += i + len(seg.Suffix)
			}
		}
	case Regexp:
		if locs := seg.expr.FindStringSubmatchIndex(path); locs != nil && locs[0] == 0 {
			params[seg.Name] = path[:locs[3]]
			return locs[1]
		}
	}

	return -1
}

// 获取相同的前缀字符串
func longestPrefix(s1, s2 string) int {
	l := len(s1)
	if len(s2) < l {
		l = len(s2)
	}

	startIndex := -10
	endIndex := -10
	state := endByte
	for i := 0; i < l; i++ {
		switch s1[i] {
		case startByte:
			startIndex = i
			state = startByte
		case endByte:
			state = endByte
			endIndex = i
		}

		if s1[i] != s2[i] {
			if state != endByte ||
				endIndex == i {
				return startIndex
			}
			return i
		}
	}

	if endIndex == l-1 {
		return startIndex
	}

	return l
}
