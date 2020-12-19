package tree

import (
	"net/http"
	"net/http/httptest"

	"github.com/issue9/assert"

	"github.com/hupf3/RestAPI/help/handlers"
	"github.com/hupf3/RestAPI/help/params"
)

func buildHandler(code int) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(code)
	})
}

type tester struct {
	tree *Tree
	a    *assert.Assertion
}

func newTester(a *assert.Assertion) *tester {
	return &tester{
		tree: New(false, false),
		a:    a,
	}
}

func (n *tester) add(method, pattern string, code int) {
	nn, err := n.tree.getNode(pattern)
	n.a.NotError(err).NotNil(nn)

	if nn.handlers == nil {
		nn.handlers = handlers.New(false, false)
	}

	nn.handlers.Add(buildHandler(code), method)
}

func (n *tester) handler(method, path string, code int) (http.Handler, params.Params) {
	hs, ps := n.tree.Handler(path)
	n.a.NotNil(ps).NotNil(hs)

	h := hs.Handler(method)
	n.a.NotNil(h)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(method, path, nil)
	h.ServeHTTP(w, r)
	n.a.Equal(w.Code, code)

	return h, ps
}

func (n *tester) matchTrue(method, path string, code int) {
	h, _ := n.handler(method, path, code)
	n.a.NotNil(h)
}

func (n *tester) paramsTrue(method, path string, code int, params map[string]string) {
	_, ps := n.handler(method, path, code)
	n.a.Equal(ps, params)
}

func (n *tester) urlTrue(pattern string, params map[string]string, url string) {
	u, err := n.tree.URL(pattern, params)
	n.a.NotError(err)
	n.a.Equal(u, url)
}
