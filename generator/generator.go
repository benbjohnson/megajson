package generator

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"io"
)

// Generates a parser for a single Go AST type.
func GenerateEncoder(pkg string, typeSpec *ast.TypeSpec, w io.Writer) error {
	var b bytes.Buffer
	name := typeSpec.Name.Name

	// Ignore type non-struct specs.
	typ, ok := typeSpec.Type.(*ast.StructType)
	if !ok {
		return nil
	}

	// Retrieve all relevant fields.
	fields := make([]*ast.Field, 0)
	for _, f := range typ.Fields.List {
		if ident, ok := f.Type.(*ast.Ident); ok {
			switch ident.Name {
			case "string", "int":
				fields = append(fields, f)
			}
		}
	}

	// Generate the encoder.
	fmt.Fprintf(&b, "package %s\n", pkg)
	fmt.Fprintln(&b, "import \"io\"")
	fmt.Fprintln(&b, "import \"github.com/benbjohnson/megajson/encoding\"")
	fmt.Fprintf(&b, "type %sJSONEncoder struct {", name)
	fmt.Fprintln(&b, "w io.Writer")
	fmt.Fprintln(&b, "}")

	fmt.Fprintf(&b, "func (e *%sJSONEncoder) Encode(v *%s) error {\n", name, name)
	for _, f := range fields {
		if ident, ok := f.Type.(*ast.Ident); ok {
			name := f.Names[0]

			switch ident.Name {
            case "string":
                fmt.Fprintf(&b, "\tif _, err := encoding.WriteString(e.w, v.%s); err != nil {\n\t\treturn err\n\t}\n", name)
            case "int":
                fmt.Fprintf(&b, "\tif _, err := encoding.WriteInt(e.w, v.%s); err != nil {\n\t\treturn err\n\t}\n", name)
            case "uint":
                fmt.Fprintf(&b, "\tif _, err := encoding.WriteUint(e.w, v.%s); err != nil {\n\t\treturn err\n\t}\n", name)
            case "bool":
                fmt.Fprintf(&b, "\tif _, err := encoding.WriteBool(e.w, v.%s); err != nil {\n\t\treturn err\n\t}\n", name)
			}
		}
	}
	fmt.Fprintf(&b, "\treturn nil\n")
	fmt.Fprintf(&b, "}\n")

	// Debugging
	fmt.Println(">>>>>>>>>>> DEBUG >>>>>>>>>>>")
	fmt.Println(b.String())
	fmt.Println("<<<<<<<<<<< DEBUG <<<<<<<<<<<")
	fmt.Println("")

	// Format source.
	bfmt, err := format.Source(b.Bytes());
	if err != nil {
		fmt.Println("ERR: ", err)
		return err
	}

	// Write to output stream.
	if _, err := w.Write(bfmt); err != nil {
		return err
	}

	return nil
}
