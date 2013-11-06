package generator

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"io"
	"text/template"
)

var decoderTemplate *template.Template

func init() {
	m := template.FuncMap{
		"types": getTypeSpecs,
		"fields": getStructFields,
		"istype": isType,
		"subtype": getSubType,
		"fieldname": getFieldName,
		"keyname": getJSONKeyName,
	}
	decoderTemplate = template.Must(template.New("decoder").Funcs(m).Parse(string(decoder_tmpl())))
}

func writeFileDecoder(w io.Writer, file *ast.File) error {
	// Ignore files without type specs.
	if types := getTypeSpecs(file); len(types) == 0 {
		return nil
	}

	var b bytes.Buffer
	if err := decoderTemplate.Execute(&b, file); err != nil {
		return err
	}

	bfmt, err := format.Source(b.Bytes())
	if err != nil {
		fmt.Println("-----\n", b.String(), "\n-----")
		return err
	}

	fmt.Println("-----\n", string(bfmt), "\n-----")

	if _, err := w.Write(bfmt); err != nil {
		return err
	}

	return nil
}
