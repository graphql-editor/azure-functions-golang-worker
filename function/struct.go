package function

import (
	"reflect"
	"strings"
	"sync"
)

type field struct {
	typ       reflect.Type
	tagged    bool
	name      string
	omitEmpty bool
	asString  bool
	index     []int
}

func getTag(field *reflect.StructField) string {
	tag := field.Tag.Get("azfunc")
	if tag == "" {
		return ""
	}
	return strings.Split(tag, ",")[0]
}

var fieldCache sync.Map

func typeFields(t reflect.Type) []field {
	current := []field{}
	next := []field{{typ: t}}
	visited := map[reflect.Type]bool{}
	fieldAt := map[string]int{}
	orphans := []int{}
	var fields []field
	var level int
	for len(next) > 0 {
		level++
		current, next = next, current[:0]
		nextCount := map[reflect.Type]bool{}
		for _, f := range current {
			if visited[f.typ] {
				continue
			}
			visited[f.typ] = true
			for i := 0; i < f.typ.NumField(); i++ {
				sf := f.typ.Field(i)
				isUnexported := sf.PkgPath != ""
				if sf.Anonymous {
					t := sf.Type
					if t.Kind() == reflect.Ptr {
						t = t.Elem()
					}
					if isUnexported && t.Kind() != reflect.Struct {
						continue
					}
				} else if isUnexported {
					continue
				}
				tag := getTag(&sf)
				if tag == "-" {
					continue
				}
				index := make([]int, len(f.index)+1)
				copy(index, f.index)
				index[len(f.index)] = i
				name := tag
				ft := sf.Type
				if ft.Name() == "" && ft.Kind() == reflect.Ptr {
					ft = ft.Elem()
				}
				if name != "" || !sf.Anonymous || ft.Kind() != reflect.Struct {
					tagged := name != ""
					if name == "" {
						name = sf.Name
					}
					if fAt, ok := fieldAt[name]; ok {
						if level > len(fields[fAt].index) {
							continue
						}
						if fields[fAt].tagged || (!tagged && !fields[fAt].tagged) {
							continue
						}
						orphans = append(orphans, fAt)
					}
					fieldAt[name] = len(fields)
					fields = append(fields, field{
						typ:    ft,
						tagged: tagged,
						name:   name,
						index:  index,
					})
					continue
				}
				if !nextCount[ft] {
					nextCount[ft] = true
					next = append(next, field{index: index, typ: ft})
				}
			}
		}
	}
	for i, orphan := range orphans {
		fields = append(fields[:orphan-i], fields[orphan-i+1:]...)
	}
	for i := range fields {
		typ := t
		for _, i := range fields[i].index {
			if typ.Kind() == reflect.Ptr {
				typ = typ.Elem()
			}
			typ = typ.Field(i).Type
		}
	}
	return fields
}

func cachedTypeFields(t reflect.Type) []field {
	if f, ok := fieldCache.Load(t); ok {
		return f.([]field)
	}
	f, _ := fieldCache.LoadOrStore(t, typeFields(t))
	return f.([]field)
}
