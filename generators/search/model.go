package search

import (
	"fmt"
	"html/template"

	"github.com/dizzyfool/genna/model"
	"github.com/dizzyfool/genna/util"
)

// TemplatePackage stores package info
type TemplatePackage struct {
	Package string

	HasImports bool
	Imports    []string

	GoPGVer string

	Entities []TemplateEntity
}

// NewTemplatePackage creates a package for template
func NewTemplatePackage(entities []model.Entity, options Options) TemplatePackage {

	imports := util.NewSet()

	var models []TemplateEntity
	for _, entity := range entities {
		mdl := NewTemplateEntity(entity, options)
		if len(mdl.Columns) == 0 {
			continue
		}

		for _, imp := range mdl.Imports {
			imports.Add(imp)
		}

		for _, col := range mdl.Columns {
			if col.Relaxed {
				imports.Add("reflect")
			}
		}

		models = append(models, mdl)
	}

	goPGVer := ""
	if options.GoPgVer >= 9 {
		goPGVer = fmt.Sprintf("/v%d", options.GoPgVer)
	}

	return TemplatePackage{
		Package: options.Package,

		HasImports: imports.Len() > 0,
		Imports:    imports.Elements(),

		GoPGVer: goPGVer,

		Entities: models,
	}
}

// TemplateEntity stores struct info
type TemplateEntity struct {
	model.Entity

	NoAlias bool
	Alias   string

	Columns []TemplateColumn

	Imports []string
}

// NewTemplateEntity creates an entity for template
func NewTemplateEntity(entity model.Entity, options Options) TemplateEntity {
	if entity.HasMultiplePKs() {
		options.KeepPK = true
	}

	imports := util.NewSet()

	var columns []TemplateColumn
	for _, column := range entity.Columns {
		if column.IsArray || column.GoType == model.TypeMapInterface || column.GoType == model.TypeMapString {
			continue
		}

		columns = append(columns, NewTemplateColumn(column, options))
		if column.Import != "" {
			imports.Add(column.Import)
		}
	}

	return TemplateEntity{
		Entity: entity,

		NoAlias: options.NoAlias,
		Alias:   util.DefaultAlias,

		Columns: columns,
		Imports: imports.Elements(),
	}
}

// TemplateColumn stores column info
type TemplateColumn struct {
	model.Column

	Relaxed bool

	HasTags bool
	Tag     template.HTML

	UseCustomRender bool
	CustomRender    template.HTML
}

// NewTemplateColumn creates a column for template
func NewTemplateColumn(column model.Column, options Options) TemplateColumn {
	if !options.KeepPK && column.IsPK {
		column.GoName = util.ID
	}

	if options.Relaxed {
		column.Type = model.TypeInterface
	} else {
		column.Type = fmt.Sprintf("*%s", column.GoType)
	}

	// add json tag
	tags := util.NewAnnotation()
	if options.AddJSONTag {
		tags.AddTag("json", util.Underscore(column.PGName))
	}

	return TemplateColumn{
		Relaxed: options.Relaxed,
		Column:  column,
		HasTags: tags.Len() > 0,
		Tag:     template.HTML(fmt.Sprintf("`%s`", tags.String())),
	}
}
