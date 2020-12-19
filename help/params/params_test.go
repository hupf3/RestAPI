package params

import (
	"context"
	"net/http"
	"testing"

	"github.com/issue9/assert"
)

func getParams(params map[string]string, a *assert.Assertion) Params {
	r, err := http.NewRequest(http.MethodGet, "/to/path", nil)
	a.NotError(err).NotNil(r)

	ctx := context.WithValue(r.Context(), ContextKeyParams, Params(params))
	r = r.WithContext(ctx)
	return Get(r)
}

func TestGetParams(t *testing.T) {
	a := assert.New(t)

	r, err := http.NewRequest(http.MethodGet, "/to/path", nil)
	a.NotError(err).NotNil(r)
	ps := Get(r)
	a.Nil(ps)

	maps := map[string]string{"key1": "1"}
	r, err = http.NewRequest(http.MethodGet, "/to/path", nil)
	a.NotError(err).NotNil(r)
	ctx := context.WithValue(r.Context(), ContextKeyParams, Params(maps))
	r = r.WithContext(ctx)
	ps = Get(r)
	a.Equal(ps, maps)
}

func TestParams_String(t *testing.T) {
	a := assert.New(t)

	ps := getParams(map[string]string{
		"key": "1",
	}, a)

	val, err := ps.String("key")
	a.NotError(err).Equal(val, "1")
	a.True(ps.Exists("key"))
	a.Equal(ps.MustString("key", "2"), "1")

}

func TestParams_Int(t *testing.T) {
	a := assert.New(t)

	ps := getParams(map[string]string{
		"key": "1",
	}, a)

	val, err := ps.Int("key")
	a.NotError(err).Equal(val, 1)
	a.Equal(ps.MustInt("key", 2), 1)
}

func TestParams_Uint(t *testing.T) {
	a := assert.New(t)

	ps := getParams(map[string]string{
		"key": "1",
	}, a)

	val, err := ps.Uint("key")
	a.NotError(err).Equal(val, 1)
	a.Equal(ps.MustUint("key", 9), 1)
}

func TestParams_Bool(t *testing.T) {
	a := assert.New(t)

	ps := getParams(map[string]string{
		"key": "true",
	}, a)

	val, err := ps.Bool("key")
	a.NotError(err).True(val)
	a.True(ps.MustBool("key", false))
}

func TestParams_Float(t *testing.T) {
	a := assert.New(t)

	ps := getParams(map[string]string{
		"key": "1",
	}, a)

	val, err := ps.Float("key")
	a.NotError(err).Equal(val, 1.0)
	a.Equal(ps.MustFloat("key", 2.0), 1.0)
}
