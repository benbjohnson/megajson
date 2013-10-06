package generator

import (
    "bytes"
    "fmt"
    "go/ast"
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
            case "string":
                fields = append(fields, f)
            }
        }
    }

    // Generate the encoder.
    fmt.Fprintf(&b, "package %s\n\n", pkg)
    fmt.Fprintf(&b, "import \"io\"\n\n")
    fmt.Fprintf(&b, "type %sJSONEncoder {\n", name)
    fmt.Fprintf(&b, "\tw io.Writer\n")
    fmt.Fprintf(&b, "}\n\n")

    fmt.Fprintf(&b, "func (e *%sJSONEncoder) Encode(v *%s) error {\n", name, name)
    for _, f := range fields {
        if ident, ok := f.Type.(*ast.Ident); ok {
            switch ident.Name {
            case "string":
                fmt.Fprintf(&b, "\tif _, err := io.WriteString(e.w, v.%s); err != nil {\n\t\treturn err\n\t}\n", f.Names[0])
            }
        }
    }
    fmt.Fprintf(&b, "\treturn nil\n")
    fmt.Fprintf(&b, "}\n")

    // Debugging
    fmt.Println(b.String())

    // Write to output stream.
    if _, err := w.Write(b.Bytes()); err != nil {
        return err
    }

    return nil
}
