package conv

import (
	"github.com/xdnc/go-tools/conv_generator/a"
	"github.com/xdnc/go-tools/conv_generator/b"
	"testing"
)

type ARep = a.Req
type BReq = b.Req

func TestGenConv(t *testing.T) {
	s := GenConv(&ARep{}, &BReq{})
	t.Log(s)
}
