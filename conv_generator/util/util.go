package util

func GetPtr[T interface{}](v T) *T {
	return &v
}

func CopyPtr[T interface{}](p *T) *T {
	if p == nil {
		return nil
	}
	v := *p
	return &v
}

func GetVal[T interface{}](p *T) T {
	if p == nil {
		var v T
		return v
	}
	return *p
}
