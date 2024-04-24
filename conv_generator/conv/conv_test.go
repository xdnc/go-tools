package conv

import (
	"github.com/xdnc/go-tools/conv_generator/a"
	"github.com/xdnc/go-tools/conv_generator/b"
	"sort"
	"testing"
)

type ARep = a.Req
type BReq = b.Req

func TestGenConv(t *testing.T) {
	s := GenConv(&ARep{}, &BReq{})
	t.Log(s)
}

func TestGenerator(t *testing.T) {
	t.Log(globalFunctionMap)
	t.Log(globalNameMap)
	t.Log()
	logMap(t, globalFunctionMap)
	t.Log()
	logMap(t, globalNameMap)
}

func logMap(t *testing.T, m map[string]string) {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		v := m[k]
		t.Logf("\"%s\": \"%s\"\n", k, v)
	}
}
