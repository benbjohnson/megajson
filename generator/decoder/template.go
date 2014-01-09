package decoder

import (
	"go/ast"
	"reflect"
	"strings"
	"text/template"
)

var tmpl *template.Template

func init() {
	tmpl = template.Must(template.New("decoder.tmpl").Funcs(template.FuncMap{
		"types":           types,
		"fields":          fields,
		"istype":          istype,
		"isprimitivetype": isprimitivetype,
		"subtype":         subtype,
		"fieldname":       fieldname,
		"keyname":         keyname,
	}).Parse(string(tmplsrc())))
}

// types retrieves all a list of all available struct type specs in a file.
func types(f *ast.File) []*ast.TypeSpec {
	s := make([]*ast.TypeSpec, 0)
	for _, decl := range f.Decls {
		if decl, ok := decl.(*ast.GenDecl); ok {
			for _, spec := range decl.Specs {
				if spec, ok := spec.(*ast.TypeSpec); ok {
					if _, ok := spec.Type.(*ast.StructType); ok {
						s = append(s, spec)
					}
				}
			}
		}
	}
	return s
}

// fields retrieves all fields from a struct type spec.
func fields(spec *ast.TypeSpec) []*ast.Field {
	s := make([]*ast.Field, 0)
	if structType, ok := spec.Type.(*ast.StructType); ok {
		for _, field := range structType.Fields.List {
			if keyname(field) != "" {
				s = append(s, field)
			}
		}
	}
	return s
}

// getType returns the name of the type of the field.
func getType(field *ast.Field) string {
	if ident, ok := field.Type.(*ast.Ident); ok {
		return ident.Name
	} else if _, ok := field.Type.(*ast.StarExpr); ok {
		return "*"
	} else if _, ok := field.Type.(*ast.ArrayType); ok {
		return "[]"
	}
	return ""
}

// istype returns true if the field is a given type.
func istype(field *ast.Field, typ string) bool {
	return getType(field) == typ
}

// isprimitivetype returns true if the field is a primitive type.
func isprimitivetype(field *ast.Field) bool {
	switch getType(field) {
	case "string", "int", "int64", "uint", "uint64", "float32", "float64", "bool":
		return true
	}
	return false
}

// subtype returns the subtype of a pointer or array.
func subtype(field *ast.Field) string {
	if typ, ok := field.Type.(*ast.StarExpr); ok {
		if ident, ok := typ.X.(*ast.Ident); ok {
			return ident.Name
		}
	} else if typ, ok := field.Type.(*ast.ArrayType); ok {
		if typ, ok := typ.Elt.(*ast.StarExpr); ok {
			if ident, ok := typ.X.(*ast.Ident); ok {
				return ident.Name
			}
		}
		if ident, ok := typ.Elt.(*ast.Ident); ok {
			return ident.Name
		}
	}
	return ""
}

// fieldname returns the first name in a field.
func fieldname(field *ast.Field) string {
	return field.Names[0].Name
}

// keyname returns the JSON key to be used for a field.
func keyname(field *ast.Field) string {
	tags := tags(field)

	if len(tags) > 0 {
		if len(tags[0]) == 0 {
			return fieldname(field)
		} else if tags[0] == "-" {
			return ""
		} else {
			return tags[0]
		}
	} else {
		return fieldname(field)
	}
}

// tags returns the JSON tags on a field.
func tags(field *ast.Field) []string {
	var tag string
	if field.Tag != nil {
		tag = field.Tag.Value[1 : len(field.Tag.Value)-1]
		tag = reflect.StructTag(tag).Get("json")
	}
	return strings.Split(tag, ",")
}
