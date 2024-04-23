package conv

import (
	"github.com/xdnc/go-tools/con_generator/a"
	"github.com/xdnc/go-tools/con_generator/b"
	"testing"
)

func TestGenConv(t *testing.T) {
	s := GenConv(&a.Req{}, &b.Req{})
	t.Log(s)
}
