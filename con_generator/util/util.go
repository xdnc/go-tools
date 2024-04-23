package util

func GetPtr[T interface{}](x T) *T {
	return &x
}

func CopyPtr[T interface{}](p *T) *T {
	v := *p
	return &v
}

func GetVal[T interface{}](p *T) T {
	return *p
}
