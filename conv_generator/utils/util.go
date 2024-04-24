package utils

import "reflect"

func DefaultValue[T interface{}]() (v T) {
	return
}
func TypeOf[T interface{}]() reflect.Type {
	return reflect.TypeOf(DefaultValue[T]())
}

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

func CopyValSlice[T interface{}](s []T) []T {
	res := make([]T, len(s))
	for i, v := range s {
		res[i] = v
	}
	return res
}

func CopyPtrSlice[T interface{}](s []*T) []*T {
	res := make([]*T, len(s))
	for i, p := range s {
		res[i] = CopyPtr(p)
	}
	return res
}

func CopyVal2PtrSSlice[T interface{}](s []T) []*T {
	res := make([]*T, len(s))
	for i, v := range s {
		res[i] = GetPtr(v)
	}
	return res
}

func CopyPtr2ValSlice[T interface{}](s []*T) []T {
	res := make([]T, len(s))
	for i, p := range s {
		res[i] = GetVal(p)
	}
	return res
}

/* -------------------------------------  number  -------------------------------------- */
type number interface {
	int | uint |
		int8 | int16 | int32 | int64 |
		uint8 | uint16 | uint32 | uint64 |
		float32 | float64
}

func GetNumPtr[TIn, TOut number](v TIn) *TOut {
	v2 := TOut(v)
	return &v2
}

func CopyNumPtr[TIn, TOut number](p *TIn) *TOut {
	if p == nil {
		return nil
	}
	v := TOut(*p)
	return &v
}

func GetNumVal[TIn, TOut number](p *TIn) TOut {
	if p == nil {
		var v TOut
		return v
	}
	return TOut(*p)
}

func CopyNumValSlice[TIn, TOut number](s []TIn) []TOut {
	res := make([]TOut, len(s))
	for i, v := range s {
		res[i] = TOut(v)
	}
	return res
}

func CopyNumPtrSlice[TIn, TOut number](s []*TIn) []*TOut {
	res := make([]*TOut, len(s))
	for i, p := range s {
		res[i] = CopyNumPtr[TIn, TOut](p)
	}
	return res
}

func CopyNumVal2PtrSSlice[TIn, TOut number](s []TIn) []*TOut {
	res := make([]*TOut, len(s))
	for i, v := range s {
		res[i] = GetNumPtr[TIn, TOut](v)
	}
	return res
}

func CopyNumPtr2ValSlice[TIn, TOut number](s []*TIn) []TOut {
	res := make([]TOut, len(s))
	for i, p := range s {
		res[i] = GetNumVal[TIn, TOut](p)
	}
	return res
}
