package generator

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"io"
	"reflect"
	"strconv"
	"strings"
)

func writeFileEncoder(w io.Writer, file *ast.File) error {
	// Write header.
	var b bytes.Buffer
	fmt.Fprintf(&b, "package %s\n", file.Name.Name)
	fmt.Fprintf(&b, "import (\n")
	fmt.Fprintf(&b, "\"io\"\n")
	fmt.Fprintf(&b, "\"github.com/benbjohnson/megajson/encoding\"\n")
	fmt.Fprintf(&b, ")\n")

	// Loop over each spec and create a encoder.
	generated := false
	for _, decl := range file.Decls {
		if decl, ok := decl.(*ast.GenDecl); ok {
			for _, spec := range decl.Specs {
				if spec, ok := spec.(*ast.TypeSpec); ok {
					err := writeTypeEncoder(&b, spec)
					if err != nil {
						return err
					}
					generated = true
				}
			}
		}
	}

	// If no types were found to encode then skip this file.
	if !generated {
		return nil
	}

	// fmt.Println("+++++\n", b.String(), "\n+++++")

	// Format source.
	bfmt, err := format.Source(b.Bytes())
	if err != nil {
		return err
	}

	fmt.Println("+++++\n", string(bfmt), "\n+++++")

	// Write to output file.
	if _, err := w.Write(bfmt); err != nil {
		return err
	}

	return nil
}

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
		tag = f.Tag.Value[1 : len(f.Tag.Value)-1]
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
	var err error
	var b bytes.Buffer
	switch typ := f.Type.(type) {
	case *ast.Ident:
		err = writePrimativeFieldEncoding(&b, fmt.Sprintf("v.%s", name), typ, tag)
	case *ast.StarExpr:
		err = writePointerFieldEncoding(&b, fmt.Sprintf("v.%s", name), typ, tag)
	case *ast.ArrayType:
		err = writeArrayFieldEncoding(&b, fmt.Sprintf("v.%s", name), typ, tag)
	default:
		return unsupportedTypeError
	}

	// Exit if the field throws an error or is ignored.
	if err != nil {
		return err
	}

	fmt.Fprintf(w, "if err := encoding.WriteString(e.w, %s); err != nil {\nreturn err\n}\n", strconv.Quote(key))
	fmt.Fprintf(w, "if err := encoding.WriteByte(e.w, ':'); err != nil {\nreturn err\n}\n")
	b.WriteTo(w)

	return nil
}

// writeFieldEncoding generates the encoder code for a single primative field.
func writePrimativeFieldEncoding(w io.Writer, varname string, typ *ast.Ident, tag string) error {
	switch typ.Name {
	case "string":
		fmt.Fprintf(w, "if err := encoding.WriteString(e.w, %s); err != nil {\nreturn err\n}\n", varname)
	case "int":
		fmt.Fprintf(w, "if err := encoding.WriteInt(e.w, %s); err != nil {\nreturn err\n}\n", varname)
	case "int64":
		fmt.Fprintf(w, "if err := encoding.WriteInt64(e.w, %s); err != nil {\nreturn err\n}\n", varname)
	case "uint":
		fmt.Fprintf(w, "if err := encoding.WriteUint(e.w, %s); err != nil {\nreturn err\n}\n", varname)
	case "uint64":
		fmt.Fprintf(w, "if err := encoding.WriteUint64(e.w, %s); err != nil {\nreturn err\n}\n", varname)
	case "float32":
		fmt.Fprintf(w, "if err := encoding.WriteFloat32(e.w, %s); err != nil {\nreturn err\n}\n", varname)
	case "float64":
		fmt.Fprintf(w, "if err := encoding.WriteFloat64(e.w, %s); err != nil {\nreturn err\n}\n", varname)
	case "bool":
		fmt.Fprintf(w, "if err := encoding.WriteBool(e.w, %s); err != nil {\nreturn err\n}\n", varname)
	default:
		return unsupportedTypeError
	}

	return nil
}

// writePointerFieldEncoding generates the encoding code for a single field with a pointer type.
func writePointerFieldEncoding(w io.Writer, varname string, typ *ast.StarExpr, tag string) error {
	x, ok := typ.X.(*ast.Ident)
	if !ok {
		return unsupportedTypeError
	}

	fmt.Fprintf(w, "if err := New%sJSONEncoder(e.w).Encode(%s); err != nil {\nreturn err\n}\n", x.Name, varname)
	return nil
}

// writeArrayFieldEncoding generates the encoding code for a single field with a array type.
func writeArrayFieldEncoding(w io.Writer, varname string, typ *ast.ArrayType, tag string) error {
	var b bytes.Buffer

	fmt.Fprintf(&b, "if err := encoding.WriteByte(e.w, '['); err != nil {\nreturn err\n}\n")

	fmt.Fprintf(&b, "for index, elem := range %s {", varname)
	fmt.Fprintf(&b, "if index > 0 { if err := encoding.WriteByte(e.w, ','); err != nil { return err } } \n")

	switch elt := typ.Elt.(type) {
	case *ast.Ident:
		if err := writePrimativeFieldEncoding(&b, "elem", elt, ""); err != nil {
			return err
		}

	case *ast.StarExpr:
		if err := writePointerFieldEncoding(&b, "elem", elt, ""); err != nil {
			return err
		}

	default:
		return unsupportedTypeError
	}

	fmt.Fprintf(&b, "}\n")
	fmt.Fprintf(&b, "if err := encoding.WriteByte(e.w, ']'); err != nil {\nreturn err\n}\n")

	_, err := b.WriteTo(w)
	return err
}
