package generator

const templateModel = `//lint:file-ignore U1000 ignore unused code, it's generated
package {{.Package}}{{if .HasImports}}

import ({{range .Imports}}
    "{{.}}"{{end}}
){{end}}

var Columns = struct { {{range .Models}}
	{{.StructName}} struct{ 
		{{range $i, $e := .Columns}}{{if $i}}, {{end}}{{.FieldName}}{{end}} string{{if .HasRelations}}

		{{range $i, $e := .Relations}}{{if $i}}, {{end}}{{.FieldName}}{{end}} string{{end}}
	}{{end}}
}{ {{range .Models}}
	{{.StructName}}: struct { 
		{{range $i, $e := .Columns}}{{if $i}}, {{end}}{{.FieldName}}{{end}} string{{if .HasRelations}}

		{{range $i, $e := .Relations}}{{if $i}}, {{end}}{{.FieldName}}{{end}} string{{end}}
	}{ {{range .Columns}}
		{{.FieldName}}: "{{.FieldDBName}}",{{end}}{{if .HasRelations}}
		{{range .Relations}}
		{{.FieldName}}: "{{.FieldName}}",{{end}}{{end}}
	},{{end}}
}

var Tables = struct { {{range .Models}}
	{{.StructName}} struct {
		Name string
	}{{end}}
}{ {{range .Models}}
	{{.StructName}}: struct {
		Name string
	}{ 
		Name: "{{.TableName}}",
	},{{end}}
}
{{range .Models}}
type {{.StructName}} struct {
	tableName struct{} {{.StructTag}}
	{{range .Columns}}
	{{.FieldName}} {{.FieldType}} {{.FieldTag}} {{.FieldComment}}{{end}}{{if .HasRelations}}
	{{range .Relations}}
	{{.FieldName}} {{.FieldType}} {{.FieldTag}} {{.FieldComment}}{{end}}{{end}}
}
{{end}}
`
