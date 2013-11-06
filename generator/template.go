package generator

import (
	"go/ast"
)

// getTypeSpecs retrieves all struct TypeSpec objects in a File.
func getTypeSpecs(f *ast.File) []*ast.TypeSpec {
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

// getStructFields retrieves all fields from a TypeSpec.
func getStructFields(spec *ast.TypeSpec) []*ast.Field {
	s := make([]*ast.Field, 0)
	if structType, ok := spec.Type.(*ast.StructType); ok {
		for _, field := range structType.Fields.List {
			s = append(s, field)
		}
	}
	return s
}

// isType returns true if the field is a given type.
func isType(field *ast.Field, typ string) bool {
	if ident, ok := field.Type.(*ast.Ident); ok {
		return ident.Name == typ
	}
	return false
}

// getFieldName returns the first name in a field.
func getFieldName(field *ast.Field) string {
	return field.Names[0].Name
}
