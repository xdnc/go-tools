package gen

import (
	"github.com/xdnc/go-tools/con_generator/a"
	"github.com/xdnc/go-tools/con_generator/b"
	"github.com/xdnc/go-tools/con_generator/util"
)

func p_a__Req2p_b__Req(in *a.Req) *b.Req {
	return &b.Req{
		Header:  a__Header2b__Header(in.GetHeader()),
		PHeader: p_a__Header2p_b__Header(in.GetPHeader()),
		Int32:   in.GetInt32(),
		PInt32:  util.GetPtr(in.GetPInt32()),
		String:  in.GetString(),
		PString: util.GetPtr(in.GetPString()),
		Strings: s_string2s_string(in.GetStrings()),
		Items:   s_a__Item2s_b__Item(in.GetItems()),
		PItems:  s_p_a__Item2s_p_b__Item(in.GetPItems()),
	}
}

func a__Header2b__Header(in a.Header) b.Header {
	return b.Header{
		User:     in.GetUser(),
		Password: util.GetPtr(in.GetPassword()),
	}
}

func p_a__Header2p_b__Header(in *a.Header) *b.Header {
	return &b.Header{
		User:     in.GetUser(),
		Password: util.GetPtr(in.GetPassword()),
	}
}

func s_string2s_string(in []string) []string {
	out := make([]string, len(in))
	for i := range in {
		out[i] = in[i]
	}
	return out
}

func s_a__Item2s_b__Item(in []a.Item) []b.Item {
	out := make([]b.Item, len(in))
	for i := range in {
		out[i] = a__Item2b__Item(in[i])
	}
	return out
}

func a__Item2b__Item(in a.Item) b.Item {
	return b.Item{
		Val: in.GetVal(),
	}
}

func s_p_a__Item2s_p_b__Item(in []*a.Item) []*b.Item {
	out := make([]*b.Item, len(in))
	for i := range in {
		out[i] = p_a__Item2p_b__Item(in[i])
	}
	return out
}

func p_a__Item2p_b__Item(in *a.Item) *b.Item {
	return &b.Item{
		Val: in.GetVal(),
	}
}
