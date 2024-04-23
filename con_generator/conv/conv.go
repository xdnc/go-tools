package conv

import (
	"fmt"
	"reflect"
	"strings"
)

func GenConv(from, to interface{}) string {
	tFrom := reflect.TypeOf(from)
	tTo := reflect.TypeOf(to)
	sb := &strings.Builder{}
	sb.WriteString("\n")
	genConvType(tFrom, tTo, sb, map[string]bool{})
	return sb.String()
}

func genConvType(from, to reflect.Type, sb *strings.Builder, existed map[string]bool) {
	fromName, toName := from.String(), to.String()
	key := fmt.Sprintf("%s => %s", fromName, toName)
	if existed[key] {
		return
	}
	existed[key] = true
	var waitConv [][2]reflect.Type
	if to.Kind() == reflect.Struct || (to.Kind() == reflect.Ptr && to.Elem().Kind() == reflect.Struct) {
		fromFields := map[string]reflect.StructField{}
		fromMethods := map[string]reflect.Method{}
		fElem := from
		for i := 0; i < fElem.NumMethod(); i++ {
			method := fElem.Method(i)
			fromMethods[method.Name] = method
		}
		if fElem.Kind() == reflect.Struct {
			fPtr := reflect.New(fElem).Type()
			for i := 0; i < fPtr.NumMethod(); i++ {
				method := fPtr.Method(i)
				fromMethods[method.Name] = method
			}
		} else if fElem.Kind() == reflect.Ptr {
			fElem = fElem.Elem()
			for i := 0; i < fElem.NumMethod(); i++ {
				method := fElem.Method(i)
				fromMethods[method.Name] = method
			}
		}
		for i := 0; i < fElem.NumField(); i++ {
			field := fElem.Field(i)
			fromFields[field.Name] = field
		}
		sb.WriteString(fmt.Sprintf("func %s2%s(in %s) %s {\n", formatTypeName(fromName), formatTypeName(toName), fromName, toName))
		sb.WriteString(fmt.Sprintf("\treturn %s{\n", strings.ReplaceAll(toName, "*", "&")))
		tElem := to
		if tElem.Kind() == reflect.Ptr {
			tElem = tElem.Elem()
		}
		for i := 0; i < tElem.NumField(); i++ {

			field := tElem.Field(i)
			fromField, okField := fromFields[field.Name]
			fromMethod, okMethod := fromMethods["Get"+field.Name]
			if !(okField || okMethod) {
				continue
			}
			toFieldType := field.Type
			if toFieldType.Kind() == reflect.Struct ||
				(toFieldType.Kind() == reflect.Ptr && toFieldType.Elem().Kind() == reflect.Struct) {
				if okMethod {
					fromMethodOutType := fromMethod.Type.Out(0)
					waitConv = append(waitConv, [2]reflect.Type{fromMethodOutType, toFieldType})
					sb.WriteString(fmt.Sprintf("\t\t%s: %s2%s(in.%s()),\n",
						field.Name, formatTypeName(fromMethodOutType.String()), formatTypeName(toFieldType.String()), fromMethod.Name))
				} else {
					fromFieldType := fromField.Type
					waitConv = append(waitConv, [2]reflect.Type{fromFieldType, toFieldType})
					sb.WriteString(fmt.Sprintf("\t\t%s: %s2%s(in.%s),\n",
						field.Name, formatTypeName(fromFieldType.String()), formatTypeName(toFieldType.String()), fromField.Name))
				}
			} else if toFieldType.Kind() == reflect.Ptr {
				if okMethod {
					fromMethodOutType := fromMethod.Type.Out(0)
					if fromMethodOutType.Kind() != toFieldType.Elem().Kind() {
						sb.WriteString(fmt.Sprintf("\t\t%s: util.GetPtr(%s(in.%s())),\n",
							field.Name, toFieldType.Elem().String(), fromMethod.Name))
					} else {
						sb.WriteString(fmt.Sprintf("\t\t%s: util.GetPtr(in.%s()),\n",
							field.Name, fromMethod.Name))
					}
				} else if fromField.Type.Kind() == reflect.Ptr {
					if fromField.Type.Elem().Kind() != toFieldType.Elem().Kind() {
						sb.WriteString(fmt.Sprintf("\t\t%s: util.GetPtr(%s(util.GetVal(in.%s))),\n",
							field.Name, toFieldType.Elem().String(), fromField.Name))
					} else {
						sb.WriteString(fmt.Sprintf("\t\t%s: util.CopyPtr(in.%s),\n",
							field.Name, fromField.Name))
					}
				} else {
					if fromField.Type.Kind() != toFieldType.Elem().Kind() {
						sb.WriteString(fmt.Sprintf("\t\t%s: util.GetPtr(%s(in.%s)),\n",
							field.Name, toFieldType.Elem().String(), fromField.Name))
					} else {
						sb.WriteString(fmt.Sprintf("\t\t%s: util.GetPtr(in.%s),\n",
							field.Name, fromField.Name))
					}
				}
			} else if toFieldType.Kind() == reflect.String ||
				(toFieldType.Kind() >= reflect.Bool && toFieldType.Kind() <= reflect.Complex128) {
				if okMethod {
					fromMethodOutType := fromMethod.Type.Out(0)
					if fromMethodOutType.Kind() != toFieldType.Kind() {
						sb.WriteString(fmt.Sprintf("\t\t%s: %s(in.%s()),\n",
							field.Name, toFieldType.String(), fromMethod.Name))
					} else {
						sb.WriteString(fmt.Sprintf("\t\t%s: in.%s(),\n",
							field.Name, fromMethod.Name))
					}
				} else if fromField.Type.Kind() == reflect.Ptr {
					if fromField.Type.Elem().Kind() != toFieldType.Kind() {
						sb.WriteString(fmt.Sprintf("\t\t%s: %s(util.GetVal(in.%s)),\n",
							field.Name, toFieldType.String(), fromField.Name))
					} else {
						sb.WriteString(fmt.Sprintf("\t\t%s: util.GetVal(in.%s),\n",
							field.Name, fromField.Name))
					}
				} else {
					if fromField.Type.Kind() != toFieldType.Kind() {
						sb.WriteString(fmt.Sprintf("\t\t%s: %s(in.%s),\n",
							field.Name, toFieldType.String(), fromField.Name))
					} else {
						sb.WriteString(fmt.Sprintf("\t\t%s: in.%s,\n",
							field.Name, fromField.Name))
					}
				}
			} else if toFieldType.Kind() == reflect.Slice {
				if okMethod {
					fromMethodOutType := fromMethod.Type.Out(0)
					waitConv = append(waitConv, [2]reflect.Type{fromMethodOutType, toFieldType})
					sb.WriteString(fmt.Sprintf("\t\t%s: %s2%s(in.%s()),\n",
						field.Name, formatTypeName(fromMethodOutType.String()), formatTypeName(toFieldType.String()), fromMethod.Name))
				} else {
					fromFieldType := fromField.Type
					waitConv = append(waitConv, [2]reflect.Type{fromFieldType, toFieldType})
					sb.WriteString(fmt.Sprintf("\t\t%s: %s2%s(in.%s),\n",
						field.Name, formatTypeName(fromFieldType.String()), formatTypeName(toFieldType.String()), fromField.Name))
				}
			}
		}
		sb.WriteString("\t}\n}\n\n")
	} else if to.Kind() == reflect.Slice {
		sb.WriteString(fmt.Sprintf("func %s2%s(in %s) %s {\n",
			formatTypeName(fromName), formatTypeName(toName), fromName, toName))
		sb.WriteString(fmt.Sprintf("\tout := make(%s, len(in))\n", toName))
		sb.WriteString("\tfor i := range in {\n")
		inElem := from.Elem()
		outElem := to.Elem()

		if outElem.Kind() == reflect.Struct ||
			(outElem.Kind() == reflect.Ptr && outElem.Elem().Kind() == reflect.Struct) {
			waitConv = append(waitConv, [2]reflect.Type{inElem, outElem})
			sb.WriteString(fmt.Sprintf("\t\tout[i] = %s2%s(in[i])\n",
				formatTypeName(inElem.String()), formatTypeName(outElem.String())))
		} else if outElem.Kind() == reflect.Ptr {
			if inElem.Kind() == reflect.Ptr {
				if inElem.Elem().Kind() != outElem.Elem().Kind() {
					sb.WriteString(fmt.Sprintf("\t\tout[i] = util.GetPtr(%s(util.GetVal(in[i])))\n",
						outElem.Elem().String()))
				} else {
					sb.WriteString("\t\tout[i] = util.CopyPtr(in)\n")
				}
			} else {
				if inElem.Kind() != outElem.Elem().Kind() {
					sb.WriteString(fmt.Sprintf("\t\tout[i] = util.GetPtr(%s(in[i]))\n",
						outElem.Elem().String()))
				} else {
					sb.WriteString("\t\tout[i] = util.GetPtr(in)\n")
				}
			}
		} else if outElem.Kind() == reflect.String ||
			(outElem.Kind() >= reflect.Bool && outElem.Kind() <= reflect.Complex128) {
			if inElem.Kind() == reflect.Ptr {
				if inElem.Elem().Kind() != outElem.Kind() {
					sb.WriteString(fmt.Sprintf("\t\tout[i] = %s(util.GetVal(in[i]))\n",
						outElem.String()))
				} else {
					sb.WriteString("\t\tout[i] = util.GetVal(in[i])\n")
				}
			} else {
				if inElem.Kind() != outElem.Kind() {
					sb.WriteString(fmt.Sprintf("\t\tout[i] = %s(in[i])\n",
						outElem.String()))
				} else {
					sb.WriteString("\t\tout[i] = in[i]\n")
				}
			}
		} else if outElem.Kind() == reflect.Slice {
			waitConv = append(waitConv, [2]reflect.Type{inElem, outElem})
			sb.WriteString(fmt.Sprintf("\t\tout[i] = %s2%s(in[i])\n",
				formatTypeName(inElem.String()), formatTypeName(outElem.String())))
		}

		sb.WriteString("\t}\n\treturn out\n}\n\n")
	}
	for _, nextConv := range waitConv {
		genConvType(nextConv[0], nextConv[1], sb, existed)
	}
}

func formatTypeName(s string) string {
	s = strings.ReplaceAll(s, ".", "__")
	s = strings.ReplaceAll(s, "[]", "s_")
	return strings.ReplaceAll(s, "*", "p_")
}
