package generator

import (
	"bytes"
	"errors"
	"fmt"
	"go/ast"
	"io"
	"reflect"
	"strings"
)

// Used for marking fields as not written. Not actually returned to user.
var unsupportedTypeError = errors.New("Unsupported type")
var ignoreFieldError = errors.New("Ignore field")

// GenerateTypeEncoder generates a struct for a single Type.
func writeTypeEncoder(w io.Writer, typeSpec *ast.TypeSpec) error {
	var b bytes.Buffer
	name := typeSpec.Name.Name

	// Ignore type non-struct specs.
	structType, ok := typeSpec.Type.(*ast.StructType)
	if !ok {
		return nil
	}

	// Generate the encoder.
	fmt.Fprintf(&b, "type %sJSONEncoder struct {", name)
	fmt.Fprintln(&b, "w io.Writer")
	fmt.Fprintln(&b, "}")

	// Generate the constructor.
	fmt.Fprintf(&b, "func New%sJSONEncoder(w io.Writer) *%sJSONEncoder {\n", name, name)
	fmt.Fprintf(&b, "return &%sJSONEncoder{w: w}\n", name)
	fmt.Fprintln(&b, "}")
	fmt.Fprintln(&b, "")

	// Generate the encode function.
	fmt.Fprintf(&b, "func (e *%sJSONEncoder) Encode(v *%s) error {\n", name, name)
	fmt.Fprintf(&b, "if v == nil {\nreturn encoding.WriteBytes(e.w, []byte(`null`))\n}\n\n")

	fmt.Fprintf(&b, "if err := encoding.WriteByte(e.w, '{'); err != nil {\nreturn err\n}\n")

	index := 0
	for _, field := range structType.Fields.List {
		for _, name := range field.Names {
			// Write to a temporary buffer to check if anything is written.
			var buf bytes.Buffer
			if err := writeFieldEncoding(&buf, name.Name, field); buf.Len() > 0 {
				// Write separating comma after the first field.
				if index > 0 {
					fmt.Fprintf(&b, "if err := encoding.WriteByte(e.w, ','); err != nil {\nreturn err\n}\n")
				}

				// Copy over to main buffer.
				buf.WriteTo(&b)
				index++

			} else if err != nil && err != unsupportedTypeError && err != ignoreFieldError {
				return err
			}
		}
	}

	fmt.Fprintf(&b, "if err := encoding.WriteByte(e.w, '}'); err != nil {\nreturn err\n}\n")
	fmt.Fprintf(&b, "return nil\n")
	fmt.Fprintf(&b, "}\n")

	// Copy buffer to writer.
	if _, err := b.WriteTo(w); err != nil {
		return err
	}

	return nil
}

// writeFieldEncoding generates the encoder code for a single field.
func writeFieldEncoding(w io.Writer, name string, f *ast.Field) error {
	var tag string
	if f.Tag != nil {
		tag = f.Tag.Value[1:len(f.Tag.Value)-1]
		tag = reflect.StructTag(tag).Get("json")
	}
	tags := strings.Split(tag, ",")

	// Extract json key name.
	var key = name
	if len(tags) > 0 {
		if tags[0] == "-" {
			return ignoreFieldError
		} else if len(tags[0]) > 0 {
			key = tags[0]
		}
	}

	// Primative and complex types need separate handling.
	switch f.Type.(type) {
	case *ast.Ident:
		return writePrimativeFieldEncoding(w, name, key, f, tag)
	case *ast.StarExpr:
		return writePointerFieldEncoding(w, name, key, f, tag)
	case *ast.ArrayType:
		fmt.Println("ARRAY!")
	}
	return unsupportedTypeError
}

// writeFieldEncoding generates the encoder code for a single primative field.
func writePrimativeFieldEncoding(w io.Writer, name string, key string, f *ast.Field, tag string) error {
	typ := f.Type.(*ast.Ident)

	fmt.Fprintf(w, "if err := encoding.WriteString(e.w, \"%s\"); err != nil {\nreturn err\n}\n", key)
	fmt.Fprintf(w, "if err := encoding.WriteByte(e.w, ':'); err != nil {\nreturn err\n}\n")

	switch typ.Name {
	case "string":
		fmt.Fprintf(w, "if err := encoding.WriteString(e.w, v.%s); err != nil {\nreturn err\n}\n", name)
	case "int":
		fmt.Fprintf(w, "if err := encoding.WriteInt(e.w, v.%s); err != nil {\nreturn err\n}\n", name)
	case "int64":
		fmt.Fprintf(w, "if err := encoding.WriteInt64(e.w, v.%s); err != nil {\nreturn err\n}\n", name)
	case "uint":
		fmt.Fprintf(w, "if err := encoding.WriteUint(e.w, v.%s); err != nil {\nreturn err\n}\n", name)
	case "uint64":
		fmt.Fprintf(w, "if err := encoding.WriteUint64(e.w, v.%s); err != nil {\nreturn err\n}\n", name)
	case "float32":
		fmt.Fprintf(w, "if err := encoding.WriteFloat32(e.w, v.%s); err != nil {\nreturn err\n}\n", name)
	case "float64":
		fmt.Fprintf(w, "if err := encoding.WriteFloat64(e.w, v.%s); err != nil {\nreturn err\n}\n", name)
	case "bool":
		fmt.Fprintf(w, "if err := encoding.WriteBool(e.w, v.%s); err != nil {\nreturn err\n}\n", name)
	default:
		return unsupportedTypeError
	}

	return nil
}

// writePointerFieldEncoding generates the encoding code for a single field with a pointer type.
func writePointerFieldEncoding(w io.Writer, name string, key string, f *ast.Field, tag string) error {
	typ := f.Type.(*ast.StarExpr)
	x, ok := typ.X.(*ast.Ident)
	if !ok {
		return unsupportedTypeError
	}

	fmt.Fprintf(w, "if err := encoding.WriteString(e.w, \"%s\"); err != nil {\nreturn err\n}\n", name)
	fmt.Fprintf(w, "if err := encoding.WriteByte(e.w, ':'); err != nil {\nreturn err\n}\n")
	fmt.Fprintf(w, "if err := New%sJSONEncoder(e.w).Encode(v.%s); err != nil {\nreturn err\n}\n", x.Name, name)
	return nil
}
