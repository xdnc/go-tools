package b

type Header struct {
	User     string
	Password *string
}

func (h *Header) GetUser() string {
	if h == nil {
		return ""
	}
	return h.User
}

func (h *Header) GetPassword() string {
	if h == nil || h.Password == nil {
		return ""
	}
	return *h.Password
}

type Item struct {
	Val int
}

func (i *Item) GetVal() int {
	if i == nil {
		return 0
	}
	return i.Val
}

type Req struct {
	Header  Header
	PHeader *Header
	Int32   int32
	PInt32  *int32
	String  string
	PString *string
	Strings []string
	Items   []Item
	PItems  []*Item
}

func (r *Req) GetHeader() Header {
	if r == nil {
		return Header{}
	}
	return r.Header
}

func (r *Req) GetPHeader() *Header {
	if r == nil {
		return nil
	}
	return r.PHeader
}

func (r *Req) GetInt32() int32 {
	if r == nil {
		return 0
	}
	return r.Int32
}

func (r *Req) GetPInt32() int32 {
	if r == nil || r.PInt32 == nil {
		return 0
	}
	return *r.PInt32
}

func (r *Req) GetString() string {
	if r == nil {
		return ""
	}
	return r.String
}

func (r *Req) GetPString() string {
	if r == nil || r.PString == nil {
		return ""
	}
	return *r.PString
}

func (r *Req) GetStrings() []string {
	if r == nil {
		return nil
	}
	return r.Strings
}

func (r *Req) GetItems() []Item {
	if r == nil {
		return nil
	}
	return r.Items
}

func (r *Req) GetPItems() []*Item {
	if r == nil {
		return nil
	}
	return r.PItems
}
