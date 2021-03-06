package syntax

import "fmt"

// 状态
type state struct {
	start     int
	end       int
	separator int
	state     byte
	err       string // 错误信息
}

// 创建新的state
func newState() *state {
	s := &state{}
	s.reset()

	return s
}

// 设置状态
func (s *state) setStart(index int) {
	if s.state != endByte {
		s.err = fmt.Sprintf("不能嵌套 %s！", string(startByte))
		return
	}

	if s.end+1 == index {
		s.err = "两个命名参数不能相邻！"
		return
	}

	s.start = index
	s.state = startByte
}

// 重置
func (s *state) reset() {
	s.start = 0
	s.end = -10
	s.separator = -10
	s.state = endByte
	s.err = ""
}

// 设置status的end
func (s *state) setEnd(index int) {
	if s.state == endByte {
		s.err = fmt.Sprintf("%s %s 必须成对出现！", string(startByte), string(endByte))
		return
	}

	if index == s.start+1 {
		s.err = "未指定参数名称！"
		return
	}

	if index == s.separator+1 {
		s.err = "未指定的正则表达式！"
		return
	}

	s.state = endByte
	s.end = index
}

// 设置status的Separator
func (s *state) setSeparator(index int) {
	if s.state != startByte {
		s.err = fmt.Sprintf("字符(%s)只能出现在 %s %s 中间！", string(separatorByte), string(startByte), string(endByte))
		return
	}

	if index == s.start+1 {
		s.err = "未指定参数名称！"
		return
	}

	s.state = separatorByte
	s.separator = index
}
