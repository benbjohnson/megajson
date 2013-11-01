package generator

import (
	"bytes"
	"errors"
	"fmt"
	"go/ast"
	"io"
)

// Used for marking fields as not written. Not actually returned to user.
var unsupportedTypeError = errors.New("Unsupported type")

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
	fmt.Fprintf(&b, "if err := encoding.WriteByte(e.w, '{'); err != nil {\nreturn err\n}\n")

	index := 0
	for _, f := range structType.Fields.List {
		for _, name := range f.Names {
			if index > 0 {
				fmt.Fprintf(&b, "if err := encoding.WriteByte(e.w, ','); err != nil {\nreturn err\n}\n")
			}

			if err := writeFieldEncoding(&b, name.Name, f); err == nil {
				index++
			} else if err != unsupportedTypeError {
				return err
			}
		}
	}
	fmt.Fprintf(&b, "if err := encoding.WriteByte(e.w, '}'); err != nil {\nreturn err\n}\n")
	fmt.Fprintf(&b, "return nil\n")
	fmt.Fprintf(&b, "}\n")

	// Write to formatted output stream.
	if _, err := b.WriteTo(w); err != nil {
		return err
	}

	return nil
}

// writeFieldEncoding generates the encoder code for a single field.
func writeFieldEncoding(w io.Writer, name string, f *ast.Field) error {
	typ, ok := f.Type.(*ast.Ident)
	if !ok {
		return nil
	}

	fmt.Fprintf(w, "if err := encoding.WriteString(e.w, \"%s\"); err != nil {\nreturn err\n}\n", name)
	fmt.Fprintf(w, "if err := encoding.WriteByte(e.w, ':'); err != nil {\nreturn err\n}\n")

	switch typ.Name {
	case "string":
		fmt.Fprintf(w, "if err := encoding.WriteString(e.w, v.%s); err != nil {\nreturn err\n}\n", name)
	case "int":
		fmt.Fprintf(w, "if err := encoding.WriteInt(e.w, v.%s); err != nil {\nreturn err\n}\n", name)
	case "uint":
		fmt.Fprintf(w, "if err := encoding.WriteUint(e.w, v.%s); err != nil {\nreturn err\n}\n", name)
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
