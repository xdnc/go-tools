package conv

import (
	"fmt"
	"github.com/xdnc/go-tools/conv_generator/utils"
	"reflect"
	"strings"
)

type Package struct {
	Path  string
	Name  string
	Alias string
}

type Generator struct {
	PackageName string

	functionMap map[string]string

	Packages   []*Package
	packageMap map[string]*Package
}

func NewGenerator(packageName string, funcMap map[string]string, packages ...*Package) *Generator {
	g := &Generator{PackageName: packageName}
	g.AddFunctions(globalFunctionMap)
	g.AddFunctions(funcMap)
	g.AddPackages(packages...)
	return g
}

func (g *Generator) AddPackages(packages ...*Package) {
	if g.packageMap == nil {
		g.packageMap = make(map[string]*Package)
	}
	for _, p := range packages {
		g.packageMap[p.Path] = p
	}
	g.Packages = make([]*Package, 0, len(g.packageMap))
	for _, p := range g.packageMap {
		g.Packages = append(g.Packages, p)
	}
}

func (g *Generator) AddFunction(k, f string) {
	if g.functionMap == nil {
		g.functionMap = make(map[string]string)
	}
	g.functionMap[k] = f
}

func (g *Generator) AddFunctions(functionMap map[string]string) {
	if g.functionMap == nil {
		g.functionMap = make(map[string]string, len(functionMap))
	}
	for k, f := range functionMap {
		g.functionMap[k] = f
	}
}

func (g *Generator) GenConv(from, to interface{}, fileName string) error {
	tFrom, tTo := reflect.TypeOf(from), reflect.TypeOf(to)
	fg := &fileGenerator{
		g:             g,
		from:          tFrom,
		to:            tTo,
		fileName:      fileName,
		stringBuilder: &strings.Builder{},
		genQueue:      [][2]reflect.Type{{tFrom, tTo}},
	}
	_ = fg
	return nil
}

type fileGenerator struct {
	g             *Generator
	from, to      reflect.Type
	fileName      string
	stringBuilder *strings.Builder
	genQueue      [][2]reflect.Type
	nameMap       map[string]string
}

func (fg *fileGenerator) gen() {
	for len(fg.genQueue) > 0 {
		top := fg.genQueue[0]
		fg.genQueue = fg.genQueue[1:]
		from, to := top[0], top[1]
		_, _ = from, to
	}
}

func (fg *fileGenerator) genConv(from, to reflect.Type) {
	if isStructOrStructPtr(to) {

	} else {

	}
}

func (fg *fileGenerator) genStructConv(from, to reflect.Type) {

}

func (fg *fileGenerator) genConvFuncName(from, to reflect.Type) string {
	return fmt.Sprintf("%s2%s",
		convTypeName(fg.formatTypeName(from)), convTypeName(fg.formatTypeName(to)))
}

func convTypeName(tName string) string {
	tName = strings.ReplaceAll(tName, ".", "__")
	tName = strings.ReplaceAll(tName, "*", "p_")
	tName = strings.ReplaceAll(tName, "[]", "s_")
	return tName
}

func (fg *fileGenerator) formatTypeName(t reflect.Type) string {
	rawName := t.String()
	if newName, ok := fg.nameMap[rawName]; ok {
		return newName
	}
	if fg.nameMap == nil {
		fg.nameMap = map[string]string{}
	}
	start := 0
	for start < len(rawName) {
		if isAlpha_(rawName[start]) {
			break
		}
		start++
	}
	for _, p := range fg.g.Packages {
		if p.Alias == "" {
			continue
		}
		if strings.HasPrefix(rawName[start:], p.Name+".") {
			newName := rawName[:start] + strings.Replace(rawName[start:], p.Name+".", p.Alias+".", 1)
			fg.nameMap[rawName] = newName
			return newName
		}
	}
	fg.nameMap[rawName] = rawName
	return rawName
}

func isAlpha_(c byte) bool {
	return c == '_' || (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z')
}

func isStructOrStructPtr(t reflect.Type) bool {
	return t.Kind() == reflect.Struct || (t.Kind() == reflect.Ptr && t.Elem().Kind() == reflect.Struct)
}

func getKey(from, to reflect.Type) string {
	return fmt.Sprintf("%s => %s", from.String(), to.String())
}

var globalFunctionMap map[string]string
var globalNameMap map[string]string

func init() {
	globalFunctionMap = map[string]string{}

	tString, tStringPtr := utils.TypeOf[string](), utils.TypeOf[*string]()

	globalFunctionMap[getKey(tString, tString)] = ""
	globalFunctionMap[getKey(tStringPtr, tStringPtr)] = "utils.CopyPtr"
	globalFunctionMap[getKey(tString, tStringPtr)] = "utils.GetPtr"
	globalFunctionMap[getKey(tStringPtr, tString)] = "utils.GetVal"

	tStringSlice, tStringPtrSlice := utils.TypeOf[[]string](), utils.TypeOf[[]*string]()

	globalFunctionMap[getKey(tStringSlice, tStringSlice)] = "utils.CopyValSlice"
	globalFunctionMap[getKey(tStringPtrSlice, tStringPtrSlice)] = "utils.CopyPtrSlice"
	globalFunctionMap[getKey(tStringSlice, tStringPtrSlice)] = "utils.CopyVal2PtrSSlice"
	globalFunctionMap[getKey(tStringPtrSlice, tStringSlice)] = "utils.CopyPtr2ValSlice"

	numTypes := []reflect.Type{
		utils.TypeOf[byte](), utils.TypeOf[int](), utils.TypeOf[uint](),
		utils.TypeOf[int8](), utils.TypeOf[int16](), utils.TypeOf[int32](), utils.TypeOf[int64](),
		utils.TypeOf[uint8](), utils.TypeOf[uint16](), utils.TypeOf[uint32](), utils.TypeOf[uint64](),
		utils.TypeOf[float32](), utils.TypeOf[float64](),
	}

	numPtrTypes := []reflect.Type{
		utils.TypeOf[*byte](), utils.TypeOf[*int](), utils.TypeOf[*uint](),
		utils.TypeOf[*int8](), utils.TypeOf[*int16](), utils.TypeOf[*int32](), utils.TypeOf[*int64](),
		utils.TypeOf[*uint8](), utils.TypeOf[*uint16](), utils.TypeOf[*uint32](), utils.TypeOf[*uint64](),
		utils.TypeOf[*float32](), utils.TypeOf[*float64](),
	}

	for i, fromType := range numTypes {
		for j, toType := range numTypes {
			if i == j {
				globalFunctionMap[getKey(fromType, toType)] = ""
			} else {
				globalFunctionMap[getKey(fromType, toType)] = toType.String()
			}
		}
	}

	for i, fromType := range numPtrTypes {
		for j, toType := range numPtrTypes {
			if i == j {
				globalFunctionMap[getKey(fromType, toType)] = "utils.CopyPtr"
			} else {
				globalFunctionMap[getKey(fromType, toType)] = fmt.Sprintf(
					"utils.CopyNumPtr[%s,%s]",
					fromType.Elem().String(), toType.Elem().String())
			}
		}
	}

	for i, fromType := range numTypes {
		for j, toType := range numPtrTypes {
			if i == j {
				globalFunctionMap[getKey(fromType, toType)] = "utils.GetPtr"
			} else {
				globalFunctionMap[getKey(fromType, toType)] = fmt.Sprintf(
					"utils.GetNumPtr[%s,%s]",
					fromType.String(), toType.Elem().String())
			}
		}
	}

	for i, fromType := range numPtrTypes {
		for j, toType := range numTypes {
			if i == j {
				globalFunctionMap[getKey(fromType, toType)] = "utils.GetVal"
			} else {
				globalFunctionMap[getKey(fromType, toType)] = fmt.Sprintf(
					"utils.GetNumVal[%s,%s]",
					fromType.Elem().String(), toType.String())
			}
		}
	}

	numSliceTypes := []reflect.Type{
		utils.TypeOf[[]byte](), utils.TypeOf[[]int](), utils.TypeOf[[]uint](),
		utils.TypeOf[[]int8](), utils.TypeOf[[]int16](), utils.TypeOf[[]int32](), utils.TypeOf[[]int64](),
		utils.TypeOf[[]uint8](), utils.TypeOf[[]uint16](), utils.TypeOf[[]uint32](), utils.TypeOf[[]uint64](),
		utils.TypeOf[[]float32](), utils.TypeOf[[]float64](),
	}

	numPtrSliceTypes := []reflect.Type{
		utils.TypeOf[[]*byte](), utils.TypeOf[[]*int](), utils.TypeOf[[]*uint](),
		utils.TypeOf[[]*int8](), utils.TypeOf[[]*int16](), utils.TypeOf[[]*int32](), utils.TypeOf[[]*int64](),
		utils.TypeOf[[]*uint8](), utils.TypeOf[[]*uint16](), utils.TypeOf[[]*uint32](), utils.TypeOf[[]*uint64](),
		utils.TypeOf[[]*float32](), utils.TypeOf[[]*float64](),
	}

	for i, fromType := range numSliceTypes {
		for j, toType := range numSliceTypes {
			if i == j {
				globalFunctionMap[getKey(fromType, toType)] = "utils.CopySlice"
			} else {
				globalFunctionMap[getKey(fromType, toType)] = fmt.Sprintf(
					"utils.CopyNumSlice[%s,%s]",
					fromType.Elem().String(), toType.String())
			}
		}
	}

	for i, fromType := range numPtrSliceTypes {
		for j, toType := range numPtrSliceTypes {
			if i == j {
				globalFunctionMap[getKey(fromType, toType)] = "utils.CopyPtrSlice"
			} else {
				globalFunctionMap[getKey(fromType, toType)] = fmt.Sprintf(
					"utils.CopyNumPtrSlice[%s,%s]",
					fromType.Elem().Elem().String(), toType.Elem().Elem().String())
			}
		}
	}

	for i, fromType := range numSliceTypes {
		for j, toType := range numPtrSliceTypes {
			if i == j {
				globalFunctionMap[getKey(fromType, toType)] = "utils.CopyVal2PtrSSlice"
			} else {
				globalFunctionMap[getKey(fromType, toType)] = fmt.Sprintf(
					"utils.CopyNumVal2PtrSSlice[%s,%s]",
					fromType.Elem().String(), toType.Elem().Elem().String())
			}
		}
	}

	for i, fromType := range numPtrSliceTypes {
		for j, toType := range numSliceTypes {
			if i == j {
				globalFunctionMap[getKey(fromType, toType)] = "utils.CopyPtr2ValSlice"
			} else {
				globalFunctionMap[getKey(fromType, toType)] = fmt.Sprintf(
					"utils.CopyNumPtr2ValSlice[%s,%s]",
					fromType.Elem().Elem().String(), toType.Elem().String())
			}
		}
	}

	globalNameMap = map[string]string{}
	for _, t := range []reflect.Type{tString, tStringPtr, tStringSlice, tStringPtrSlice} {
		globalNameMap[t.String()] = t.String()
	}
	for _, types := range [][]reflect.Type{numTypes, numPtrTypes, numSliceTypes, numPtrSliceTypes} {
		for _, t := range types {
			globalNameMap[t.String()] = t.String()
		}
	}
}
