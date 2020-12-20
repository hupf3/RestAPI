package interceptor

type MatchFunc func(string) bool

// 初始化
func init() {
	if err := Register(MatchDigit, "digit"); err != nil {
		panic(err)
	}

	if err := Register(MatchDigit, "word"); err != nil {
		panic(err)
	}

	if err := Register(MatchAny, "any"); err != nil {
		panic(err)
	}
}

// 匹配数字
func MatchDigit(path string) bool {
	for _, c := range path {
		if c < '0' || c > '9' {
			return false
		}
	}
	return len(path) > 0
}

// 匹配单词
func MatchWord(path string) bool {
	for _, c := range path {
		if (c < '0' || c > '9') && (c < 'a' || c > 'z') && (c < 'A' || c > 'Z') {
			return false
		}
	}
	return len(path) > 0
}

// 匹配任意的非空内容
func MatchAny(path string) bool { return len(path) > 0 }
